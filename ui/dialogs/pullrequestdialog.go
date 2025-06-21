package dialogs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/repomancer/repomancer/internal"
	"github.com/repomancer/repomancer/ui/widgets"
	"strings"
)

func PullRequestDialog(window fyne.Window, project *internal.Project, onAdded func(title, description string)) (*dialog.FormDialog, *widgets.EscapeEntry) {
	title := widgets.NewEscapeEntry()
	title.SetPlaceHolder("Title")
	if strings.TrimSpace(project.PullRequestTitle) != "" {
		title.SetText(project.PullRequestTitle)
	}
	title.MaxLength = 60
	titleItem := widget.NewFormItem("Title", title)
	title.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return fmt.Errorf("title is required")
		}
		return nil
	}

	description := widgets.NewEscapeEntry()
	description.SetPlaceHolder("Describe your changes and reasoning in detail, and match the commits with the proposed change explanation...")
	description.MultiLine = true
	description.Wrapping = fyne.TextWrapWord
	description.SetMinRowsVisible(15)
	if strings.TrimSpace(project.PullRequestDescription) != "" {
		description.SetText(project.PullRequestDescription)
	}
	description.Validator = func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return fmt.Errorf("description is required")
		}
		return nil
	}
	descriptionItem := widget.NewFormItem("Description", description)

	d := dialog.NewForm("Pull Request",
		"Create Pull Request",
		"Cancel",
		[]*widget.FormItem{titleItem, descriptionItem},
		func(b bool) {
			if b {
				onAdded(strings.TrimSpace(title.Text), strings.TrimSpace(description.Text))
			}
		},
		window)

	title.OnEscape = func() {
		d.Hide()
	}
	description.OnEscape = func() {
		d.Hide()
	}

	d.Resize(fyne.NewSize(900, 600))
	return d, title
}
