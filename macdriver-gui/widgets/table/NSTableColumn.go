package table

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTableColumn struct {
	objc.Object
}

var nsTableColumn = objc.Get("NSTableColumn")

func NewNSTableColumn() NSTableColumn {
	return NSTableColumn{Object: nsTableColumn.Alloc().Init()}
}

func (c NSTableColumn) SetTitle(title string) {
	c.Set("title:", core.String(title))
}

func (c NSTableColumn) SetHeaderCell(header NSTableHeaderCell) {
	c.Set("headerCell:", header)
}
