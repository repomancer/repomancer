package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// ShortcutHandlingEntry is an Entry field that also handles length, allowed characters
// and handles Cmd-W to hide/close the current window
// HandleShortcut the function that will be called for Cmd-W
// AllowedCharacters string containing all valid characters. If "", everything is allowed
// MaxLength maximum number of characters, or -1 for unlimited
type ShortcutHandlingEntry struct {
	widget.Entry
	HandleShortcut    func(*desktop.CustomShortcut)
	AllowedCharacters string
	MaxLength         int
}

func (m *ShortcutHandlingEntry) TypedShortcut(s fyne.Shortcut) { //for local
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		return
	} else {
		t := s.(*desktop.CustomShortcut)
		if m.HandleShortcut != nil {
			m.HandleShortcut(t)
		}
	}
}

func (m *ShortcutHandlingEntry) TypedRune(r rune) {
	if m.MaxLength == -1 || len(m.Text) < m.MaxLength {
		if m.AllowedCharacters == "" {
			m.Entry.TypedRune(r)
		} else {
			if strings.ContainsRune(m.AllowedCharacters, r) {
				m.Entry.TypedRune(r)
			}
		}
	}
}

func NewShortcutHandlingEntry(window fyne.Window, isMainWindow bool) *ShortcutHandlingEntry {
	item := &ShortcutHandlingEntry{}
	item.ExtendBaseWidget(item)
	item.HandleShortcut = func(s *desktop.CustomShortcut) {
		if s.KeyName == fyne.KeyW && s.Modifier == desktop.SuperModifier {
			if isMainWindow {
				window.Close()
			} else {
				window.Hide()
			}
		}
	}
	item.MaxLength = -1
	item.AllowedCharacters = ""
	return item
}
