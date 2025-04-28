package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/jashort/repomancer/internal"
	"image/color"
	"log"
	"os"
)

func checkRequirements() (string, error) {
	stdout, stderr, err := internal.RunCommand("", 3, "zsh", "-c", "-i", "gh --version")
	if err != nil {
		return stderr, err
	}
	return "Found gh command\n" + stdout, err
}

func GotoStartScreen(app fyne.App, w fyne.Window) {
	w.Resize(fyne.NewSize(500, 600))

	newBtn := widget.NewButton("New Project", nil)
	openBtn := widget.NewButton("Open Project", nil)
	settingsBtn := widget.NewButton("Settings", nil)
	quitBtn := widget.NewButton("Quit", nil)
	status := widget.NewLabel("")

	top := canvas.NewText("Repomancer", color.White)
	top.Alignment = fyne.TextAlignCenter
	top.TextSize = 16

	newBtn.Disable()
	openBtn.Disable()
	settingsBtn.Disable()

	newBtn.OnTapped = func() {
		d, focus := NewAddProjectDialog(w, func(name, location string) {
			project, err := internal.CreateProject(name, "", location)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			GotoProjectScreen(w, project)
		})
		d.Show()
		w.Canvas().Focus(focus)
	}
	openBtn.OnTapped = func() {
		d := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
			}
			if reader == nil {
				// Nothing was chosen
				return
			}
			project, err := internal.OpenProject(reader.Path())
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				GotoProjectScreen(w, project)
			}
		}, w)

		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		existing := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
		testData := storage.NewFileURI(existing)
		var dir fyne.ListableURI
		dir, err = storage.ListerForURI(testData)
		if err != nil {
			log.Printf("Failed to open directory: %s", err)
			testData = storage.NewFileURI(dirname)
			dir, err = storage.ListerForURI(testData)
			if err != nil {
				dialog.ShowError(err, w)
			}
		}
		d.SetLocation(dir)
		d.SetFileName("config.json")
		d.Show()
	}
	settingsBtn.OnTapped = func() {
		settingsWindow := NewSettingsWindow(app)
		settingsWindow.Show()
	}
	quitBtn.OnTapped = func() {
		w.Close()
	}

	content := container.New(layout.NewVBoxLayout(), newBtn, openBtn, settingsBtn, quitBtn, layout.NewSpacer(), status)

	go func() {
		msg, err := checkRequirements()
		if err != nil {
			status.SetText(fmt.Sprintf("%s\n%s", msg, "gh must be available on the path. Install it from https://cli.github.com/ and restart Repomancer"))
			newBtn.Disable()
			openBtn.Disable()
			settingsBtn.Disable()
		} else {
			status.SetText(msg)
			newBtn.Enable()
			openBtn.Enable()
			settingsBtn.Enable()
		}
	}()

	screen := container.NewBorder(top, nil, nil, nil, content)
	w.SetContent(screen)
}
