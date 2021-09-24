package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/andygrunwald/vdf"
)

func get_config_path() string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	return filepath.Join(dir, ".steam/root/config/config.vdf")
}

func lookup(m map[string]interface{}, x ...string) (map[string]interface{}, error) {
	y := m

	for _, s := range x {
		if y[s] == nil {
			return nil, errors.New("Key not found: " + s)
		} else {
			y = y[s].(map[string]interface{})
		}
	}

	return y, nil
}

func main() {
	path := get_config_path()

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	x, err := lookup(m, "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
	if err != nil {
		log.Fatal(err)
	}

	for id := range x {
		fmt.Println(id)
	}
}
