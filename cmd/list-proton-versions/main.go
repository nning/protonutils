package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/andygrunwald/vdf"
	"github.com/nning/list_proton_versions/appid"

	. "github.com/nning/list_proton_versions"
)

func get_config_path() string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	return path.Join(dir, ".steam", "root", "config", "config.vdf")
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
	PanicOnError(err)

	p := vdf.NewParser(f)
	m, err := p.Parse()
	PanicOnError(err)

	x, err := lookup(m, "InstallConfigStore", "Software", "Valve", "Steam", "CompatToolMapping")
	PanicOnError(err)

	versions := make(map[string][]string)
	appid := appid.New()

	for id, cfg := range x {
		if id == "0" {
			continue
		}

		v := cfg.(map[string]interface{})["name"].(string)

		if v == "" {
			v = "Default"
		}

		if versions[v] == nil {
			versions[v] = make([]string, 0)
		}

		name := appid.GetName(id)
		if name != "💩" {
			versions[v] = append(versions[v], name+" ("+id+")")
		}
	}

	for version, games := range versions {
		fmt.Println(version)

		for _, game := range games {
			fmt.Println("\t" + game)
		}

		fmt.Println()
	}

	appid.Write()
}
