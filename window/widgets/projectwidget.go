package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"repomancer/internal"
	"strings"
)

func ShowLogWindow(repository *internal.Repository) {
	w2 := fyne.CurrentApp().NewWindow(repository.Title())
	cmdW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierSuper}

	w2.Canvas().AddShortcut(cmdW, func(shortcut fyne.Shortcut) {
		w2.Hide()
	})

	logText := widget.NewListWithData(repository.GetLogBinding(),
		func() fyne.CanvasObject {
			w := widget.NewLabel("")
			return w
		}, func(item binding.DataItem, object fyne.CanvasObject) {
			w := object.(*widget.Label)
			i := item.(binding.String)
			w.Bind(i)

			val, _ := i.Get()
			if strings.Contains(val, "Command:") {
				w.Importance = widget.HighImportance
			} else {
				w.Importance = widget.MediumImportance
			}
		})
	logText.HideSeparators = true

	w2.Resize(fyne.NewSize(800, 800))
	w2.SetContent(logText)
	w2.Show()
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

func (pw *ProjectWidget) PushChanges() {
	selectedCount := pw.project.SelectedRepositoryCount()
	for i := 0; i < pw.project.RepositoryCount(); i++ {
		r := pw.project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			result, err := internal.PushChanges(r, pw.project)
			if err != nil {
				log.Printf("Error pushing changes on %s:  %v", r.Name, err)
				return
			}
			log.Printf("Pushed changes on %s: %s", r.Name, result)
		}
	}
}

func (pw *ProjectWidget) CreatePR() {
	selectedCount := pw.project.SelectedRepositoryCount()
	for i := 0; i < pw.project.RepositoryCount(); i++ {
		r := pw.project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			result, err := internal.CreatePullRequest(r, pw.project)
			if err != nil {
				log.Printf("Error creating pull request on %s:  %v", r.Name, err)
				return
			}
			log.Printf("Created pull request on %s: %s", r.Name, result)
		}
	}
}

func (pw *ProjectWidget) CheckPRStatus() {
	selectedCount := pw.project.SelectedRepositoryCount()
	for i := 0; i < pw.project.RepositoryCount(); i++ {
		r := pw.project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			err := internal.UpdatePullRequestInfo(r)
			if err != nil {
				r.LastCommandResult = err
			}
			pw.Refresh()
		}
	}
}

func (pw *ProjectWidget) Refresh() {
	pw.list.Refresh()
	msg := fmt.Sprintf("%d/%d Selected", pw.project.SelectedRepositoryCount(), pw.project.RepositoryCount())
	pw.statusLabel.SetText(msg)
}

func (pw *ProjectWidget) ExecuteJobQueue() {
	var jobsToRun []*internal.Job

	for i := 0; i < pw.project.RepositoryCount(); i++ {
		for j := 0; j < pw.project.GetRepository(i).JobCount(); j++ {
			if !pw.project.GetRepository(i).GetJob(j).Finished {
				jobsToRun = append(jobsToRun, pw.project.GetRepository(i).GetJob(j))
			}
		}
	}
	log.Printf("Found %d jobs to run", len(jobsToRun))

	go func() {
		for i := 0; i < len(jobsToRun); i++ {
			pw.statusLabel.SetText(fmt.Sprintf("Running %d/%d", i+1, len(jobsToRun)))
			jobsToRun[i].Run()
			pw.Refresh()
		}
	}()

	pw.statusLabel.SetText(fmt.Sprintf("%d jobs finished", len(jobsToRun)))
}

func (pw *ProjectWidget) LoadProject(project *internal.Project) {
	pw.project = project
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

	pw.list = widget.NewList(
		func() int { return pw.project.RepositoryCount() },
		func() fyne.CanvasObject {
			return NewRepositoryWidget("", "")
		},
		func(listItemId widget.ListItemID, obj fyne.CanvasObject) {
			repo := pw.project.GetRepository(listItemId)
			rw := obj.(*RepositoryWidget)
			rw.Name.SetText(repo.Title())
			rw.Status.Bind(binding.BindString(&repo.Status))
			if repo.LastCommandResult != nil {
				rw.LastCommandResult.SetText(fmt.Sprintf("%s", repo.LastCommandResult))
			} else {
				rw.LastCommandResult.SetText("")
			}
			if repo.PullRequest != nil {
				rw.PullRequestInfo.SetText(fmt.Sprintf("%s (%s) %s", repo.PullRequest.Url, repo.PullRequest.Status, repo.PullRequest.StatusCheckRollupState))
			} else {
				rw.PullRequestInfo.SetText("")
			}

			if repo.Selected {
				rw.Selected.SetIcon(theme.CheckButtonCheckedIcon())
			} else {
				rw.Selected.SetIcon(theme.CheckButtonIcon())
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
				pw.statusLabel.SetText(fmt.Sprintf("%d/%d Selected", pw.project.SelectedRepositoryCount(), pw.project.RepositoryCount()))
			}
			queued := repo.QueuedJobs()
			if queued > 1 {
				rw.CommandsCount.SetText(fmt.Sprintf("%d jobs pending", queued))
			} else if queued == 1 {
				rw.CommandsCount.SetText(fmt.Sprintf("%d job pending", queued))
			} else {
				rw.CommandsCount.SetText("")
			}
			rw.LogsBtn.OnTapped = func() {
				log.Printf("Viewing logs for %s", repo.Name)
				ShowLogWindow(repo)
			}
			if repo.JobCount() == 0 {
				rw.LogsBtn.Disable()
			} else {
				rw.LogsBtn.Enable()
			}
			rw.OpenBtn.OnTapped = func() {
				log.Printf("Opening %s", repo.Name)
				cmd := exec.Command("open", repo.BaseDir)
				err := cmd.Run()
				if err != nil {
					log.Printf("Error opening %s: %s", repo.BaseDir, err)
				}
			}
		},
	)

	pw.statusLabel.Importance = widget.LowImportance

	return pw
}
