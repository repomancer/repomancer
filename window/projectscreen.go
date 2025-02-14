package screens

import (
	"fmt"
	"repomancer/internal"
	"repomancer/window/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"strings"
)

type SelectRange int

const (
	All SelectRange = iota
	None
	Errors
	TenMore
)

func PushChanges() {
	selectedCount := project.SelectedRepositoryCount()
	for i := 0; i < project.RepositoryCount(); i++ {
		r := project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			result, err := internal.PushChanges(r, project)
			if err != nil {
				log.Printf("Error pushing changes on %s:  %v", r.Name, err)
				return
			}
			log.Printf("Pushed changes on %s: %s", r.Name, result)
		}
	}
}

func CreatePR() {
	selectedCount := project.SelectedRepositoryCount()
	for i := 0; i < project.RepositoryCount(); i++ {
		r := project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			result, err := internal.CreatePullRequest(r, project)
			if err != nil {
				log.Printf("Error creating pull request on %s:  %v", r.Name, err)
				return
			}
			log.Printf("Created pull request on %s: %s", r.Name, result)
		}
	}
}

func CheckPRStatus() {
	selectedCount := project.SelectedRepositoryCount()
	for i := 0; i < project.RepositoryCount(); i++ {
		r := project.GetRepository(i)
		if r.Selected || selectedCount == 0 {
			err := internal.UpdatePullRequestInfo(r)
			if err != nil {
				r.LastCommandResult = err
			}
			list.Refresh()
		}
	}
}

func Select(selectRange SelectRange) {
	cnt := 0
	if selectRange == All {
		for i := 0; i < project.RepositoryCount(); i++ {
			project.GetRepository(i).Selected = true
			cnt++
		}
	} else if selectRange == None {
		for i := 0; i < project.RepositoryCount(); i++ {
			project.GetRepository(i).Selected = false
		}
	} else if selectRange == Errors {
		for i := 0; i < project.RepositoryCount(); i++ {
			if project.GetRepository(i).LastCommandResult != nil {
				project.GetRepository(i).Selected = true
				cnt++
			} else {
				project.GetRepository(i).Selected = false
			}
		}
	} else if selectRange == TenMore {
		added := 0
		for i := 0; i < project.RepositoryCount(); i++ {
			if !project.GetRepository(i).Selected {
				if added < 10 {
					project.GetRepository(i).Selected = true
					added++
					cnt++
				}
			}
		}
	}
	list.Refresh()
	var msg string
	if cnt > 0 {
		msg = fmt.Sprintf("%d/%d Selected (%d added)", project.SelectedRepositoryCount(), project.RepositoryCount(), cnt)
	} else {
		msg = fmt.Sprintf("%d/%d Selected", project.SelectedRepositoryCount(), project.RepositoryCount())
	}
	statusLabel.SetText(msg)
}

var project *internal.Project
var list *widget.List
var statusLabel *widget.Label
var progressBar *widget.ProgressBar
var CommandInput *widget.Entry
var runBtn *widget.Button

func RefreshRepoList() {
	list.Refresh()
}

func NewProjectScreen(state *internal.State) fyne.Window {
	window := state.App.NewWindow("")
	window.Resize(fyne.NewSize(800, 600))
	//project = p
	CommandInput = widget.NewEntry()
	CommandInput.OnSubmitted = func(s string) {
		cmd := strings.TrimSpace(CommandInput.Text)
		log.Println("cmd:", cmd)
		CommandInput.SetText("")
		CommandInput.Refresh()

		AddJobToRepositories(cmd)
		ExecuteJobQueue()
	}

	runBtn = widget.NewButton("Run", func() {
		cmd := strings.TrimSpace(CommandInput.Text)
		CommandInput.SetText("")
		CommandInput.Refresh()
		AddJobToRepositories(cmd)
		ExecuteJobQueue()
	})

	statusLabel = widget.NewLabel("")
	statusLabel.Importance = widget.LowImportance

	progressBar = widget.NewProgressBar()
	progressBar.SetValue(0.0)
	progressBar.Hide()

	addBtn := widget.NewButton("Add Repository", func() {
		//ShowAddRepositoryWindow(project)
	})
	addMultipleBtn := widget.NewButton("Add Multiple", func() {
		//ShowAddMultipleRepositoryWindow(project)
	})
	addBtn.Importance = widget.HighImportance

	toolbar := container.NewHBox(
		addBtn,
		addMultipleBtn,
		widgets.NewContextMenuButton("Select...",
			fyne.NewMenu("",
				fyne.NewMenuItem("All", func() { Select(All) }),
				fyne.NewMenuItem("None", func() { Select(None) }),
				fyne.NewMenuItem("Errors", func() { Select(Errors) }),
				fyne.NewMenuItem("10 More", func() { Select(TenMore) }),
			)),

		// Todo: Not implemented
		//widgets.NewContextMenuButton("Sort...",
		//	fyne.NewMenu("",
		//		fyne.NewMenuItem("Checked", func() { Select(All) }),
		//	)),

		widgets.NewContextMenuButton("GitHub...",
			fyne.NewMenu("",
				fyne.NewMenuItem("Commit", func() {
					//ShowCommitWindow(project)
				}),
				fyne.NewMenuItem("Push", func() { PushChanges() }),
				fyne.NewMenuItem("Create Pull Request", func() { CreatePR() }),
				fyne.NewMenuItem("Refresh PR status", func() { CheckPRStatus() }),
			),
		),
	)

	list = widget.NewList(
		func() int { return project.RepositoryCount() },
		func() fyne.CanvasObject {
			return widgets.NewRepositoryWidget("", "")
		},
		func(listItemId widget.ListItemID, obj fyne.CanvasObject) {
			repo := project.GetRepository(listItemId)
			rw := obj.(*widgets.RepositoryWidget)
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
				statusLabel.SetText(fmt.Sprintf("%d/%d Selected", project.SelectedRepositoryCount(), project.RepositoryCount()))
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
				//ShowLogWindow(repo)
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

	header := container.NewBorder(toolbar, nil, nil, runBtn, CommandInput)
	footer := container.NewGridWithColumns(2, statusLabel, progressBar)

	l := container.NewBorder(header, footer, nil, nil, list)

	//Select(None)

	//menu := fyne.NewMainMenu(FileMenu(window), EditMenu(), ProjectMenu(window, project), ViewMenu())
	//window.SetMainMenu(menu)

	window.SetContent(l)
	return window
}

func AddJobToRepositories(cmd string) {
	if project.SelectedRepositoryCount() == 0 {
		// Nothing selected, run everywhere
		for i := 0; i < project.RepositoryCount(); i++ {
			j := internal.NewJob(project.GetRepository(i), cmd)
			project.GetRepository(i).AddJob(j)
		}
	} else {
		// Only run on selected repos
		for i := 0; i < project.RepositoryCount(); i++ {
			if project.GetRepository(i).Selected {
				j := internal.NewJob(project.GetRepository(i), cmd)
				project.GetRepository(i).AddJob(j)
			}
		}
	}
	list.Refresh()
}

func ExecuteJobQueue() {
	var jobsToRun []*internal.Job

	for i := 0; i < project.RepositoryCount(); i++ {
		for j := 0; j < project.GetRepository(i).JobCount(); j++ {
			if !project.GetRepository(i).GetJob(j).Finished {
				jobsToRun = append(jobsToRun, project.GetRepository(i).GetJob(j))
			}
		}
	}
	log.Printf("Found %d jobs to run", len(jobsToRun))

	progressBar.SetValue(0.0)
	progressBar.Min = 0
	progressBar.Max = float64(len(jobsToRun))
	progressBar.Show()

	go func() {
		for i := 0; i < len(jobsToRun); i++ {
			statusLabel.SetText(fmt.Sprintf("Running %d/%d", i+1, len(jobsToRun)))
			jobsToRun[i].Run()
			progressBar.SetValue(float64(i + 1))
			progressBar.Refresh()
			list.Refresh()
		}
	}()

	progressBar.Hide()
	statusLabel.SetText(fmt.Sprintf("%d jobs finished", len(jobsToRun)))
}
