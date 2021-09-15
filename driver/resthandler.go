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
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
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
		ctx := context.WithValue(r.Context(), handlerContextKey, handler)
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

	if deviceResource.Properties.MediaType != "" {
		contentType := request.Header.Get(common.ContentType)

		handler.logger.Debugf("Content Type is '%s' & Media Type is '%s' and Type is '%s'",
			contentType, deviceResource.Properties.MediaType, deviceResource.Properties.ValueType)

		if contentType != deviceResource.Properties.MediaType {
			handler.logger.Errorf("Incoming reading ignored. Content Type '%s' doesn't match %s resource's Media Type '%s'",
				contentType, resourceName, deviceResource.Properties.MediaType)

			http.Error(writer, "Wrong Content-Type", http.StatusBadRequest)
			return
		}
	}

	var reading interface{}
	if deviceResource.Properties.ValueType == common.ValueTypeBinary {
		reading, err = handler.readBodyAsBinary(writer, request)
	} else {
		reading, err = handler.readBodyAsString(writer, request)
	}

	if err != nil {
		handler.logger.Errorf("Incoming reading ignored. Unable to read request body: %s", err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := handler.newCommandValue(resourceName, reading, deviceResource.Properties.ValueType)
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

func (handler RestHandler) readBodyAsString(writer http.ResponseWriter, request *http.Request) (string, error) {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	if len(body) == 0 {
		return "", fmt.Errorf("no request body provided")
	}

	return string(body), nil
}

func (handler RestHandler) readBodyAsBinary(writer http.ResponseWriter, request *http.Request) ([]byte, error) {
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
		writer.Write([]byte("Bad context pass to handler"))
		return
	}

	handler.processAsyncRequest(writer, request)
}

func (handler RestHandler) newCommandValue(resourceName string, reading interface{}, valueType string) (*models.CommandValue, error) {
	var err error
	var result = &models.CommandValue{}
	castError := "fail to parse %v reading, %v"

	if !checkValueInRange(valueType, reading) {
		err = fmt.Errorf("parse reading fail. Reading %v is out of the value type(%v)'s range", reading, valueType)
		handler.logger.Error(err.Error())
		return result, err
	}

	var val interface{}
	switch valueType {
	case common.ValueTypeBinary:
		var ok bool
		val, ok = reading.([]byte)
		if !ok {
			return nil, fmt.Errorf(castError, resourceName, "not []byte")
		}
	case common.ValueTypeBool:
		val, err = cast.ToBoolE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeString:
		val, err = cast.ToStringE(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeUint8:
		val, err = cast.ToUint8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeUint16:
		val, err = cast.ToUint16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeUint32:
		val, err = cast.ToUint32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeUint64:
		val, err = cast.ToUint64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeInt8:
		val, err = cast.ToInt8E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeInt16:
		val, err = cast.ToInt16E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeInt32:
		val, err = cast.ToInt32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeInt64:
		val, err = cast.ToInt64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeFloat32:
		val, err = cast.ToFloat32E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	case common.ValueTypeFloat64:
		val, err = cast.ToFloat64E(reading)
		if err != nil {
			return nil, fmt.Errorf(castError, resourceName, err)
		}
	default:
		err = fmt.Errorf("return result fail, none supported value type: %v", valueType)
	}

	result, err = models.NewCommandValue(resourceName, valueType, val)
	if err != nil {
		return nil, err
	}

	result.Origin = time.Now().UnixNano()
	return result, nil
}

func checkValueInRange(valueType string, reading interface{}) bool {
	isValid := false

	if valueType == common.ValueTypeString || valueType == common.ValueTypeBool || valueType == common.ValueTypeBinary {
		return true
	}

	if valueType == common.ValueTypeInt8 || valueType == common.ValueTypeInt16 ||
		valueType == common.ValueTypeInt32 || valueType == common.ValueTypeInt64 {
		val := cast.ToInt64(reading)
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == common.ValueTypeUint8 || valueType == common.ValueTypeUint16 ||
		valueType == common.ValueTypeUint32 || valueType == common.ValueTypeUint64 {
		val := cast.ToUint64(reading)
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == common.ValueTypeFloat32 || valueType == common.ValueTypeFloat64 {
		val := cast.ToFloat64(reading)
		isValid = checkFloatValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType string, val uint64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeUint8:
		if val <= math.MaxUint8 {
			isValid = true
		}
	case common.ValueTypeUint16:
		if val <= math.MaxUint16 {
			isValid = true
		}
	case common.ValueTypeUint32:
		if val <= math.MaxUint32 {
			isValid = true
		}
	case common.ValueTypeUint64:
		maxiMum := uint64(math.MaxUint64)
		if val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType string, val int64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeInt8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case common.ValueTypeInt16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case common.ValueTypeInt32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case common.ValueTypeInt64:
		isValid = true
	}
	return isValid
}

func checkFloatValueRange(valueType string, val float64) bool {
	var isValid = false
	switch valueType {
	case common.ValueTypeFloat32:
		if math.Abs(val) >= math.SmallestNonzeroFloat32 && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case common.ValueTypeFloat64:
		if math.Abs(val) >= math.SmallestNonzeroFloat64 && math.Abs(val) <= math.MaxFloat64 {
			isValid = true
		}
	}
	return isValid
}
