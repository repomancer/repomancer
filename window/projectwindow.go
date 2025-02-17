package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"log"
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

	pw.Toolbar.AddRepository.OnTapped = func() {
		d, entry := AddRepositoryDialog(w, project)
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
		log.Println("Commit not implemented")
	}
	pw.Toolbar.GitPush.Action = func() {
		cmd := fmt.Sprintf("git push origin '%s'", project.Name)

		project.AddInternalJobToRepositories(cmd, func(job *internal.Job) {
			log.Printf("Job complete: %s", job.Duration())
		})
		pw.Refresh()
		pw.ExecuteJobQueue()
	}
	pw.Toolbar.GitOpenPullRequest.Action = func() {
		log.Println("Open Pull Request not implemented")
	}
	pw.Toolbar.GitRefreshStatus.Action = func() {
		log.Println("Refresh Status not implemented")
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

	return w
}
