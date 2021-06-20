package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTextField struct {
	cocoa.NSView
}

var NSTextField_ = objc.Get("NSTextField")

func NewNSTextField(frame core.NSRect) NSTextField {
	return NSTextField{NSView: cocoa.NSView{Object: NSTextField_.Alloc().Send("initWithFrame:", frame)}}
}

func (t NSTextField) SetBackgroundColor(color cocoa.NSColor) {
	t.Set("backgroundColor:", &color)
}

func (t NSTextField) SetIsBordered(isBordered bool) {
	t.Set("bordered:", isBordered)
}

func (t NSTextField) SetStringValue(val string) {
	t.Set("stringValue:", core.String(val))
}

//func (t NSTextField) BackgroundColor() cocoa.NSColor {
//	return t.Get("backgroundColor:")
//}
