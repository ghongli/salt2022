# Common make targets for Docker images.

HUB_REGISTRY ?= ghongli/cncamp-cloudnative
#TARGET_OS     ?= $(GOOS)
#TARGET_ARCH   ?= $(GOARCH)
TARGET_OS     ?= linux
TARGET_ARCH   ?= amd64
REL_VERSION   ?= latest
ifeq ($(REL_VERSION),edge)
	REL_VERSION := latest
endif

# Add latest tag if LATEST_RELEASE is true
LATEST_RELEASE ?=

# Docker image build and push setting
DOCKER:=docker
DOCKERFILE_DIR ?= ./docker
DOCKERFILE ?= Dockerfile
DOCKER_MULTI_ARCH=linux-amd64 linux-arm linux-arm64 windows-amd64

# build docker image for linux
BIN_PATH ?= $(OUT_DIR)/$(TARGET_OS)_$(TARGET_ARCH)

ifeq ($(TARGET_OS), windows)
  DOCKERFILE:=Dockerfile-windows
  BIN_PATH := $(BIN_PATH)/release
else ifeq ($(origin DEBUG), undefined)
#  BIN_PATH := $(BIN_PATH)/release
else ifeq ($(DEBUG),0)
#  BIN_PATH := $(BIN_PATH)/release
else
  DOCKERFILE:=Dockerfile-debug
  BIN_PATH := $(BIN_PATH)/debug
endif

ifeq ($(TARGET_ARCH),arm)
  DOCKER_IMAGE_PLATFORM:=$(TARGET_OS)/$(TARGET_ARCH)/v7
else ifeq ($(TARGET_ARCH),arm64)
  DOCKER_IMAGE_PLATFORM:=$(TARGET_OS)/$(TARGET_ARCH)/v8
else
  DOCKER_IMAGE_PLATFORM:=$(TARGET_OS)/$(TARGET_ARCH)
endif

# To use buildx: https://github.com/docker/buildx#docker-ce
export DOCKER_CLI_EXPERIMENTAL=enabled

################################################################################
# Target: check-docker-env                                                     #
################################################################################
#.PHONY: check-docker-env # check the docker required environment variables.
check-docker-env:
ifeq ($(HUB_REGISTRY),)
	$(error current `HUB_REGISTRY` environment variable must be set)
else
	$(info current env `HUB_REGISTRY`: $(HUB_REGISTRY))
endif
ifeq ($(DOCKER_IMAGE_PREFIX),)
	$(error current env `DOCKER_IMAGE_PREFIX` environment variable must be set)
else
	$(info current env `DOCKER_IMAGE_PREFIX`: $(DOCKER_IMAGE_PREFIX))
endif

################################################################################
# Target: check-arch                                                           #
################################################################################
#.PHONY: check-arch # check the docker required os,arch environment variables.
check-arch:
ifeq ($(TARGET_OS),)
	$(error `TARGET_OS` environment variable must be set)
else
	$(info env `TARGET_OS`: $(TARGET_OS))
endif
ifeq ($(TARGET_ARCH),)
	$(error `TARGET_ARCH` environment variable must be set)
else
	$(info env `TARGET_ARCH`: $(TARGET_ARCH))
endif

################################################################################
# Target: build image                                                          #
################################################################################
.PHONY: docker-build # build application's image.
BUILD_APPS:=$(foreach ITEM,$(APPS),build-$(ITEM))
docker-build: check-docker-env check-arch $(BUILD_APPS)

# Generate docker image build targets
define genDockerImageBuild
.PHONY: build-$(1)
build-$(1):
ifeq ($(TARGET_ARCH),amd64)
	$(DOCKER) build --build-arg PKG_FILES=* -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(HUB_REGISTRY)-$(1):$(REL_VERSION)-$(TARGET_OS)-$(TARGET_ARCH)
else
	-$(DOCKER) buildx create --use --name appbuild
	-$(DOCKER) run --rm --privileged multiarch/qemu-user-static --reset -p yes
	$(DOCKER) build --build-arg PKG_FILES=* --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(HUB_REGISTRY)-$(1):$(REL_VERSION)-$(TARGET_OS)-$(TARGET_ARCH)
endif
endef

# Generate docker image build targets
$(foreach ITEM,$(APPS),$(eval $(call genDockerImageBuild,$(ITEM))))

################################################################################
# Target: push image                                                           #
################################################################################
.PHONY: docker-push # publish application's image to the registry.
PUSH_APPS:=$(foreach ITEM,$(APPS),push-$(ITEM))
docker-push: docker-build $(PUSH_APPS)

# Generate docker image push targets
define genDockerImagePush
.PHONY: push-$(1)
push-$(1):
ifeq ($(TARGET_ARCH),amd64)
	$(DOCKER) push $(HUB_REGISTRY)-$(1):$(REL_VERSION)-$(TARGET_OS)-$(TARGET_ARCH)
else
	-$(DOCKER) buildx create --use --name appbuild
	-$(DOCKER) run --rm --privileged multiarch/qemu-user-static --reset -p yes
	$(DOCKER) build --build-arg PKG_FILES=* --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(HUB_REGISTRY)-$(1):$(REL_VERSION)-$(TARGET_OS)-$(TARGET_ARCH) --push
endif
endef

# Generate docker image push targets
$(foreach ITEM,$(APPS),$(eval $(call genDockerImagePush,$(ITEM))))

################################################################################
# Target: create docker manifest                                               #
################################################################################
.PHONY: manifest-create # create docker manifest.
CREATE_MANIFEST_APPS:=$(foreach ITEM,$(APPS),manifest-create-$(ITEM))
manifest-create: check-docker-env $(CREATE_MANIFEST_APPS)

# Generate docker manifest create
define genDockerManifestCreate
.PHONY: manifest-create-$(1)
manifest-create-$(1):
	$(DOCKER) manifest create $(HUB_REGISTRY)-$(1):$(REL_VERSION) $(DOCKER_MULTI_ARCH:%=$(HUB_REGISTRY)-$(1):$(REL_VERSION)-%)
ifeq ($(LATEST_RELEASE),true)
	$(DOCKER) manifest create $(HUB_REGISTRY)-$(1):$(LATEST_TAG) $(DOCKER_MULTI_ARCH:%=$(HUB_REGISTRY)-$(1):$(LATEST_TAG)-%)
endif
endef

# Generate docker manifest create
$(foreach ITEM,$(APPS),$(eval $(call genDockerManifestCreate,$(ITEM))))

################################################################################
# Target: push docker manifest                                                 #
################################################################################
.PHONY: manifest-push # publish docker manifest.
PUSH_MANIFEST_APPS:=$(foreach ITEM,$(APPS),manifest-push-$(ITEM))
manifest-push: manifest-create $(PUSH_MANIFEST_APPS)

# Generate docker manifest create
define genDockerManifestPush
.PHONY: manifest-push-$(1)
manifest-push-$(1):
	$(DOCKER) manifest push $(HUB_REGISTRY)-$(1):$(REL_VERSION)
ifeq ($(LATEST_RELEASE),true)
	$(DOCKER) manifest push $(HUB_REGISTRY)-$(1):$(LATEST_TAG)
endif
endef

# Generate docker manifest create
$(foreach ITEM,$(APPS),$(eval $(call genDockerManifestPush,$(ITEM))))

################################################################################
# Target: package docker image                                                 #
################################################################################
.PHONY: image_package # package docker image (build binaries, build images, push images)
image_package: docker-push manifest-push