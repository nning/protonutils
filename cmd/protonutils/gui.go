package main

import (
	_ "embed"

	_ "github.com/andlabs/ui/winmanifest"

	"github.com/andlabs/ui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start GUI",
	Run:   guiRun,
}

var mainwin *ui.Window

func init() {
	rootCmd.AddCommand(guiCmd)
	guiCmd.Flags().BoolVarP(&ignoreCache, "ignore-cache", "c", false, "Ignore app ID/name cache")
	guiCmd.Flags().StringVarP(&user, "user", "u", "", "Steam user name (or SteamID3)")
}

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	hbox.Append(ui.NewButton("Button"), false)
	hbox.Append(ui.NewCheckbox("Checkbox"), false)

	vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("Entries")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	entryForm.Append("Entry", ui.NewEntry(), false)
	entryForm.Append("Password Entry", ui.NewPasswordEntry(), false)
	entryForm.Append("Search Entry", ui.NewSearchEntry(), false)
	entryForm.Append("Multiline Entry", ui.NewMultilineEntry(), true)
	entryForm.Append("Multiline Entry No Wrap", ui.NewNonWrappingMultilineEntry(), true)

	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("libui Control Gallery", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Test 1", makeBasicControlsPage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func guiRun(cmd *cobra.Command, args []string) {
	ui.Main(setupUI)
}
