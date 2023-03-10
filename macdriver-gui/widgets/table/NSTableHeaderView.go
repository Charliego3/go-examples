package table

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/objc"
)

type NSTableHeaderView struct {
	cocoa.NSView
}

var nsTableHeaderView = objc.Get("NSTableHeaderViews")

func NewNSTableHeaderView(table NSTableView) NSTableHeaderView {
	header := nsTableHeaderView.Alloc().Init()
	header.Set("tableView:", &table)
	return NSTableHeaderView{cocoa.NSView{}}
}
