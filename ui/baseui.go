package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/repomancer/repomancer/internal"
	"log"
)

type WindowType string

const (
	Start    WindowType = "start"
	Project             = "project"
	Settings            = "settings"
	About               = "about"
)

type BaseUI struct {
	App           fyne.App
	windows       map[WindowType]fyne.Window
	projectWindow fyne.Window
	startWindow   fyne.Window
	aboutWindow   fyne.Window
	mainMenu      *fyne.MainMenu
	ctx           context.Context
}

func NewBaseUI(app fyne.App) *BaseUI {
	b := &BaseUI{
		App:     app,
		ctx:     context.Background(),
		windows: map[WindowType]fyne.Window{},
	}
	return b
}

func (b *BaseUI) NewWindow(title string) fyne.Window {
	w := b.App.NewWindow(title)
	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
		w.Close()
	})
	return w
}

func (b *BaseUI) ShowWindow(name WindowType) {
	log.Printf("Showing window: %s", name)
	_, ok := b.windows[name]
	if !ok {
		switch name {
		case Start:
			b.windows[name] = NewStartScreen(b)
		case Settings:
			b.windows[name] = NewSettingsWindow(b)
		case About:
			b.windows[name] = NewAboutWindow(b)
		case Project:
			log.Fatalf("Use ShowProjectWindow() to open project window, not ShowWindow()")
		default:
			log.Fatalf("Unknown window name: %s", name)
		}
	}
	if b.mainMenu == nil {
		// main menu hasn't been created before, create and add it to this window
		b.mainMenu = createMainMenu(b)
		b.windows[name].SetMainMenu(b.mainMenu)
	}
	b.windows[name].Show()
	b.windows[name].RequestFocus()
	b.windows[name].SetOnClosed(func() {
		delete(b.windows, name)
	})
}

func (b *BaseUI) ShowProjectWindow(project *internal.Project) {
	log.Printf("Showing project window for project: %s", project.Name)
	if b.projectWindow == nil {
		w := NewProjectWindow(b, project)
		b.projectWindow = *w
	}
	b.projectWindow.SetOnClosed(func() {
		b.projectWindow = nil
	})
	b.projectWindow.Show()
	b.projectWindow.RequestFocus()
	b.CloseWindow(Start)
}

// CloseWindow closes window with given name if it is open
func (b *BaseUI) CloseWindow(name WindowType) {
	log.Printf("Closing window: %s", name)
	w, ok := b.windows[name]
	if ok {
		w.Close()
	}
}
