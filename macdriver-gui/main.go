package main

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets"
	"runtime"
	"time"
)

func main() {
	textField()
}

func textField() {
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
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

		textField := widgets.NewNSTextField(core.Rect(10, 10, 100, 22))

		rootView.Send("addSubview:", &textField)

		win.SetContentView(rootView)
		//win.SetTitleVisibility(cocoa.NSWindowTitleHidden)
		win.SetIgnoresMouseEvents(false)
		win.SetMovableByWindowBackground(false)
		win.SetLevel(0)
		win.MakeKeyAndOrderFront(rootView)
		win.SetCollectionBehavior(cocoa.NSWindowCollectionBehaviorCanJoinAllSpaces)
		win.Center()
	})

	app.Retain()
	app.Run()
}

func statusBarApp() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		golog.Errorf("Param: %+v", notification)

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
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)

	})
	app.Run()
}
