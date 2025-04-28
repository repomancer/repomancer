package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jashort/repomancer/screen/widgets"
	"log"
	"os"
	"strings"
)

func NewAddProjectDialog(window fyne.Window, onAdded func(name string, location string)) (*dialog.FormDialog, *widgets.ShortcutHandlingEntry) {
	name := widgets.NewShortcutHandlingEntry(window)
	name.SetPlaceHolder("my-project")
	name.AllowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRQSTUVWXYZ0123456789_-/"
	name.MaxLength = 50

	projectName := widget.NewFormItem("Name", name)
	projectName.HintText = "Used as branch name for pull requests"

	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	basePath := fyne.CurrentApp().Preferences().StringWithFallback("baseDirectory", dirname)

	location := widgets.NewEscapeEntry()
	location.Text = basePath
	location.Validator = func(s string) error {
		if s == basePath {
			return nil
		}
		_, err := os.Stat(s)
		if !os.IsNotExist(err) {
			return fmt.Errorf("location must not exist")
		}
		return nil
	}

	projectLocation := widget.NewFormItem("Location", location)

	name.OnChanged = func(value string) {
		location.SetText(basePath + strings.Replace(value, "/", "_", -1))
		location.Refresh()
	}

	d := dialog.NewForm("New Project",
		"Create",
		"Cancel",
		[]*widget.FormItem{projectName, projectLocation},
		func(b bool) {
			if b && onAdded != nil {
				onAdded(name.Text, location.Text)
			}
		},
		window)

	name.OnSubmitted = func(s string) {
		d.Submit()
	}

	return d, name
}
