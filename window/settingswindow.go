package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"repomancer/internal"
	"repomancer/window/widgets"
	"strings"
)

type SettingsWindow struct {
	locationLbl   *widgets.LabelWithHelp
	locationEntry *widgets.ShortcutHandlingEntry
	okButton      *widget.Button
	cancelButton  *widget.Button
}

func NewSettingsWindow(state *internal.State) fyne.Window {
	w := state.NewHideableWindow("Settings")

	p := SettingsWindow{
		locationLbl:   widgets.NewLabelWithHelpWidget("Default Location", "Directory where new projects will be created", w),
		locationEntry: widgets.NewShortcutHandlingEntry(w, false),
		okButton:      widget.NewButton("Save", nil),
		cancelButton:  widget.NewButton("Cancel", nil),
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	existing := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
	p.locationEntry.Text = existing

	p.okButton.OnTapped = func() {
		p.okButton.Disable()

		newLocation := strings.TrimSpace(p.locationEntry.Text)
		if !strings.HasSuffix(newLocation, "/") {
			newLocation = newLocation + "/"
		}

		fyne.CurrentApp().Preferences().SetString("baseDirectory", newLocation)
		log.Printf("Setting baseDirectory to: %s", newLocation)
		p.okButton.Enable()
		w.Hide()
	}

	p.cancelButton.OnTapped = func() { w.Hide() }

	form := container.New(layout.NewFormLayout(), p.locationLbl, p.locationEntry)

	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(container.NewVBox(form, widget.NewSeparator(), container.NewGridWithColumns(2, p.okButton, p.cancelButton)))
	return w
}
