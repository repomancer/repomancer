package screens

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"repomancer/screen/widgets"
	"strconv"
	"strings"
)

var settingsWindow fyne.Window

func NewSettingsWindow(a fyne.App) fyne.Window {
	if settingsWindow != nil {
		return settingsWindow
	}
	w := a.NewWindow("Settings")
	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
		w.Close()
	})
	w.SetOnClosed(func() {
		settingsWindow = nil
		log.Println("Settings screen closed")
	})

	locationLbl := widgets.NewLabelWithHelpWidget("Default Location", "Directory where new projects will be created", w)
	locationEntry := widget.NewEntry()
	workerCountLbl := widgets.NewLabelWithHelpWidget("Workers", "Number of concurrent workers used for running shell commands\nTakes effect when Project is opened", w)
	workerCount := widget.NewEntry()
	okButton := widget.NewButton("Save", nil)
	cancelButton := widget.NewButton("Cancel", nil)

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape {
			w.Close()
		}
		if key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
			saveSettings(okButton, locationEntry, workerCount, w)
		}
	})

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	existing := a.Preferences().StringWithFallback("baseDirectory", dirname)
	locationEntry.Text = existing

	existingWorkerCount := a.Preferences().IntWithFallback("workerCount", 5)
	workerCount.SetText(strconv.Itoa(existingWorkerCount))
	workerCount.Validator = func(s string) error {
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

	okButton.OnTapped = func() {
		saveSettings(okButton, locationEntry, workerCount, w)
	}

	cancelButton.OnTapped = func() {
		w.Close()
	}

	form := container.New(layout.NewFormLayout(), locationLbl, locationEntry, workerCountLbl, workerCount)

	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(container.NewVBox(form, widget.NewSeparator(), container.NewGridWithColumns(2, okButton, cancelButton)))
	settingsWindow = w
	return w
}

func saveSettings(okButton *widget.Button, locationEntry *widget.Entry, workerCount *widget.Entry, w fyne.Window) {
	okButton.Disable()

	newLocation := strings.TrimSpace(locationEntry.Text)
	if !strings.HasSuffix(newLocation, "/") {
		newLocation = newLocation + "/"
	}

	fyne.CurrentApp().Preferences().SetString("baseDirectory", newLocation)
	log.Printf("Setting baseDirectory to: %s", newLocation)

	workers, _ := strconv.Atoi(workerCount.Text)
	fyne.CurrentApp().Preferences().SetInt("workerCount", workers)
	okButton.Enable()
	w.Close()
}
