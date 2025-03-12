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

package terminator

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new terminator API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for terminator API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	CreateTerminator(params *CreateTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateTerminatorCreated, error)

	DeleteTerminator(params *DeleteTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteTerminatorOK, error)

	DetailTerminator(params *DetailTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DetailTerminatorOK, error)

	ListTerminators(params *ListTerminatorsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListTerminatorsOK, error)

	PatchTerminator(params *PatchTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchTerminatorOK, error)

	UpdateTerminator(params *UpdateTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*UpdateTerminatorOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  CreateTerminator creates a terminator resource

  Create a terminator resource. Requires admin access.
*/
func (a *Client) CreateTerminator(params *CreateTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*CreateTerminatorCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateTerminatorParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "createTerminator",
		Method:             "POST",
		PathPattern:        "/terminators",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &CreateTerminatorReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*CreateTerminatorCreated)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for createTerminator: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  DeleteTerminator deletes a terminator

  Delete a terminator by id. Requires admin access.
*/
func (a *Client) DeleteTerminator(params *DeleteTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DeleteTerminatorOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteTerminatorParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "deleteTerminator",
		Method:             "DELETE",
		PathPattern:        "/terminators/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteTerminatorReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteTerminatorOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for deleteTerminator: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  DetailTerminator retrieves a single terminator

  Retrieves a single terminator by id. Requires admin access.
*/
func (a *Client) DetailTerminator(params *DetailTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DetailTerminatorOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDetailTerminatorParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "detailTerminator",
		Method:             "GET",
		PathPattern:        "/terminators/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DetailTerminatorReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DetailTerminatorOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for detailTerminator: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  ListTerminators lists terminators

  Retrieves a list of terminator resources; supports filtering, sorting, and pagination. Requires admin access.

*/
func (a *Client) ListTerminators(params *ListTerminatorsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListTerminatorsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListTerminatorsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listTerminators",
		Method:             "GET",
		PathPattern:        "/terminators",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListTerminatorsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListTerminatorsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for listTerminators: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PatchTerminator updates the supplied fields on a terminator

  Update the supplied fields on a terminator. Requires admin access.
*/
func (a *Client) PatchTerminator(params *PatchTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*PatchTerminatorOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchTerminatorParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "patchTerminator",
		Method:             "PATCH",
		PathPattern:        "/terminators/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PatchTerminatorReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PatchTerminatorOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for patchTerminator: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  UpdateTerminator updates all fields on a terminator

  Update all fields on a terminator by id. Requires admin access.
*/
func (a *Client) UpdateTerminator(params *UpdateTerminatorParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*UpdateTerminatorOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateTerminatorParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "updateTerminator",
		Method:             "PUT",
		PathPattern:        "/terminators/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &UpdateTerminatorReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*UpdateTerminatorOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for updateTerminator: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
