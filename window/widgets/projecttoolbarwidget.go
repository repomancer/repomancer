package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProjectToolbarWidget struct {
	widget.BaseWidget
	AddRepository *widget.Button
	//AddMultipleRepositories *widget.Button
	SelectMenu         *ContextMenuButton
	SelectAll          *fyne.MenuItem
	SelectNone         *fyne.MenuItem
	SelectErrors       *fyne.MenuItem
	SelectTenMore      *fyne.MenuItem
	GitMenu            *ContextMenuButton
	GitCommit          *fyne.MenuItem
	GitPush            *fyne.MenuItem
	GitOpenPullRequest *fyne.MenuItem
	GitRefreshStatus   *fyne.MenuItem
}

func NewProjectToolbarWidget() *ProjectToolbarWidget {
	item := &ProjectToolbarWidget{
		AddRepository: widget.NewButton("Add Repository", nil),
		//AddMultipleRepositories: widget.NewButton("Add Multiple", nil),
		SelectAll:          fyne.NewMenuItem("Select All", nil),
		SelectNone:         fyne.NewMenuItem("Select None", nil),
		SelectErrors:       fyne.NewMenuItem("Select Errors", nil),
		SelectTenMore:      fyne.NewMenuItem("Select Next 10", nil),
		GitCommit:          fyne.NewMenuItem("Commit", nil),
		GitPush:            fyne.NewMenuItem("Push", nil),
		GitOpenPullRequest: fyne.NewMenuItem("Open Pull Request", nil),
		GitRefreshStatus:   fyne.NewMenuItem("Refresh Status", nil),
	}
	item.ExtendBaseWidget(item)

	item.SelectMenu = NewContextMenuButton("Select...",
		fyne.NewMenu("Select",
			item.SelectAll, item.SelectNone, item.SelectErrors, item.SelectTenMore))

	item.GitMenu = NewContextMenuButton("GitHub...", fyne.NewMenu("GitHub",
		item.GitCommit, item.GitPush, item.GitOpenPullRequest, item.GitRefreshStatus))

	return item
}

func (item *ProjectToolbarWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.AddRepository, item.SelectMenu, item.GitMenu)
	return widget.NewSimpleRenderer(c)
}
