.PHONY: all
all: battery

REPO := github.com/LEI/battery

.PHONY: battery
battery: vendor test install

.PHONY: vendor
vendor:
	dep ensure -vendor-only

.PHONY: test
test:
	go test -race -v ./...

.PHONY: install
install:
	go install $(REPO)

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: snapshot
snapshot:
	goreleaser --rm-dist --snapshot
