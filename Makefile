.PHONY: clean test lint

PREFIX = ~/.local/bin
MAN_PREFIX = ~/.local/share/man
ZSH_PREFIX = ~/.local/share/zsh/functions

SOURCES = $(shell find . -name \*.go)

UTILS_BIN_DIR = cmd/protonutils
UTILS_BIN_FILE = protonutils
UTILS_BIN = $(UTILS_BIN_DIR)/$(UTILS_BIN_FILE)
COMPLETION_ZSH_SRC = completion.zsh
MAN_SRC = man1

VERSION = $(shell ./build/version.sh)
BUILDTIME = $(shell date -u +"%Y%m%d%H%M%S")

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS += -ldflags "$(GOLDFLAGS)"

build: $(UTILS_BIN)

$(UTILS_BIN): $(SOURCES)
	cd $(UTILS_BIN_DIR); go build $(GOFLAGS)

$(UTILS_BIN_FILE): $(UTILS_BIN)

clean:
	rm -f $(UTILS_BIN) $(COMPLETION_ZSH_SRC)
	rm -rf $(MAN_SRC)

run: run_utils

run_utils: $(UTILS_BIN)
	./$(UTILS_BIN) $(args)

test:
	go test -cover -coverprofile .coverage ./...

coverage: test
	go tool cover -html .coverage

lint:
	golint ./...

build_pie: GOLDFLAGS += -s -w -linkmode external -extldflags \"$(LDFLAGS)\"
build_pie: GOFLAGS += -trimpath -buildmode=pie -mod=readonly -modcacherw
build_pie: build

completion: build
	$(UTILS_BIN) completion zsh > $(COMPLETION_ZSH_SRC)

man: build
	$(UTILS_BIN) -m $(MAN_SRC)

release: build_pie
	upx -qq --best $(UTILS_BIN)
	ls -lh $(UTILS_BIN)

install: build_pie completion man
	mkdir -p $(PREFIX) $(ZSH_PREFIX) $(MAN_PREFIX)
	cp $(UTILS_BIN) $(PREFIX)
	cp $(COMPLETION_ZSH_SRC) $(ZSH_PREFIX)/_protonutils
	cp -r $(MAN_SRC) $(MAN_PREFIX)/
