BINARY := energy-prometheus-exporter
BUILD_DIR := $(shell pwd)/build
CHART_NAME := $(shell grep 'name:' deployment/helm-chart/Chart.yaml | awk '{print $$2}')
CHART_VERSION := $(shell grep 'version:' deployment/helm-chart/Chart.yaml | awk '{print $$2}')
IMAGE := quay.io/touchardv/energy-prometheus-exporter
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
SOURCES := $(shell find . -name '*.go')
TAG := latest
TARGET ?= $(shell uname -m)

ifeq ($(GOARCH), arm64)
 DOCKER_BUILDX_PLATFORM := linux/arm64/v8
else ifeq ($(GOARCH), amd64)
 DOCKER_BUILDX_PLATFORM := linux/amd64
endif

.DEFAULT_GOAL := build
.PHONY: build
build: $(BUILD_DIR)/$(BINARY)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/$(BINARY): $(SOURCES)
	go mod tidy
	go build -o $(BUILD_DIR)/$(BINARY) .

$(BUILD_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz:
	helm package deployment/helm-chart -d $(BUILD_DIR)

$(BINARY)-linux-$(GOARCH): $(BUILD_DIR) $(SOURCES)
	go mod tidy
	GOOS=linux GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY)-linux-$(GOARCH) .

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	go clean

.PHONY: package
package: package-helm-chart package-image

.PHONY: package-helm-chart
package-helm-chart: $(BUILD_DIR)/$(CHART_NAME)-$(CHART_VERSION).tgz

.PHONY: package-image
package-image: $(BINARY)-linux-$(GOARCH)
	docker buildx build --progress plain \
	--platform $(DOCKER_BUILDX_PLATFORM) \
	--tag $(IMAGE):$(TAG) --load -f deployment/Dockerfile .

.PHONY: run
run: $(BUILD_DIR)/$(BINARY)
	source local.env && $(BUILD_DIR)/$(BINARY) serve

.PHONY: run-image
run-image: package-image
	docker run -it --rm -e HOMEWIZARD_PROMETHEUS_EXPORTER_URL=http://foo.bar --entrypoint=/bin/sh $(IMAGE):$(TAG)

.PHONY: test
test:
	go test -v -cover -timeout 10s ./...
	helm lint deployment/helm-chart
