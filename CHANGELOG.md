
<a name="EdgeX REST Device Service (found in device-rest-go) Changelog"></a>
## EdgeX REST Device Service
[Github repository](https://github.com/edgexfoundry/device-rest-go)

### Change Logs for EdgeX Dependencies
- [device-sdk-go](https://github.com/edgexfoundry/device-sdk-go/blob/main/CHANGELOG.md)
- [go-mod-core-contracts](https://github.com/edgexfoundry/go-mod-core-contracts/blob/main/CHANGELOG.md)
- [go-mod-bootstrap](https://github.com/edgexfoundry/go-mod-bootstrap/blob/main/CHANGELOG.md)  (indirect dependency)
- [go-mod-messaging](https://github.com/edgexfoundry/go-mod-messaging/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-registry](https://github.com/edgexfoundry/go-mod-registry/blob/main/CHANGELOG.md)  (indirect dependency)
- [go-mod-secrets](https://github.com/edgexfoundry/go-mod-secrets/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-configuration](https://github.com/edgexfoundry/go-mod-configuration/blob/main/CHANGELOG.md) (indirect dependency)

## [4.0.0] Odessa - 2025-03-12 (Only compatible with the 4.x releases)

### ✨  Features

- Enable PIE support for ASLR and full RELRO ([be40ce9…](https://github.com/edgexfoundry/device-rest-go/commit/be40ce9c5fc91884215293b6bf7bbda8bd3c0d2f))

### ♻ Code Refactoring

- Update module to v4 ([6bfdd3f…](https://github.com/edgexfoundry/device-rest-go/commit/6bfdd3ff5f3e792eafad2b4c95b01495d5837e2e))
```text

BREAKING CHANGE: update go module to v4

```

### 🐛 Bug Fixes

- Only one ldflags flag is allowed ([ec1ec23…](https://github.com/edgexfoundry/device-rest-go/commit/ec1ec23295ae906f939fdcce9f1e5a4eef1babde))

### 👷 Build

- Upgrade to go-1.23, Linter1.61.0 and Alpine 3.20 ([21d3830…](https://github.com/edgexfoundry/device-rest-go/commit/21d3830f1720a37f1895b62ff4ca15d755f1aed1))

## [v3.1.0] Napa - 2023-11-15 (Only compatible with the 3.x releases)


### ✨  Features

- Remove snap packaging ([#353](https://github.com/edgexfoundry/device-rest-go/issues/353)) ([758b2be…](https://github.com/edgexfoundry/device-rest-go/commit/758b2bed15a3d13a7bac802538a1e1661ebd32c9))
```text

BREAKING CHANGE: Remove snap packaging ([#353](https://github.com/edgexfoundry/device-rest-go/issues/353))

```
- Replace gorilla/mux with labstack/echo ([#334](https://github.com/edgexfoundry/device-rest-go/issues/334)) ([cf6f5c6…](https://github.com/edgexfoundry/device-rest-go/commit/cf6f5c6b2bf2b83036e025c5ddfa4d5590d4fdb7))


### ♻ Code Refactoring

- Remove obsolete comments from config file ([#355](https://github.com/edgexfoundry/device-rest-go/issues/355)) ([71e268e…](https://github.com/edgexfoundry/device-rest-go/commit/71e268e2bbeb5f32b791e9fe51a95144d8b6bbea))
- Remove github.com/pkg/errors from Attribution.txt ([bbfcc02…](https://github.com/edgexfoundry/device-rest-go/commit/bbfcc0226598174b6127e805ff9b931fc1acdcd9))


### 🐛 Bug Fixes

- Range check for float64/32 equal to 0.0 ([#330](https://github.com/edgexfoundry/device-rest-go/issues/330)) ([89173ea…](https://github.com/edgexfoundry/device-rest-go/commit/89173ea65a943aba0f4628b8150ce015ac4d5016))


### 📖 Documentation

- [Device REST] Add Swagger file ([#360](https://github.com/edgexfoundry/device-rest-go/issues/360)) ([47318dc…](https://github.com/edgexfoundry/device-rest-go/commit/47318dce372b4a1f6fa2293d118d34d1e168baa0))


### 👷 Build

- Upgrade to go-1.21, Linter1.54.2 and Alpine 3.18 ([#340](https://github.com/edgexfoundry/device-rest-go/issues/340)) ([7661139…](https://github.com/edgexfoundry/device-rest-go/commit/7661139dd1c3d848a7247d8835c5c2fde1868ee5))


### 🤖 Continuous Integration

- Add automated release workflow on tag creation ([4617776…](https://github.com/edgexfoundry/device-rest-go/commit/4617776163951c6dac5c14d5cbb20ddb9bf695df))


## [v3.0.0] Minnesota - 2023-05-31 (Only compatible with the 3.x releases)

### Features ✨
- Add device validation function ([#289](https://github.com/edgexfoundry/device-rest-go/pull/289))
    ```text
    BREAKING CHANGE: Implement `ValidateDevice` function to validate device protocol properties for core-metadata  
    ```
- Support two-way REST communication ([#229](https://github.com/edgexfoundry/device-rest-go/issues/229)) ([#83d0d95](https://github.com/edgexfoundry/device-rest-go/commits/83d0d95))
- Update for common config ([#251](https://github.com/edgexfoundry/device-rest-go/pull/251))
    ```text
    BREAKING CHANGE: Configuration file is changed to remove common config settings
    ```
- Use latest SDK for MessageBus Request API ([#247](https://github.com/edgexfoundry/device-rest-go/pull/247))
    ```text
    BREAKING CHANGE: Commands via MessageBus topic configuration are changed
    ```
- Remove ZeroMQ MessageBus capability ([#235](https://github.com/edgexfoundry/device-rest-go/pull/235))
    ```text
    BREAKING CHANGE: ZeroMQ MessageBus capability no longer available
    ```

### Bug Fixes 🐛
- **snap:** Refactor to avoid conflicts with readonly config provider directory ([#275](https://github.com/edgexfoundry/device-rest-go/issues/275)) ([#5399131](https://github.com/edgexfoundry/device-rest-go/commits/5399131))

### Code Refactoring ♻
- Change configuration and devices files format to YAML ([#281](https://github.com/edgexfoundry/device-rest-go/pull/281))
    ```text
    BREAKING CHANGE: Configuration files are now in YAML format, Default file name is now configuration.yaml
    ```
- **snap:** Update command and metadata sourcing ([#266](https://github.com/edgexfoundry/device-rest-go/issues/266)) ([#2ded8c8](https://github.com/edgexfoundry/device-rest-go/commits/2ded8c8))
- **snap:** Drop the support for legacy snap env options ([#238](https://github.com/edgexfoundry/device-rest-go/issues/238))
    ```text
    BREAKING CHANGE:
    - Drop the support for legacy snap options with env. prefix
    - Upgrade edgex-snap-hooks to v3
    - Upgrade edgex-snap-testing Github action to v3
    - Add snap's Go module to dependabot
    - Other minor refactoring
    ```

### Documentation 📖
- Add main branch Warning ([#316](https://github.com/edgexfoundry/device-rest-go/issues/316)) ([#c2b4310](https://github.com/edgexfoundry/device-rest-go/commits/c2b4310))

### Build 👷
- Update to Go 1.20, Alpine 3.17 and linter v1.51.2 ([#265](https://github.com/edgexfoundry/device-rest-go/issues/265)) ([#eb45550](https://github.com/edgexfoundry/device-rest-go/commits/eb45550))

## [v2.3.0] Levski - 2022-11-09  (Only compatible with the 2.x releases)

### Features ✨

- Add Service Metrics configuration ([#34f8305](https://github.com/edgexfoundry/device-rest-go/commits/34f8305))
- Add NATS configuration and build option ([#74fc3e8](https://github.com/edgexfoundry/device-rest-go/commits/74fc3e8))
- Add commanding via message configuration ([#9080408](https://github.com/edgexfoundry/device-rest-go/commits/9080408))
- Add go-winio to attribution (new SPIFFE dependency) ([#2bcd884](https://github.com/edgexfoundry/device-rest-go/commits/2bcd884))
- **snap:** add config interface with unique identifier ([#219](https://github.com/edgexfoundry/device-rest-go/issues/219)) ([#336f3e7](https://github.com/edgexfoundry/device-rest-go/commits/336f3e7))

### Bug Fixes 🐛

- **snap:** typo in install hook error log ([#192](https://github.com/edgexfoundry/device-rest-go/issues/192)) ([#6db4960](https://github.com/edgexfoundry/device-rest-go/commits/6db4960))

### Code Refactoring ♻

- **snap:** edgex-snap-hooks related upgrade ([#202](https://github.com/edgexfoundry/device-rest-go/issues/202)) ([#deaabcb](https://github.com/edgexfoundry/device-rest-go/commits/deaabcb))

### Build 👷

- Upgrade to Go 1.18 and alpine 3.16 ([#56a929b](https://github.com/edgexfoundry/device-rest-go/commits/56a929b))

## [v2.2.0] Kamakura - 2022-05-11  (Only compatible with the 2.x releases)

### Features ✨
- Enable security hardening ([#cc8c89e](https://github.com/edgexfoundry/device-rest-go/commits/cc8c89e))
- Addressed review issue of value type checking in check uint and int methods ([#747ebf9](https://github.com/edgexfoundry/device-rest-go/commits/747ebf9))

### Bug Fixes 🐛
- **snap:** expose parent directory in device-config plug ([#f967545](https://github.com/edgexfoundry/device-rest-go/commits/f967545))

### Code Refactoring ♻
- **snap:** remove obsolete passthrough usage ([#fd6a719](https://github.com/edgexfoundry/device-rest-go/commits/fd6a719))
- **snap:** remove redundant content indentifier ([#67a8320](https://github.com/edgexfoundry/device-rest-go/commits/67a8320))

### Build 👷
- Update to latest SDK w/o ZMQ on windows ([#222d1f3](https://github.com/edgexfoundry/device-rest-go/commits/222d1f3))
    ```
    BREAKING CHANGE:
    ZeroMQ no longer supported on native Windows for EdgeX
    MessageBus
    ```
- **snap:** source snap metadata from external repo ([#45d15e4](https://github.com/edgexfoundry/device-rest-go/commits/45d15e4))
- **snap:** Add go tidy compat 1.7 flag ([#9ff0b86](https://github.com/edgexfoundry/device-rest-go/commits/9ff0b86))

### Continuous Integration 🔄
- gomod changes related for Go 1.17 ([#5bc210e](https://github.com/edgexfoundry/device-rest-go/commits/5bc210e))
- Go 1.17 related changes ([#06c0ff5](https://github.com/edgexfoundry/device-rest-go/commits/06c0ff5))

## [v2.1.0] Jakarta - 2021-11-18  (Only compatible with the 2.x releases)

### Features ✨
- Add support for Object ValueType ([#96e184b](https://github.com/edgexfoundry/device-rest-go/commits/96e184b))
- Update configuration for new CORS and Secrets File settings ([#5acad06](https://github.com/edgexfoundry/device-rest-go/commits/5acad06))

### Bug Fixes 🐛
- Update device resource RW permission to read-only ([#9eb04a7](https://github.com/edgexfoundry/device-rest-go/commits/9eb04a7))
- Remove the code: condition that is always true ([#4225c14](https://github.com/edgexfoundry/device-rest-go/commits/4225c14))
- Update all TOML to use quote and not single-quote ([#93dcba7](https://github.com/edgexfoundry/device-rest-go/commits/93dcba7))

### Documentation 📖
- Add snap section in README.md ([#31e4a8a](https://github.com/edgexfoundry/device-rest-go/commits/31e4a8a))
- Update build status badge ([#ad5ae9a](https://github.com/edgexfoundry/device-rest-go/commits/ad5ae9a))
- **snap:** Update snap/README's format ([#0360482](https://github.com/edgexfoundry/device-rest-go/commits/0360482))
- **snap:** Update snap/README (part 2) ([#1941e36](https://github.com/edgexfoundry/device-rest-go/commits/1941e36))
- **snap:** Update snap/README ([#55adbc2](https://github.com/edgexfoundry/device-rest-go/commits/55adbc2))

### Build 👷
- Update to latest SDK and released go-mods ([#754aee6](https://github.com/edgexfoundry/device-rest-go/commits/754aee6))
- Update to latest SDK ([#97198e6](https://github.com/edgexfoundry/device-rest-go/commits/97198e6))
- Update alpine base to 3.14 ([#721085b](https://github.com/edgexfoundry/device-rest-go/commits/721085b))
- **snap:** Update snap packaging ([#cbfdaf9](https://github.com/edgexfoundry/device-rest-go/commits/cbfdaf9))
- **snap:** Update `base` to core20 ([#7435c4d](https://github.com/edgexfoundry/device-rest-go/commits/7435c4d))
- **snap:** Update README and add snap/local/hooks/go.sum ([#ae9526d](https://github.com/edgexfoundry/device-rest-go/commits/ae9526d))

### Continuous Integration 🔄
- Remove need for CI specific Dockerfile ([#b871a07](https://github.com/edgexfoundry/device-rest-go/commits/b871a07))

## [v2.0.0] Ireland - 2021-06-30  (Only compatible with the 2.x releases)

### Features ✨
- Enable using MessageBus as the default ([#01c2e73](https://github.com/edgexfoundry/device-rest-go/commits/01c2e73))
- Add Registry/Config Access token capability ([#03a48d5](https://github.com/edgexfoundry/device-rest-go/commits/03a48d5))
- Remove Logging configuration ([#c973575](https://github.com/edgexfoundry/device-rest-go/commits/c973575))
### Bug Fixes 🐛
- use correct service key in SecretStore paths ([#23b2ca7](https://github.com/edgexfoundry/device-rest-go/commits/23b2ca7))
- Add Type='vault' to [SecretStore] config ([#99e6da9](https://github.com/edgexfoundry/device-rest-go/commits/99e6da9))
### Code Refactoring ♻
- remove unimplemented InitCmd/RemoveCmd configuration ([#d82b524](https://github.com/edgexfoundry/device-rest-go/commits/d82b524))
- Change PublishTopicPrefix value to be 'edgex/events/device' ([#13945f3](https://github.com/edgexfoundry/device-rest-go/commits/13945f3))
- Update configuration for change to common ServiceInfo struct Moved non-common settings under [Device] section ([#7571376](https://github.com/edgexfoundry/device-rest-go/commits/7571376))
    ```
    BREAKING CHANGE:
    Service configuration has changed
    ```
- Update to assign and uses new Port Assignments ([#2f1c2cc](https://github.com/edgexfoundry/device-rest-go/commits/2f1c2cc))
    ```
    BREAKING CHANGE:
    Device Rest default port number has changed to 59986
    ```
- Added go mod tidy under test target ([#dd01544](https://github.com/edgexfoundry/device-rest-go/commits/dd01544))
- Update for new service key names and overrides for hyphen to underscore ([#2ecd16f](https://github.com/edgexfoundry/device-rest-go/commits/2ecd16f))
    ```
    BREAKING CHANGE:
    Service key names used in configuration have changed.
    ```
- use v2 device-sdk ([#8b511d7](https://github.com/edgexfoundry/device-rest-go/commits/8b511d7))
### Documentation 📖
- update README for v2 ([#f51f5ca](https://github.com/edgexfoundry/device-rest-go/commits/f51f5ca))
- Add badges to readme ([#972f9a5](https://github.com/edgexfoundry/device-rest-go/commits/972f9a5))
### Build 👷
- update build files for v2 ([#a01389d](https://github.com/edgexfoundry/device-rest-go/commits/a01389d))
- **snap:** set release name to 'ireland' ([#903fe29](https://github.com/edgexfoundry/device-rest-go/commits/903fe29))
- update go.mod to go 1.16 ([#0dd2d84](https://github.com/edgexfoundry/device-rest-go/commits/0dd2d84))
- update Dockerfiles to use go 1.16 ([#2544f5c](https://github.com/edgexfoundry/device-rest-go/commits/2544f5c))
- **snap:** update snap v2 support ([#b99a89d](https://github.com/edgexfoundry/device-rest-go/commits/b99a89d))
- **snap:** update environment overrides for device and profile dir ([#5707fd0](https://github.com/edgexfoundry/device-rest-go/commits/5707fd0))
- **snap:** update epoch for Ireland release ([#629973d](https://github.com/edgexfoundry/device-rest-go/commits/629973d))
- **snap:** fix regression due to v2 build changes ([#a2ffdda](https://github.com/edgexfoundry/device-rest-go/commits/a2ffdda))
- **snap:** update go to 1.16 ([#fc4971f](https://github.com/edgexfoundry/device-rest-go/commits/fc4971f))
- **snap:** '-go' suffix removed from device name ([#8b5b60a](https://github.com/edgexfoundry/device-rest-go/commits/8b5b60a))
- **snap:** run 'go mod tidy' ([#e393ce8](https://github.com/edgexfoundry/device-rest-go/commits/e393ce8))
### Continuous Integration 🔄
- update local docker image names ([#2c710f7](https://github.com/edgexfoundry/device-rest-go/commits/2c710f7))

<a name="v1.2.1"></a>
## [v1.2.1] - 2021-02-02
### Features ✨
- **snap:** add startup-duration and startup-interval configure options ([#4b44503](https://github.com/edgexfoundry/device-rest-go/commits/4b44503))
### Build 👷
- **deps:** Bump github.com/edgexfoundry/device-sdk-go ([#70](https://github.com/edgexfoundry/device-rest-go/issues/70)) ([#abd24f1](https://github.com/edgexfoundry/device-rest-go/commits/abd24f1))
### Continuous Integration 🔄
- add semantic.yml for commit linting, update PR template to latest ([#c3aa815](https://github.com/edgexfoundry/device-rest-go/commits/c3aa815))
- standardize dockerfiles ([#998a81b](https://github.com/edgexfoundry/device-rest-go/commits/998a81b))

<a name="v1.2.0"></a>
## [v1.2.0] - 2020-11-18
### Doc
- correct build instructions ([#36](https://github.com/edgexfoundry/device-rest-go/issues/36)) ([#a96498e](https://github.com/edgexfoundry/device-rest-go/commits/a96498e))
### Bug Fixes 🐛
- **snap:** Update snap versioning logic ([#ad0a8ed](https://github.com/edgexfoundry/device-rest-go/commits/ad0a8ed))
### Code Refactoring ♻
- Upgrade SDK to v1.2.4-dev.34 ([#54](https://github.com/edgexfoundry/device-rest-go/issues/54)) ([#4f6fe4f](https://github.com/edgexfoundry/device-rest-go/commits/4f6fe4f))
- update dockerfile to appropriately use ENTRYPOINT and CMD, closes[#34](https://github.com/edgexfoundry/device-rest-go/issues/34) ([#46301eb](https://github.com/edgexfoundry/device-rest-go/commits/46301eb))
### Build 👷
- upgrade device-sdk-go ([#42](https://github.com/edgexfoundry/device-rest-go/issues/42)) ([#0a79c20](https://github.com/edgexfoundry/device-rest-go/commits/0a79c20))
- Upgrade to Go1.15 ([#069cb69](https://github.com/edgexfoundry/device-rest-go/commits/069cb69))
- **all:** Enable use of DependaBot to maintain Go dependencies ([#755b338](https://github.com/edgexfoundry/device-rest-go/commits/755b338))
- **deps:** Bump github.com/edgexfoundry/device-sdk-go ([#5430346](https://github.com/edgexfoundry/device-rest-go/commits/5430346))
- **deps:** Bump github.com/spf13/cast from 1.3.0 to 1.3.1 ([#72307df](https://github.com/edgexfoundry/device-rest-go/commits/72307df))

<a name="v1.1.2"></a>
## [v1.1.2] - 2020-08-19
### Features ✨
- **device-rest:** Add snap directory for device-rest ([#6a789b2](https://github.com/edgexfoundry/device-rest-go/commits/6a789b2))
### Documentation 📖
- Add standard PR template ([#d097784](https://github.com/edgexfoundry/device-rest-go/commits/d097784))

<a name="v1.1.1"></a>
## [v1.1.1] - 2020-06-15
### Bug Fixes 🐛
- Update package name in main.go to match the one in version.go ([#fb37ef4](https://github.com/edgexfoundry/device-rest-go/commits/fb37ef4))

<a name="v1.1.0"></a>
## [v1.1.0] - 2020-05-13
### Code Refactoring ♻
- Set default logging level to INFO ([#d5d9203](https://github.com/edgexfoundry/device-rest-go/commits/d5d9203))
- upgrade SDK to v1.2.0 ([#32a9f9d](https://github.com/edgexfoundry/device-rest-go/commits/32a9f9d))
### Build 👷
- Switch to Go 1.13 ([#2cc5958](https://github.com/edgexfoundry/device-rest-go/commits/2cc5958))

<a name="v1.0.0"></a>
## v1.0.0 - 2020-02-18
### Features ✨
- **rest ds:** Implement new REST Device service ([#5c6b288](https://github.com/edgexfoundry/device-rest-go/commits/5c6b288))
### Bug Fixes 🐛
- Update to latest release of SDK V1.1.2 for mediaType fix ([#49bf546](https://github.com/edgexfoundry/device-rest-go/commits/49bf546))
