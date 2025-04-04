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
	a := app.NewWithID("com.sheersky.repomancer")
	a.SetIcon(data.FyneLogo)
	screens.LogLifecycle(a)
	w := a.NewWindow("Fyne Demo")
	topWindow = w

	w.SetMainMenu(screens.MakeMenu(a, w))
	w.SetMaster()

	if len(os.Args) < 2 {
		w.SetContent(screens.NewStartScreen(a, w))
	} else {
		project, err := internal.OpenProject(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		w.SetContent(screens.NewProjectWindow(w, project))

	}
	w.ShowAndRun()

}
