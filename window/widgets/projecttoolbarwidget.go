package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProjectToolbarWidget struct {
	widget.BaseWidget
	AddRepository            *widget.Button
	AddMultipleRepositories  *widget.Button
	SelectMenu               *ContextMenuButton
	SelectAll                *fyne.MenuItem
	SelectNone               *fyne.MenuItem
	SelectErrors             *fyne.MenuItem
	SelectTenMore            *fyne.MenuItem
	SelectWithPullRequest    *fyne.MenuItem
	SelectWithoutPullRequest *fyne.MenuItem
	GitMenu                  *ContextMenuButton
	GitCommit                *fyne.MenuItem
	GitPush                  *fyne.MenuItem
	GitOpenPullRequest       *fyne.MenuItem
	GitRefreshStatus         *fyne.MenuItem
	CopyMenu                 *ContextMenuButton
	CopyRepositoryList       *fyne.MenuItem
	CopyRepositoryStatus     *fyne.MenuItem
}

func NewProjectToolbarWidget() *ProjectToolbarWidget {
	item := &ProjectToolbarWidget{
		AddRepository:            widget.NewButton("Add Repository", nil),
		AddMultipleRepositories:  widget.NewButton("Add Multiple", nil),
		SelectAll:                fyne.NewMenuItem("All", nil),
		SelectNone:               fyne.NewMenuItem("None", nil),
		SelectErrors:             fyne.NewMenuItem("Errors", nil),
		SelectTenMore:            fyne.NewMenuItem("Next 10", nil),
		SelectWithPullRequest:    fyne.NewMenuItem("Repos With PullRequest", nil),
		SelectWithoutPullRequest: fyne.NewMenuItem("Repos Without Pull Request", nil),
		// TODO: Add more variants. Select Everything with unmerged PRs?
		// Select Merged PRs?
		// Etc
		GitCommit:            fyne.NewMenuItem("Commit", nil),
		GitPush:              fyne.NewMenuItem("Push", nil),
		GitOpenPullRequest:   fyne.NewMenuItem("Open Pull Request", nil),
		GitRefreshStatus:     fyne.NewMenuItem("Refresh Status", nil),
		CopyRepositoryList:   fyne.NewMenuItem("Repository List", nil),
		CopyRepositoryStatus: fyne.NewMenuItem("Pull Request Status", nil),
	}
	item.ExtendBaseWidget(item)

	item.SelectMenu = NewContextMenuButton("Select...",
		fyne.NewMenu("Select",
			item.SelectAll, item.SelectNone, item.SelectErrors, item.SelectTenMore, item.SelectWithPullRequest, item.SelectWithoutPullRequest))

	item.GitMenu = NewContextMenuButton("GitHub...",
		fyne.NewMenu("GitHub",
			item.GitCommit, item.GitPush, item.GitOpenPullRequest, item.GitRefreshStatus))

	item.CopyMenu = NewContextMenuButton("Copy...",
		fyne.NewMenu("Copy", item.CopyRepositoryList, item.CopyRepositoryStatus))
	return item
}

func (item *ProjectToolbarWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.AddRepository, item.AddMultipleRepositories, item.SelectMenu, item.GitMenu, item.CopyMenu)
	return widget.NewSimpleRenderer(c)
}
