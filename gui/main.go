package gui

import (
	"embed"
	"encoding/json"

	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/steam"
	"github.com/wailsapp/wails"
)

type versionResult struct {
	Version   string `json:"version"`
	Buildtime string `json:"buildtime"`
	URL       string `json:"url"`
}

//go:embed public/build
var files embed.FS

type context struct {
	user        string
	cfg         *config.Config
	ignoreCache bool
}

var ctx *context

func errorToJSON(err error) string {
	return "{error: \"" + err.Error() + "\"}"
}

func list() string {
	s, err := steam.New(ctx.user, ctx.cfg.SteamRoot, ctx.ignoreCache)
	if err != nil {
		return errorToJSON(err)
	}

	err = s.ReadCompatToolVersions()
	if err != nil {
		return errorToJSON(err)
	}

	j, err := json.MarshalIndent(s.CompatToolVersions, "", "  ")
	if err != nil {
		return errorToJSON(err)
	}

	err = s.SaveCache()
	if err != nil {
		return errorToJSON(err)
	}

	return string(j)
}

func Run(user string, cfg *config.Config, ignoreCache bool) {
	name := "public/build/bundle"

	js, err := files.ReadFile(name + ".js")
	if err != nil {
		panic(err)
	}

	css, err := files.ReadFile(name + ".css")
	if err != nil {
		panic(err)
	}

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "protonutils",
		JS:     string(js),
		CSS:    string(css),
		Colour: "#ffffff",
	})

	ctx = &context{user, cfg, ignoreCache}

	app.Bind(list)

	app.Run()
}
