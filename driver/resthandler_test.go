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
	"os"
	"testing"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	sdk "github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var handler *RestHandler

func TestMain(m *testing.M) {
	service := &sdk.DeviceService{}
	logger := logger.NewMockClient()
	asyncValues := make(chan<- *models.AsyncValues)

	handler = NewRestHandler(service, logger, asyncValues)
	os.Exit(m.Run())
}

func TestNewCommandValue(t *testing.T) {

	tests := []struct {
		Name          string
		Value         interface{}
		Expected      interface{}
		Type          string
		ErrorExpected bool
	}{
		{"Test A Binary", []byte{1, 0, 0, 1}, []byte{1, 0, 0, 1}, v2.ValueTypeBinary, false},
		{"Test A String JSON", `{"name" : "My JSON"}`, `{"name" : "My JSON"}`, v2.ValueTypeString, false},
		{"Test A String Text", "Random Text", "Random Text", v2.ValueTypeString, false},
		{"Test A Bool true", "true", true, v2.ValueTypeBool, false},
		{"Test A Bool false", "false", false, v2.ValueTypeBool, false},
		{"Test A Bool error", "bad", nil, v2.ValueTypeBool, true},
		{"Test A Float32+", "123.456", float32(123.456), v2.ValueTypeFloat32, false},
		{"Test A Float32-", "-123.456", float32(-123.456), v2.ValueTypeFloat32, false},
		{"Test A Float32 error", "-123.junk", nil, v2.ValueTypeFloat32, true},
		{"Test A Float64+", "456.123", 456.123, v2.ValueTypeFloat64, false},
		{"Test A Float64-", "-456.123", -456.123, v2.ValueTypeFloat64, false},
		{"Test A Float64 error", "Random", nil, v2.ValueTypeFloat64, true},
		{"Test A Uint8", "255", uint8(255), v2.ValueTypeUint8, false},
		{"Test A Uint8 error", "FF", nil, v2.ValueTypeUint8, true},
		{"Test A Uint16", "65535", uint16(65535), v2.ValueTypeUint16, false},
		{"Test A Uint16 error", "FFFF", nil, v2.ValueTypeUint16, true},
		{"Test A Uint32", "4294967295", uint32(4294967295), v2.ValueTypeUint32, false},
		{"Test A Uint32 error", "FFFFFFFF", nil, v2.ValueTypeUint32, true},
		{"Test A Uint64", "6744073709551615", uint64(6744073709551615), v2.ValueTypeUint64, false},
		{"Test A Uint64 error", "FFFFFFFFFFFFFFFF", nil, v2.ValueTypeUint64, true},
		{"Test A Int8+", "101", int8(101), v2.ValueTypeInt8, false},
		{"Test A Int8-", "-101", int8(-101), v2.ValueTypeInt8, false},
		{"Test A Int8 error", "-101.98", nil, v2.ValueTypeInt8, true},
		{"Test A Int16+", "2001", int16(2001), v2.ValueTypeInt16, false},
		{"Test A Int16-", "-2001", int16(-2001), v2.ValueTypeInt16, false},
		{"Test A Int16 error", "-FF", nil, v2.ValueTypeInt16, true},
		{"Test A Int32+", "32000", int32(32000), v2.ValueTypeInt32, false},
		{"Test A Int32-", "-32000", int32(-32000), v2.ValueTypeInt32, false},
		{"Test A Int32 error", "-32.456", nil, v2.ValueTypeInt32, true},
		{"Test A Int64+", "214748364800", int64(214748364800), v2.ValueTypeInt64, false},
		{"Test A Int64-", "-214748364800", int64(-214748364800), v2.ValueTypeInt64, false},
		{"Test A Int64 error", "-21474.99", nil, v2.ValueTypeInt64, true},
	}

	for _, currentTest := range tests {
		t.Run(currentTest.Name, func(t *testing.T) {
			cmdVal, err := handler.newCommandValue("test", currentTest.Value, currentTest.Type)
			if currentTest.ErrorExpected {
				assert.Error(t, err, "Expected an Error")
			} else {
				require.NoError(t, err, "Unexpected an Error")
				assert.Equal(t, cmdVal.Value, currentTest.Expected)
			}
		})
	}
}
