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
	rm -f cmd/munin-miflora/miflorad
	rm -f cmd/munin-miflora/munin-miflora
	rm -f cmd/munin-miflora/munin-miflora-gatt

.PHONY: build
build: cmd/munin-miflora/miflorad cmd/munin-miflora/munin-miflora cmd/munin-miflora/munin-miflora-gatt

.PHONY: test
test: build
	cd common && go test -v -race && cd ..

.PHONY: remote-run
remote-run: clean
	cd cmd/$(RUN_COMMAND) && CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w" && cd ../..
	file cmd/$(RUN_COMMAND)/$(RUN_COMMAND)
	scp cmd/$(RUN_COMMAND)/$(RUN_COMMAND) extzero:$(RUN_COMMAND)
	ssh extzero "./$(RUN_COMMAND) $(RUN_OPTIONS)"

.PHONY: cmd/munin-miflora/miflorad
cmd/munin-miflora/miflorad:
	cd cmd/miflorad && CGO_ENABLED=0 go build -buildmode=pie -ldflags "-X main.version=$(MIFLORAD_VERSION)" && cd ../..

.PHONY: cmd/munin-miflora/munin-miflora
cmd/munin-miflora/munin-miflora:
	cd cmd/munin-miflora && CGO_ENABLED=0 go build -buildmode=pie && cd ../..

.PHONY: cmd/munin-miflora/munin-miflora-gatt
cmd/munin-miflora/munin-miflora-gatt:
	cd cmd/munin-miflora-gatt && CGO_ENABLED=0 go build -buildmode=pie && cd ../..
