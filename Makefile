.PHONY: clean

SOURCES = $(shell find . -name \*.go)

LIST_BIN_DIR = "cmd/list-proton-versions"
LIST_BIN = "$(LIST_BIN_DIR)/list-proton-versions"

build: $(LIST_BIN)

$(LIST_BIN): $(SOURCES)
	cd $(LIST_BIN_DIR); go build

clean:
	rm -f $(LIST_BIN)

run: $(LIST_BIN)
	./$(LIST_BIN)
