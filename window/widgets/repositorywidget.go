package widgets

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"os/exec"
	"repomancer/internal"
)

type RepositoryWidget struct {
	widget.BaseWidget
	Name              *widget.Hyperlink
	Status            *widget.Label
	LastCommandResult *widget.Label
	Selected          *ToggleIconWidget
	OpenBtn           *widget.Button
	LogsBtn           *widget.Button
	CommandsCount     *widget.Label
	PullRequestUrl    *widget.Hyperlink
	PullRequestInfo   *widget.Label
}

func NewRepositoryWidget(title, comment string) *RepositoryWidget {
	item := &RepositoryWidget{
		Name:              widget.NewHyperlink(title, nil),
		Status:            widget.NewLabel(comment),
		LastCommandResult: widget.NewLabel(""),
		CommandsCount:     widget.NewLabel("0"),
		Selected:          NewToggleWidget(nil),
		OpenBtn:           widget.NewButton("Open", nil),
		LogsBtn:           widget.NewButton("Logs", nil),
		PullRequestUrl:    widget.NewHyperlink("https://github.com/organization/repository", nil),
		PullRequestInfo:   widget.NewLabel(""),
	}

	item.Status.Truncation = fyne.TextTruncateEllipsis
	item.PullRequestUrl.Truncation = fyne.TextTruncateOff

	item.ExtendBaseWidget(item)

	return item
}

func (rw *RepositoryWidget) CreateRenderer() fyne.WidgetRenderer {
	statusLine := container.NewBorder(nil,
		nil,
		container.NewHBox(rw.CommandsCount, rw.LastCommandResult),
		container.NewHBox(rw.PullRequestUrl, rw.PullRequestInfo), rw.Status)
	c := container.NewBorder(nil, statusLine, rw.Selected, container.NewHBox(rw.LogsBtn, rw.OpenBtn), rw.Name)
	return widget.NewSimpleRenderer(c)
}

func (rw *RepositoryWidget) Update(repo *internal.Repository) {
	rw.Name.SetText(repo.Title())
	rw.Name.URL = repo.GetUrl()
	rw.Status.Bind(binding.BindString(&repo.Status))
	if repo.LastCommandResult != nil {
		rw.LastCommandResult.SetText(fmt.Sprintf("%s", repo.LastCommandResult))
	} else {
		rw.LastCommandResult.SetText("")
	}
	if repo.PullRequest != nil && repo.PullRequest.Url != "" {
		rw.PullRequestUrl.SetText(repo.PullRequest.Url)
		_ = rw.PullRequestUrl.SetURLFromString(repo.PullRequest.Url)
		rw.PullRequestUrl.Show()
		rw.PullRequestInfo.SetText(fmt.Sprintf("(%s) %s", repo.PullRequest.Status, repo.PullRequest.StatusCheckRollupState))
	} else {
		rw.PullRequestInfo.SetText("")
		rw.PullRequestUrl.SetText("")
		rw.PullRequestUrl.Refresh()
		rw.PullRequestUrl.Hidden = true
		rw.PullRequestInfo.SetText("No Pull Request found")
	}

	if repo.Selected {
		rw.Selected.SetIcon(theme.CheckButtonCheckedIcon())
	} else {
		rw.Selected.SetIcon(theme.CheckButtonIcon())
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
		err := ShowLogWindow(repo)
		if err != nil {
			log.Fatal(err)
		}
	}

	// If there's a log file for this repository, enable the View Logs button
	if _, err := os.Stat(repo.LogFile); err == nil {
		rw.LogsBtn.Enable()
	} else if errors.Is(err, os.ErrNotExist) {
		rw.LogsBtn.Disable()
	} else {
		// Schrödinger: file may or may not exist.
		rw.LogsBtn.Disable()
	}

	rw.OpenBtn.OnTapped = func() {
		log.Printf("Opening %s", repo.Name)
		cmd := exec.Command("open", repo.BaseDir)
		err := cmd.Run()
		if err != nil {
			log.Printf("Error opening %s: %s", repo.BaseDir, err)
		}
	}
}
