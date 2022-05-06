# device-rest-go
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/device-rest-go/branch/main/graph/badge.svg?token=fmbJjqOyk4)](https://codecov.io/gh/edgexfoundry/device-rest-go) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-rest-go)](https://goreportcard.com/report/github.com/edgexfoundry/device-rest-go) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-rest-go/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-rest-go)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-rest-go) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-rest-go-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/commits)

EdgeX device service for REST protocol

This device service provides easy way for 3'rd party applications, such as Point of Sale, CV Analytics, etc., to push data into EdgeX via the REST protocol. 

The current implementation is meant for one-way communication into EdgeX via async readings. If future use cases determine a need for`commanding`, i.e. two-communication, it can be added then.

## Runtime Prerequisite    

- core-data
  - Redis DB
- core-metadata

## REST Endpoints

This device service creates the additional parametrized `REST` endpoint:

```
/api/v2/resource/{deviceName}/{resourceName}
```

- `deviceName` refers to a `device` managed by the REST device service.
- `resourceName`refers to the `device resource` defined in the `device profile` associated with the given `deviceName` .

The data, `text` or `binary`,  posted to this endpoint is type validated and type casted (text data only) to the type defined by the specified `device resource`. The resulting value is then sent into EdgeX via the Device SDK's `async values` channel. 

> *Note: When binary data is used the EdgeX event/reading is `CBOR` encoded by the `Device SDK` and the binary value in the reading is`NOT` be stored in the database by `Core Data`. The `CBOR` encoded event/reading, with the binary value, is published to the `Message Bus` for `Application Services` to consume*

> *Note: All non-binary data is consumed as text. The text is casted to the specific type of the specified `device resource` once it passes type validation.*

## Configuration

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

## Device Profile

As with all device services the `device profile` is where the **Device Name**, **Device Resources** and **Device Commands** are define. The parameterized REST endpoint described above references these definitions. Each `Device` has it's own device profile. There are three sample device profiles that define the devices referenced in the above sample configuration.

- [**sample-image-device**](cmd/res/profiles/sample-image-device.yaml)
- [**sample-json-device**](cmd/res/profiles/sample-json-device.yaml)
- [**sample-numeric-device**](cmd/res/profiles/sample-numeric-device.yaml)

> *Note: The`isHidden` field is set to true since this device service does not support Commanding.* 

## Testing/Simulation

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

## AutoEvents

Since `Commanding` is not implemented, specifying `AutoEvents` in the configuration will result in errors. Thus `AutoEvents` should not be specified in the configuration.

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
