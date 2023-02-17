package main

import (
	"github.com/caseymrm/menuet"
)

func main() {
	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		//Title: "Tools",
		Image: "hammer",
	})
	app.Label = "com.github.charlie.tools"
	app.Children = menuItems

	app.RunApplication()
}

func menuItems() []menuet.MenuItem {
	var items []menuet.MenuItem
	items = append(items, proxyItem(items))
	items = append(items, generatePassword())
	return items
}

func notify(title, stitle, msg string) {
	menuet.App().Alert(menuet.Alert{
		MessageText:     title,
		InformativeText: msg,
		Buttons:         []string{"OK"},
	})
}
