package main

import (
	"fyne.io/fyne/v2/app"
	"log"
	"os"
	"repomancer/internal"
	"repomancer/window"
)

func main() {

	state := internal.State{}
	state.App = app.NewWithID("com.sheersky.repomancer")
	state.StartWindow = screens.NewStartScreen(&state)
	state.NewProjectWindow = screens.NewAddProjectScreen(&state)
	state.SettingsWindow = screens.NewPreferenceScreen(&state)

	if len(os.Args) < 2 {
		state.StartWindow.ShowAndRun()
	} else {
		project, err := internal.OpenProject(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		window := screens.NewProjectWindow(&state, project)
		window.ShowAndRun()
	}
}
