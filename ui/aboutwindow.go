package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func NewAboutWindow(b *BaseUI) fyne.Window {
	w := b.NewWindow("About")
	msg := fmt.Sprintf(`
Repomancer %s (Build %d)

Repository changes at scale
`, fyne.CurrentApp().Metadata().Version, fyne.CurrentApp().Metadata().Build)

	homepage := "https://github.com/repomancer/repomancer"
	u, _ := url.Parse(homepage)
	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyEscape {
			w.Close()
		}
	})

	w.SetContent(container.NewVBox(widget.NewLabel(msg), widget.NewHyperlink(homepage, u), widget.NewButton("Close", func() { w.Close() })))
	return w
}
