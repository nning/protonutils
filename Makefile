.PHONY: clean test lint

SOURCES = $(shell find . -name \*.go)

LIST_BIN_DIR = cmd/list-proton-versions
LIST_BIN_FILE = list-proton-versions
LIST_BIN = $(LIST_BIN_DIR)/$(LIST_BIN_FILE)

VERSION := $(shell ./build/version.sh)
BUILDTIME := $(shell date -u +"%Y%m%d%H%M%S")

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

build: $(LIST_BIN)

$(LIST_BIN): $(SOURCES)
	cd $(LIST_BIN_DIR); go build $(GOFLAGS)

clean:
	rm -f $(LIST_BIN) ./list-proton-versions-*

run: $(LIST_BIN)
	./$(LIST_BIN)

test:
	go test ./...

lint:
	golint ./...

release: build
	cp $(LIST_BIN) $(LIST_BIN_FILE)-$(VERSION)
