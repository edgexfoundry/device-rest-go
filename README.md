# device-rest-go
[![Build Status](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/badge/icon)](https://jenkins.edgexfoundry.org/view/EdgeX%20Foundry%20Project/job/edgexfoundry/job/device-rest-go/job/main/) [![Code Coverage](https://codecov.io/gh/edgexfoundry/device-rest-go/branch/main/graph/badge.svg?token=fmbJjqOyk4)](https://codecov.io/gh/edgexfoundry/device-rest-go) [![Go Report Card](https://goreportcard.com/badge/github.com/edgexfoundry/device-rest-go)](https://goreportcard.com/report/github.com/edgexfoundry/device-rest-go) [![GitHub Latest Dev Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?include_prereleases&sort=semver&label=latest-dev)](https://github.com/edgexfoundry/device-rest-go/tags) ![GitHub Latest Stable Tag)](https://img.shields.io/github/v/tag/edgexfoundry/device-rest-go?sort=semver&label=latest-stable) [![GitHub License](https://img.shields.io/github/license/edgexfoundry/device-rest-go)](https://choosealicense.com/licenses/apache-2.0/) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/edgexfoundry/device-rest-go) [![GitHub Pull Requests](https://img.shields.io/github/issues-pr-raw/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/pulls) [![GitHub Contributors](https://img.shields.io/github/contributors/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/contributors) [![GitHub Committers](https://img.shields.io/badge/team-committers-green)](https://github.com/orgs/edgexfoundry/teams/device-rest-go-committers/members) [![GitHub Commit Activity](https://img.shields.io/github/commit-activity/m/edgexfoundry/device-rest-go)](https://github.com/edgexfoundry/device-rest-go/commits)

> **Warning**  
> The **main** branch of this repository contains work-in-progress development code for the upcoming release, and is **not guaranteed to be stable or working**.
> It is only compatible with the [main branch of edgex-compose](https://github.com/edgexfoundry/edgex-compose) which uses the Docker images built from the **main** branch of this repo and other repos.
>
> **The source for the latest release can be found at [Releases](https://github.com/edgexfoundry/device-rest-go/releases).**

EdgeX device service for REST protocol

This device service supports two-way communication for commanding. REST endpoints provides easy way for 3'rd party applications, such as Point of Sale, CV Analytics, etc., to push async data into EdgeX via the REST protocol. Commands allows EdgeX to send GET and PUT request to end device. 

## Documentation

For latest documentation please visit https://docs.edgexfoundry.org/latest/microservices/device/services/device-rest/Purpose

The locally built Docker image can then be used in place of the published Docker image in your compose file.
See [Compose Builder](https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder#gen) `nat-bus` option to generate compose file for NATS and local dev images.

## Build Instructions

1. Clone the device-rest-go repo with the following command:

        git clone https://github.com/edgexfoundry/device-rest-go.git

2. Build a docker image by using the following command:

        make docker

3. Alternatively the device service can be built natively:

        make build

## Build with NATS Messaging
Currently, the NATS Messaging capability (NATS MessageBus) is opt-in at build time.
This means that the published Docker images do not include the NATS messaging capability.

The following make commands will build the local binary or local Docker image with NATS messaging
capability included.
```makefile
make build-nats
make docker-nats
```
 
## Packaging

This component is packaged as docker images.
Please refer to the [Dockerfile] and [Docker Compose Builder] scripts.

[Dockerfile]: Dockerfile
[Docker Compose Builder]: https://github.com/edgexfoundry/edgex-compose/tree/main/compose-builder
