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
	nameEntry      *widgets.ShortcutHandlingEntry
	locationLbl    *widgets.LabelWithHelp
	locationEntry  *widgets.ShortcutHandlingEntry
	prMessageLbl   *widget.Label
	prMessageEntry *widgets.ShortcutHandlingEntry
	statusMessage  *widget.Label
	okButton       *widget.Button
	cancelButton   *widget.Button
}

func (p *AddProjectScreen) Validate() bool {
	var errors []string
	_, err := os.Stat(p.locationEntry.Text)
	if !os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("%s already exists", p.locationEntry.Text))
	}

	if p.nameEntry.Text == "" {
		errors = append(errors, "Name is required")
	}
	if p.locationEntry.Text == "" {
		errors = append(errors, "Location is required")
	}
	if p.prMessageEntry.Text == "" {
		errors = append(errors, "Pull Request Message is required")
	}

	if len(errors) == 0 {
		p.okButton.Enable()
		p.statusMessage.Text = ""
		p.statusMessage.Refresh()
		return true
	} else {
		p.okButton.Disable()
		p.statusMessage.Text = strings.Join(errors, "\n")
		p.statusMessage.Refresh()
		return false
	}
}

func NewAddProjectScreen(state *internal.State) fyne.Window {
	w := state.NewHideableWindow("New Project")
	p := AddProjectScreen{
		nameLbl:        widgets.NewLabelWithHelpWidget("Name", "Project Name. Also used for the name of the git branch.\nValid characters: [A-Za-z0-9_-]", w),
		nameEntry:      widgets.NewShortcutHandlingEntry(w, false),
		locationLbl:    widgets.NewLabelWithHelpWidget("Location", "Where project data and cloned repositories will be stored. Must not exist.", w),
		locationEntry:  widgets.NewShortcutHandlingEntry(w, false),
		prMessageLbl:   widget.NewLabel("Pull Request\nMessage"),
		prMessageEntry: widgets.NewShortcutHandlingEntry(w, false),
		statusMessage:  widget.NewLabel(""),
		okButton:       widget.NewButton("Create", nil),
		cancelButton:   widget.NewButton("Cancel", nil),
	}
	p.nameEntry.AllowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789_-"
	p.nameEntry.MaxLength = 50
	p.statusMessage.Wrapping = fyne.TextWrapWord
	p.prMessageEntry.Wrapping = fyne.TextWrapWord
	p.prMessageEntry.MultiLine = true
	p.okButton.Disable()

	p.nameEntry.OnChanged = func(s string) {
		p.Validate()
	}
	p.locationEntry.OnChanged = func(s string) {
		p.Validate()
	}
	p.prMessageEntry.OnChanged = func(s string) {
		p.Validate()
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	basePath := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
	p.locationEntry.Text = basePath

	p.nameEntry.OnChanged = func(value string) {
		p.locationEntry.SetText(basePath + strings.Replace(value, "/", "_", -1))
		p.locationEntry.Refresh()
	}

	p.okButton.OnTapped = func() {
		p.statusMessage.Text = p.nameEntry.Text
		p.okButton.Disable()
		//project, err := internal.CreateProject(p.nameEntry.Text, p.prMessageEntry.Text, p.locationEntry.Text)
		if err != nil {
			p.statusMessage.Text = err.Error()
			p.statusMessage.Refresh()
			p.okButton.Enable()
		} else {
			//AddProjectScreen(window, project)
		}
	}

	p.cancelButton.OnTapped = func() { w.Close() }

	grid := container.New(layout.NewFormLayout(), p.nameLbl, p.nameEntry, p.locationLbl, p.locationEntry, p.prMessageLbl, p.prMessageEntry)

	w.Resize(fyne.NewSize(600, 600))
	w.SetContent(container.NewVBox(grid, p.statusMessage, p.okButton, p.cancelButton))
	w.Canvas().Focus(p.nameEntry)
	return w
}
