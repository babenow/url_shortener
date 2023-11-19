PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
CONFIG_PATH = $(PROJECT_DIR)/config/local.yaml
PATH := $(PROJECT_BIN):$(PATH)
PATH := $(CONFIG_PATH):$(PATH)

BINARY = urlsh
GOLANGCI_LINT_VERSION = v1.55.0

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint
ifeq ($(OS), Windown_NT)
	BINARY := $(BINARY).exe
GOLANGCI_LINT := $(GOLANGCI_LINT).exe
endif

$(shell [ -f $(PROJECT_BIN) ] || mkdir -p $(PROJECT_BIN))



.PHONY:build
build:lint
	go build -o $(PROJECT_BIN)/$(BINARY) ./cmd/urlsh

.PHONY:run
run:build
	CONFIG_PATH=$(CONFIG_PATH) $(BINARY)


.PHONY:test
test:
	CONFIG_PATH=$(CONFIG_PATH) go test -v --timeout 30s ./...

.PHONY:.install-linter
.install-linter:
	$(shell [ -f $(GOLANGCI_LINT) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) $(GOLANGCI_LINT_VERSION))

.PHONY:lint
lint:.install-linter
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY:lint-fast
lint-fast:
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml


.DEFAULT_GOAL := build



