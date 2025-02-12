package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ShortcutHandlingEntry struct {
	widget.Entry
	HandleShortcut func(*desktop.CustomShortcut)
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

	return item
}
