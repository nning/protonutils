.PHONY: clean test lint

PREFIX = ~/.local/bin
MAN_PREFIX = ~/.local/share/man
ZSH_PREFIX = ~/.local/share/zsh/functions

SOURCES = $(shell find . -name \*.go -o -wholename \*build/\*.js -o -wholename \*build/\*.css | grep -v node_modules)
GUI_SOURCES = $(shell find gui/src -name \*.js -o -name \*.svelte)

UTILS_BIN_DIR = cmd/protonutils
UTILS_BIN_FILE = protonutils
UTILS_BIN = $(UTILS_BIN_DIR)/$(UTILS_BIN_FILE)
COMPLETION_ZSH_SRC = completion.zsh
MAN_SRC = man1
GUI_BUNDLE_DIR = gui/public/build
GUI_BUNDLE = $(GUI_BUNDLE_DIR)/bundle.js $(GUI_BUNDLE_DIR)/bundle.css

VERSION = $(shell ./build/version.sh)
BUILDTIME = $(shell date -u +"%Y%m%d%H%M%S")

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS += -ldflags "$(GOLDFLAGS)"

build: $(UTILS_BIN)

$(GUI_BUNDLE): $(GUI_SOURCES)
	cd gui && npm install && npm run build

$(UTILS_BIN): $(SOURCES) $(GUI_BUNDLE)
	cd $(UTILS_BIN_DIR) && go build $(GOFLAGS)

clean:
	rm -f $(UTILS_BIN) $(COMPLETION_ZSH_SRC)
	rm -rf $(MAN_SRC)
	rm -rf $(GUI_BUNDLE_DIR)

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
