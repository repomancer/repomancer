package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"log"
	"repomancer/internal"
	"repomancer/window/widgets"
)

type ProjectWindow struct {
	fyne.Window
	pw *widgets.ProjectWidget
}

func (w *ProjectWindow) LoadProject(project *internal.Project) {
	w.pw.LoadProject(project)
}

func NewProjectWindow(state *internal.State, project *internal.Project) fyne.Window {
	w := state.NewQuitWindow(project.Name)

	pw := widgets.NewProjectWidget()
	pw.LoadProject(project)

	pw.Toolbar.AddRepository.OnTapped = func() {
		log.Println("New Repository")
		dialog.ShowInformation("New Repository", "Add repository", w)
	}
	w.Resize(fyne.NewSize(1000, 800))
	w.SetMaster()
	w.SetContent(pw)

	pw.Toolbar.AddRepository.OnTapped = func() {
		dialog.ShowInformation("Add Repository", "Add Repository", w)
	}
	pw.Toolbar.AddMultipleRepositories.OnTapped = func() {
		dialog.ShowInformation("Add Multiple", "Multiple", w)
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

	pw.Toolbar.GitCommit.Action = func() {
		log.Println("Commit not implemented")
	}
	pw.Toolbar.GitPush.Action = func() {
		log.Println("Push not implemented")
	}
	pw.Toolbar.GitOpenPullRequest.Action = func() {
		log.Println("Open Pull Request not implemented")
	}
	pw.Toolbar.GitRefreshStatus.Action = func() {
		log.Println("Refresh Status not implemented")
	}

	return w
}
