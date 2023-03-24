package main

import (
	"github.com/caseymrm/menuet"
)

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
	items = append(items, generatePassword())
	items = append(items, openProject())
	items = append(items, openTerminal())
	return items
}

func notify(title, msg string) {
	menuet.App().Alert(menuet.Alert{
		MessageText:     title,
		InformativeText: msg,
		Buttons:         []string{"OK"},
	})
}

func warning(msg string) {
	notify("Oops!", msg)
}
