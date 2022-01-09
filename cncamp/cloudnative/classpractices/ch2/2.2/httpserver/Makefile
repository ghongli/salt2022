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
BASE_PACKAGE_NAME := github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver
RELEASE_NAME ?= server
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

.PHONY: tidy # up the go modules
tidy:
	$(GO) mod tidy

.PHONY: check-diff # check no changes
check-diff:
	git checkout -- go.mod go.sum # can differ between go1.15 and go1.16
	#git diff --exit-code
	git diff --exit-code $(PROJECT_ROOT_DIR)/go.mod
	git diff --exit-code $(PROJECT_ROOT_DIR)/go.sum

.PHONY: yarn-ensure # Install the pinned version of yarn.
yarn-ensure:
	@./tools/install-yarn.sh

PHONY: preflight-checks # check if environment has all required dependencies.
preflight-checks:
	@tools/preflight-checks.sh

.PHONY: api-lint # Lint the generated API assets.
api-lint:
	tools/compile-protos.sh -c "$(PROJECT_ROOT_DIR)/apis" -l

.PHONY: api-lint-fix # Lint and fix the generated API assets.
api-lint-fix:
	tools/compile-protos.sh -c "$(PROJECT_ROOT_DIR)/apis" -lf

.PHONY: api-verify # Verify API proto changes including generated backend assets.
api-verify:
	find pkg/api -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;
	$(MAKE) api
	tools/ensure-no-diff.sh pkg/api

.PHONY: backend # Build the standalone backend.
backend: preflight-checks build

.PHONY: backend-lint # Lint the backend code.
backend-lint: preflight-checks
	tools/golangci-lint.sh run

.PHONY: backend-lint-fix # Lint and fix the backend code.
backend-lint-fix: preflight-checks
	tools/golangci-lint.sh run --fix
	$(GO) mod tidy && $(GO) mod verify

.PHONY: backend-test # Run unit tests for the backend code.
backend-test: preflight-checks
	$(GO) test -race -covermode=atomic ./...

.PHONY: backend-verify # Verify go modules' requirements files are clean.
backend-verify:
	$(GO) mod tidy && $(GO) mod verify
	tools/ensure-no-diff.sh pkg

.PHONY: lint # Lint all of the code.
lint: api-lint backend-lint

.PHONY: lint-fix # Lint and fix all of the code.
lint-fix: api-lint-fix backend-lint-fix

.PHONY: test # Unit test all of the code.
test: backend-test

.PHONY: verify # Verify all of the code.
verify: api-verify backend-verify

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