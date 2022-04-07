.PHONY: build test unittest lint clean prepare update docker

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

# see https://shibumi.dev/posts/hardening-executables
CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2"
CGO_CFLAGS="-O2 -pipe -fno-plt"
CGO_CXXFLAGS="-O2 -pipe -fno-plt"
CGO_LDFLAGS="-Wl,-O1,–sort-common,–as-needed,-z,relro,-z,now"

MICROSERVICES=cmd/device-rest

.PHONY: $(MICROSERVICES)

ARCH=$(shell uname -m)

DOCKERS=docker_device_rest_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-rest-go.Version=$(VERSION)" -trimpath -mod=readonly
CGOFLAGS=-ldflags "-linkmode=external -X github.com/edgexfoundry/device-rest-go.Version=$(VERSION)" -trimpath -mod=readonly -buildmode=pie

tidy:
	go mod tidy

build: $(MICROSERVICES)

cmd/device-rest:
	$(GOCGO) build $(CGOFLAGS) -o $@ ./cmd

unittest:
	$(GOCGO) test ./... -coverprofile=coverage.out ./...

lint:
	@which golangci-lint >/dev/null || echo "WARNING: go linter not installed. To install, run\n  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.42.1"
	@if [ "z${ARCH}" = "zx86_64" ] && which golangci-lint >/dev/null ; then golangci-lint run --config .golangci.yml ; else echo "WARNING: Linting skipped (not on x86_64 or linter not installed)"; fi

test: unittest lint
	$(GOCGO) vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICES)

update:
	$(GO) mod download

docker: $(DOCKERS)

docker_device_rest_go:
	docker build \
		--build-arg http_proxy \
		--build-arg https_proxy \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-rest:$(GIT_SHA) \
		-t edgexfoundry/device-rest:$(VERSION)-dev \
		.

vendor:
	$(GO) mod vendor
