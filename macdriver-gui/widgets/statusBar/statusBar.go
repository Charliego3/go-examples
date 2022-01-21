package statusBar

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"runtime"
)

type StatusMenuBarApplication struct {
	cocoa.NSApplication
	Menu cocoa.NSMenu
	obj  cocoa.NSStatusItem
}

type SubMenu struct {
	SubTitle string
	Action   func(object objc.Object)
}

func NewStatusBarApp(title string, fn func(item cocoa.NSStatusItem)) StatusMenuBarApplication {
	cocoa.TerminateAfterWindowsClose = false
	runtime.LockOSThread()
	menu := cocoa.NSMenu_New()
	var obj cocoa.NSStatusItem
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		obj = cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle(title)
		obj.SetMenu(menu)

		fn(obj)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	return StatusMenuBarApplication{NSApplication: app, Menu: menu, obj: obj}
}

func (a StatusMenuBarApplication) SetTitle(title string) {
	a.obj.Button().SetTitle(title)
}

func (a StatusMenuBarApplication) AddSubMenu(title string, menus ...SubMenu) {
	subItem := cocoa.NSMenuItem_New()
	subItem.SetTitle(title)
	subMenu := cocoa.NSMenu_New()
	subItem.SetSubmenu(subMenu)

	for _, menu := range menus {
		t1 := cocoa.NSMenuItem_New()
		t1.SetTitle(menu.SubTitle)
		object, selector := core.Callback(menu.Action)
		t1.SetTarget(object)
		t1.SetAction(selector)
		subMenu.AddItem(t1)
	}

	a.Menu.AddItem(subItem)
}

func (a StatusMenuBarApplication) AddMenuItem(title string, action func(object objc.Object)) cocoa.NSMenuItem {
	obj, sel := core.Callback(action)
	item := cocoa.NSMenuItem_New()
	item.SetTitle(title)
	item.SetAction(sel)
	item.SetTarget(obj)
	a.Menu.AddItem(item)
	return item
}

func (a StatusMenuBarApplication) AddMenuItemWithSelector(title string, sel objc.Selector) cocoa.NSMenuItem {
	item := cocoa.NSMenuItem_New()
	item.SetTitle(title)
	item.SetAction(sel)
	a.Menu.AddItem(item)
	return item
}

func (a StatusMenuBarApplication) AddTerminateItem(title ...string) {
	itemTitle := "Quit"
	if len(title) > 0 {
		itemTitle = title[0]
	}
	a.AddMenuItemWithSelector(itemTitle, objc.Sel("terminate:"))
}

func (a StatusMenuBarApplication) AddItemSeparator() {
	a.Menu.AddItem(cocoa.NSMenuItem_Separator())
}

func (a StatusMenuBarApplication) Run() {
	a.NSApplication.Run()
}
