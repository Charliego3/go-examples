package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTextField struct {
	objc.Object
}

var NSTextField_ = objc.Get("NSTextField")

func NSTextField_Init(frame core.NSRect) NSTextField {
	return NSTextField{NSTextField_.Alloc().Send("initWithFrame:", &frame)}
}

func (t NSTextField) SetBackgroundColor(color cocoa.NSColor) {
	t.Set("backgroundColor:", &color)
}

func (t NSTextField) SetStringValue(val string) {
	t.Set("stringValue:", core.String(val))
}

//func (t NSTextField) BackgroundColor() cocoa.NSColor {
//	return t.Get("backgroundColor:")
//}
