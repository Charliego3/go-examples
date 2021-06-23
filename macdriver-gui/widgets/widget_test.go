package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"github.com/whimthen/temp/macdriver-gui/widgets/statusBar"
	"github.com/whimthen/temp/macdriver-gui/widgets/table"
	"testing"
)

func TestTableView(t *testing.T) {
	app := statusBar.NewStatusBarApp("🛠", cocoa.NSSquareStatusItemLength)
	app.AddMenuItem("TableView Example", tableView)
	app.AddMenuItem("Show Alert", showAlert)
	app.AddMenuItem("Open file selection", openFileSelection)
	app.AddMenuItem("Open new window", openWindow)
	app.AddItemSeparator()
	app.AddTerminateItem()
	app.Run()
}

func tableView(objc.Object) {
	window := cocoa.NSWindow_Init(
		core.Rect(0, 0, 600, 665),
		cocoa.NSClosableWindowMask|
			cocoa.NSResizableWindowMask|
			cocoa.NSMiniaturizableWindowMask|
			cocoa.NSFullSizeContentViewWindowMask|
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

	nsAlert := alert.NewNSAlert()
	nsAlert.SetAlertStyle(alert.Critical)
	nsAlert.SetMessageText("Alert message")
	nsAlert.SetInformativeText("Detailed description of nsAlert message")
	nsAlert.AddButtonWithTitle("Default")
	nsAlert.AddButtonWithTitle("Alternative")
	//nsAlert.SetShowsSuppressionButton(true)
	//nsAlert.SetSuppressionButtonTitle("Don't Set again.......")
	//nsAlert.AddButtonWithTitle("Other")
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
