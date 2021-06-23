package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"runtime"
)

func init() {
}

func main() {
	cocoa.TerminateAfterWindowsClose = false
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		initStatusMenuBar()
		//initEmptyWin()
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
	obj.SetMenu(menu)

	menu.AddItem(itemAlert)
	menu.AddItem(html2csv())
	menu.AddItem(openWindow())
	menu.AddItem(cocoa.NSMenuItem_Separator())
	menu.AddItem(itemQuit)
}

func openWindow() cocoa.NSMenuItem {
	call, selector := core.Callback(func(objc.Object) {
		//window := cocoa.NSWindow_New()
		window := cocoa.NSWindow_Init(
			core.Rect(0, 0, 600, 665),
			cocoa.NSClosableWindowMask|
				cocoa.NSResizableWindowMask|
				cocoa.NSMiniaturizableWindowMask|
				cocoa.NSFullSizeContentViewWindowMask|
				cocoa.NSTexturedBackgroundWindowMask|
				cocoa.NSTitledWindowMask,
			cocoa.NSBackingStoreBuffered,
			false,
		)
		window.SetHasShadow(true)
		window.SetTitlebarAppearsTransparent(true)

		rect := core.Rect(0, 0, 600, 665)
		view := cocoa.NSView_Init(rect)

		window.SetContentView(view)
		window.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		window.SetIgnoresMouseEvents(false)
		window.SetMovableByWindowBackground(false)
		window.SetLevel(0)
		window.MakeKeyAndOrderFront(view)
		window.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
		window.Center()

		nsAlert := alert.NewNSAlert()
		nsAlert.SetAlertStyle(alert.Critical)
		nsAlert.SetMessageText("Alert message")
		nsAlert.SetInformativeText("Detailed description of nsAlert message")
		nsAlert.AddButtonWithTitle("Default")
		nsAlert.AddButtonWithTitle("Alternative")
		nsAlert.AddButtonWithTitle("Other")
		//nsAlert.Send("showSheetModalForWindow")
		//nsAlert.Show()
		//sel := objc.Sel("buttonAction:")
		nsAlert.Send("beginSheetModalForWindow:completionHandler:", &window, core.String("buttonAction:"))
	})
	openWin := cocoa.NSMenuItem_New()
	openWin.SetTitle("Open new Window")
	openWin.SetAction(selector)
	openWin.SetTarget(call)
	return openWin
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
	nsAlert := alert.NewNSAlert()
	nsAlert.SetAlertStyle(alert.Informational)
	nsAlert.SetShowsHelp(true)
	nsAlert.Set("helpAnchor:", core.String("www.baidu.com"))
	nsAlert.SetMessageText("Alert message")
	nsAlert.SetInformativeText("Detailed description of nsAlert message")
	nsAlert.AddButtonWithTitle("Default")
	nsAlert.AddButtonWithTitle("Alternative")
	nsAlert.AddButtonWithTitle("Other")
	nsAlert.Show()
}
