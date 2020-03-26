PACKAGES=$(shell go list ./... )

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BINDIR ?= $(GOPATH)/bin

GORELEASER = $(BINDIR)/goreleaser

export GO111MODULE = on

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=WrkOracle \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=wrkoracle \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=wrkoracle \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: lint install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/wrkoracle

build: clean go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o build/wrkoracle ./cmd/wrkoracle

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

lint:
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	go mod verify

test:
	@go test -mod=readonly ./...

clean:
	rm -rf build/

snapshot: goreleaser
	goreleaser --snapshot --skip-publish --rm-dist

release: goreleaser
	goreleaser --rm-dist

goreleaser: $(GORELEASER)
$(GORELEASER):
	@echo "Installing goreleaser..."
	@(cd /tmp && go get github.com/goreleaser/goreleaser)
