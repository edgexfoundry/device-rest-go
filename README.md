# device-rest-go
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/device-rest-go/branch/main/graph/badge.svg?token=fmbJjqOyk4)](https://codecov.io/gh/edgexfoundry/device-rest-go) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-rest-go)](https://goreportcard.com/report/github.com/edgexfoundry/device-rest-go) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-rest-go/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-rest-go)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-rest-go) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-rest-go-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/commits)

EdgeX device service for REST protocol

This device service supports two-way communication. REST server provides easy way for 3'rd party applications, such as Point of Sale, CV Analytics, etc., to push data into EdgeX via the REST protocol. REST client allows EdgeX to send GET and PUT request to end device. 

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

## REST Server

This device service creates the additional parametrized `REST` endpoint:
End device will be pushing data to these endpoints. For this end device will be running with REST client.

```
/api/v2/resource/{deviceName}/{resourceName}
```

- `deviceName` refers to a `device` managed by the REST device service.
- `resourceName`refers to the `device resource` defined in the `device profile` associated with the given `deviceName` .

The data, `text` or `binary`,  posted to this endpoint is type validated and type casted (text data only) to the type defined by the specified `device resource`. The resulting value is then sent into EdgeX via the Device SDK's `async values` channel. 

> *Note: When binary data is used the EdgeX event/reading is `CBOR` encoded by the `Device SDK` and the binary value in the reading is`NOT` be stored in the database by `Core Data`. The `CBOR` encoded event/reading, with the binary value, is published to the `Message Bus` for `Application Services` to consume*

> *Note: All non-binary data is consumed as text. The text is casted to the specific type of the specified `device resource` once it passes type validation.*

### Configuration

This device service use the standard configuration defined by the **Device SDK**. 

The `DeviceList` configuration is standard except that the `DeviceList.Protocols` can be empty. The following is a sample `DeviceList` that works with the sample device profiles referenced below.

```toml
[[DeviceList]]
  Name = "sample-json"
  ProfileName = "sample-json"
  Description = "RESTful Device that sends in JSON data"
  Labels = [ "rest", "json" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]
[[DeviceList]]
  Name = "sample-image"
  ProfileName = "sample-image"
  Description = "RESTful Device that sends in binary image data"
  Labels = [ "rest", "binary", "image" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]    
[[DeviceList]]
  Name = "sample-numeric"
  ProfileName = "sample-numeric"
  Description = "RESTful Device that sends in numeric data"
  Labels = [ "rest", "numeric", "float", "int" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]
```

### Device Profile

As with all device services the `device profile` is where the **Device Name**, **Device Resources** and **Device Commands** are define. The parameterized REST endpoint described above references these definitions. Each `Device` has it's own device profile. There are three sample device profiles that define the devices referenced in the above sample configuration.

- [**sample-image-device**](cmd/res/profiles/sample-image-device.yaml)
- [**sample-json-device**](cmd/res/profiles/sample-json-device.yaml)
- [**sample-numeric-device**](cmd/res/profiles/sample-numeric-device.yaml)

### Testing/Simulation of REST server

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

## REST Client

This device service allows EdgeX to send GET/PUT request to end device.
End device will respond to GET and PUT requests received from EdgeX. For this end device will be running REST server.

### Configuration

To use this device service, device profile and device file needs to be configured. 
This device service use the standard configuration defined by the **Device SDK**. 

The `DeviceList` configuration is standard except that it is mandatory to provide end device parameters in the `DeviceList.Protocols.EndDevice_Params` structure. The following is a sample `DeviceList` that works with the sample device profiles referenced below. `ED_URI_Prefix` parameter is optional.

```toml
[[DeviceList]]
  Name = "ED_ID1"
  ProfileName = "ED_ID1"
  Description = "RESTful Device that sends data"
  Labels = [ "rest", "ED_ID1" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.EndDevice_Params]
	  ED_IP = "127.0.0.1"
	  ED_PORT = "5000"
	  ED_URI_Prefix = "api"
  [[DeviceList.AutoEvents]]
    Interval = "20s"
    OnChange = false
    SourceName = "jsonRes"
```

### Device Profile

This section shows sample device profile defined for REST client functionality.

- [**ED_ID1**](cmd/res/profiles/ED_ID1.yaml)

### Testing/Simulation of REST client

This device service supports commanding functionality for the data types as shown in below table.

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

Using `curl` command-line utility or `PostMan` we can test commanding functionaity of the REST device service.

**GET Request**

Example GET request to `int8` device resource using curl command-line utility is as shown below.
```
   $ curl --request GET http://localhost:59882/api/v2/device/name/ED_ID1/int8
```
`ED_ID1` is the device name as defined in the device file.
Example expected success response from the end device is as shown below.

```
   {
   "apiVersion" : "v2",
   "event" : {
      "apiVersion" : "v2",
      "deviceName" : "ED_ID1",
      "id" : "46baf3d5-98fd-4073-b52e-801660b01ce6",
      "origin" : 1670506568209119757,
      "profileName" : "ED_ID1",
      "readings" : [
         {
            "deviceName" : "ED_ID1",
            "id" : "c7d4d4fe-13f5-423a-8d62-0e57f8dbc063",
            "origin" : 1670506568209111164,
            "profileName" : "ED_ID1",
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

**PUT Request**

Example PUT request to `int8` device resource using curl command-line utility is as shown below.
```
   $ curl -i -X PUT -H "Content-Type: application/json" -d '{"int8":12}' http://localhost:59882/api/v2/device/name/ED_ID1/int8
```
`ED_ID1` is the device name as defined in the device list.

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
