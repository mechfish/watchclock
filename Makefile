PKGS := $(shell go list ./... | grep -v /vendor)
BINARY := watchclock
VERSION := $(shell git describe --always --dirty)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: test
test:
	go test $(PKGS)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p dist/$(os)
	env GOOS=$(os) GOARCH=amd64 go build -ldflags="-X main.version=${VERSION}" -o dist/$(os)/$(BINARY) ./cmd/$(BINARY)
	if [ $(os) = 'windows' ]; then mv dist/$(os)/$(BINARY) dist/$(os)/$(BINARY).exe; fi

build: $(PLATFORMS)

.PHONY: clean
clean:
	rm -rf dist/*
