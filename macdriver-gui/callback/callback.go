package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/button"
	"runtime"
)

var (
	nsAlert = objc.Get("NSAlert")

	nsAlertController      objc.Object
	nsAlertControllerClass objc.Class
)

func init() {
	nsAlertControllerClass = objc.NewClass("DefaultNSAlertDelegate", "NSAlert")
	nsAlertControllerClass.AddMethod("beginSheetModalForWindow:", func(alert, handler objc.Object) {
		println(alert)
	})
	objc.RegisterClass(nsAlertControllerClass)
	nsAlertController = objc.Get("DefaultNSAlertDelegate").Alloc().Init()
}

type NSAlert struct {
	objc.Object
}

func NewNSAlert() NSAlert {
	alert := NSAlert{nsAlert.Alloc().Init()}
	alert.SetDelegate(nsAlertController)
	return alert
}

func (i NSAlert) SetDelegate(delegate objc.Object) {
	i.Send("setDelegate:", delegate)
}

func (i NSAlert) Delegate() objc.Object {
	return i.Send("delegate")
}

func (i NSAlert) MessageText() string {
	return i.Get("messageText").String()
}

func (i NSAlert) SetMessageText(s string) {
	i.Set("messageText:", core.String(s))
}

func (i NSAlert) InformativeText() string {
	return i.Get("informativeText").String()
}

func (i NSAlert) SetInformativeText(s string) {
	i.Set("informativeText:", core.String(s))
}

func (i NSAlert) AddButtonWithTitle(s string) {
	i.Send("addButtonWithTitle:", core.String(s))
}

func (i NSAlert) BeginSheetModalForWindow(win cocoa.NSWindow) objc.Object {
	return i.Send("beginSheetModalForWindow:completionHandler:", &win, nil)
}

func main() {
	runtime.LockOSThread()
	app := cocoa.NSApp_WithDidLaunch(wenAppLaunch)
	app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyAccessory)
	app.ActivateIgnoringOtherApps(true)
	app.Run()
}

func wenAppLaunch(notification objc.Object) {
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

	rect := core.Rect(0, 0, 600, 665)
	view := cocoa.NSView_Init(rect)
	view.Set("setTranslatesAutoresizingMaskIntoConstraints:", false)

	nsButton := button.NSButton{NSView: cocoa.NSView{Object: objc.Get("NSButton").Alloc().Init()}}
	nsButton.Set("title:", core.String("titled button"))
	nsButton.Set("")
	view.Send("addSubview:", &nsButton)

	btn := button.NewButtonWithFrame(core.Rect(100, 100, 200, 22))
	btn.SetTitle("Show Alert with sheet modal")
	btn.SetType()
	btn.SetAction(func(object objc.Object) {
		showAlertWithSheet(window)
	})
	view.Send("addSubview:", &btn)

	window.SetTitle("Test sheet modal alert")
	window.SetContentView(view)
	window.MakeKeyAndOrderFront(view)
	window.Center()
}

func showAlertWithSheet(window cocoa.NSWindow) {
	nsAlert := NewNSAlert()
	nsAlert.SetMessageText("Alert test sheet message")
	nsAlert.SetInformativeText("Detailed description of nsAlert message")
	nsAlert.AddButtonWithTitle("OK")
	nsAlert.AddButtonWithTitle("Second")
	nsAlert.BeginSheetModalForWindow(window)
}
