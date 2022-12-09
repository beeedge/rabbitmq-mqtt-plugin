.PHONY: all
all: format test build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/beeedge/beethings
VERSION_PACKAGE=github.com/beeedge/beethings/pkg/version

# ==============================================================================
# Includes

include build/lib/common.mk
include build/lib/golang.mk
include build/lib/image.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="device-management monitor"
               Example: make build.multiarch BINS="device-management monitor" PLATFORMS="linux_amd64 linux_arm64 linux_arm" VERSION="v1.0.0"
  IMAGES       Backend images to make. Default is all of cmd.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image IMAGES="device-management monitor"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64、linux_arm64 and linux_arm
               This option is available when using: make image.multiarch
               Example: make image.multiarch IMAGES="device-management monitor" PLATFORMS="linux_amd64 linux_arm64" VERSION="v1.0.0"
  VERSION      The version information compiled into binaries.
               The default is obtained from git.
  V            Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

.PHONY: format
format:
	go fmt ./pkg/... ./cmd/...

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"