package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"log"
	"os"
	"repomancer/internal"
	"repomancer/screen"
)

func main() {
	a := app.NewWithID("com.sheersky.repomancer")
	a.SetIcon(data.FyneLogo)
	screens.LogLifecycle(a)
	w := a.NewWindow("Repomancer")

	w.SetMainMenu(screens.MakeMenu(a, w))
	w.SetMaster()

	if len(os.Args) < 2 {
		screens.GotoStartScreen(a, w)
	} else {
		project, err := internal.OpenProject(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		screens.GotoProjectScreen(w, project)

	}
	w.ShowAndRun()

}
