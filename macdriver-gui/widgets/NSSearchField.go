package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSSearchField struct {
	cocoa.NSView
}

var nsSearchField = objc.Get("NSSearchField")

func NewNSSearchField(frame core.NSRect) NSSearchField {
	return NSSearchField{NSView: cocoa.NSView{}}
}

func (s NSSearchField) SetRecentSearches(recent ...string) {
	s.Set("recentSearches:", NewStringNSArray(recent...))
}
