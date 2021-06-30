package table

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/type_alias"
)

type NSTableColumn struct {
	objc.Object
}

var nsTableColumn = objc.Get("NSTableColumn")

func NewNSTableColumn(identifier string) NSTableColumn {
	return NSTableColumn{Object: nsTableColumn.Alloc().Send("initWithIdentifier:", type_alias.NewNSUserInterfaceItemIdentifier(identifier))}
}

func (c NSTableColumn) SetTitle(title string) {
	c.Set("title:", core.String(title))
}

func (c NSTableColumn) SetHeaderCell(header NSTableHeaderCell) {
	c.Set("headerCell:", header)
}
