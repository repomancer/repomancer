package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"log"
	"os"
	"repomancer/internal"
	"repomancer/window"
)

var topWindow fyne.Window

func main() {

	state := internal.State{}
	state.App = app.NewWithID("com.sheersky.repomancer")

	a := app.NewWithID("com.sheersky.repomancer")
	a.SetIcon(data.FyneLogo)
	screens.LogLifecycle(a)
	w := a.NewWindow("Fyne Demo")
	topWindow = w

	w.SetMainMenu(screens.MakeMenu(a, w))
	w.SetMaster()

	if len(os.Args) < 2 {
		w.SetContent(screens.NewStartScreen(a, w))
		w.Resize(fyne.NewSize(400, 600))
	} else {
		project, err := internal.OpenProject(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		window := screens.NewProjectWindow(&state, project)
		window.ShowAndRun()
	}
	w.ShowAndRun()

}
