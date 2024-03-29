package table

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

var nsTableView = objc.Get("NSTableView")

type NSTableView struct {
	cocoa.NSView

	coder      NSCoder
	dataSource NSTableViewDataSource
	delegate   NSTableViewDelegate
}

func NewNSTableView(frame core.NSRect) NSTableView {
	tableView := NSTableView{
		NSView: cocoa.NSView{},
		coder:  NewNSCoder(),
	}
	dataSource := NewNSTableViewDataSource()
	delegate := NewNSTableViewDelegate()
	tableView.dataSource = dataSource
	tableView.delegate = delegate
	tableView.Set("dataSource:", dataSource.Object)
	tableView.Set("delegate:", delegate.Object)
	return tableView
}

func NewNSTableViewWithCoder() NSTableView {
	coder := NewNSCoder()
	return NSTableView{
		NSView: cocoa.NSView{},
		coder:  coder,
	}
}

func (t NSTableView) AddTableColumn(columns ...NSTableColumn) {
	for _, column := range columns {
		t.Send("addTableColumn:", column)
	}
}

func (t NSTableView) SetSelectionHighlightStyle(style NSTableViewSelectionHighlightStyle) {
	t.Set("selectionHighlightStyle:", int(style))
}

func (t NSTableView) SetRowHeight(height float64) {
	t.Set("rowHeight:", height)
}

func (t NSTableView) SetRowSizeStyle(style NSTableViewRowSizeStyle) {
	t.Set("rowSizeStyle:", int(style))
}

func (t NSTableView) SetStyle(style NSTableViewStyle) {
	t.Set("style:", style)
}

func (t NSTableView) SetGridStyleMask(style NSTableViewGridLineStyle) {
	t.Set("gridStyleMask:", style)
}

func (t NSTableView) SetGridColor(color cocoa.NSColor) {
	t.Set("gridColor:", color)
}
