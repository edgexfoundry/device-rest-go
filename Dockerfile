#
# Copyright (c) 2022 Intel Corporation
# Copyright (c) 2021 IOTech Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

ARG BASE=golang:1.18-alpine3.16
FROM ${BASE} AS builder

ARG MAKE='make build'
ARG ALPINE_PKG_BASE="make git openssh-client gcc libc-dev zeromq-dev libsodium-dev"
ARG ALPINE_PKG_EXTRA=""

RUN sed -e 's/dl-cdn[.]alpinelinux.org/dl-4.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

WORKDIR /device-rest-go

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
RUN $MAKE

FROM alpine:3.16

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2022: Intel'

LABEL Name=device-rest-go Version=${VERSION}

# dumb-init needed for injected secure bootstrapping entrypoint script when run in secure mode.
RUN apk add --update --no-cache zeromq dumb-init

COPY --from=builder /device-rest-go/cmd /
COPY --from=builder /device-rest-go/LICENSE /
COPY --from=builder /device-rest-go/Attribution.txt /

EXPOSE 59986

ENTRYPOINT ["/device-rest"]
CMD ["--cp=consul://edgex-core-consul:8500", "--confdir=/res", "--registry"]