package screens

import (
	"fyne.io/fyne/v2"
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

func NewProjectWindow(state *internal.State) fyne.Window {
	pw := widgets.NewProjectWidget()
	w := &ProjectWindow{state.NewQuitWindow("test"), pw}

	w.SetMaster()
	w.SetContent(pw)
	return w
}
