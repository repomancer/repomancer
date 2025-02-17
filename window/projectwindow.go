package screens

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"repomancer/internal"
	"repomancer/window/widgets"
	"strings"
	"time"
)

func NewProjectWindow(state *internal.State, project *internal.Project) fyne.Window {
	w := state.NewQuitWindow(project.Name)

	pw := widgets.NewProjectWidget()
	pw.LoadProject(project)

	w.Resize(fyne.NewSize(1000, 800))
	w.SetMaster()
	w.SetContent(pw)

	pw.Toolbar.AddRepository.OnTapped = func() {
		d, entry := AddRepositoryDialog(w, project, func() { pw.Refresh() })
		d.Show()
		w.Canvas().Focus(entry)
	}
	//pw.Toolbar.AddMultipleRepositories.OnTapped = func() {
	//	dialog.ShowInformation("Add Multiple", "Multiple", w)
	//}
	pw.Toolbar.SelectAll.Action = func() {
		project.Select(internal.All)
		pw.Refresh()
	}
	pw.Toolbar.SelectNone.Action = func() {
		project.Select(internal.None)
		pw.Refresh()
	}
	pw.Toolbar.SelectErrors.Action = func() {
		project.Select(internal.Errors)
		pw.Refresh()
	}
	pw.Toolbar.SelectTenMore.Action = func() {
		project.Select(internal.TenMore)
		pw.Refresh()
	}

	pw.Toolbar.GitCommit.Action = func() {
		message := widget.NewMultiLineEntry()
		message.Wrapping = fyne.TextWrapWord
		message.SetPlaceHolder("Title\n\nThis commit...")

		content := []*widget.FormItem{widget.NewFormItem("Commit Message", message)}
		d := dialog.NewForm("Commit message", "Commit", "Cancel", content, func(b bool) {
			if b {
				// TODO: shell escaping problem. Pipe message in to StdIn or write it to a file?
				cmd := fmt.Sprintf("git add . && git commit -m '%s'", message.Text)
				project.AddInternalJobToRepositories(cmd, nil)
				pw.Refresh()
				pw.ExecuteJobQueue()
			}
		}, w)
		d.Resize(fyne.NewSize(500, 300))
		d.Show()
		w.Canvas().Focus(message)
	}
	pw.Toolbar.GitPush.Action = func() {
		selected := project.SelectedRepositories()
		for _, repo := range selected {
			job := internal.NewPushJob(repo, project)
			repo.AddJob(job)
		}
		pw.Refresh()
		pw.ExecuteJobQueue()
	}
	pw.Toolbar.GitOpenPullRequest.Action = func() {

		// TODO: Check for pullRequestBody.md in project root, prompt/fail if it doesn't exist
		selected := project.SelectedRepositories()

		for _, repo := range selected {
			job := internal.NewPullRequestJob(repo, project)
			repo.AddJob(job)
		}
		pw.Refresh()
		pw.ExecuteJobQueue()
	}
	pw.Toolbar.GitRefreshStatus.Action = func() {
		cmd := "gh pr status --json number,url,state,statusCheckRollup"
		project.AddInternalJobToRepositories(cmd, func(job *internal.Job) {

			var resp internal.GitHubPrResponse
			err := json.Unmarshal([]byte(strings.Join(job.StdOut, "\n")), &resp)
			if err != nil {
				log.Printf("Error unmarshalling GitHub PR response: %s", err)
				return
			}

			if resp.CurrentBranch.Number == 0 { // No PR for the current branch
				job.Repository.PullRequest = nil
			} else {
				prInfo := internal.PullRequest{
					Number:      resp.CurrentBranch.Number,
					Url:         resp.CurrentBranch.URL,
					Status:      resp.CurrentBranch.State,
					LastChecked: time.Now(),
				}

				job.Repository.PullRequest = &prInfo
				job.Repository.RepositoryStatus.PullRequestCreated = true
			}
		})
		pw.Refresh()
		pw.ExecuteJobQueue()
	}

	pw.CommandInput.OnSubmitted = func(s string) {
		cmd := strings.TrimSpace(pw.CommandInput.Text)
		log.Println("cmd:", cmd)
		pw.CommandInput.SetText("")
		pw.CommandInput.Refresh()
		project.AddJobToRepositories(cmd)
		pw.ExecuteJobQueue()
		pw.Refresh()
	}
	pw.RunBtn.OnTapped = func() {
		cmd := strings.TrimSpace(pw.CommandInput.Text)
		log.Println("cmd:", cmd)
		pw.CommandInput.SetText("")
		pw.CommandInput.Refresh()
		project.AddJobToRepositories(cmd)
		pw.ExecuteJobQueue()
		pw.Refresh()
	}

	w.SetOnClosed(func() {
		err := project.SaveProject()
		if err != nil {
			log.Printf("Error saving project: %s", err)
		}
	})

	return w
}
