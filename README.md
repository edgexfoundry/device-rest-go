# device-rest-go
EdgeX device service for REST protocol

This device service provides easy way for 3'rd party applications, such as Point of Sale, CV Analytics, etc., to push data into EdgeX via the REST protocol. 

## Runtime Requisites

- core-data
  - Mongo or Redis DB
- core-metadata

## REST Endpoints

This device service creates the additional parametrized `REST` endpoint:

```
/ap1/vi/resource/{deviceName}/{resourceName}
```

- `deviceName` refers to the `device` defined in a `device profile` and the `configuration.toml`.
- `resourceName`refers to the `device resource` defined in the `device profile` that `deviceName` references.

The data, `text` or `binary`,  posted to this endpoint is type validated and type casted (text data only) to the type defined by the specified `device resource`. The resulting value is then sent into EdgeX via the Device SDK's `async values` channel. 

> *Note: When binary data is used the EdgeX event/reading is `CBOR` encoded by the `Device SDK` and the binary value in the reading is`NOT` be stored in the database by `Core Data`. The `CBOR` encoded event/reading, with the binary value, is published to the `Message Bus` for `Application Services` to consume*

> *Note: All non-binary data is consumed as text. The text is casted to the specific type of the specified `device resource` once it passes type validation.*

## Configuration

This device service use the standard configuration defined by the **Device SDK**. 

The `DeviceList` configuration is standard except that the `DeviceList.Protocols` can be empty. The following is a sample `DeviceList` that works with the sample device profiles referenced below.

```toml
[[DeviceList]]
  Name = "sample-json"
  Profile = "sample-json"
  Description = "RESTful Device that sends in JSON data"
  Labels = [ "rest", "json" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]
[[DeviceList]]
  Name = "sample-binary"
  Profile = "sample-binary"
  Description = "RESTful Device that sends in binary data"
  Labels = [ "rest", "binary" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]    
[[DeviceList]]
  Name = "sample-numbers"
  Profile = "sample-numbers"
  Description = "RESTful Device that sends in number data"
  Labels = [ "rest", "float", "int" ]
  [DeviceList.Protocols]
    [DeviceList.Protocols.other]
```

## Device Profile

As with all device services the `device profile` is where the **Device Name**, **Device Resources** and **Device Commands** are define. The parameterized REST endpoint described above references these definitions. Each `Device` has it's own device profile. There are three sample device profiles that define the devices referenced in the above sample configuration.

- **[sample-binary-device](./cmd/res/)**
- [**sample-json-device**](./cmd/res/sample-json-device.yaml)
- [**sample-numbers-device**](./cmd/res/sample-numbers-device.yaml)

> *Note: The`coreCommands` section is omitted since this device service does not support Commanding. See below for details.* 

> *Note: The `deviceCommands` section only requires the `get` operations.*

## Commanding

The current implementation is meant for one-way communication into EdgeX. If future use cases determine that `commanding`, i.e. two-communication, is desirable it can be added then.

## AutoEvents

Since `Commanding` is not implemented, specifying `AutoEvents` in the configuration will result in errors. Thus `AutoEvents` should not be specified in the configuration.

## Installation and Execution

```bash
make build
make run
```

## Build docker image

```bash
make docker
```

