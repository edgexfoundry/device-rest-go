//
// Copyright (c) 2019 Intel Corporation
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
//
// CONTRIBUTORS              COMPANY
//===============================================================
// 1. Sathya Durai           HCL Technologies
// 2. Sudhamani Bijivemula   HCL Technologies
// 3. Vijay Annamalaisamy    HCL Technologies
//

package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"strings"
	"time"
)

type RestDriver struct {
	logger logger.LoggingClient
}

// Initialize performs protocol-specific initialization for the device
// service.
func (driver *RestDriver) Initialize(logger logger.LoggingClient, asyncValues chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	driver.logger = logger
	handler := NewRestHandler(sdk.RunningService(), logger, asyncValues)
	return handler.Start()
}

// Validate data that is being received as part of the http request and response
func (driver *RestDriver) validateData(resource models.DeviceResource, reading interface{}, valueType string, contentType string) (interface{}, error) {
	var err error
	castError := "failed to validate %v reading, %v"

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

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (driver *RestDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (responses []*dsModels.CommandValue, err error) {

	driver.logger.Debugf("In HandleReadCommands function")
	var reading interface{}
	var val interface{}
	var result = &dsModels.CommandValue{}
	var uri string

	// Declare response structure,
	// This will be filled with response received from HTTP GET request
	responses = make([]*dsModels.CommandValue, len(reqs))

	// To send request to any end device, first we need to know end device details.
	// Such as end device IP address, port number on which REST server is running etc.
	// First get all these details from the device file
	protocolParams, paramsExists := protocols[ED_PARAMS]
	if !paramsExists {
		return nil, fmt.Errorf("No End device parameters defined in the protocol list")
	}

	// Get end device IP address
	address, ok := protocolParams[ED_IP]
	if !ok {
		return nil, fmt.Errorf("ED_IP not found")
	}
	driver.logger.Debugf("ED_IP:%v", address)

	// Get end device port number
	port, ok := protocolParams[ED_PORT]
	if !ok {
		return nil, fmt.Errorf("ED_PORT not found")
	}
	driver.logger.Debugf("ED_PORT:%v", port)

	// Get end device URI prefix, This parameter will be empty for the end devices
	// which do not have any prefix
	uriPrefix, ok := protocolParams[ED_URI_PREFIX]
	if !ok {
		return nil, fmt.Errorf("ED_URI_PREFIX not found")
	}
	driver.logger.Debugf("ED_URI_PREFIX:%v", uriPrefix)

	// Now, we have got required end device information, its time to create GET request
	for i, req := range reqs {
		// First get device resource instance, needed during validation of the
		// response data received from the end device later.
		// RunningService returns the Service instance which is running.
		// service.DeviceResource retrieves the specific DeviceResource instance
		// from cache according to the Device name and Device Resource name
		resourceName := req.DeviceResourceName

		service := sdk.RunningService()

		deviceResource, ok := service.DeviceResource(deviceName, resourceName)
		if !ok {
			driver.logger.Errorf("Incoming reading ignored. Resource '%s' not found", resourceName)
			return nil, fmt.Errorf("Resource not found")
		}

		// Now get the query parameters received in the request.
		// These needs to be sent to the end device
		// Query parameters needs to be converted to string to append to the uri
		driver.logger.Debugf("URLRawQuery :%v", req.Attributes["urlRawQuery"])
		reqParam := fmt.Sprint(req.Attributes["urlRawQuery"])

		// Form URI from the end device parameters and request parameters and
		// query parameters. Omit uri prefix if it is empty
		if uriPrefix != "" {
			uri = "http://" + address + ":" + port + "/" + uriPrefix + "/" + req.DeviceResourceName + "?" + reqParam + ""
		} else {
			uri = "http://" + address + ":" + port + "/" + req.DeviceResourceName + "?" + reqParam + ""
		}
		driver.logger.Debugf("uri = %v", uri)

		// Now we have end device informationa and uri. This is enough to create
		// GET request. For this first create http client instance.
		// Then create http new request, this will not initiate request to end device
		client := &http.Client{}
		request, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			// handle error
			return nil, fmt.Errorf("GET request creation failed")
		}

		// Now, we have client instance and GET request instance
		// Initiate GET request to end device
		resp, err := client.Do(request)
		if err != nil {
			// handle error
			return nil, fmt.Errorf("Get request failed")
		}

		// GET request to end device success, Its time to parse the response received
		// Return immediately if status code is > 299
		// Ref: https://pkg.go.dev/net/http
		if resp.StatusCode > 299 {
			return nil, fmt.Errorf("Get response failed with status code: %v", resp.StatusCode)
		}

		// Reached here, as the success response is received. Let's get
		// response body to return as response to this read command request.
		// Close response body once read from it
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return responses, fmt.Errorf("Read command failed. Cmd:%v err:%v \n", req.DeviceResourceName, err)
		}
		driver.logger.Debugf("Resp Body=%s", body)

		// We are going to validate received content type against the expected
		// content type of device resource. For doing this get content type from
		// GET response header.  Take response body as it is if device
		// resource data type is binary or object. For other data types convert
		// response body to string to use during validation of reponse
		if deviceResource.Properties.ValueType == common.ValueTypeBinary ||
			deviceResource.Properties.ValueType == common.ValueTypeObject {
			reading = body
		} else {
			reading = string(body)
		}
		contentType := resp.Header.Get("Content-Type")

		val, err = driver.validateData(deviceResource, reading, deviceResource.Properties.ValueType, contentType)
		if err != nil {
			return nil, fmt.Errorf("Recevice response data is not valid")
		}

		// Now, we have valid response data. This needs to be sent as response to the read command request. Create a CommandValue according to the data type
		result, err = dsModels.NewCommandValue(deviceResource.Name, deviceResource.Properties.ValueType, val)
		if err != nil {
			return nil, err
		}
		result.Origin = time.Now().UnixNano()

		responses[i] = result
	}

	return responses, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the
// individual command.
func (driver *RestDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	driver.logger.Debugf("In HandleWriteCommands function")

	// Create http request variable to be used for creating new request
	var request *http.Request
	var err error
	var uri string

	// To send request to any end device, first we need to know end device details.
	// Such as end device IP address, port number on which REST server is running etc.
	// First get all these details from the device file
	protocolParams, paramsExists := protocols[ED_PARAMS]
	if !paramsExists {
		return fmt.Errorf("No End device parameters defined in the protocol list")
	}

	// Get End device address
	address, ok := protocolParams[ED_IP]
	if !ok {
		return fmt.Errorf("ED_IP not found")
	}
	driver.logger.Debugf("ED_IP:%v", address)

	// Get End device port number
	port, ok := protocolParams[ED_PORT]
	if !ok {
		return fmt.Errorf("ED_PORT not found")
	}
	driver.logger.Debugf("ED_PORT:%v", port)

	// Get End devie URI prefix, this can be empty
	uriPrefix, ok := protocolParams[ED_URI_PREFIX]
	if !ok {
		return fmt.Errorf("ED_URI_PREFIX not found")
	}
	driver.logger.Debugf("ED_URI_PREFIX:%v", uriPrefix)

	for i, req := range reqs {
		// First get device resource instance, needed during validation of the
		// data received in the write command request
		// RunningService returns the Service instance which is running
		// service.DeviceResource retrieves the specific DeviceResource instance
		// from cache according to the Device name and Device Resource name
		resourceName := req.DeviceResourceName

		service := sdk.RunningService()

		deviceResource, ok := service.DeviceResource(deviceName, resourceName)
		if !ok {
			driver.logger.Errorf("Incoming Writing ignored. Resource '%s' not found", resourceName)
			return nil
		}

		// Now get the query parameters received in the request.
		// These needs to be sent to the end device
		// Query parameters needs to be converted to string to append to the uri
		driver.logger.Debugf("URLRawQuery :%v", req.Attributes["urlRawQuery"])
		reqParam := fmt.Sprint(req.Attributes["urlRawQuery"])

		// Form URI from the end device parameters and request parameters and
		// query parameters. Omit uri prefix if it is empty
		if uriPrefix != "" {
			uri = "http://" + address + ":" + port + "/" + uriPrefix + "/" + req.DeviceResourceName + "?" + reqParam + ""
		} else {
			uri = "http://" + address + ":" + port + "/" + req.DeviceResourceName + "?" + reqParam + ""
		}

		// Its time to form payload to be sent to end device.
		// For this fisrt get the data received in the write command request
		// This data is validated against the expected value type of device resource
		// With the data and uri create new http PUT request
		// And, set the content type header for the PUT request
		driver.logger.Debugf("params = %v", params[i])
		reading := params[i].Value
		valueType := deviceResource.Properties.ValueType
		switch valueType {
		case common.ValueTypeObject:
			buf, _ := json.Marshal(reading)
			driver.logger.Debugf("body = %v", buf)
			if !json.Valid([]byte(buf)) {
				driver.logger.Debugf("Invalid JSON string")
			}

			// Create new PUT request, this will not send request to end device
			request, err = http.NewRequest(http.MethodPut, uri, bytes.NewReader(buf))
			if err != nil {
				// handle error
				return fmt.Errorf("PUT request creation failed")
			}
			// Set content type as application/json
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		case common.ValueTypeBool, common.ValueTypeString, common.ValueTypeUint8,
			common.ValueTypeUint16, common.ValueTypeUint32, common.ValueTypeUint64,
			common.ValueTypeInt8, common.ValueTypeInt16, common.ValueTypeInt32,
			common.ValueTypeInt64, common.ValueTypeFloat32, common.ValueTypeFloat64:
			// All other types
			contentType := "text/plain"
			_, err = driver.validateData(deviceResource, reading, deviceResource.Properties.ValueType, contentType)
			if err != nil {
				// handle error
				return fmt.Errorf("PUT request data is not valid")
			}
			// Create new PUT request
			request, err = http.NewRequest(http.MethodPut, uri, strings.NewReader(cast.ToString(reading)))
			if err != nil {
				// handle error
				return fmt.Errorf("PUT request creation failed")
			}
			// Set content type as text/plain
			request.Header.Set("Content-Type", "text/plain; charset=UTF-8")

		default:
			return fmt.Errorf("Unsupported value type: %v", valueType)
		}

		// Now we have created http PUT request instance with uri, and payload. This
		// is enough to initiate PUT request to end device.
		// First create new http client and initiate PUT request
		client := &http.Client{}

		resp, err := client.Do(request)
		if err != nil {
			// handle error
			return fmt.Errorf("PUT request failed")
		}

		// Htpp status codes till 299 fall under informational/ success category
		/* 1xx Informational
		   2xx Success
		   3xx Redirection
		   4xx Client Error
		   5xx Server Error
		*/
		// Return immediately if status code is > 299
		// Ref: https://pkg.go.dev/net/http
		if resp.StatusCode > 299 {
			return fmt.Errorf("PUT request failed with status code: %v", resp.StatusCode)
		}
	}

	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (driver *RestDriver) Stop(force bool) error {
	driver.logger.Debugf("RestDriver.Stop called: force=%v", force)
	// Nothing required to do for Stop
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (driver *RestDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	// Nothing required to do for AddDevice since new devices will be available
	// when data is posted to REST endpoint
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (driver *RestDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	// Nothing required to do for UpdateDevice since device update will be available
	// when data is posted to REST endpoint.
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (driver *RestDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	// Nothing required to do for RemoveDevice since removed device will no longer be available
	// when data is posted to REST endpoint.
	return nil
}
