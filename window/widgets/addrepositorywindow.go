package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"repomancer/internal"
)

func NewAddRepositoryDialog(project *internal.Project, window fyne.Window) {

}

// repositoryLooksValid tries to check if a string looks like a valid GitHub remote.
// Right now it only supports a limited number of formats; github.com/org/repo
// Adding a match for other formats would be good. Note that we only want remote repositories
// because we're expecting to make pull requests, and we'll need to convert ssh/etc URLs in
// to HTTPS (or normalize to no protocol)
// List of remote repo types, minus local repositories:
// https://stackoverflow.com/a/2514986
func repositoryLooksValid(str string) bool {
	_, err := internal.NormalizeGitUrl(str)
	if err != nil {
		return false
	}
	return true
}

func ShowAddRepositoryWindow(project *internal.Project, onClose func()) {
	w2 := fyne.CurrentApp().NewWindow("Add Repository")

	label1 := widget.NewLabel("Name")
	value1 := NewShortcutHandlingEntry(w2, false)
	value1.SetPlaceHolder("github.com/organization/repository")

	output := widget.NewLabel("")
	output.Wrapping = fyne.TextWrapWord

	// TODO: Extract a common method from value1.OnEnter and AddBtn
	value1.OnSubmitted = func(str string) {
		output.SetText(fmt.Sprintf("Cloning %s and switching to branch %s...", value1.Text, project.Name))
		info, err := internal.GetRepositoryInfo(value1.Text)
		if err != nil {
			output.SetText("Error: " + err.Error())
			return
		}
		err = project.AddRepositoryFromUrl(info.URL)
		if err != nil {
			output.SetText("Error: " + err.Error())
		} else {
			onClose()
			w2.Close()
		}
	}

	addBtn := widget.NewButton("Add", func() {
		info, err := internal.GetRepositoryInfo(value1.Text)
		if err != nil {
			output.SetText("Error: " + err.Error())
			return
		}
		err = project.AddRepositoryFromUrl(info.URL)
		if err != nil {
			output.SetText("Error: " + err.Error())
		} else {
			onClose()
			w2.Close()
		}
	})
	addBtn.Importance = widget.HighImportance
	addBtn.Disable()
	cancelBtn := widget.NewButton("Cancel", func() { w2.Close() })
	testBtn := widget.NewButton("Test", func() {
		info, err := internal.GetRepositoryInfo(value1.Text)
		if err != nil {
			output.SetText("Error: " + err.Error())
			addBtn.Disable()
		} else {
			output.SetText(fmt.Sprintf("Found %s\nLast Commit: %s\nURL: %s", info.Name, info.PushedAt, info.URL))
			addBtn.Enable()
		}
	})
	testBtn.Disable()

	value1.OnChanged = func(s string) {
		if repositoryLooksValid(s) {
			addBtn.Enable()
			testBtn.Enable()
		} else {
			addBtn.Disable()
			testBtn.Disable()
		}
	}

	grid := container.New(layout.NewFormLayout(), label1, value1, widget.NewLabel(""), testBtn)
	c := container.NewVBox(grid, output, layout.NewSpacer(), container.NewGridWithColumns(2, addBtn, cancelBtn))

	w2.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w2.Close()
		}
	})

	escapeShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyEscape}
	w2.Canvas().AddShortcut(escapeShortcut, func(shortcut fyne.Shortcut) {
		log.Println("We tapped Escape")
	})

	otherShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyF1}
	w2.Canvas().AddShortcut(otherShortcut, func(shortcut fyne.Shortcut) {
		log.Println("We tapped Escape")
	})

	w2.SetContent(c)

	w2.Resize(fyne.NewSize(500, 400))
	w2.Canvas().Focus(value1)
	w2.Show()
}
