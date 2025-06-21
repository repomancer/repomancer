package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"github.com/repomancer/repomancer/internal"
	"log"
	"os"

	"github.com/repomancer/repomancer/ui"
)

func main() {
	a := app.NewWithID("com.sheersky.repomancer")
	a.SetIcon(data.FyneLogo)
	myUi := ui.NewBaseUI(a)
	if len(os.Args) < 2 {
		myUi.ShowWindow(ui.Start)
	} else {
		project, err := internal.OpenProject(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		myUi.ShowProjectWindow(project)
	}
	a.Run()

}
