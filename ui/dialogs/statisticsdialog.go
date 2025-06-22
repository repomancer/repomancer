package dialogs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/repomancer/repomancer/internal"
	"sort"
	"strings"
)

func NewStatisticsDialog(w fyne.Window, project *internal.Project) *dialog.CustomDialog {
	prStatusMap := make(map[string]int)

	for i := 0; i < project.RepositoryCount(); i++ {
		r := project.GetRepository(i)
		if r.PullRequest == nil {
			continue
		}
		prStatusMap[r.PullRequest.Status]++
	}

	var msg []string
	msg = append(msg, fmt.Sprintf("Repositories: %d\n", project.RepositoryCount()))
	msg = append(msg, "Pull Requests:")

	var keys []string
	for k := range prStatusMap {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		msg = append(msg, "None")
	} else {
		sort.Strings(keys)
		for _, k := range keys {
			msg = append(msg, fmt.Sprintf("    %s\t %3d", k, prStatusMap[k]))
		}
	}
	lbl := widget.NewLabel(strings.Join(msg, "\n"))
	lbl.TextStyle.Monospace = true
	d := dialog.NewCustom("Project Statistics", "Ok", lbl, w)
	return d
}
