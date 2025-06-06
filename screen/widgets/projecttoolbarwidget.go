package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProjectToolbarWidget struct {
	widget.BaseWidget
	RepositoryMenu           *ContextMenuButton
	AddRepository            *fyne.MenuItem
	AddMultipleRepositories  *fyne.MenuItem
	DeleteRepository         *fyne.MenuItem
	DeleteLogs               *fyne.MenuItem
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
	StatisticsMenu           *ContextMenuButton
	ProjectStatistics        *fyne.MenuItem
}

func NewProjectToolbarWidget() *ProjectToolbarWidget {
	item := &ProjectToolbarWidget{
		AddRepository:            fyne.NewMenuItem("Add", nil),
		AddMultipleRepositories:  fyne.NewMenuItem("Add Multiple", nil),
		DeleteRepository:         fyne.NewMenuItem("Delete Selected Repositories", nil),
		DeleteLogs:               fyne.NewMenuItem("Delete Selected Logs", nil),
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
		ProjectStatistics:    fyne.NewMenuItem("Project Statistics", nil),
	}
	item.ExtendBaseWidget(item)
	item.RepositoryMenu = NewContextMenuButton("Repository...", fyne.NewMenu("Repository",
		item.AddRepository, item.AddMultipleRepositories, fyne.NewMenuItemSeparator(), item.DeleteRepository, item.DeleteLogs))

	item.SelectMenu = NewContextMenuButton("Select...",
		fyne.NewMenu("Select",
			item.SelectAll, item.SelectNone, item.SelectErrors, item.SelectTenMore, item.SelectWithPullRequest, item.SelectWithoutPullRequest))

	item.GitMenu = NewContextMenuButton("Git...",
		fyne.NewMenu("Git",
			item.GitCommit, item.GitPush, item.GitOpenPullRequest, item.GitRefreshStatus))

	item.CopyMenu = NewContextMenuButton("Copy...",
		fyne.NewMenu("Copy", item.CopyRepositoryList, item.CopyRepositoryStatus))

	item.StatisticsMenu = NewContextMenuButton("Statistics...",
		fyne.NewMenu("Project Statistics", item.ProjectStatistics))
	return item
}

func (item *ProjectToolbarWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.RepositoryMenu, item.SelectMenu, item.GitMenu, item.CopyMenu, item.StatisticsMenu)
	return widget.NewSimpleRenderer(c)
}
