package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"repomancer/window/widgets"
	"strings"
)

type PreferenceScreen struct {
	locationLbl   *widgets.LabelWithHelp
	locationEntry *widgets.ShortcutHandlingEntry
	okButton      *widget.Button
	cancelButton  *widget.Button
}

func NewPreferenceScreen() {
	w := fyne.CurrentApp().NewWindow("Preferences")

	p := PreferenceScreen{
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

		if err != nil {
			p.okButton.Enable()
		} else {
			w.Close()
		}
	}

	p.cancelButton.OnTapped = func() { w.Close() }

	form := container.New(layout.NewFormLayout(), p.locationLbl, p.locationEntry)

	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(container.NewVBox(form, widget.NewSeparator(), container.NewGridWithColumns(2, p.okButton, p.cancelButton)))
	w.Show()
}
