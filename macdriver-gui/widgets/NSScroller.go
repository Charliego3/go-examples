package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
)

type NSScroller struct {
	cocoa.NSView
}

func NewNSScroller() NSScroller {
	//scroll := objc.Get("NSScroller").Alloc().Init()
	//scroll.Send("scrollerWidthForControlSize:scrollerStyle:",
	//	core.NSUInteger(NSControlSizeSmall), core.NSUInteger(NSScrollerStyleOverlay))
	return NSScroller{cocoa.NSView{}}
}

func (s NSScroller) SetScrollerStyle(style NSScrollerStyle) {
	s.Set("scrollerStyle:", core.NSUInteger(style))
}
