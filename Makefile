UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
target := linux
endif
ifeq ($(UNAME), Darwin)
target := darwin
endif

COMMIT_HASH=$$(git rev-list -1 HEAD)
TAG_VERSION=$$(git tag --sort=committerdate | tail -1)

.PHONY: test
test:
	go test -count=1 -race -cover -v $(shell go list ./... | grep -v /vendor/)

.PHONY: ci
ci: lint test

.PHONY: go-lint-install
go-lint-install:
	go get -u golang.org/x/lint/golint
	cp $$(go list -f {{.Target}} golang.org/x/lint/golint) /usr/local/bin

.PHONY: lint
lint:
	golint -set_exit_status `go list ./...`

.PHONY: build
build:
	GOOS=$(target) go build -o="bin/exporter" -ldflags="\
	-X 'main.commitHash=$(COMMIT_HASH)' \
	-X 'main.version=$(TAG_VERSION)' \
	-X 'main.date=$$(date)'" \
	cmd/main.go cmd/config.go cmd/version.go

.PHONY: build-linux
build-linux:
	GOOS=linux go build -o="bin/exporter" -ldflags="\
	-X 'main.commitHash=$(COMMIT_HASH)' \
	-X 'main.version=$(TAG_VERSION)' \
	-X 'main.date=$$(date)'" \
	cmd/main.go cmd/config.go cmd/version.go

vendor:
	go mod vendor

tidy:
	go mod tidy