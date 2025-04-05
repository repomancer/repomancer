package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"repomancer/internal"
	"repomancer/screen/widgets"
	"strings"
)

func AddMultipleRepositoryDialog(window fyne.Window, project *internal.Project, onAdded func()) (*dialog.FormDialog, *widgets.EscapeEntry) {
	entry := widgets.NewEscapeEntry()
	entry.MultiLine = true
	entry.SetMinRowsVisible(15)
	entry.SetPlaceHolder("github.com/org/repository\ngithub.com/org/repository2")
	entry.Validator = func(s string) error {
		var errors []string
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				_, err := internal.NormalizeGitUrl(line)
				if err != nil {
					errors = append(errors, line)
				}
			}
		}
		if len(errors) > 0 {
			return fmt.Errorf("Invalid: %s", strings.Join(errors, ", "))
		}
		return nil
	}
	formItem := widget.NewFormItem("Repository URLs", entry)
	d := dialog.NewForm("Add Repository",
		"Add",
		"Cancel",
		[]*widget.FormItem{formItem},
		func(b bool) {
			if b {
				urls := strings.Split(entry.Text, "\n")
				var errors []error
				for _, url := range urls {
					if strings.TrimSpace(url) != "" {
						err := project.AddRepositoryFromUrl(url)
						if err != nil {
							errors = append(errors, err)
						}
						onAdded()
					}
				}

				if len(errors) > 0 {
					var msg []string
					for _, e := range errors {
						msg = append(msg, e.Error())
					}
					dialog.ShowInformation("Clone Errors", strings.TrimSpace(strings.Join(msg, "\n")), window)
				}

				onAdded()

			}
		},
		window)

	entry.OnEscape = func() {
		d.Hide()
	}
	entry.OnSubmitted = func(s string) {
		entry.SetText(s + "\n")
	}

	d.Resize(fyne.NewSize(500, 400))
	return d, entry
}
