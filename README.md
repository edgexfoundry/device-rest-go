# device-rest-go
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/device-rest-go/branch/main/graph/badge.svg?token=fmbJjqOyk4)](https://codecov.io/gh/edgexfoundry/device-rest-go) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-rest-go)](https://goreportcard.com/report/github.com/edgexfoundry/device-rest-go) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-rest-go/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-rest-go)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-rest-go) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-rest-go-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/commits)

> **Warning**  
> The **main** branch of this repository contains work-in-progress development code for the upcoming release, and is **not guaranteed to be stable or working**.
> It is only compatible with the [main branch of edgex-compose](https://github.com/edgexfoundry/edgex-compose) which uses the Docker images built from the **main** branch of this repo and other repos.
>
> **The source for the latest release can be found at [Releases](https://github.com/edgexfoundry/device-rest-go/releases).**

EdgeX device service for REST protocol

This device service supports two-way communication for commanding. REST endpoints provides easy way for 3'rd party applications, such as Point of Sale, CV Analytics, etc., to push async data into EdgeX via the REST protocol. Commands allows EdgeX to send GET and PUT request to end device. 

## Runtime Prerequisite    

- core-data
  - Redis DB
- core-metadata
- core-command


## Build with NATS Messaging
Currently, the NATS Messaging capability (NATS MessageBus) is opt-in at build time.
This means that the published Docker image and Snaps do not include the NATS messaging capability.

The following make commands will build the local binary or local Docker image with NATS messaging
capability included.
```makefile
make build-nats
make docker-nats
```

The locally built Docker image can then be used in place of the published Docker image in your compose file.
See [Compose Builder](https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder#gen) `nat-bus` option to generate compose file for NATS and local dev images.

## REST Endpoints

This device service creates the additional parametrized `REST` endpoints

### Async

End device will be pushing async data to these endpoints.

```
/api/v2/resource/{deviceName}/{resourceName}
```

- `deviceName` refers to a `device` managed by the REST device service.
- `resourceName`refers to the `device resource` defined in the `device profile` associated with the given `deviceName` .

The data, `text` or `binary`,  posted to this endpoint is type validated and type casted (text data only) to the type defined by the specified `device resource`. The resulting value is then sent into EdgeX via the Device SDK's `async values` channel. 

> *Note: When binary data is used the EdgeX event/reading is `CBOR` encoded by the `Device SDK` and the binary value in the reading is`NOT` be stored in the database by `Core Data`. The `CBOR` encoded event/reading, with the binary value, is published to the `Message Bus` for `Application Services` to consume*

> *Note: All non-binary data is consumed as text. The text is casted to the specific type of the specified `device resource` once it passes type validation.*

### Commands

Use the following core command service API format to execute commands against the rest device service.

```
/api/v2/device/name/{deviceName}/{resourceName}
```

- `deviceName` refers to a `device` managed by the REST device service.
- `resourceName`refers to the `device resource` defined in the `device profile` associated with the given `deviceName` .

Upon receiving the command, device service has to forward this command to end device. For this device service reads the end device protocol parameters from the device list protocol properties. Then constructs uri based on the end device protocol parameters, query parameters & device resource.

For GET command creates new http GET request using uri and send the GET request to end device. The response received from end device is type validated and sent as response to the GET command.

For PUT command, the PUT data, `text` or `binary`,  posted to this endpoint is type validated and type casted (text data only) to the type defined by the specified `device resource`. Then creates new http PUT request using uri, validated PUT data and then send the http PUT request to end device. The end device response status code is sent in response to the PUT command. 

## Configuration

To use this device service, device profile and device file needs to be configured. 
This device service use the standard configuration defined by the **Device SDK**. 

The `DeviceList` configuration is standard except that it is mandatory to provide end device parameters in the `DeviceList.Protocols.EndDevice_Params` structure for 2way-rest-device functionality. The following is a sample `DeviceList` that works with the sample device profiles referenced below. `path` parameter is optional.

```yaml
deviceList:
  - name: sample-json
    profileName: sample-json
    description: RESTful Device that sends in JSON data
    labels:
      - rest
      - json
    protocols:
      other: {}
  - name: sample-image
    profileName: sample-image
    description: RESTful Device that sends in binary image data
    labels:
      - rest
      - binary
      - image
    protocols:
      other: {}
  - name: sample-numeric
    profileName: sample-numeric
    description: RESTful Device that sends in numeric data
    labels:
      - rest
      - numeric
      - float
      - int
    protocols:
      other: {}
  - name: 2way-rest-device
    profileName: sample-2way-rest-device
    description: RESTful Device that sends data
    labels:
      - rest
      - 2way-rest-device
    protocols:
      REST:
        Host: 127.0.0.1
        Port: '5000'
        Path: api
    autoEvents:
      - Interval: 20s
        OnChange: false
        SourceName: json

```

## Device Profile

As with all device services the `device profile` is where the **Device Name**, **Device Resources** and **Device Commands** are define. The parameterized REST endpoints described above references these definitions. Each `Device` has it's own device profile. There are four sample device profiles that define the devices referenced in the above sample configuration.

- [**sample-image-device**](cmd/res/profiles/sample-image-device.yaml)
- [**sample-json-device**](cmd/res/profiles/sample-json-device.yaml)
- [**sample-numeric-device**](cmd/res/profiles/sample-numeric-device.yaml)
- [**sample-2way-rest-device**](cmd/res/profiles/sample-2way-rest-device.yaml)

## Testing/Simulation

### Async

The best way to test this service with simulated data is to use **PostMan** to send data to the following endpoints defined for the above device profiles.

- http://localhost:59986/api/v2/resource/sample-image/jpeg

  - POSTing a JPEG binary image file will result in the `BinaryValue` of the `Reading` being set to the JPEG image data posted.
  - Example test JPEG to post:
    - Select any JPEG file from your computer or the internet

- http://localhost:59986/api/v2/resource/sample-image/png

  - POSTing a PNG binary image file will result in the `BinaryValue` of the `Reading` being set to the PNG image data posted.
  - Example test PNG to post:
    - Select any PNG file from your computer or the internet

- http://localhost:59986/api/v2/resource/sample-json/json

  - POSTing a JSON string value will result in the  `Value` of the `Reading` being set to the JSON string value posted.

    *Note: Currently there isn't a JSON data type, thus there is no validation that the string value is valid JSON. It is up to the Application Service using the JSON to first validate it.*
    
  - Example test JSON value to post:

    ```json
    {
        "id" : "1234",
        "name" : "test data",
        "payload" : "test payload"
    }
    ```

- http://localhost:59986/api/v2/resource/sample-numeric/int
  - POSTing a text integer value will result in the  `Value` of the `Reading` being set to the string representation of the value as an `Int64`. The POSTed value is verified to be a valid `Int64` value. 
  
  - A 400 error will be returned if the POSTed value fails the `Int64` type verification.
  
  - Example test `int` value to post:
  
    ```
    1001
    ```
  
- http://localhost:59986/api/v2/resource/sample-numeric/float
  - POSTing a text float value will result in the  `Value` of the `Reading` being set to the string representation of the value as an `Float64`. The POSTed value is verified to be a valid `Float64` value. 
  
  - A 400 error will be returned if the POSTed value fails the `Float64` type verification.
  
  - Example test `float` value to post:
  
    ```
    500.568
    ```

### Commands

This device service supports commanding functionality with a sample profile for the data types as shown in below table.

| Data Type | GET                | PUT                |
|---------	|------------------- |------------------- |
| Binary	| :heavy_check_mark: | :X				  |
| Object	| :heavy_check_mark: | :heavy_check_mark: |
| Bool	   	| :heavy_check_mark: | :heavy_check_mark: |
| String	| :heavy_check_mark: | :heavy_check_mark: |
| Uint8		| :heavy_check_mark: | :heavy_check_mark: |
| Uint16	| :heavy_check_mark: | :heavy_check_mark: |
| Uint32	| :heavy_check_mark: | :heavy_check_mark: |
| Uint64	| :heavy_check_mark: | :heavy_check_mark: |
| Int8	    | :heavy_check_mark: | :heavy_check_mark: |
| Int16	    | :heavy_check_mark: | :heavy_check_mark: |
| Int32	    | :heavy_check_mark: | :heavy_check_mark: |
| Int64	    | :heavy_check_mark: | :heavy_check_mark: |
| Float32   | :heavy_check_mark: | :heavy_check_mark: |
| Float64   | :heavy_check_mark: | :heavy_check_mark: |

Using `curl` command-line utility or `PostMan` we can send GET/PUT request to EdgeX. These commands are explained in `GET Command` section below. End device can be anything, For example `nodejs based REST emulator` is used as end device for testing commanding functionaity of the REST device service. Example end device code is mentioned in `End Device` section below.

**End Device**

Example end device code using `nodejs` is as shown below. This example code is having endpoint for `int8` resource. To test GET/PUT commands for other resources, this code needs to be expanded in the same way for other device resources also.

```
///////////////////BUILD AND RUN INSTRUCTIONS/////////////////////
// Install node, npm, express module in target machine
// Run using "node end-device.js"
/////////////////////////////////////////////////////////////////

var express = require('express');
var bodyParser = require('body-parser')
var app = express();

var textParser = bodyParser.text({type: '*/*'})

//-128 to 127
var int8 = 111

// GET int8
app.get('/api/int8', function (req, res) {
console.log("Get int8 request");
res.end(int8.toString());
})

// PUT int8
app.put('/api/int8', textParser, function (req, res) {
console.log("Put int8 request");
console.log(req.body);
int8 = req.body;
res.end(int8);
})

var server = app.listen(5000, function () {
var host = server.address().address
var port = server.address().port
console.log("Server listening at http://%s:%s", host, port)
})
```

**GET Command**

Example GET request to `int8` device resource using curl command-line utility is as shown below.
```
   $ curl --request GET http://localhost:59882/api/v2/device/name/2way-rest-device/int8
```
Example GET request to `int8` device resource using **PostMan** is as shown below.

- http://localhost:59882/api/v2/device/name/2way-rest-device/int8

`2way-rest-device` is the device name as defined in the device file.
Example expected success response from the end device is as shown below.

```
   {
   "apiVersion" : "v2",
   "event" : {
      "apiVersion" : "v2",
      "deviceName" : "2way-rest-device",
      "id" : "46baf3d5-98fd-4073-b52e-801660b01ce6",
      "origin" : 1670506568209119757,
      "profileName" : "sample-2way-rest-device",
      "readings" : [
         {
            "deviceName" : "2way-rest-device",
            "id" : "c7d4d4fe-13f5-423a-8d62-0e57f8dbc063",
            "origin" : 1670506568209111164,
            "profileName" : "sample-2way-rest-device",
            "resourceName" : "int8",
            "value" : "111",
            "valueType" : "Int8"
         }
      ],
      "sourceName" : "int8"
   },
   "statusCode" : 200
   } 
```


**PUT Command**

Example PUT request to `int8` device resource using curl command-line utility is as shown below.
```
   $ curl -i -X PUT -H "Content-Type: application/json" -d '{"int8":12}' http://localhost:59882/api/v2/device/name/2way-rest-device/int8
```
Example PUT request to `int8` device resource using **PostMan** is as shown below.
- http://localhost:59882/api/v2/device/name/2way-rest-device/int8
  - PUTting a text integer value will result in the  `Value` of the `Command` being set to the string representation of the value as an `Int8`. The PUT value is verified to be a valid `Int8` value. 
  
  - A 400 error will be returned if the PUTted value fails the `Int8` type verification.
  
  - Example test `int8` value to post:
  
    ```
    12
    ```

`2way-rest-device` is the device name as defined in the device list.

Example expected success response from the end device is as shown below.
```
HTTP/1.1 200 OK
Content-Type: application/json
X-Correlation-Id: d208c432-0ee4-4d7e-b819-378bec45cbf6
Date: Thu, 08 Dec 2022 14:02:14 GMT
Content-Length: 37

{"apiVersion":"v2","statusCode":200}
```

## AutoEvents

Auto events are supported for the resources mentioned in the device profile for example `int8` device resource.
To enable autoevents functionality, It is mandatory to provide `DeviceList.AutoEvents` structure in the device file. Reference autoevents configuration is mentioned in the `Configuration` section of REST client.

## Build Instructions

1. Clone the device-rest-go repo with the following command:

        git clone https://github.com/edgexfoundry/device-rest-go.git

2. Build a docker image by using the following command:

        make docker

3. Alternatively the device service can be built natively:

        make build
        
## Packaging

This component is packaged as docker image and snap.

For docker, please refer to the [Dockerfile] and [Docker Compose Builder] scripts.

For the snap, refer to the [snap] directory.

[Dockerfile]: Dockerfile
[Docker Compose Builder]: https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder
[snap]: snap
