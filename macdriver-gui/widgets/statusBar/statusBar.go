package statusBar

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"runtime"
)

type StatusMenuBarApplication struct {
	app  cocoa.NSApplication
	menu cocoa.NSMenu
}

func NewStatusBarApp(title string, length float64) StatusMenuBarApplication {
	cocoa.TerminateAfterWindowsClose = false
	runtime.LockOSThread()
	menu := cocoa.NSMenu_New()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(length)
		obj.Retain()
		obj.Button().SetTitle(title)
		obj.SetMenu(menu)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	return StatusMenuBarApplication{app: app, menu: menu}
}

func (a StatusMenuBarApplication) AddMenuItem(title string, action func(object objc.Object)) {
	obj, sel := core.Callback(action)
	item := cocoa.NSMenuItem_New()
	item.SetTitle(title)
	item.SetAction(sel)
	item.SetTarget(obj)
	a.menu.AddItem(item)
}

func (a StatusMenuBarApplication) AddMenuItemWithSelector(title string, sel objc.Selector) {
	item := cocoa.NSMenuItem_New()
	item.SetTitle(title)
	item.SetAction(sel)
	a.menu.AddItem(item)
}

func (a StatusMenuBarApplication) AddTerminateItem(title ...string) {
	itemTitle := "Quit"
	if len(title) > 0 {
		itemTitle = title[0]
	}
	a.AddMenuItemWithSelector(itemTitle, objc.Sel("terminate:"))
}

func (a StatusMenuBarApplication) AddItemSeparator() {
	a.menu.AddItem(cocoa.NSMenuItem_Separator())
}

func (a StatusMenuBarApplication) Run() {
	a.app.Run()
}
