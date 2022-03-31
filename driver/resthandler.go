//
// Copyright (c) 2019 Intel Corporation
// Copyright (c) 2021 IOTech Ltd
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
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	model "github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
)

const (
	apiResourceRoute  = common.ApiBase + "/resource/{" + common.DeviceName + "}/{" + common.ResourceName + "}"
	handlerContextKey = "RestHandler"
)

type RestHandler struct {
	service     *sdk.DeviceService
	logger      logger.LoggingClient
	asyncValues chan<- *models.AsyncValues
}

func NewRestHandler(service *sdk.DeviceService, logger logger.LoggingClient, asyncValues chan<- *models.AsyncValues) *RestHandler {
	handler := RestHandler{
		service:     service,
		logger:      logger,
		asyncValues: asyncValues,
	}

	return &handler
}

func (handler RestHandler) Start() error {
	if err := handler.service.AddRoute(apiResourceRoute, handler.addContext(deviceHandler), http.MethodPost); err != nil {
		return fmt.Errorf("unable to add required route: %s: %s", apiResourceRoute, err.Error())
	}

	handler.logger.Infof("Route %s added.", apiResourceRoute)

	return nil
}

func (handler RestHandler) addContext(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	// Add the context with the handler so the endpoint handling code can get back to this handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), handlerContextKey, handler) //nolint
		next(w, r.WithContext(ctx))
	})
}

func (handler RestHandler) processAsyncRequest(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	deviceName := vars[common.DeviceName]
	resourceName := vars[common.ResourceName]

	handler.logger.Debugf("Received POST for Device=%s Resource=%s", deviceName, resourceName)

	_, err := handler.service.GetDeviceByName(deviceName)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Device '%s' not found", deviceName)
		http.Error(writer, fmt.Sprintf("Device '%s' not found", deviceName), http.StatusNotFound)
		return
	}

	deviceResource, ok := handler.service.DeviceResource(deviceName, resourceName)
	if !ok {
		handler.logger.Errorf("Incoming reading ignored. Resource '%s' not found", resourceName)
		http.Error(writer, fmt.Sprintf("Resource '%s' not found", resourceName), http.StatusNotFound)
		return
	}

	contentType := request.Header.Get(common.ContentType)

	var reading interface{}

	data, err := handler.readBody(request)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to read request body: %s", err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if deviceResource.Properties.ValueType == common.ValueTypeBinary || deviceResource.Properties.ValueType == common.ValueTypeObject {
		reading = data
	} else {
		reading = string(data)
	}

	value, err := handler.newCommandValue(deviceResource, reading, deviceResource.Properties.ValueType, contentType)
	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to create Command Value for Device=%s Command=%s: %s",
			deviceName, resourceName, err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	asyncValues := &models.AsyncValues{
		DeviceName:    deviceName,
		CommandValues: []*models.CommandValue{value},
	}

	handler.logger.Debugf("Incoming reading received: Device=%s Resource=%s", deviceName, resourceName)

	handler.asyncValues <- asyncValues
}

func (handler RestHandler) readBody(request *http.Request) ([]byte, error) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("no request body provided")
	}

	return body, nil
}

func deviceHandler(writer http.ResponseWriter, request *http.Request) {
	handler, ok := request.Context().Value(handlerContextKey).(RestHandler)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write([]byte("Bad context pass to handler"))
		if err != nil {
			handler.logger.Debugf("problem in writer of byte array: '%s'", err.Error())
		}
		return
	}

	handler.processAsyncRequest(writer, request)
}

func (handler RestHandler) newCommandValue(resource model.DeviceResource, reading interface{}, valueType string, contentType string) (*models.CommandValue, error) {
	var err error
	var result = &models.CommandValue{}
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

	result, err = models.NewCommandValue(resource.Name, valueType, val)
	if err != nil {
		return nil, err
	}

	result.Origin = time.Now().UnixNano()
	return result, nil
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
		valFloat := val.(float32)
		if math.Abs(float64(valFloat)) >= math.SmallestNonzeroFloat32 && math.Abs(float64(valFloat)) <= math.MaxFloat32 {
			return nil
		}
	case common.ValueTypeFloat64:
		valFloat := val.(float64)
		if math.Abs(valFloat) >= math.SmallestNonzeroFloat64 && math.Abs(valFloat) <= math.MaxFloat64 {
			return nil
		}
	}

	return fmt.Errorf("value %v for %s type is out of range", val, valueType)
}
