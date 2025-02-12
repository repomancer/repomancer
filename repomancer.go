package main

import (
	"fyne.io/fyne/v2/app"
	"repomancer/window/splash"
)

func main() {

	a := app.NewWithID("com.sheersky.repomancer")
	w := a.NewWindow("Repomancer")

	screen := splash.NewStartScreen(w)

	w.SetContent(screen)
	w.ShowAndRun()
}
