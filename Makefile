VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

FLAGS_INTEL64 = GOARCH=amd64
NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(FLAGS_INTEL) GOOS=linux
FLAGS_MAC = $(FLAGS_INTEL) GOOS=darwin
FLAGS_MAC_ARM = $GOARCH=arm64 GOOS=darwin
FLAGS_WIN = $(FLAGS_INTEL64) GOOS=windows
FLAGS_LD=-ldflags "-s -w -X github.com/sfgrp/lognsq.Build=${DATE} \
                  -X github.com/sfgrp/lognsq.Version=${VERSION}"
GOCMD = go
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get

RELEASE_DIR ?= "/tmp"
BUILD_DIR ?= "."
CLIB_DIR ?= "."

all: install

test: deps install
	$(FLAG_INTEL) go test -race ./...

test-build: deps build

deps:
	$(GOCMD) mod download;

tools: deps
	@echo Installing tools from tools.go
	@cat lognsq/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build:
	cd lognsq; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOBUILD) -o $(BUILD_DIR)

install:
	cd lognsq; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOINSTALL)

release:
	cd lognsq; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/lognsq-$(VER)-linux.tar.gz lognsq; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/lognsq-$(VER)-mac.tar.gz lognsq; \
	$(GOCLEAN); \
	$(FLAGS_MAC_ARM) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/lognsq-$(VER)-mac-arm64.tar.gz lognsq; \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(NO_C) $(GOBUILD); \
	zip -9 $(RELEASE_DIR)/lognsq-$(VER)-win-64.zip lognsq.exe; \
	$(GOCLEAN);
