package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSClipView struct {
	cocoa.NSView
}

func NewNSClipView() NSClipView {
	return NSClipView{cocoa.NSView{}}
}

func (c NSClipView) ScrollTo(point core.NSPoint) {
	c.Send("scrollTo:", point)
}

func (c NSClipView) SetDocumentView(view objc.Object) {
	c.Set("documentView:", view)
}
