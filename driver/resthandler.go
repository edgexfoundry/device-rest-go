//
// Copyright (c) 2019 Intel Corporation
// Copyright (c) 2021-2023 IOTech Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edgexfoundry/device-sdk-go/v3/pkg/interfaces"
	"github.com/edgexfoundry/device-sdk-go/v3/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/common"
	model "github.com/edgexfoundry/go-mod-core-contracts/v3/models"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

const (
	apiResourceRoute  = common.ApiBase + "/resource/:deviceName/:resourceName"
	handlerContextKey = "RestHandler"
)

type RestHandler struct {
	service     interfaces.DeviceServiceSDK
	logger      logger.LoggingClient
	asyncValues chan<- *models.AsyncValues
}

func NewRestHandler(sdk interfaces.DeviceServiceSDK) *RestHandler {
	handler := RestHandler{
		service:     sdk,
		logger:      sdk.LoggingClient(),
		asyncValues: sdk.AsyncValuesChannel(),
	}

	return &handler
}

func (handler RestHandler) Start() error {
	if err := handler.service.AddCustomRoute(apiResourceRoute, interfaces.Authenticated, handler.addContext(deviceHandler), http.MethodPost); err != nil {
		return fmt.Errorf("unable to add required route: %s: %s", apiResourceRoute, err.Error())
	}

	handler.logger.Infof("Route %s added.", apiResourceRoute)

	return nil
}

func (handler RestHandler) addContext(next echo.HandlerFunc) echo.HandlerFunc {
	// Add the context with the handler so the endpoint handling code can get back to this handler
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), handlerContextKey, handler) //nolint
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func (handler RestHandler) processAsyncRequest(c echo.Context) error {
	deviceName := c.Param(common.DeviceName)
	resourceName := c.Param(common.ResourceName)

	handler.logger.Debugf("Received POST for Device=%s Resource=%s", deviceName, resourceName)

	_, err := handler.service.GetDeviceByName(deviceName)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Device '%s' not found", deviceName)
		return c.String(http.StatusNotFound, fmt.Sprintf("Device '%s' not found", deviceName))
	}

	deviceResource, ok := handler.service.DeviceResource(deviceName, resourceName)
	if !ok {
		handler.logger.Errorf("Incoming reading ignored. Resource '%s' not found", resourceName)
		return c.String(http.StatusNotFound, fmt.Sprintf("Resource '%s' not found", resourceName))
	}

	contentType := c.Request().Header.Get(common.ContentType)

	var reading interface{}

	data, err := handler.readBody(c.Request())
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to read request body: %s", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	if deviceResource.Properties.ValueType == common.ValueTypeBinary || deviceResource.Properties.ValueType == common.ValueTypeObject {
		reading = data
	} else {
		reading = string(data)
	}

	value, err := validateCommandValue(deviceResource, reading, deviceResource.Properties.ValueType, contentType)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to validate Command Value for Device=%s Command=%s: %s",
			deviceName, resourceName, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	result, err := models.NewCommandValue(deviceResource.Name, deviceResource.Properties.ValueType, value)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to create Command Value for Device=%s Command=%s: %s",
			deviceName, resourceName, err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}
	result.Origin = time.Now().UnixNano()

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{result},
	}

	handler.logger.Debugf("Incoming reading received: Device=%s Resource=%s", deviceName, resourceName)

	handler.asyncValues <- asyncValues

	return nil
}

func (handler RestHandler) readBody(request *http.Request) ([]byte, error) {
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("no request body provided")
	}

	return body, nil
}

func deviceHandler(c echo.Context) error {
	handler, ok := c.Request().Context().Value(handlerContextKey).(RestHandler)
	if !ok {
		return c.String(http.StatusBadRequest, "Bad context pass to handler")
	}

	return handler.processAsyncRequest(c)
}

func validateCommandValue(resource model.DeviceResource, reading interface{}, valueType string, contentType string) (interface{}, error) {
	var err error
	castError := "failed to parse %v reading, %v"

	var val interface{}
	switch valueType {
	case common.ValueTypeBinary:
		var ok bool
		val, ok = reading.([]byte)
		if !ok {
			return nil, fmt.Errorf(castError, resource.Name, "not []byte")
		}
		if contentType != resource.Properties.MediaType {
			return nil, fmt.Errorf("wrong Content-Type: expected '%s' but received '%s'", resource.Properties.MediaType, contentType)
		}
	case common.ValueTypeObject:
		if contentType != common.ContentTypeJSON {
			return nil, fmt.Errorf("wrong Content-Type: expected '%s' but received '%s'", common.ContentTypeJSON, contentType)
		}

		data, ok := reading.([]byte)
		if !ok {
			return nil, fmt.Errorf(castError, resource.Name, "not []byte")
		}

		val = map[string]interface{}{}
		if err := json.Unmarshal(data, &val); err != nil {
			return nil, errors.New("unable to marshal JSON data to type Object")
		}
	case common.ValueTypeBool:
		val, err = cast.ToBoolE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case common.ValueTypeString:
		val, err = cast.ToStringE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
	case common.ValueTypeUint8:
		val, err = cast.ToUint8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkUintValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeUint16:
		val, err = cast.ToUint16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkUintValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeUint32:
		val, err = cast.ToUint32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkUintValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeUint64:
		val, err = cast.ToUint64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkUintValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeInt8:
		val, err = cast.ToInt8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkIntValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeInt16:
		val, err = cast.ToInt16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkIntValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeInt32:
		val, err = cast.ToInt32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkIntValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeInt64:
		val, err = cast.ToInt64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkIntValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeFloat32:
		val, err = cast.ToFloat32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkFloatValueRange(valueType, val); err != nil {
			return nil, err
		}
	case common.ValueTypeFloat64:
		val, err = cast.ToFloat64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resource.Name, err)
		}
		if err := checkFloatValueRange(valueType, val); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("return result fail, unsupported value type: %v", valueType)
	}

	return val, nil
}

func checkUintValueRange(valueType string, val interface{}) error {
	switch valueType {
	case common.ValueTypeUint8:
		_, ok := val.(uint8)
		if ok {
			return nil
		}
	case common.ValueTypeUint16:
		_, ok := val.(uint16)
		if ok {
			return nil
		}
	case common.ValueTypeUint32:
		_, ok := val.(uint32)
		if ok {
			return nil
		}
	case common.ValueTypeUint64:
		_, ok := val.(uint64)
		if ok {
			return nil
		}
	}
	return fmt.Errorf("value %v for %s type is out of range", val, valueType)
}

func checkIntValueRange(valueType string, val interface{}) error {
	switch valueType {
	case common.ValueTypeInt8:
		_, ok := val.(int8)
		if ok {
			return nil
		}
	case common.ValueTypeInt16:
		_, ok := val.(int16)
		if ok {
			return nil
		}
	case common.ValueTypeInt32:
		_, ok := val.(int32)
		if ok {
			return nil
		}
	case common.ValueTypeInt64:
		_, ok := val.(int64)
		if ok {
			return nil
		}
	}
	return fmt.Errorf("value %v for %s type is out of range", val, valueType)
}

func checkFloatValueRange(valueType string, val interface{}) error {
	switch valueType {
	case common.ValueTypeFloat32:
		_, ok := val.(float32)
		if ok {
			return nil
		}
	case common.ValueTypeFloat64:
		_, ok := val.(float64)
		if ok {
			return nil
		}
	}

	return fmt.Errorf("value %v for %s type is out of range", val, valueType)
}
