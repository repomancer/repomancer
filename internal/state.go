package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"log"
)

type State struct {
	App              fyne.App
	StartWindow      fyne.Window
	OpenWindow       fyne.Window
	NewProjectWindow fyne.Window
	SettingsWindow   fyne.Window
	ProjectWindow    fyne.Window
	Project          *Project
}

func (state *State) ShowSettingsWindow() {
	state.SettingsWindow.Show()
}

func (state *State) ShowProjectWindow() {
	state.ProjectWindow.Show()
	state.StartWindow.Hide()
}

func (state *State) NewHideableWindow(title string) fyne.Window {
	w := state.App.NewWindow(title)

	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}

	w.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
		log.Println("We tapped Cmd+W")
		w.Hide()
	})
	return w
}

func (state *State) NewQuitWindow(title string) fyne.Window {
	w := state.App.NewWindow(title)

	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}

	w.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
		w.Close()
	})
	return w
}
