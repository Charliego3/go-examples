package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		initStatusMenuBar()
		initEmptyWin()
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func initEmptyWin() {
	rect := core.Rect(0, 0, 0, 0)
	w := cocoa.NSWindow_Init(
		rect,
		cocoa.NSClosableWindowMask,
		cocoa.NSBackingStoreRetained,
		false,
	)

	view := cocoa.NSView_Init(rect)
	w.SetContentView(view)
	w.MakeKeyAndOrderFront(view)
	w.SetLevel(-1)
	w.Center()
}

func initStatusMenuBar() {
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	obj.Button().SetTitle("🛠")

	alertObj, alertSel := core.Callback(ShowAlert)
	itemAlert := cocoa.NSMenuItem_New()
	itemAlert.SetTitle("Show Alert")
	itemAlert.SetAction(alertSel)
	itemAlert.SetTarget(alertObj)

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemAlert)
	menu.AddItem(html2csv())
	menu.AddItem(cocoa.NSMenuItem_Separator())
	menu.AddItem(itemQuit)
	obj.SetMenu(menu)
}

func html2csv() cocoa.NSMenuItem {
	call, selector := core.Callback(func(objc.Object) {
		OpenFileSelection()
	})
	csv := cocoa.NSMenuItem_New()
	csv.SetTitle("HTML table to csv")
	csv.SetAction(selector)
	csv.SetTarget(call)
	return csv
}

func OpenFileSelection() {
	panel := widgets.NewNSOpenPanel()
	panel.SetAllowsMultipleSelection(false)
	panel.SetCanChooseDirectories(false)
	panel.SetCanCreateDirectories(false)
	panel.SetCanChooseFiles(true)
	panel.SetMessage("Please choose a html table content file")
	resp := panel.RunModalForDirectory()
	println(resp)
}

func ShowAlert(objc.Object) {
	alert := widgets.NewNSAlert()
	alert.SetMessageText("Alert message")
	alert.SetInformativeText("Detailed description of alert message")
	alert.AddButtonWithTitle("Default")
	alert.AddButtonWithTitle("Alternative")
	alert.AddButtonWithTitle("Other")
	alert.Show()
}
