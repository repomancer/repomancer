package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/repomancer/repomancer/internal"
	"github.com/repomancer/repomancer/ui/widgets"
	"log"
)

func AddRepositoryDialog(window fyne.Window, project *internal.Project, onAdded func()) (*dialog.FormDialog, *widgets.EscapeEntry) {
	entry := widgets.NewEscapeEntry()
	entry.SetPlaceHolder("github.com/org/repository")
	entry.Validator = func(s string) error {
		_, err := internal.NormalizeGitUrl(s)
		return err
	}
	formItem := widget.NewFormItem("Repository URL", entry)
	d := dialog.NewForm("Add Repository",
		"Add",
		"Cancel",
		[]*widget.FormItem{formItem},
		func(b bool) {
			log.Printf("Add Repository Dialog running: %t", b)
			if b {
				info, err := internal.GetRepositoryInfo(entry.Text)
				if err != nil {
					dialog.NewError(err, window).Show()
					return
				}
				err = project.AddRepositoryFromUrl(info.URL)
				if err != nil {
					dialog.NewError(err, window).Show()
				}
				onAdded()
			}
		},
		window)

	entry.OnEscape = func() {
		d.Hide()
	}
	entry.OnSubmitted = func(s string) {
		d.Submit()
	}

	d.Resize(fyne.NewSize(500, 400))
	return d, entry
}
