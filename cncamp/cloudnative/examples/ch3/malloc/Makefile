SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MM_SHELL ?= bash -c
MAKEFLAGS += --no-print-directory

.PHONY: help # print this help message.
help:
	@grep -E '^\.PHONY: [a-zA-Z_-]+ .*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: |#)"}; {printf "%-30s %s\n", $$2, $$3}'

.PHONY: build # builds binaries for the target.
build:
	CGO_ENABLED=1 GOOS=linux CGO_LDFLAGS="-static" go build