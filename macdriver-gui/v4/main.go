package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func main() {
	app := cocoa.NSApp_WithDidLaunch(func(notification objc.Object) {
		rect := core.Rect(0, 0, 400, 500)
		window := cocoa.NSWindow_New()
		window.SetStyleMask(
			cocoa.NSClosableWindowMask |
				cocoa.NSTitledWindowMask |
				cocoa.NSMiniaturizableWindowMask |
				cocoa.NSResizableWindowMask |
				cocoa.NSBorderlessWindowMask |
				1<<12 |
				cocoa.NSFullSizeContentViewWindowMask,
		)

		window.SetBackingType(core.NSUInteger(2))
		window.SetTitle("Macdriver Demo")
		window.SetMinSize(core.Size(400, 500))

		contentController := cocoa.NSViewController_New()
		contentController.SetView(cocoa.NSView_Init(rect))
		window.SetContentViewController(contentController)
		window.Center()
		window.MakeKeyAndOrderFront(nil)
	})

	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}
