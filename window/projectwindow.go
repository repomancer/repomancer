package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"path"
	"repomancer/internal"
	"repomancer/window/widgets"
	"strings"
)

func NewProjectWindow(state *internal.State, project *internal.Project) fyne.Window {
	w := state.NewQuitWindow(project.Name)

	pw := widgets.NewProjectWidget()
	pw.LoadProject(project)

	w.Resize(fyne.NewSize(1000, 800))
	w.SetMaster()
	w.SetContent(pw)

	pw.Toolbar.AddRepository.Action = func() {
		d, entry := AddRepositoryDialog(w, project, func() { pw.Refresh() })
		d.Show()
		w.Canvas().Focus(entry)
	}
	pw.Toolbar.AddMultipleRepositories.Action = func() {
		d, entry := AddMultipleRepositoryDialog(w, project, func() {
			pw.Refresh()
		})

		d.Show()
		w.Canvas().Focus(entry)
	}
	pw.Toolbar.DeleteRepository.Action = func() {
		count := project.SelectedRepositoryCount()
		msg := ""
		if count == 0 {
			dialog.ShowInformation("Delete Repositories", "No repositories selected", w)
			return
		} else if count == 1 {
			msg = "Delete 1 repository?"
		} else {
			msg = fmt.Sprintf("Delete %d repositories?", count)
		}

		c := dialog.NewConfirm("Delete Repositories",
			fmt.Sprintf("%s\nThis will also delete local files for selected\nrepositories but not the remote branch,\nif pushed", msg), func(confirm bool) {
				if confirm {
					project.DeleteSelectedRepositories()
					pw.Refresh()
					err := project.SaveProject()
					if err != nil {
						dialog.ShowError(err, w)
					}
				}
			}, w)
		c.SetConfirmText("Delete")
		c.SetConfirmImportance(widget.DangerImportance)
		c.Show()
	}

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
	pw.Toolbar.SelectWithPullRequest.Action = func() {
		project.Select(internal.ReposWithPullRequest)
		pw.Refresh()
	}
	pw.Toolbar.SelectWithoutPullRequest.Action = func() {
		project.Select(internal.ReposWithoutPullRequest)
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
				// TODO: Create something similar to NewPushJob
				project.AddInternalJobToRepositories(cmd, nil)
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
		pw.ExecuteJobQueue()
	}
	pw.Toolbar.GitOpenPullRequest.Action = func() {
		d, entry := PullRequestDialog(w, project, func(title, description string) {
			project.PullRequestTitle = title
			project.PullRequestDescription = description
			err := project.SaveProject()
			if err != nil {
				log.Println(err)
			}
			err = os.WriteFile(path.Join(project.ProjectDir, "PullRequest.md"), []byte(project.PullRequestDescription), 0644)
			if err != nil {
				log.Println(err)
				dialog.NewError(err, w).Show()
				return
			}
			selected := project.SelectedRepositories()

			for _, repo := range selected {
				job := internal.NewPullRequestJob(repo, project)
				job.OnComplete = func(job *internal.Job) {
					// TODO: there's probably an edge case where something doesn't run here
					// Need to switch to a single job queue/runner
					refreshJob := internal.NewPRStatusJob(repo)
					repo.AddJob(refreshJob)
				}
				repo.AddJob(job)
			}
			pw.ExecuteJobQueue()
		})
		d.Show()
		w.Canvas().Focus(entry)
	}
	pw.Toolbar.GitRefreshStatus.Action = func() {
		selected := project.SelectedRepositories()

		for _, repo := range selected {
			job := internal.NewPRStatusJob(repo)
			repo.AddJob(job)
		}
		pw.ExecuteJobQueue()
	}
	pw.Toolbar.CopyRepositoryList.Action = func() {
		var repos []string
		for _, repo := range project.SelectedRepositories() {
			repos = append(repos, repo.GetUrl().String())
		}
		w.Clipboard().SetContent(strings.Join(repos, "\n"))
	}

	pw.Toolbar.CopyRepositoryStatus.Action = func() {
		var repos []string
		for _, repo := range project.SelectedRepositories() {
			status := ""
			if repo.PullRequest != nil {
				status = repo.PullRequest.Status
			}
			line := strings.Join([]string{repo.GetUrl().String(), status}, ",")
			repos = append(repos, line)
		}
		w.Clipboard().SetContent(strings.Join(repos, "\n"))
	}

	pw.CommandInput.OnSubmitted = func(s string) {
		cmd := strings.TrimSpace(pw.CommandInput.Text)
		log.Println("cmd:", cmd)
		pw.CommandInput.SetText("")
		pw.CommandInput.Refresh()
		project.AddJobToRepositories(cmd)
		pw.ExecuteJobQueue()
	}
	pw.RunBtn.OnTapped = func() {
		cmd := strings.TrimSpace(pw.CommandInput.Text)
		log.Println("cmd:", cmd)
		pw.CommandInput.SetText("")
		pw.CommandInput.Refresh()
		project.AddJobToRepositories(cmd)
		pw.ExecuteJobQueue()
	}

	w.SetOnClosed(func() {
		err := project.SaveProject()
		if err != nil {
			log.Printf("Error saving project: %s", err)
		}
	})

	return w
}
