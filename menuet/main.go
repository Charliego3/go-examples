package main

import "github.com/caseymrm/menuet"

func main() {
	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		// Title: "Tools",
		Image: "hammer",
	})
	app.Label = "com.github.charlie.tools"
	app.Children = menuItems

	app.RunApplication()
}

func menuItems() []menuet.MenuItem {
	var items []menuet.MenuItem
	items = append(items, proxyItem(items))
	items = append(items, menuet.MenuItem{Text: "empty"})
	return items
}

func notification(title, stitle, msg string) {
	menuet.App().Notification(menuet.Notification{
		Title:                        title,
		Subtitle:                     stitle,
		Message:                      msg,
		ActionButton:                 "OK",
		CloseButton:                  "Close",
		RemoveFromNotificationCenter: true,
	})
}
