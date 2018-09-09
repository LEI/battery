SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

PKGS := $(go list ./... | grep -vF /vendor/)

REPO := github.com/LEI/battery

.PHONY: default
default: vendor test format install

.PHONY: get-deps
get-deps:
	# dep, goreleaser...
	go get golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	# go get -u golang.org/x/lint/golint honnef.co/go/tools/cmd/megacheck

# .PHONY: build
# build:
# 	go fmt ./...
# 	BATTERY_BUILD_PLATFORMS=$(PLATFORM) BATTERY_BUILD_ARCHS=$(ARCH) ./build.sh
# 	cp ./dist/battery-$(PLATFORM)-$(ARCH) battery

.PHONY: vendor
vendor:
	dep ensure -vendor-only

.PHONY: test
test:
	go test -race -v ./...
	go vet $(PKGS)
	golint $(PKGS)
	# megacheck -unused.exported $(PKGS)
	# -ignore "github.com/golang/dep/internal/test/test.go:U1000 github.com/golang/dep/gps/prune.go:U1000 github.com/golang/dep/manifest.go:U1000"

.PHONY: format
format:
	go imports ./...

# .PHONY: test
# test: build
# 	./battery check

# .PHONY: install
# install: build
# 	cp ./battery $(GOBIN)

.PHONY: install
install:
	go install $(REPO)

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: snapshot
snapshot:
	goreleaser --rm-dist --snapshot
