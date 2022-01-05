package gui

import (
	_ "embed"
	"encoding/json"

	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/steam"
	"github.com/wailsapp/wails"
)

type result struct {
	Text   string `json:"text"`
	Number uint64 `json:"number"`
}

//go:embed public/build/bundle.js
var js string

//go:embed public/build/bundle.css
var css string

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
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "protonutils",
		JS:     js,
		CSS:    css,
		Colour: "#ffffff",
	})

	ctx = &context{user, cfg, ignoreCache}

	app.Bind(list)

	app.Run()
}
