.PHONY: clean test check build.local build.linux build.osx build.docker build.push

BINARY        ?= n26
VERSION       ?= $(shell git describe --tags --always --dirty)
TAG           ?= $(VERSION)
SOURCES       = $(shell find . -name '*.go')
GOPKGS        = $(shell go list ./...)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X main.version=$(VERSION) -w -s

default: build.local

clean:
	rm -rf build
	rm -rf $(GENERATED)

test: $(GENERATED)
	go test -v $(GOPKGS)

check:
	golint $(GOPKGS)
	go vet -v $(GOPKGS)

lint:
	golangci-lint run ./...

build.local: build/$(BINARY)
build.linux: build/linux/$(BINARY)
build.osx: build/osx/$(BINARY)

build/$(BINARY): $(GENERATED) $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/linux/$(BINARY): $(GENERATED) $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/linux/$(BINARY) -ldflags "$(LDFLAGS)" .

build/osx/$(BINARY): $(GENERATED) $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o build/osx/$(BINARY) -ldflags "$(LDFLAGS)" .
