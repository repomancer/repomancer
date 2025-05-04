package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"github.com/repomancer/repomancer/internal"
	"github.com/repomancer/repomancer/screen"
	"log"
	"os"
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
