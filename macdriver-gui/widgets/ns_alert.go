package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSAlert struct {
	objc.Object
}

var NSAlert_ = objc.Get("NSAlert")

func NewNSAlert() NSAlert {
	return NSAlert{NSAlert_.Alloc().Init()}
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

func (i NSAlert) Show() objc.Object {
	return i.Send("runModal")
}

func (i NSAlert) BeginSheetModalForWindow(win cocoa.NSWindow, callback func(resp objc.Object)) objc.Object {
	cb, se := core.Callback(callback)
	_, _ = cb, se
	return i.Send("beginSheetModalForWindow:completionHandler:", &win, nil)
}
