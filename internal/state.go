package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type State struct {
	App fyne.App
}

func (state *State) NewHideableWindow(title string) fyne.Window {
	w := state.App.NewWindow(title)

	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}

	w.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
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
