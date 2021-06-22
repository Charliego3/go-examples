package main

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"runtime"
	"time"
)

func main() {
	golog.SetLevel("debug")

	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		initStatusMenuBar()

		tableView := NewNSTableView(core.Rect(0, 0, 100, 200))
		//tableView.Send("draw:", nil)
		_ = tableView

		wRect := core.Rect(0, 0, 600, 665)
		w := cocoa.NSWindow_Init(
			wRect,
			cocoa.NSClosableWindowMask|
				cocoa.NSResizableWindowMask|
				cocoa.NSMiniaturizableWindowMask|
				//cocoa.NSTexturedBackgroundWindowMask|
				cocoa.NSFullSizeContentViewWindowMask|
				cocoa.NSTitledWindowMask,
			cocoa.NSBackingStoreRetained,
			false,
		)

		tfRect := core.Rect(10, 200, 200, 21)
		textField := widgets.NewNSTextField(tfRect)
		//textField.SetBackgroundColor(cocoa.Color(206, 85, 33, 1))
		//textField.Set("bezelStyle:", objc.Get("NSTextFieldRoundedBezel"))
		//textField.Set("drawsBackground:", true)
		//textField.SetWantsLayer(true)
		//textField.Layer().SetCornerRadius(5)
		//textField.Layer().Set("borderWidth:", 2)
		//textField.Layer().Set("masksToBounds:", true)
		//inView := cocoa.NSView_Init(tfRect)
		//inView.Set("")
		//textField.Send("drawWithFrame:inView:", tfRect, nil)

		sfRect := core.Rect(10, 100, 200, 21)
		searchField := widgets.NewNSSearchField(sfRect)
		searchField.SetRecentSearches("abc", "ddd", "TextField")

		textField.SetStringValue("TextFieldTest")
		textField.SetIsBordered(true)
		textField.Set("placeholderString:", core.String("PlaceholderString"))
		textField.Set("drawsBackground:", true)

		textView := cocoa.NSTextView_Init(core.Rect(0, 320, 300, 100))
		textView.SetSelectable(true)
		textView.SetString("test")

		rect := core.Rect(0, 0, 600, 665)
		view := cocoa.NSView_Init(rect)
		view.SetWantsLayer(true)
		view.Send("addSubview:", &tableView)
		view.Send("addSubview:", &textView)
		view.Send("addSubview:", &searchField)
		view.Send("addSubview:", &textField)
		//view.SetWantsLayer(true)
		//view.Layer().SetCornerRadius(32.0)
		//view.AddSubviewPositionedRelativeTo(textField, cocoa.NSWindowBelow, nil)
		//view.AddSubviewPositionedRelativeTo(textView, cocoa.NSWindowBelow, nil)
		//view.AddSubviewPositionedRelativeTo(tableView, cocoa.NSWindowAbove, nil)
		//view.Send("draw:", &rect)
		//view.AddSubviewPositionedRelativeTo(tableView, 1, w)

		w.Send("setMinSize:", core.Size(300, 300))
		w.Send("setContentMinSize:", core.Size(300, 300))
		w.SetContentView(view)
		w.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		w.SetTitlebarAppearsTransparent(true)
		w.SetIgnoresMouseEvents(false)
		w.SetMovableByWindowBackground(false)
		//w.SetBackgroundColor(cocoa.NSColor_Init(46, 81, 133, 1))
		//w.SetBackgroundColor(cocoa.NSColor_Clear())
		w.SetLevel(0)
		w.SetTitle("NSTableView")
		w.MakeKeyAndOrderFront(view)
		w.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
		w.Center()

		//alert := widgets.NewNSAlert()
		//alert.SetMessageText("Alert message")
		//alert.SetInformativeText("Detailed description of alert message")
		//alert.AddButtonWithTitle("Default")
		//alert.AddButtonWithTitle("Alternative")
		//alert.AddButtonWithTitle("Other")
		//alert.BeginSheetModalForWindow(w, func(resp objc.Object) {
		//	println(resp)
		//})
		//alert.Show()
	})

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemQuit)
	app.SetMainMenu(menu)

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func ShowAlert(win objc.Object) {
	alert := alert.NewNSAlert()
	alert.SetMessageText("Alert message")
	alert.SetInformativeText("Detailed description of alert message")
	alert.AddButtonWithTitle("Default")
	alert.AddButtonWithTitle("Alternative")
	alert.AddButtonWithTitle("Other")

	//alert.Show()
	alert.BeginSheetModalForWindow(win.(cocoa.NSWindow))
}

func initStatusMenuBar() {
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	obj.Button().SetTitle("▶️ Ready")

	nextClicked := make(chan bool)
	go func() {
		state := -1
		timer := 1500
		countdown := false
		for {
			select {
			case <-time.After(1 * time.Second):
				if timer > 0 && countdown {
					timer = timer - 1
				}
				if timer <= 0 && state%2 == 1 {
					state = (state + 1) % 4
				}
			case <-nextClicked:
				state = (state + 1) % 4
				timer = map[int]int{
					0: 1500,
					1: 1500,
					2: 0,
					3: 300,
				}[state]
				if state%2 == 1 {
					countdown = true
				} else {
					countdown = false
				}
			}
			labels := map[int]string{
				0: "▶️ Ready %02d:%02d",
				1: "✴️ Working %02d:%02d",
				2: "✅ Finished %02d:%02d",
				3: "⏸️ Break %02d:%02d",
			}
			obj.Button().SetTitle(fmt.Sprintf(labels[state], timer/60, timer%60))
		}
	}()
	nextClicked <- true

	itemNext := cocoa.NSMenuItem_New()
	itemNext.SetTitle("Next")
	itemNext.SetAction(objc.Sel("nextClicked:"))
	cocoa.DefaultDelegateClass.AddMethod("nextClicked:", func(_ objc.Object) {
		nextClicked <- true
	})

	itemAlert := cocoa.NSMenuItem_New()
	itemAlert.SetTitle("Show Alert")
	itemAlert.SetAction(objc.Sel("showAlert:"))
	cocoa.DefaultDelegateClass.AddMethod("showAlert:", ShowAlert)

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemNext)
	menu.AddItem(itemAlert)
	menu.AddItem(cocoa.NSMenuItem_Separator())
	menu.AddItem(itemQuit)
	obj.SetMenu(menu)
}
