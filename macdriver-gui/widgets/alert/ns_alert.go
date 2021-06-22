package alert

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

var (
	nsAlert_ = objc.Get("NSAlert")

	nsAlertController objc.Object
	nsAlertControllerClass objc.Class
)

//func init() {
//	nsAlertControllerClass = objc.NewClass("NSAlertController", "NSAlert")
//	objc.RegisterClass(nsAlertControllerClass)
//	nsAlertController = objc.Get("NSAlertController").Alloc().Init()
//}

type NSAlert struct {
	objc.Object
}

func NewNSAlert() NSAlert {
	alert := NSAlert{nsAlert_.Alloc().Init()}
	//alert.SetDelegate(nsAlertController)
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

func (i NSAlert) SetAlertStyle(style NSAlertStyle) {
	i.Set("alertStyle:", core.NSUInteger(style))
}

func (i NSAlert) AlertStyle() NSAlertStyle {
	return NSAlertStyle(i.Get("alertStyle").Int())
}

func (i NSAlert) SetShowsHelp(show bool) {
	i.Set("showsHelp:", show)
}

func (i NSAlert) ShowsHelp() objc.Object {
	return i.Get("showsHelp")
}

func (i NSAlert) SetHelpAnchor(anchor string) {
	i.Set("helpAnchor:", core.String(anchor))
}

func (i NSAlert) HelpAnchor() string {
	return i.Get("helpAnchor").String()
}

func (i NSAlert) AddButtonWithTitle(s string) {
	i.Send("addButtonWithTitle:", core.String(s))
}

func (i NSAlert) Show() objc.Object {
	return i.Send("runModal")
}

func (i NSAlert) BeginSheetModalForWindow(win cocoa.NSWindow) objc.Object {
	callback, selector := core.Callback(func(resp objc.Object) {
		println(resp)
	})
	_ = callback
	_ = selector
	return i.Delegate().Send("beginSheetModalForWindow:", &win)
}

func NewNSAlert_WithSheetModal(win cocoa.NSWindow, callback func(resp objc.Object)) NSAlert {
	alert := NSAlert{nsAlert_.Alloc().Init()}
	//nsAlertControllerClass.AddMethod("showSheetModalForWindow:", func(a objc.Object) {
	//	if nsAlert, ok := a.(NSAlert); ok {
	//		//obj := nsAlert.BeginSheetModalForWindow(win)
	//		//callback(obj)
	//		fmt.Printf("%v\n", nsAlert)
	//	}
	//})
	//alert.SetDelegate(nsAlertController)
	alert.Object.Class().AddMethod("buttonAction:", func(r objc.Object) {
		println("buttonAction:::::", r)
	})
	return alert
}
