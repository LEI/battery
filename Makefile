SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

PKGS := $(shell go list ./... | grep -vF /vendor/)

REPO := github.com/LEI/battery

.PHONY: default
default: test format install

.PHONY: dep
DEP := $(shell command -v dep 2>/dev/null)
dep:
ifeq ($(DEP),)
	# Darwin: brew install dep
	curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: test-deps
test-deps:
# go get golang.org/x/tools/cmd/goimports
# go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/lint/golint honnef.co/go/tools/cmd/megacheck

# .PHONY: build
# build:
# 	go fmt ./...
# 	BATTERY_BUILD_PLATFORMS=$(PLATFORM) BATTERY_BUILD_ARCHS=$(ARCH) ./build.sh
# 	cp ./dist/battery-$(PLATFORM)-$(ARCH) battery

.PHONY: vendor
vendor: dep
	dep ensure -vendor-only

.PHONY: test
# gometalinter --fast # --errors
test: test-deps vendor
	go test -race -v ./...
	go vet $(PKGS)
	golint $(PKGS)
	megacheck # -unused.exported -ignore "github.com/LEI/dot/battery.go:U1000" $(PKGS)
# -ignore "github.com/golang/dep/internal/test/test.go:U1000 \
# github.com/golang/dep/gps/prune.go:U1000 \
# github.com/golang/dep/manifest.go:U1000"

.PHONY: format
format:
	go fmt ./...
# goimports -l *.go

# .PHONY: test
# test: build
# 	./battery check

# .PHONY: install
# install: build
# 	cp ./battery $(GOBIN)

.PHONY: install
install: vendor
	go install $(REPO)

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: snapshot
snapshot:
	goreleaser --rm-dist --snapshot
