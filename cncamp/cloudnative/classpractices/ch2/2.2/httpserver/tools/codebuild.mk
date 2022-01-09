# Common make targets for build and archive.

export GO111MODULE ?= on
export GOPROXY ?= https://goproxy.cn,direct
export GOSUMDB ?= sum.golang.google.cn

GO := GO111MODULE=on go
# default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0

### os, arch, tag, test, ext vars
LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL ?= windows
endif
export GOOS ?= $(TARGET_OS_LOCAL)

ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL = amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL = arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL = arm
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),arm64)
	TARGET_ARCH_LOCAL = arm64
else
	TARGET_ARCH_LOCAL = amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

ifeq ($(GOARCH),amd64)
	LATEST_TAG = latest
else
	LATEST_TAG = latest-$(GOARCH)
endif

ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
export ARCHIVE_EXT = .zip
else
BINARY_EXT_LOCAL:=
export ARCHIVE_EXT = .tar.gz
endif

export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

ifeq ($(origin DEBUG), undefined)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else ifeq ($(DEBUG),0)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else
  BUILDTYPE_DIR:=debug
  GCFLAGS:=-gcflags="all=-N -l"
  LDFLAGS:="$(DEFAULT_LDFLAGS)"
  $(info Build with debugger information)
endif

BUILDTAGS:=-tags $(DEFAULT_BUILDTAGS)

OUT_DIR ?= ./build
APP_OUT_DIR := $(OUT_DIR)/$(GOOS)_$(GOARCH)/$(BUILDTYPE_DIR)
APP_LINUX_OUT_DIR := $(OUT_DIR)/linux_$(GOARCH)/$(BUILDTYPE_DIR)
ifneq ($(OUT_DIR), ./build)
  APP_OUT_DIR := $(OUT_DIR)/$(BUILDTYPE_DIR)
  APP_LINUX_OUT_DIR := $(OUT_DIR)/linux/$(BUILDTYPE_DIR)
endif

################################################################################
# Target: build                                                                #
################################################################################
.PHONY: build # builds binaries for the target.
APP_BINS:=$(foreach ITEM,$(BINARIES),$(APP_OUT_DIR)/$(ITEM)$(BINARY_EXT))
build: $(APP_BINS)

# Generate builds binaries for the target
# Params:
# $(1): the binary name for the target
# $(2): the binary main directory
# $(3): the target os
# $(4): the target arch
# $(5): the output directory
define genBinariesForTarget
.PHONY: $(5)/$(1)
$(5)/$(1):
	CGO_ENABLED=$(CGO) GOOS=$(3) GOARCH=$(4) $(GO) build $(BUILDTAGS) $(GCFLAGS) -ldflags=$(LDFLAGS) \
	-o $(5)/$(1) $(2)/;
endef

# Generate binary targets
# ./cmd/$(ITEM)
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM)$(BINARY_EXT),.,$(GOOS),$(GOARCH),$(APP_OUT_DIR))))

################################################################################
# Target: build-linux                                                          #
################################################################################
.PHONY: build-linux # build linux binaries for the target.
BUILD_LINUX_BINS:=$(foreach ITEM,$(BINARIES),$(APP_LINUX_OUT_DIR)/$(ITEM))
build-linux: $(BUILD_LINUX_BINS)

# Generate linux binaries targets to build linux docker image
# ./cmd/$(ITEM)
ifneq ($(GOOS), linux)
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM),.,linux,$(GOARCH),$(APP_LINUX_OUT_DIR))))
endif

################################################################################
# Target: archive                                                              #
################################################################################
.PHONY: archive # archive files for each binary.
ARCHIVE_OUT_DIR ?= $(APP_OUT_DIR)
ARCHIVE_FILE_EXTS:=$(foreach ITEM,$(BINARIES),archive-$(ITEM)$(ARCHIVE_EXT))
archive: build $(ARCHIVE_FILE_EXTS)

# Generate archive files for each binary
# $(1): the binary name to be archived
# $(2): the archived file output directory
define genArchiveBinary
ifeq ($(GOOS),windows)
archive-$(1).zip:
	7z.exe a -tzip "$(2)\\$(1)_$(GOOS)_$(GOARCH)$(ARCHIVE_EXT)" "$(APP_OUT_DIR)\\$(1)$(BINARY_EXT)"
else
archive-$(1).tar.gz:
	tar czf "$(2)/$(1)_$(GOOS)_$(GOARCH)$(ARCHIVE_EXT)" -C "$(APP_OUT_DIR)" "$(1)$(BINARY_EXT)"
endif
endef

# Generate archive-*.[zip|tar.gz] targets
$(foreach ITEM,$(BINARIES),$(eval $(call genArchiveBinary,$(ITEM),$(ARCHIVE_OUT_DIR))))
