package main

import (
	"fmt"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"runtime"
	"time"
)

func main() {
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		initStatusMenuBar()

		tableView := NSTableView_Init(core.Rect(0, 0, 100, 200))
		//tableView.Send("draw:", nil)

		w := cocoa.NSWindow_Init(
			core.Rect(0, 0, 500, 300),
			cocoa.NSClosableWindowMask|
				cocoa.NSResizableWindowMask|
				cocoa.NSMiniaturizableWindowMask|
				cocoa.NSFullSizeContentViewWindowMask|
				cocoa.NSTitledWindowMask,
			cocoa.NSBackingStoreRetained,
			false,
		)

		view := cocoa.NSView_Init(core.Rect(0, 0, 500, 300))
		view.AddSubviewPositionedRelativeTo(objc.Get("NSTextField").Alloc().Init(), 0, nil)
		view.AddSubviewPositionedRelativeTo(tableView, 1, nil)

		textView := cocoa.NSTextView_Init(core.Rect(0, 0, 500, w.Frame().Size.Height-29))
		textView.SetSelectable(true)
		textView.SetString("test")

		w.SetContentView(tableView)
		w.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		w.SetTitlebarAppearsTransparent(true)
		w.SetIgnoresMouseEvents(false)
		w.SetMovableByWindowBackground(false)
		//w.SetLevel(0)
		w.SetBackgroundColor(cocoa.NSColor_Init(46, 81, 133, 1))
		w.SetTitle("NSTableView")
		w.MakeKeyAndOrderFront(tableView)
		w.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorDefault)
		w.Center()
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

	itemQuit := cocoa.NSMenuItem_New()
	itemQuit.SetTitle("Quit")
	itemQuit.SetAction(objc.Sel("terminate:"))

	menu := cocoa.NSMenu_New()
	menu.AddItem(itemNext)
	menu.AddItem(cocoa.NSMenuItem_Separator())
	menu.AddItem(itemQuit)
	obj.SetMenu(menu)
}
