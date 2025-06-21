package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ToggleIconWidget struct {
	widget.Icon
	OnTapped func()
	Selected bool
}

func NewToggleWidget(fn func()) *ToggleIconWidget {
	entry := &ToggleIconWidget{}
	entry.Resource = theme.CheckButtonIcon()
	entry.OnTapped = fn
	entry.Selected = false
	entry.ExtendBaseWidget(entry)
	return entry
}

func (t *ToggleIconWidget) SetIcon(icon fyne.Resource) {
	t.Icon.SetResource(icon)
}
func (t *ToggleIconWidget) Tapped(_ *fyne.PointEvent) {
	if t.OnTapped == nil {
		return
	}
	t.Selected = !t.Selected
	t.OnTapped()
}
