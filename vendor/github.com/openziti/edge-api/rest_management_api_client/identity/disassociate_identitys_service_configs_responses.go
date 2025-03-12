// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// __          __              _
// \ \        / /             (_)
//  \ \  /\  / /_ _ _ __ _ __  _ _ __   __ _
//   \ \/  \/ / _` | '__| '_ \| | '_ \ / _` |
//    \  /\  / (_| | |  | | | | | | | | (_| | : This file is generated, do not edit it.
//     \/  \/ \__,_|_|  |_| |_|_|_| |_|\__, |
//                                      __/ |
//                                     |___/

package identity

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openziti/edge-api/rest_model"
)

// DisassociateIdentitysServiceConfigsReader is a Reader for the DisassociateIdentitysServiceConfigs structure.
type DisassociateIdentitysServiceConfigsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DisassociateIdentitysServiceConfigsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDisassociateIdentitysServiceConfigsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewDisassociateIdentitysServiceConfigsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewDisassociateIdentitysServiceConfigsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDisassociateIdentitysServiceConfigsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewDisassociateIdentitysServiceConfigsTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 503:
		result := NewDisassociateIdentitysServiceConfigsServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDisassociateIdentitysServiceConfigsOK creates a DisassociateIdentitysServiceConfigsOK with default headers values
func NewDisassociateIdentitysServiceConfigsOK() *DisassociateIdentitysServiceConfigsOK {
	return &DisassociateIdentitysServiceConfigsOK{}
}

/* DisassociateIdentitysServiceConfigsOK describes a response with status code 200, with default header values.

Base empty response
*/
type DisassociateIdentitysServiceConfigsOK struct {
	Payload *rest_model.Empty
}

func (o *DisassociateIdentitysServiceConfigsOK) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsOK  %+v", 200, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsOK) GetPayload() *rest_model.Empty {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.Empty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDisassociateIdentitysServiceConfigsBadRequest creates a DisassociateIdentitysServiceConfigsBadRequest with default headers values
func NewDisassociateIdentitysServiceConfigsBadRequest() *DisassociateIdentitysServiceConfigsBadRequest {
	return &DisassociateIdentitysServiceConfigsBadRequest{}
}

/* DisassociateIdentitysServiceConfigsBadRequest describes a response with status code 400, with default header values.

The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information
*/
type DisassociateIdentitysServiceConfigsBadRequest struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *DisassociateIdentitysServiceConfigsBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsBadRequest  %+v", 400, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsBadRequest) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDisassociateIdentitysServiceConfigsUnauthorized creates a DisassociateIdentitysServiceConfigsUnauthorized with default headers values
func NewDisassociateIdentitysServiceConfigsUnauthorized() *DisassociateIdentitysServiceConfigsUnauthorized {
	return &DisassociateIdentitysServiceConfigsUnauthorized{}
}

/* DisassociateIdentitysServiceConfigsUnauthorized describes a response with status code 401, with default header values.

The supplied session does not have the correct access rights to request this resource
*/
type DisassociateIdentitysServiceConfigsUnauthorized struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *DisassociateIdentitysServiceConfigsUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsUnauthorized  %+v", 401, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsUnauthorized) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDisassociateIdentitysServiceConfigsNotFound creates a DisassociateIdentitysServiceConfigsNotFound with default headers values
func NewDisassociateIdentitysServiceConfigsNotFound() *DisassociateIdentitysServiceConfigsNotFound {
	return &DisassociateIdentitysServiceConfigsNotFound{}
}

/* DisassociateIdentitysServiceConfigsNotFound describes a response with status code 404, with default header values.

The requested resource does not exist
*/
type DisassociateIdentitysServiceConfigsNotFound struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *DisassociateIdentitysServiceConfigsNotFound) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsNotFound  %+v", 404, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsNotFound) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDisassociateIdentitysServiceConfigsTooManyRequests creates a DisassociateIdentitysServiceConfigsTooManyRequests with default headers values
func NewDisassociateIdentitysServiceConfigsTooManyRequests() *DisassociateIdentitysServiceConfigsTooManyRequests {
	return &DisassociateIdentitysServiceConfigsTooManyRequests{}
}

/* DisassociateIdentitysServiceConfigsTooManyRequests describes a response with status code 429, with default header values.

The resource requested is rate limited and the rate limit has been exceeded
*/
type DisassociateIdentitysServiceConfigsTooManyRequests struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *DisassociateIdentitysServiceConfigsTooManyRequests) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsTooManyRequests  %+v", 429, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsTooManyRequests) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDisassociateIdentitysServiceConfigsServiceUnavailable creates a DisassociateIdentitysServiceConfigsServiceUnavailable with default headers values
func NewDisassociateIdentitysServiceConfigsServiceUnavailable() *DisassociateIdentitysServiceConfigsServiceUnavailable {
	return &DisassociateIdentitysServiceConfigsServiceUnavailable{}
}

/* DisassociateIdentitysServiceConfigsServiceUnavailable describes a response with status code 503, with default header values.

The request could not be completed due to the server being busy or in a temporarily bad state
*/
type DisassociateIdentitysServiceConfigsServiceUnavailable struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *DisassociateIdentitysServiceConfigsServiceUnavailable) Error() string {
	return fmt.Sprintf("[DELETE /identities/{id}/service-configs][%d] disassociateIdentitysServiceConfigsServiceUnavailable  %+v", 503, o.Payload)
}
func (o *DisassociateIdentitysServiceConfigsServiceUnavailable) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *DisassociateIdentitysServiceConfigsServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
