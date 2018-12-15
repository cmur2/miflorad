.POSIX:

MIFLORA_ADDR?=00:00:00:00:00:00

.PHONY: all
all: clean build test

.PHONY: run
run: clean build test
	sudo cmd/munin-miflora/munin-miflora xyz $(MIFLORA_ADDR)

.PHONY: clean
clean:
	rm -f cmd/munin-miflora/munin-miflora
	rm -f cmd/munin-miflora/munin-miflora-gatt

.PHONY: build
build: cmd/munin-miflora/munin-miflora cmd/munin-miflora/munin-miflora-gatt

.PHONY: test
test: build
	cd common && go test -v -race && cd ..

.PHONY: cmd/munin-miflora/munin-miflora
cmd/munin-miflora/munin-miflora:
	cd cmd/munin-miflora && CGO_ENABLED=0 go build && cd ../..

.PHONY: cmd/munin-miflora/munin-miflora-gatt
cmd/munin-miflora/munin-miflora-gatt:
	cd cmd/munin-miflora-gatt && CGO_ENABLED=0 go build && cd ../..
