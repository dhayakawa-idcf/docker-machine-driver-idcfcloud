VERSION=0.4.1

GOVERSION ?= 1.12.7
BUILD_DIR := bin
CMD_DIR := docker-machine-driver-idcfcloud
PKG ?= github.com/dhayakawa-idcf/docker-machine-driver-idcfcloud

TARGETS=linux darwin windows

ifeq ($(shell uname), Darwin)
CURRENT_OS=darwin
else ifeq ($(shell uname), Linux)
CURRENT_OS=linux
else
CURRENT_OS=windows
endif

LDFLAGS ?= -X $(PKG)/pkg/version.Version=$(VERSION)


.PHONY: all
all: vendor mkdir build

.PHONY: build
build: build-$(CURRENT_OS)

.PHONY: build-all
build-all: $(TARGETS:%=build-%)

build-%:
	@echo "==> Building the driver - $*"
	@docker run -v $(PWD):/go/src/$(PKG) \
		-w /go/src/$(PKG) \
		-e GOOS=$* -e GOARCH=amd64 -e CGO_ENABLED=0 -e GOFLAGS=-mod=vendor golang:$(GOVERSION)-alpine3.10 \
		go build -o $(BUILD_DIR)/$(CMD_DIR)-$*-amd64-$(VERSION) -ldflags "$(LDFLAGS)" $(PKG)/cmd/$(CMD_DIR)/

.PHONY: clean-all
clean-all: $(TARGETS:%=clean-%)
	@rm -fr bin/*

clean-%:
	@echo "==> Cleaning releases"
	@GOOS=$* go clean -i -x ./...

.PHONY: install
install:
	cp $(BUILD_DIR)/$(CMD_DIR)-$(CURRENT_OS)-amd64-$(VERSION) /usr/local/bin/$(CMD_DIR)

.PHONY: release
release:
	ghr $(VERSION) bin/

.PHONY: vendor
vendor:
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

mkdir:
	@echo "==> Creating build directories"
	@mkdir -p $(BUILD_DIR)
