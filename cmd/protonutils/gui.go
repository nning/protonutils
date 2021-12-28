package main

import (
	"log"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start GUI",
	Run:   gui,
}

func init() {
	rootCmd.AddCommand(guiCmd)
}

func gui(cmd *cobra.Command, args []string) {
	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "protonutils",
		BaseDirectoryPath: "example",
	})
	exitOnError(err)
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	err = a.Start()
	exitOnError(err)

	d := func(src string) ([]byte, error) {
		l.Println("TEST", src)
		return []byte{}, nil
	}

	astilectron.NewDisembedderProvisioner(d, "example", "example", l)

	// New window
	var w *astilectron.Window
	w, err = a.NewWindow("example/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	})
	exitOnError(err)

	// Create windows
	err = w.Create()
	exitOnError(err)

	// Blocking pattern
	a.Wait()
}
