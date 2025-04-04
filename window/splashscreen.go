package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"repomancer/internal"
)

type StartScreen struct {
	newBtn      *widget.Button
	openBtn     *widget.Button
	settingsBtn *widget.Button
	quitBtn     *widget.Button
	status      *widget.Label
}

func (s *StartScreen) Log(message string) {
	s.status.SetText(message)
	s.status.Refresh()
	log.Println(message)
}

func (s *StartScreen) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	s.Log(msg)
}

func checkRequirements() (string, error) {
	stdout, stderr, err := internal.RunCommand("", 3, "gh", "--version")
	if err != nil {
		return stderr, err
	}
	return "Found gh command\n" + stdout, err
}

func NewStartScreen(app fyne.App, w fyne.Window) *fyne.Container {
	s := &StartScreen{
		newBtn:      widget.NewButton("New Project", nil),
		openBtn:     widget.NewButton("Open Project", nil),
		settingsBtn: widget.NewButton("Settings", nil),
		quitBtn:     widget.NewButton("Quit", nil),
		status:      widget.NewLabel(""),
	}
	top := canvas.NewText("Repomancer", color.White)
	top.Alignment = fyne.TextAlignCenter
	top.TextSize = 16

	s.newBtn.Disable()
	s.openBtn.Disable()
	s.settingsBtn.Disable()

	//s.newBtn.OnTapped = func() {
	//	d, focus := NewAddProjectDialog(w, func(name, location string) {
	//		project, err := internal.CreateProject(name, "", location)
	//		if err != nil {
	//			dialog.ShowError(err, w)
	//			return
	//		}
	//		projectWidget := NewProjectWindow(state, project)
	//		projectWidget.Show()
	//		w.Hide()
	//	})
	//	d.Show()
	//	w.Canvas().Focus(focus)
	//}
	//s.openBtn.OnTapped = func() {
	//	d := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
	//		if err != nil {
	//			dialog.ShowError(err, w)
	//		}
	//		if reader == nil {
	//			// Nothing was chosen
	//			return
	//		}
	//		project, err := internal.OpenProject(reader.Path())
	//		if err != nil {
	//			dialog.ShowError(err, *w)
	//		} else {
	//			//window := NewProjectWindow(state, project)
	//			//window.Show()
	//			//w.Hide()
	//		}
	//	}, w)
	//
	//	dirname, err := os.UserHomeDir()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	existing := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)
	//	testData := storage.NewFileURI(existing)
	//	var dir fyne.ListableURI
	//	dir, err = storage.ListerForURI(testData)
	//	if err != nil {
	//		log.Printf("Failed to open directory: %s", err)
	//		testData = storage.NewFileURI(dirname)
	//		dir, err = storage.ListerForURI(testData)
	//		if err != nil {
	//			dialog.ShowError(err, w)
	//		}
	//	}
	//	d.SetLocation(dir)
	//	d.SetFileName("config.json")
	//	d.Show()
	//}
	s.settingsBtn.OnTapped = func() {
		settingsWindow := NewSettingsWindow(app)
		settingsWindow.Show()
	}
	s.quitBtn.OnTapped = func() {
		w.Close()
	}

	content := container.New(layout.NewVBoxLayout(), s.newBtn, s.openBtn, s.settingsBtn, s.quitBtn, layout.NewSpacer(), s.status)

	go func() {
		msg, err := checkRequirements()
		if err != nil {
			s.Logf("%s\n%s", msg, "gh must be available on the path. Install it from https://cli.github.com/ and restart Repomancer ")
			s.newBtn.Disable()
			s.openBtn.Disable()
			s.settingsBtn.Disable()
		} else {
			s.Logf("%s", msg)
			s.newBtn.Enable()
			s.openBtn.Enable()
			s.settingsBtn.Enable()
		}
	}()

	screen := container.NewBorder(top, nil, nil, nil, content)
	return screen
}
