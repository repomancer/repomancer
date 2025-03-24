package screens

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"repomancer/internal"
	"repomancer/window/widgets"
	"strconv"
	"strings"
)

type SettingsWindow struct {
	locationLbl    *widgets.LabelWithHelp
	locationEntry  *widgets.ShortcutHandlingEntry
	workerCountLbl *widgets.LabelWithHelp
	workerCount    *widgets.ShortcutHandlingEntry
	okButton       *widget.Button
	cancelButton   *widget.Button
}

func NewSettingsWindow(state *internal.State) fyne.Window {
	w := state.NewHideableWindow("Settings")

	p := SettingsWindow{
		locationLbl:    widgets.NewLabelWithHelpWidget("Default Location", "Directory where new projects will be created", w),
		locationEntry:  widgets.NewShortcutHandlingEntry(w, false),
		workerCountLbl: widgets.NewLabelWithHelpWidget("Workers", "Number of concurrent workers used for running shell commands\nTakes effect when Project is opened", w),
		workerCount:    widgets.NewShortcutHandlingEntry(w, false),
		okButton:       widget.NewButton("Save", nil),
		cancelButton:   widget.NewButton("Cancel", nil),
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	existing := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
	p.locationEntry.Text = existing

	existingWorkerCount := fyne.CurrentApp().Preferences().IntWithFallback("workerCount", 5)
	p.workerCount.SetText(strconv.Itoa(existingWorkerCount))
	p.workerCount.Validator = func(s string) error {
		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		} else {
			if val < 1 || val > 50 {
				return errors.New("worker count must be between 1 and 50")
			}
		}
		return nil
	}

	p.okButton.OnTapped = func() {
		p.okButton.Disable()

		newLocation := strings.TrimSpace(p.locationEntry.Text)
		if !strings.HasSuffix(newLocation, "/") {
			newLocation = newLocation + "/"
		}

		fyne.CurrentApp().Preferences().SetString("baseDirectory", newLocation)
		log.Printf("Setting baseDirectory to: %s", newLocation)

		workers, _ := strconv.Atoi(p.workerCount.Text)
		fyne.CurrentApp().Preferences().SetInt("workerCount", workers)
		p.okButton.Enable()
		w.Hide()
	}

	p.cancelButton.OnTapped = func() { w.Hide() }

	form := container.New(layout.NewFormLayout(), p.locationLbl, p.locationEntry, p.workerCountLbl, p.workerCount)

	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(container.NewVBox(form, widget.NewSeparator(), container.NewGridWithColumns(2, p.okButton, p.cancelButton)))
	return w
}
