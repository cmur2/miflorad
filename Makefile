.POSIX:

MIFLORA_ADDR?=00:00:00:00:00:00
MIFLORAD_VERSION?=master

RUN_COMMAND=miflorad
RUN_OPTIONS=$(MIFLORA_ADDR)

.PHONY: all
all: clean build test

.PHONY: run
run: clean build test
	sudo cmd/$(RUN_COMMAND)/$(RUN_COMMAND) $(RUN_OPTIONS)

.PHONY: clean
clean:
	rm -f cmd/miflorad/miflorad
	rm -f cmd/munin-miflora/munin-miflora
	rm -f cmd/munin-miflora-gatt/munin-miflora-gatt

.PHONY: build
build: cmd/miflorad/miflorad cmd/munin-miflora/munin-miflora cmd/munin-miflora-gatt/munin-miflora-gatt

.PHONY: test
test: build
	cd cmd/miflorad && go test -v -race && cd ../..
	cd common && go test -v -race && cd ..

.PHONY: remote-run
remote-run: clean
	cd cmd/$(RUN_COMMAND) && CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" && cd ../..
	file cmd/$(RUN_COMMAND)/$(RUN_COMMAND)
	scp cmd/$(RUN_COMMAND)/$(RUN_COMMAND) extzero:$(RUN_COMMAND)
	ssh extzero "./$(RUN_COMMAND) $(RUN_OPTIONS)"

.PHONY: release
release:
	mkdir -p pkg
	cd cmd/miflorad && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "../../pkg/miflorad-$(MIFLORAD_VERSION)-linux-amd64" -ldflags="-s -w -X main.version=$(MIFLORAD_VERSION)" && cd ../..
	cd cmd/miflorad && CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o "../../pkg/miflorad-$(MIFLORAD_VERSION)-linux-arm" -ldflags="-s -w -X main.version=$(MIFLORAD_VERSION)" && cd ../..
#	github-release "v$(MIFLORAD_VERSION)" pkg/miflorad-$(MIFLORAD_VERSION)-* --commit "master" --tag "v$(MIFLORAD_VERSION)" --prerelease --github-repository "cmur2/miflorad"
	github-release "v$(MIFLORAD_VERSION)" pkg/miflorad-$(MIFLORAD_VERSION)-* --commit "master" --tag "v$(MIFLORAD_VERSION)" --github-repository "cmur2/miflorad"

.PHONY: cmd/miflorad/miflorad
cmd/miflorad/miflorad:
	cd cmd/miflorad && CGO_ENABLED=0 go build -buildmode=pie -ldflags "-X main.version=$(MIFLORAD_VERSION)" && cd ../..

.PHONY: cmd/munin-miflora/munin-miflora
cmd/munin-miflora/munin-miflora:
	cd cmd/munin-miflora && CGO_ENABLED=0 go build -buildmode=pie && cd ../..

.PHONY: cmd/munin-miflora-gatt/munin-miflora-gatt
cmd/munin-miflora-gatt/munin-miflora-gatt:
	cd cmd/munin-miflora-gatt && CGO_ENABLED=0 go build -buildmode=pie && cd ../..
