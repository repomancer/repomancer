package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RepositoryWidget struct {
	widget.BaseWidget
	Name              *widget.Label
	Status            *widget.Label
	LastCommandResult *widget.Label
	Selected          *ToggleIconWidget
	OpenBtn           *widget.Button
	LogsBtn           *widget.Button
	CommandsCount     *widget.Label
	PullRequestInfo   *widget.Label
}

func NewRepositoryWidget(title, comment string) *RepositoryWidget {
	item := &RepositoryWidget{
		Name:              widget.NewLabel(title),
		Status:            widget.NewLabel(comment),
		LastCommandResult: widget.NewLabel(""),
		CommandsCount:     widget.NewLabel("0"),
		Selected:          NewToggleWidget(nil),
		OpenBtn:           widget.NewButton("Open", nil),
		LogsBtn:           widget.NewButton("Logs", nil),
		PullRequestInfo:   widget.NewLabel(""),
	}

	item.Status.Truncation = fyne.TextTruncateEllipsis
	item.ExtendBaseWidget(item)

	return item
}

func (item *RepositoryWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, container.NewGridWithColumns(4, item.CommandsCount, item.Status, item.PullRequestInfo, item.LastCommandResult), item.Selected, container.NewHBox(item.LogsBtn, item.OpenBtn), item.Name)
	return widget.NewSimpleRenderer(c)
}
