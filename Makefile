B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d)

# get current user name
USER=$(shell whoami)
# get current user group
GROUP=$(shell id -gn)

.PHONY: build run test deploy status remove
build: info
	- go build -v --ldflags="-X main.version=$(REV)" -o ./bin/ ./cmd/api

run: build
	- ./bin/api

test:
	go test -v -race -mod=vendor ./...

info:
	- @echo "revision $(REV)"

.DEFAULT_GOAL: build