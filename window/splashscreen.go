package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"os"
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
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	stdout, stderr, err := internal.ShellOut("gh --version", home)
	if err != nil {
		return stderr, err
	}
	return "Found gh command\n" + stdout, err
}

func NewStartScreen(window fyne.Window) fyne.CanvasObject {
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

	s.newBtn.OnTapped = func() {
		w2 := fyne.CurrentApp().NewWindow("New Project")
		NewProjectScreen(w2)
		w2.Show()
	}
	s.openBtn.OnTapped = func() {
		dialog.ShowFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				s.Logf("%v", err)
			}
			if reader == nil {
				// Nothing was chosen
				return
			}
			s.Logf("Open Project: %s", reader.Path())
			// TODO: Open Project
			//project, err := internal.OpenProject(reader.Path())
			if err != nil {
				s.Logf("Failed to open project: %s", err)
			} else {
				//log.Printf("Open Project: %s", project)
				//MainScreen(window, project)
			}
		}, window)
	}
	s.settingsBtn.OnTapped = func() {
		log.Println("Settings")
		NewPreferenceScreen()
	}
	s.quitBtn.OnTapped = func() {
		window.Close()
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
	window.Resize(fyne.NewSize(500, 600))
	return screen
}
