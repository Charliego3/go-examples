package table

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

var (
	nsTableView_ = objc.Get("NSTableView")

	tableViewDelegateClass objc.Class
	tableViewDelegate      objc.Object
)

func init() {
	tableViewDelegateClass = objc.NewClass("NSTableViewDelegate", "NSObject")
	tableViewDelegateClass.AddMethod("", func(obj objc.Object) {

	})

}

type NSTableView struct {
	cocoa.NSView
}

func NewNSTableView(frame core.NSRect) NSTableView {
	return NSTableView{cocoa.NSView{nsTableView_.Alloc().Send("initWithFrame:", frame)}}
}
