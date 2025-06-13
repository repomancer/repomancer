package screens

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/url"
)

func MakeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New Project", func() {
		log.Println("New Project")
		// TODO: New project
	})
	openItem := fyne.NewMenuItem("Open Project", func() {
		log.Println("Open Project")
		// TODO: Open project
	})
	//checkedItem := fyne.NewMenuItem("Checked", nil)
	//checkedItem.Checked = true
	//disabledItem := fyne.NewMenuItem("Disabled", nil)
	//disabledItem.Disabled = true
	//otherItem := fyne.NewMenuItem("Other", nil)

	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
	dirItem.Icon = theme.FolderIcon()

	openSettings := func() {
		w := NewSettingsWindow(a)
		w.Show()
	}
	showAbout := func() {
		w := a.NewWindow("About")
		msg := fmt.Sprintf("Repomancer %s (Build %d)", fyne.CurrentApp().Metadata().Version, fyne.CurrentApp().Metadata().Build)
		msg += "\n\n"
		msg += "Repository changes at scale"
		msg += fyne.CurrentApp().Metadata().Custom["Repository"]
		ca := fyne.CurrentApp().Metadata()
		log.Printf("%v", ca)
		homepage := "https://github.com/repomancer/repomancer"
		u, _ := url.Parse(homepage)
		w.SetContent(container.NewVBox(widget.NewLabel(msg), widget.NewHyperlink(homepage, u), widget.NewButton("Close", func() { w.Close() })))
		w.Show()
	}
	aboutItem := fyne.NewMenuItem("About", showAbout)
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		ShortcutFocused(cutShortcut, w)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		ShortcutFocused(copyShortcut, w)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		ShortcutFocused(pasteShortcut, w)
	})
	pasteItem.Shortcut = pasteShortcut
	//performFind := func() { fmt.Println("Menu Find") }
	//findItem := fyne.NewMenuItem("Find", performFind)
	//findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	//w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
	//	performFind()
	//})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://github.com/repomancer/repomancer")
			_ = a.OpenURL(u)
		}),
	)

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, openItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	file.Items = append(file.Items, aboutItem)
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator()),
		helpMenu,
	)
	return main
}

func ShortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func LogLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		//log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		//log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		//log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		//log.Println("Lifecycle: Exited Foreground")
	})
}
