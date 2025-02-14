package main

import (
	"fyne.io/fyne/v2/app"
	"repomancer/internal"
	"repomancer/window"
)

func main() {

	state := internal.State{}
	state.App = app.NewWithID("com.sheersky.repomancer")
	state.StartWindow = screens.NewStartScreen(&state)
	state.NewProjectWindow = screens.NewAddProjectScreen(&state)
	state.SettingsWindow = screens.NewPreferenceScreen(&state)
	state.ProjectWindow = screens.NewProjectScreen(&state)
	state.StartWindow.ShowAndRun()

}
