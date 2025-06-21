package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"net/url"
)

func createMainMenu(b *BaseUI) *fyne.MainMenu {
	return fyne.NewMainMenu(createFileMenu(b), createHelpMenu(b))
}

func createFileMenu(b *BaseUI) *fyne.Menu {
	settingsItem := fyne.NewMenuItem("Settings", func() {
		b.ShowWindow(Settings)
	})
	settingsItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyComma,
		Modifier: fyne.KeyModifierSuper,
	}
	closeItem := fyne.NewMenuItem("Close Project", func() {
		b.projectWindow.Close()
		b.ShowWindow(Start)
	})
	aboutItem := fyne.NewMenuItem("About", func() { b.ShowWindow(About) })

	fileMenu := fyne.NewMenu("File", closeItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		fileMenu.Items = append(fileMenu.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	fileMenu.Items = append(fileMenu.Items, aboutItem)
	return fileMenu
}

func createHelpMenu(b *BaseUI) *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://github.com/repomancer/repomancer")
			_ = b.App.OpenURL(u)
		}),
	)
}
