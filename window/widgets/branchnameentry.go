package widgets

import "fyne.io/fyne/v2/widget"

type BranchNameEntry struct{ widget.Entry }

func NewBranchNameEntry() *BranchNameEntry {
	entry := &BranchNameEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *BranchNameEntry) TypedRune(r rune) {
	if len(e.Text) < 45 {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '/' {
			e.Entry.TypedRune(r)
		}
	}
}
