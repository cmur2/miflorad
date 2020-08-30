
# https://tech.davis-hansson.com/p/make/
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eux -o pipefail -c
.DELETE_ON_ERROR:
.SILENT:
.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

MIFLORA_ADDR?=00:00:00:00:00:00
MIFLORAD_VERSION?=master

RUN_COMMAND=miflorad
RUN_OPTIONS=$(MIFLORA_ADDR)

.PHONY: all
all: clean build test ## Run clean, build and test (default goal)

.PHONY: run
run: clean build test ## Run clean, build, test and finally launch the $RUN_COMMAND as root
	sudo cmd/$(RUN_COMMAND)/$(RUN_COMMAND) $(RUN_OPTIONS)

.PHONY: clean
clean: ## Remove all produced executables
	rm -f cmd/miflorad/miflorad
	rm -f cmd/munin-miflora/munin-miflora
	rm -f cmd/munin-miflora-gatt/munin-miflora-gatt

.PHONY: build
build: cmd/miflorad/miflorad cmd/munin-miflora/munin-miflora cmd/munin-miflora-gatt/munin-miflora-gatt ## Build all produced executables

.PHONY: test
test: build ## Run all tests
	pushd cmd/miflorad
	go test -v -race
	popd
	pushd common
	go test -v -race
	popd

.PHONY: remote-run
remote-run: clean ## Run clean, build $RUN_COMMAND for Linux on ARM and launch it via SSH on extzero
	pushd cmd/$(RUN_COMMAND)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w"
	popd
	file cmd/$(RUN_COMMAND)/$(RUN_COMMAND)
	scp cmd/$(RUN_COMMAND)/$(RUN_COMMAND) extzero:$(RUN_COMMAND)
	ssh extzero "./$(RUN_COMMAND) $(RUN_OPTIONS)"

.PHONY: release
release: ## Build and upload release version of miflorad to Github
	mkdir -p pkg
	pushd cmd/miflorad
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "../../pkg/miflorad-$(MIFLORAD_VERSION)-linux-amd64" -ldflags="-s -w -X main.version=$(MIFLORAD_VERSION)"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o "../../pkg/miflorad-$(MIFLORAD_VERSION)-linux-arm" -ldflags="-s -w -X main.version=$(MIFLORAD_VERSION)"
	popd
#	github-release "v$(MIFLORAD_VERSION)" pkg/miflorad-$(MIFLORAD_VERSION)-* --commit "master" --tag "v$(MIFLORAD_VERSION)" --prerelease --github-repository "cmur2/miflorad"
	github-release "v$(MIFLORAD_VERSION)" pkg/miflorad-$(MIFLORAD_VERSION)-* --commit "master" --tag "v$(MIFLORAD_VERSION)" --github-repository "cmur2/miflorad"

.PHONY: cmd/miflorad/miflorad
cmd/miflorad/miflorad:
	pushd cmd/miflorad
	CGO_ENABLED=0 go build -buildmode=pie -ldflags "-X main.version=$(MIFLORAD_VERSION)"

.PHONY: cmd/munin-miflora/munin-miflora
cmd/munin-miflora/munin-miflora:
	pushd cmd/munin-miflora
	CGO_ENABLED=0 go build -buildmode=pie

.PHONY: cmd/munin-miflora-gatt/munin-miflora-gatt
cmd/munin-miflora-gatt/munin-miflora-gatt:
	pushd cmd/munin-miflora-gatt
	CGO_ENABLED=0 go build -buildmode=pie

.PHONY: help
help: ## Print this help text
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
