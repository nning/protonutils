.PHONY: clean
SOURCES = $(shell find . -name \*.go)

build: cmd/list-proton-versions/list-proton-versions

cmd/list-proton-versions/list-proton-versions: $(SOURCES)
	cd cmd/list-proton-versions; go build

clean:
	rm -f cmd/list-proton-versions/list-proton-versions
