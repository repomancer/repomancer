package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"repomancer/internal"
	"strings"
)

// ShowLogWindow launch the default system viewer for the repository's .log file
// On macOS, usually the "console" app
func ShowLogWindow(repository *internal.Repository) error {
	cmd := exec.Command("open", repository.LogFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

type ProjectWidget struct {
	widget.BaseWidget
	list         *widget.List
	statusLabel  *widget.Label
	CommandInput *widget.Entry
	RunBtn       *widget.Button
	project      *internal.Project
	Toolbar      *ProjectToolbarWidget
}

func (pw *ProjectWidget) Refresh() {
	pw.list.Refresh()
	selectedCount := pw.project.SelectedRepositoryCount()
	msg := fmt.Sprintf("%d/%d Selected", selectedCount, pw.project.RepositoryCount())
	pw.statusLabel.SetText(msg)
	if selectedCount == 0 {
		pw.Toolbar.DeleteRepository.Label = fmt.Sprintf("Delete All Repositories")
		pw.Toolbar.DeleteLogs.Label = fmt.Sprintf("Clear All Logs")
	} else if selectedCount == 1 {
		pw.Toolbar.DeleteRepository.Label = "Delete 1 Repository"
		pw.Toolbar.DeleteLogs.Label = "Clear 1 Log"
	} else {
		pw.Toolbar.DeleteRepository.Label = fmt.Sprintf("Delete %d Repositories", selectedCount)
		pw.Toolbar.DeleteLogs.Label = fmt.Sprintf("Clear %d Logs", selectedCount)
	}
	if strings.TrimSpace(pw.CommandInput.Text) == "" || pw.project.RepositoryCount() == 0 {
		pw.RunBtn.Disable()
	} else {
		pw.RunBtn.Enable()
	}
}

func (pw *ProjectWidget) ExecuteJobQueue() {
	repositories := pw.project.Repositories

	for _, repo := range repositories {
		pw.project.WorkerChannel <- repo
	}
}

func (pw *ProjectWidget) LoadProject(project *internal.Project) {
	pw.project = project
	pw.Refresh()
}

func (pw *ProjectWidget) CreateRenderer() fyne.WidgetRenderer {

	header := container.NewBorder(pw.Toolbar, nil, nil, pw.RunBtn, pw.CommandInput)
	footer := container.NewGridWithColumns(1, pw.statusLabel)
	c := container.NewBorder(header, footer, nil, nil, pw.list)

	return widget.NewSimpleRenderer(c)
}

func NewProjectWidget() *ProjectWidget {
	pw := &ProjectWidget{
		BaseWidget:   widget.BaseWidget{},
		statusLabel:  widget.NewLabel(""),
		CommandInput: widget.NewEntry(),
		RunBtn:       widget.NewButton("Run", nil),
		Toolbar:      NewProjectToolbarWidget(),
	}
	pw.ExtendBaseWidget(pw)
	pw.list = widget.NewList(
		func() int { return pw.project.RepositoryCount() },
		func() fyne.CanvasObject {
			return NewRepositoryWidget("", "")
		},
		func(listItemId widget.ListItemID, obj fyne.CanvasObject) {
			repo := pw.project.GetRepository(listItemId)
			rw := obj.(*RepositoryWidget)
			rw.Update(repo)
			repo.OnUpdated = func(_ *internal.Repository) {
				fyne.Do(pw.list.Refresh)
			}
			rw.Selected.OnTapped = func() {
				repo.Selected = !repo.Selected
				log.Printf("%s checked: %t", repo.Name, repo.Selected)
				if repo.Selected {
					rw.Selected.SetIcon(theme.CheckButtonCheckedIcon())
				} else {
					rw.Selected.SetIcon(theme.CheckButtonIcon())
				}
				rw.Selected.Refresh()
				pw.Refresh()
			}

		},
	)

	pw.statusLabel.Importance = widget.LowImportance

	return pw
}
