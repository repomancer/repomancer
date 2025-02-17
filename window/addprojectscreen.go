package screens

import (
	"fmt"
	"repomancer/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"repomancer/window/widgets"
	"strings"
)

type AddProjectScreen struct {
	nameLbl        *widgets.LabelWithHelp
	NameEntry      *widgets.ShortcutHandlingEntry
	locationLbl    *widgets.LabelWithHelp
	LocationEntry  *widgets.ShortcutHandlingEntry
	prMessageLbl   *widget.Label
	PrMessageEntry *widgets.ShortcutHandlingEntry
	statusMessage  *widget.Label
	OkButton       *widget.Button
	CancelButton   *widget.Button
}

func (p *AddProjectScreen) Validate() bool {
	var errors []string
	_, err := os.Stat(p.LocationEntry.Text)
	if !os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("%s already exists", p.LocationEntry.Text))
	}

	if p.NameEntry.Text == "" {
		errors = append(errors, "Name is required")
	}
	if p.LocationEntry.Text == "" {
		errors = append(errors, "Location is required")
	}
	if p.PrMessageEntry.Text == "" {
		errors = append(errors, "Pull Request Message is required")
	}

	if len(errors) == 0 {
		p.OkButton.Enable()
		p.statusMessage.Text = ""
		p.statusMessage.Refresh()
		return true
	} else {
		p.OkButton.Disable()
		p.statusMessage.Text = strings.Join(errors, "\n")
		p.statusMessage.Refresh()
		return false
	}
}

func NewAddProjectScreen(state *internal.State) fyne.Window {
	w := state.NewHideableWindow("New Project")
	p := AddProjectScreen{
		nameLbl:        widgets.NewLabelWithHelpWidget("Name", "Project Name. Also used for the name of the git branch.\nValid characters: [A-Za-z0-9_-]", w),
		NameEntry:      widgets.NewShortcutHandlingEntry(w, false),
		locationLbl:    widgets.NewLabelWithHelpWidget("Location", "Where project data and cloned repositories will be stored. Must not exist.", w),
		LocationEntry:  widgets.NewShortcutHandlingEntry(w, false),
		prMessageLbl:   widget.NewLabel("Pull Request\nMessage"),
		PrMessageEntry: widgets.NewShortcutHandlingEntry(w, false),
		statusMessage:  widget.NewLabel(""),
		OkButton:       widget.NewButton("Create", nil),
		CancelButton:   widget.NewButton("Cancel", nil),
	}
	p.NameEntry.AllowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789_-"
	p.NameEntry.MaxLength = 50
	p.statusMessage.Wrapping = fyne.TextWrapWord
	p.PrMessageEntry.Wrapping = fyne.TextWrapWord
	p.PrMessageEntry.MultiLine = true
	p.OkButton.Disable()

	p.NameEntry.OnChanged = func(s string) {
		p.Validate()
	}
	p.LocationEntry.OnChanged = func(s string) {
		p.Validate()
	}
	p.PrMessageEntry.OnChanged = func(s string) {
		p.Validate()
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	basePath := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
	p.LocationEntry.Text = basePath

	p.NameEntry.OnChanged = func(value string) {
		p.LocationEntry.SetText(basePath + strings.Replace(value, "/", "_", -1))
		p.LocationEntry.Refresh()
	}

	p.OkButton.OnTapped = func() {
		p.statusMessage.Text = p.NameEntry.Text
		p.OkButton.Disable()
		project, err := internal.CreateProject(p.NameEntry.Text, p.PrMessageEntry.Text, p.LocationEntry.Text)
		if err != nil {
			p.statusMessage.Text = err.Error()
			p.statusMessage.Refresh()
			p.OkButton.Enable()
		} else {
			window := NewProjectWindow(state, project)
			window.Show()
			//AddProjectScreen(window, project)
		}
	}

	p.CancelButton.OnTapped = func() { w.Close() }

	grid := container.New(layout.NewFormLayout(), p.nameLbl, p.NameEntry, p.locationLbl, p.LocationEntry, p.prMessageLbl, p.PrMessageEntry)

	w.Resize(fyne.NewSize(600, 600))
	w.SetContent(container.NewVBox(grid, p.statusMessage, p.OkButton, p.CancelButton))
	w.Canvas().Focus(p.NameEntry)
	return w
}
