package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type EscapeEntry struct {
	widget.Entry
	OnEscape  func()
	MaxLength int
}

func (m *EscapeEntry) TypedKey(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyEscape && m.OnEscape != nil {
		m.OnEscape()
	}
	if event.Name == fyne.KeyEnter || event.Name == fyne.KeyReturn && m.OnSubmitted != nil {
		m.OnSubmitted(m.Text)
	}
	m.Entry.TypedKey(event)
}

func (m *EscapeEntry) TypedRune(r rune) {
	if m.MaxLength == -1 || len(m.Text) < m.MaxLength {
		m.Entry.TypedRune(r)
	}
}

func NewEscapeEntry() *EscapeEntry {
	item := &EscapeEntry{}
	item.ExtendBaseWidget(item)
	item.MaxLength = -1
	return item
}
