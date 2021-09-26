package widgets

import (
	"github.com/kataras/golog"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"github.com/whimthen/temp/macdriver-gui/widgets/statusBar"
	"github.com/whimthen/temp/macdriver-gui/widgets/table"
	"github.com/whimthen/temp/macdriver-gui/widgets/type_alias"
	"testing"
)

func TestConstraint(t *testing.T) {
	//t.Log(objc.Get("NSCommandKeyMask"))
	//
	//
	//class := objc.Get("NSLayoutConstraint")
	//t.Logf("ConstraintClass: %#v", class)
	//alloc := class.Alloc()
	//t.Logf("ConstraintAlloc: %#v", alloc)
	//instance := alloc.Init()
	//t.Logf("Constraint: %#v", instance)
	rootView := cocoa.NSView_Init(core.Rect(0, 0, 100, 100))
	subView := cocoa.NSView_Init(core.Rect(0, 0, 200, 300))
	subView.SetBackgroundColor(cocoa.Color(255, 255, 0, 1))
	subView.SetWantsLayer(true)
	subView.Set("translatesAutoresizingMaskIntoConstraints:", core.False)
	rootView.Send("addSubview:", subView)

	dict := core.NSDictionary_Init(rootView, core.String("view1"), subView, core.String("view2"))
	t.Log(dict)

	t.Logf("Constraint: %#v", NewNSLayoutConstraintWithFormat(rootView, subView))
	constraint(nil)
}

func TestViews(t *testing.T) {
	app := statusBar.NewStatusBarApp("🛠", cocoa.NSSquareStatusItemLength)
	app.AddSubMenu("Window Example",
		statusBar.SubMenu{
			SubTitle: "Layout Constraint",
			Action:   layoutConstraint,
		},
		statusBar.SubMenu{
			SubTitle: "Test1",
			Action: func(object objc.Object) {

			},
		},
	)

	app.AddMenuItem("TableView Example", tableView)
	app.AddMenuItem("Show Alert", showAlert)
	app.AddMenuItem("Open file selection", openFileSelection)
	app.AddMenuItem("Open new window", openWindow)
	app.AddMenuItem("Constraint", constraint)
	app.AddItemSeparator()
	app.AddTerminateItem()
	app.Run()
}

func constraint(objc.Object) {
	win := cocoa.NSWindow_Init(
		core.Rect(0, 0, 600, 665),
		cocoa.NSClosableWindowMask|
			cocoa.NSResizableWindowMask|
			cocoa.NSMiniaturizableWindowMask|
			cocoa.NSFullSizeContentViewWindowMask|
			cocoa.NSTitledWindowMask,
		cocoa.NSBackingStoreBuffered,
		false,
	)
	win.SetHasShadow(true)
	//win.SetTitlebarAppearsTransparent(true)

	rootView := cocoa.NSView_Init(win.Frame())

	subView := cocoa.NSView_Init(core.Rect(0, 0, 200, 300))
	subView.SetBackgroundColor(cocoa.Color(255, 255, 0, 1))
	subView.SetWantsLayer(true)
	//subView.Layer().SetCornerRadius(5)
	subView.Set("translatesAutoresizingMaskIntoConstraints:", false)

	rootView.Send("addSubview:", subView)
	//constraint := objc.Get("NSLayoutConstraint").Alloc().
	//	Send("constraintWithItem:attribute:relatedBy:toItem:attribute:multiplier:constant:",
	//		&subView, NSLayoutAttributeTop, NSLayoutRelationEqual, &rootView, NSLayoutAttributeLeft, 1.0, 100.0)

	//topConstraint := NewNSLayoutConstraint()
	//topConstraint.SetConstraintWithItem(subView, NSLayoutAttributeLeft, NSLayoutRelationEqual, rootView, NSLayoutAttributeLeft, 1.0, 40)
	//topConstraint := NewNSLayoutConstraintWithAttr(subView,
	//	NSLayoutAttributeLeft,
	//	NSLayoutRelationEqual,
	//	rootView,
	//	NSLayoutAttributeLeft,
	//	1.0, 50.0,
	//)
	//constraint := objc.Get("NSLayoutConstraint").Alloc()
	//constraint.Set("constant:", 200.0)
	rootView.Send("addConstraints:", NewNSLayoutConstraintWithFormat(rootView, subView))

	win.SetContentView(rootView)
	//win.SetTitleVisibility(cocoa.NSWindowTitleHidden)
	win.SetIgnoresMouseEvents(false)
	win.SetMovableByWindowBackground(false)
	win.SetLevel(0)
	win.MakeKeyAndOrderFront(rootView)
	win.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
	win.Center()
}

func layoutConstraint(object objc.Object) {
	win := cocoa.NSWindow_Init(
		core.Rect(0, 0, 600, 665),
		cocoa.NSClosableWindowMask|
			cocoa.NSResizableWindowMask|
			cocoa.NSMiniaturizableWindowMask|
			cocoa.NSFullSizeContentViewWindowMask|
			cocoa.NSTitledWindowMask,
		cocoa.NSBackingStoreBuffered,
		false,
	)
	win.SetHasShadow(true)
	win.SetTitlebarAppearsTransparent(true)

	rootView := cocoa.NSView_Init(win.Frame())

	subView := cocoa.NSView{Object: objc.Get("NSView").Alloc().Init()}
	subView.SetBackgroundColor(cocoa.Color(255, 255, 0, 1))
	subView.SetWantsLayer(true)
	subView.Layer().SetCornerRadius(10)
	subView.Set("translatesAutoresizingMaskIntoConstraints:", false)

	rootView.Send("addSubview:", subView)
	rootView.Send("addConstraint:", NewNSLayoutConstraintWithAttr(subView,
		NSLayoutAttributeLeft,
		NSLayoutRelationEqual,
		rootView,
		NSLayoutAttributeLeft,
		1.0, 10.0,
	))
	rootView.Send("addConstraint:", NewNSLayoutConstraintWithAttr(subView,
		NSLayoutAttributeRight,
		NSLayoutRelationEqual,
		rootView,
		NSLayoutAttributeRight,
		1.0, 10.0,
	))
	rootView.Send("addConstraint:", NewNSLayoutConstraintWithAttr(subView,
		NSLayoutAttributeTop,
		NSLayoutRelationEqual,
		rootView,
		NSLayoutAttributeTop,
		1.0, 10.0,
	))
	rootView.Send("addConstraint:", NewNSLayoutConstraintWithAttr(subView,
		NSLayoutAttributeHeight,
		NSLayoutRelationEqual,
		rootView,
		NSLayoutAttributeHeight,
		1.0, 200.0,
	))

	win.SetContentView(rootView)
	win.SetTitleVisibility(cocoa.NSWindowTitleHidden)
	win.SetIgnoresMouseEvents(false)
	win.SetMovableByWindowBackground(false)
	win.SetLevel(0)
	win.MakeKeyAndOrderFront(rootView)
	win.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
	win.Center()
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

	identifier := type_alias.NewNSUserInterfaceItemIdentifier("tablecell")
	_ = identifier
	golog.Errorf("Identifier: %+v", identifier)

	rect := core.Rect(0, 0, 500, 500)
	sv := NewNSScrollView(rect)
	sv.Set("verticalLineScroll:", float64(10))
	clipView := NewNSClipView()
	tableView := table.NewNSTableView(rect)
	c1 := table.NewNSTableColumn("Column1")
	c1.SetTitle("Column1 Title")
	c1.Set("minWidth:", float64(150))
	c1.SetHeaderCell(table.NewNSTableHeaderCell("Column1 HeaderCell1"))
	c2 := table.NewNSTableColumn("Column2")
	c2.SetTitle("Column2 Title")
	c2.Set("editable:", true)
	c2.Set("headerToolTip:", core.String("Header ToolTip"))
	c2.SetHeaderCell(table.NewNSTableHeaderCell("Column2 HeaderCell2"))
	c3 := table.NewNSTableColumn("Number")
	c3.SetHeaderCell(table.NewNSTableHeaderCell("序号"))
	tableView.AddTableColumn(c3, c1, c2)
	tableView.SetSelectionHighlightStyle(table.NSTableViewSelectionHighlightStyleRegular)
	tableView.SetRowHeight(16)
	tableView.SetRowSizeStyle(table.NSTableViewRowSizeStyleCustom)
	//tableView.SetStyle(table.NSTableViewStyleAutomatic)
	tableView.SetGridStyleMask(table.NSTableViewSolidHorizontalGridLineMask)
	//tableView.SetGridColor(cocoa.Color(104, 104, 53, 1))
	clipView.SetDocumentView(tableView)
	sv.SetContentView(clipView)
	sv.SetHorizontalScroller(NewNSScroller())
	sv.SetVerticalScroller(NewNSScroller())
	sv.SetBorderType(NoBorderType)

	window.SetContentView(sv)
	window.SetTitleVisibility(cocoa.NSWindowTitleHidden)
	window.SetIgnoresMouseEvents(false)
	window.SetMovableByWindowBackground(false)
	window.SetLevel(0)
	window.MakeKeyAndOrderFront(sv)
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
	resp := panel.Show()
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
