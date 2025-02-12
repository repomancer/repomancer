package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"repomancer/internal"
)

type LogWidget struct {
	widget.BaseWidget
	Timestamp *widget.TextGrid
	Command   *widget.TextGrid
	StdOut    *widget.TextGrid
	StdErr    *widget.TextGrid
	ErrorText *widget.Label
}

func (w *LogWidget) Set(job *internal.Job) {
	w.Timestamp.SetText(fmt.Sprintf("%s (Duration: %s)", job.StartTime.Format("2006-01-02 15:04:05"), job.Duration()))

	w.Command.SetText(fmt.Sprintf("Command: %s\n", job.Command))

	if len(job.StdOut) > 0 {
		w.StdOut.SetText(job.StdOut)
		w.StdOut.Show()
	} else {
		w.StdOut.SetText("")
		w.StdOut.Hide()
	}
	if len(job.StdErr) > 0 {
		w.StdErr.SetText(fmt.Sprintf("StdErr:\n%s", job.StdErr))
		w.StdErr.Show()
	} else {
		w.StdErr.SetText("")
		w.StdErr.Hide()
	}
	if job.Error != nil {
		w.ErrorText.Text = fmt.Sprintf("Error: %s", job.Error)
		w.ErrorText.Show()
	} else {
		w.ErrorText.Text = ""
		w.ErrorText.Hide()
	}
	w.Refresh()
}

func NewLogWidget() *LogWidget {
	item := &LogWidget{
		Timestamp: widget.NewTextGrid(),
		Command:   widget.NewTextGrid(),
		StdOut:    widget.NewTextGrid(),
		StdErr:    widget.NewTextGrid(),
		ErrorText: widget.NewLabel("Error goes here"),
	}
	item.ExtendBaseWidget(item)

	item.StdOut.ShowLineNumbers = true
	item.StdErr.ShowLineNumbers = true

	return item
}

func (item *LogWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVBox(item.Timestamp, item.Command, item.ErrorText, item.StdOut, item.StdErr)
	return widget.NewSimpleRenderer(c)
}
