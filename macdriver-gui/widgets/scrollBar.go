package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSScroller struct {
	cocoa.NSView
}

func NewNSScroller() NSScroller {
	return NSScroller{cocoa.NSView{objc.Get("NSScroller").Alloc().Init()}}
}

func (s NSScroller) SetScrollerStyle(style NSScrollerStyle) {
	s.Set("scrollerStyle:", core.NSUInteger(style))
}
