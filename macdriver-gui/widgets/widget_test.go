package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"github.com/whimthen/temp/macdriver-gui/widgets/table"
	"runtime"
	"testing"
)

func TestTableView(t *testing.T) {
	registerView(newMenuItem("TableView Example", tableView))
}

func tableView(objc.Object) {
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

	scrollView := NewNSScrollView(core.Rect(0, 0, 500, 500))
	clipView := NewNSClipView()
	tableView := table.NewNSTableView(core.Rect(0, 0, 300, 300))
	clipView.SetDocumentView(tableView.NSView)
	scrollView.SetContentView(clipView)
	scrollView.SetHorizontalScroller(NewNSScroller())
	scrollView.SetBorderType(BezelBorderType)

	window.SetContentView(scrollView)
	window.SetTitleVisibility(cocoa.NSWindowTitleHidden)
	window.SetIgnoresMouseEvents(false)
	window.SetMovableByWindowBackground(false)
	window.SetLevel(0)
	window.MakeKeyAndOrderFront(scrollView)
	window.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
	window.Center()
}

func registerView(items ...cocoa.NSMenuItem) {
	cocoa.TerminateAfterWindowsClose = false
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		initStatusMenuBar(items...)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func initStatusMenuBar(items ...cocoa.NSMenuItem) {
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	obj.Button().SetTitle("🛠")

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))

	menu := cocoa.NSMenu_New()

	for _, item := range items {
		menu.AddItem(item)
	}

	menu.AddItem(newMenuItem("Show Alert", showAlert))
	menu.AddItem(newMenuItem("Open file selection", openFileSelection))
	menu.AddItem(newMenuItem("Open new window", openWindow))
	menu.AddItem(cocoa.NSMenuItem_Separator())
	menu.AddItem(itemQuit)
	obj.SetMenu(menu)
}

func newMenuItem(title string, callback func(object objc.Object)) cocoa.NSMenuItem {
	alertObj, alertSel := core.Callback(callback)
	itemAlert := cocoa.NSMenuItem_New()
	itemAlert.SetTitle(title)
	itemAlert.SetAction(alertSel)
	itemAlert.SetTarget(alertObj)
	return itemAlert
}

func openWindow(_ objc.Object) {
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

	nsAlert := alert.NewNSAlert_WithSheetModal(window, func(resp objc.Object) {
		println(resp)
	})
	nsAlert.SetAlertStyle(alert.Critical)
	nsAlert.SetMessageText("Alert message")
	nsAlert.SetInformativeText("Detailed description of nsAlert message")
	nsAlert.AddButtonWithTitle("Default")
	nsAlert.AddButtonWithTitle("Alternative")
	nsAlert.AddButtonWithTitle("Other")
	nsAlert.Send("beginSheetModalForWindow:completionHandler:", &window, nil)
}

func openFileSelection(_ objc.Object) {
	panel := NewNSOpenPanel()
	panel.SetAllowsMultipleSelection(false)
	panel.SetCanChooseDirectories(false)
	panel.SetCanCreateDirectories(false)
	panel.SetCanChooseFiles(true)
	panel.SetMessage("Please choose a html table content file")
	resp := panel.RunModalForDirectory()
	println(resp)
}

func showAlert(objc.Object) {
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
