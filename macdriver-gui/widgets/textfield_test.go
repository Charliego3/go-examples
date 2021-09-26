package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"testing"
)

func TestNSTextField(t *testing.T) {
	t.Logf("%#v", objc.Get("NSDictionaryOfVariableBindings"))

	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		win := cocoa.NSWindow_Init(
			core.Rect(0, 0, 600, 665),
			cocoa.NSClosableWindowMask|
				cocoa.NSResizableWindowMask|
				cocoa.NSMiniaturizableWindowMask|
				//cocoa.NSTexturedBackgroundWindowMask|
				cocoa.NSFullSizeContentViewWindowMask|
				cocoa.NSTitledWindowMask,
			cocoa.NSBackingStoreRetained,
			false,
		)
		win.SetHasShadow(true)
		//win.SetTitlebarAppearsTransparent(true)

		tfRect := core.Rect(10, 200, 200, 21)
		textField := NewNSTextField(tfRect)
		textField.Set("placeholderString:", core.String("PlaceholderString"))
		textField.Set("drawsBackground:", true)
		t.Logf("TextField: %#v", textField)

		sfRect := core.Rect(10, 100, 200, 21)
		searchField := NewNSSearchField(sfRect)
		searchField.SetRecentSearches("abc", "ddd", "TextField")

		view := cocoa.NSView_Init(win.Frame())
		view.SetWantsLayer(true)
		view.Send("addSubview:", &textField)
		view.Send("addSubview:", &searchField)

		win.SetTitle("NSTextField")
		win.SetContentView(view)
		win.Send("setMinSize:", core.Size(300, 300))
		win.SetIgnoresMouseEvents(false)
		win.SetMovableByWindowBackground(false)
		win.SetLevel(0)
		win.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		win.MakeKeyAndOrderFront(view)
		win.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
		win.Center()
	})

	mainMenu := cocoa.NSMenu_New()
	mainMenu.SetTitle("MainMenu")
	rootMenu := cocoa.NSMenu_New()
	rootMenu.SetTitle("root_menu")

	obj, sel := core.Callback(func(object objc.Object) {})
	item := cocoa.NSMenuItem_New()
	item.SetTitle("menu1")
	item.SetAction(sel)
	item.SetTarget(obj)
	//mainMenu.Send("setSubmenu:forItem:", rootMenu, item)
	mainMenu.AddItem(item)

	t.Logf("NSApp.mainMenu: %v", mainMenu)
	app.Set("windowsMenu:", mainMenu)

	app.Retain()
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.Run()
}
