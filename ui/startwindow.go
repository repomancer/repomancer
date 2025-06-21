package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/repomancer/repomancer/internal"
	"github.com/repomancer/repomancer/ui/dialogs"
	"log"
	"os"
)

type StartWindow struct {
	fyne.Window
	status  *widget.Label
	newBtn  *widget.Button
	openBtn *widget.Button
}

func (s *StartWindow) CheckStatus() {
	go func() {
		stdout, stderr, err := internal.RunCommand("", 3, "gh --version")
		fyne.Do(func() {
			if err != nil {
				msg := fmt.Sprintf(`
Error running gh --version:
%s
%v
- Is gh installed and available on the PATH?
- Check shell and shell arguments in settings`, stderr, err)
				s.status.SetText(msg)
				s.newBtn.Disable()
				s.openBtn.Disable()
			} else {
				s.status.SetText("Found gh command\n" + stdout)
				s.newBtn.Enable()
				s.openBtn.Enable()
			}
			s.status.Refresh()
			s.newBtn.Refresh()
			s.openBtn.Refresh()
		})
	}()
}

func NewStartScreen(b *BaseUI) StartWindow {
	w := StartWindow{
		Window:  b.NewWindow("Repomancer"),
		status:  widget.NewLabel(""),
		newBtn:  widget.NewButton("New Project", nil),
		openBtn: widget.NewButton("Open Project", nil),
	}
	w.newBtn.OnTapped = func() {
		d, focus := dialogs.NewAddProjectDialog(w, func(name, location string) {
			project, err := internal.CreateProject(name, "", location)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			b.ShowProjectWindow(project)
		})
		d.Show()
		w.Window.Canvas().Focus(focus)
	}

	w.openBtn.OnTapped = func() {
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
				b.ShowProjectWindow(project)
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

	w.Resize(fyne.NewSize(500, 600))
	w.status.Wrapping = fyne.TextWrapWord
	settingsBtn := widget.NewButton("Settings", func() {
		b.ShowWindow(Settings)
		b.windows[Settings].SetCloseIntercept(func() {
			w.CheckStatus()
			w.Close()
		})
	})
	quitBtn := widget.NewButton("Quit", func() {
		b.App.Quit()
	})
	w.SetContent(container.NewVBox(w.newBtn, w.openBtn, settingsBtn, quitBtn, w.status))
	w.CheckStatus()
	return w
}
