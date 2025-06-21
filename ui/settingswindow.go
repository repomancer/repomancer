package ui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/repomancer/repomancer/ui/widgets"
	"log"
	"os"
	"strconv"
	"strings"
)

func NewSettingsWindow(ui *BaseUI) fyne.Window {
	w := ui.NewWindow("Settings")
	locationLbl := widgets.NewLabelWithHelpWidget("Default Location", "Directory where new projects will be created", w)
	locationEntry := widget.NewEntry()
	workerCountLbl := widgets.NewLabelWithHelpWidget("Workers", "Number of concurrent workers used for running shell commands\nTakes effect when Project is opened", w)
	workerCount := widget.NewEntry()
	shell := widget.NewEntry()
	shellLbl := widgets.NewLabelWithHelpWidget("Shell", "Commands will be executed inside this shell", w)
	shellArgs := widget.NewEntry()
	shellArgsLbl := widgets.NewLabelWithHelpWidget("Shell Arguments", "Shell arguments", w)
	okButton := widget.NewButton("Save", nil)
	cancelButton := widget.NewButton("Cancel", nil)
	errorLabel := widget.NewLabel("")

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape {
			w.Close()
		}
		if key.Name == fyne.KeyEnter || key.Name == fyne.KeyReturn {
			saveSettings(okButton, locationEntry, workerCount, shell, shellArgs, errorLabel, w)
		}
	})

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	existing := ui.App.Preferences().StringWithFallback("baseDirectory", dirname)
	locationEntry.Text = existing

	existingWorkerCount := ui.App.Preferences().IntWithFallback("workerCount", 5)
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

	shell.SetText(ui.App.Preferences().StringWithFallback("shell", "/bin/zsh"))
	shell.Validator = func(s string) error {
		_, err := os.Stat(s)
		if err != nil {
			return err
		}
		return nil
	}
	shellArgs.SetText(ui.App.Preferences().StringWithFallback("shellArguments", "--login --interactive"))
	okButton.OnTapped = func() {
		saveSettings(okButton, locationEntry, workerCount, shell, shellArgs, errorLabel, w)
	}

	cancelButton.OnTapped = func() {
		w.Close()
	}

	form := container.New(layout.NewFormLayout(), locationLbl, locationEntry, workerCountLbl, workerCount, shellLbl, shell, shellArgsLbl, shellArgs)

	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(container.NewVBox(form, widget.NewSeparator(), container.NewGridWithColumns(2, okButton, cancelButton), errorLabel))
	return w
}

func saveSettings(okButton *widget.Button, locationEntry *widget.Entry, workerCount *widget.Entry, shell *widget.Entry, shellArgs *widget.Entry, errorLabel *widget.Label, w fyne.Window) {
	err1 := locationEntry.Validate()
	err2 := workerCount.Validate()
	err3 := shell.Validate()
	err4 := shellArgs.Validate()
	if err1 != nil {
		errorLabel.SetText(fmt.Sprint(err1))
		return
	}
	if err2 != nil {
		errorLabel.SetText(fmt.Sprint(err2))
		return
	}
	if err3 != nil {
		errorLabel.SetText(fmt.Sprint(err3))
		return
	}
	if err4 != nil {
		errorLabel.SetText(fmt.Sprint(err4))
		return
	}
	newLocation := strings.TrimSpace(locationEntry.Text)
	if !strings.HasSuffix(newLocation, "/") {
		newLocation = newLocation + "/"
	}

	fyne.CurrentApp().Preferences().SetString("baseDirectory", newLocation)
	log.Printf("Setting baseDirectory to: %s", newLocation)

	workers, _ := strconv.Atoi(workerCount.Text)
	fyne.CurrentApp().Preferences().SetInt("workerCount", workers)

	fyne.CurrentApp().Preferences().SetString("shell", shell.Text)
	fyne.CurrentApp().Preferences().SetString("shellArguments", shellArgs.Text)

	okButton.Enable()
	w.Close()
}
