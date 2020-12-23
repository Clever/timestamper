include golang.mk
include lambda.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
export PATH := $(PWD)/bin:$(PATH)
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /tools)
REPONAME := $(notdir $(shell pwd))
CMD ?= handler
APP_NAME ?= $(REPONAME)

.PHONY: test build run $(PKGS) install_deps

$(eval $(call golang-version-check,1.13))

test: generate $(PKGS)

build: generate
	$(call lambda-build-go,./cmd/$(CMD),$(APP_NAME))

build-local: generate
	$(call golang-build,./cmd/$(CMD),$(APP_NAME))

run: build-local
	IS_LOCAL=true bin/$(APP_NAME)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

generate:
	go generate ./cmd/$(CMD)

install_deps:
	go mod vendor
	go build -o bin/mockgen    ./vendor/github.com/golang/mock/mockgen
	go build -o bin/go-bindata ./vendor/github.com/kevinburke/go-bindata/go-bindata
