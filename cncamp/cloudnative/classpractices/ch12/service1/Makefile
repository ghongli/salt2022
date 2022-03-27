SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MM_SHELL ?= bash -c
MAKEFLAGS += --no-print-directory

GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)

REL_VERSION            ?= latest
ifeq ($(REL_VERSION),edge)
	REL_VERSION := latest
endif

DOCS_DEPLOY_USE_SSH ?= true
DOCS_DEPLOY_GIT_USER ?= git

################################################################################
# Go build details                                                             #
################################################################################
BASE_PACKAGE_NAME := github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch12/service1
RELEASE_NAME ?= httpserverv1
BINARIES ?= server
VERSION ?= edge
ifdef REL_VERSION
	VERSION := $(REL_VERSION)
endif
#VERSION := $(shell git describe --tags)
DEFAULT_LDFLAGS += -X "$(BASE_PACKAGE_NAME)/main.date=$(shell date +"%Y%m%d_%H%M%S")"
DEFAULT_LDFLAGS += -X "$(BASE_PACKAGE_NAME)/main.version=$(VERSION)"
DEFAULT_LDFLAGS += -X "$(BASE_PACKAGE_NAME)/main.commit=$(GIT_COMMIT)"
DEFAULT_LDFLAGS += -X "$(BASE_PACKAGE_NAME)/main.goVersion=$(shell go version | sed -E 's/go version go(.*)\ .*/\1/')"
DEFAULT_BUILDTAGS ?= app

YARN:=./build/bin/yarn.sh
PROJECT_ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: help # Print this help message.
 help:
	@grep -E '^\.PHONY: [a-zA-Z_-]+ .*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: |#)"}; {printf "%-30s %s\n", $$2, $$3}'

################################################################################
# Target: code build and archive                                               #
################################################################################
include tools/codebuild.mk

################################################################################
# Target: docker                                                               #
################################################################################
DOCKER_IMAGE_PREFIX ?=cncamp-
APPS                ?=$(RELEASE_NAME)

include tools/docker.mk