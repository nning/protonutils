.PHONY: clean

cmd/list-proton-versions/list-proton-versions: cmd/list-proton-versions/main.go
	cd cmd/list-proton-versions; go build

clean:
	rm -f cmd/list-proton-versions/list-proton-versions
