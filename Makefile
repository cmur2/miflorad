.POSIX:

.PHONY: all
all: clean build test

.PHONY: run
run: clean build test
	sudo cmd/munin-miflora/munin-miflora xyz 00:00:00:00:00:00

.PHONY: clean
clean:
	rm -f cmd/munin-miflora/munin-miflora

.PHONY: build
build: cmd/munin-miflora/munin-miflora

.PHONY: test
test: build

cmd/munin-miflora/munin-miflora:
	cd cmd/munin-miflora && go build && cd ../..
