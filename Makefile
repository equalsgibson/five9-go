.PHONY: full build test test-go lint lint-go fix fix-go watch clean docs-go

SHELL=/bin/bash -o pipefail
$(shell git config core.hooksPath ops/git-hooks)
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := /usr/local/bin:$(GO_PATH)/bin:$(PATH)

full: clean lint test build

## Build the project
build:

## Test the project
test: test-go

test-go:
	@mkdir -p five9/ops/docs/coverage/
	@go install github.com/boumenot/gocover-cobertura@latest
	go test -p 1 -count=1 -cover -coverprofile five9/ops/docs/coverage/coverage-profile.txt ./...
	@go tool cover -func five9/ops/docs/coverage/coverage-profile.txt | awk '/^total/{print "{\"total\":\""$$3"\"}"}' > five9/ops/docs/coverage/coverage.json
	@go tool cover -html five9/ops/docs/coverage/coverage-profile.txt -o five9/ops/docs/coverage/coverage.html
	@gocover-cobertura < five9/ops/docs/coverage/coverage-profile.txt > five9/ops/docs/coverage/coverage-cobertura.xml

## Lint the project
lint: lint-go

lint-go:
	go get -d ./...
	go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

## Fix the project
fix: fix-go

fix-go:
	go mod tidy
	gofmt -s -w .

## Watch the project
watch:

## Clean the project
clean:
	git clean -Xdff --exclude="!.env*local"

## Run the docs server for the project
docs-go:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "listening on http://127.0.0.1:6060/pkg/github.com/equalsgibson/five9-go"
	@godoc -http=127.0.0.1:6060
