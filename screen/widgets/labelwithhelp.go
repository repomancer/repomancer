package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type tappableIcon struct {
	widget.Icon
	onTapped func()
}

func newTappableIcon(res fyne.Resource, onTapped func()) *tappableIcon {
	icon := &tappableIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	icon.onTapped = onTapped
	return icon
}

func (t *tappableIcon) Tapped(_ *fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

type LabelWithHelp struct {
	widget.BaseWidget
	Label *widget.Label
	Icon  *tappableIcon
}

func NewLabelWithHelpWidget(title, comment string, window fyne.Window) *LabelWithHelp {
	item := &LabelWithHelp{
		Label: widget.NewLabel(title),
		Icon: newTappableIcon(theme.HelpIcon(), func() {
			i := dialog.NewInformation(title, comment, window)
			i.Show()
		}),
	}
	item.Label.Alignment = fyne.TextAlignTrailing
	item.Label.Truncation = fyne.TextTruncateOff
	item.ExtendBaseWidget(item)

	return item
}

func (item *LabelWithHelp) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(item.Label, item.Icon)
	return widget.NewSimpleRenderer(c)
}
