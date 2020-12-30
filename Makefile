export CGO_ENABLED = 0
export GOARCH = amd64
export BINARY = cfn-diff

.DEFAULT_GOAL := help
.PHONY: help
help: ## help
	@echo '  see: https://github.com/kmrok/cfn-diff'
	@echo ''
	@grep -E '^[%/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: install
install: ## install
	go install $(shell pwd)/cmd/cfn-diff

.PHONY: build
build: make build-linux

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux: ## build for Linux
	GOOS=linux go build -o ./cmd/bin/${BINARY}-linux-${GOARCH} ./cmd/cfn-diff/*

.PHONY: build-darwin
build-darwin: ## build for MacOS
	GOOS=darwin go build -o ./cmd/bin/${BINARY}-darwin-${GOARCH} ./cmd/cfn-diff/*

.PHONY: build-windows
build-windows: ## build for Windows
	GOOS=windows go build -o ./cmd/bin/${BINARY}-windows-${GOARCH} ./cmd/cfn-diff/*
