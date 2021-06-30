package table

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTableHeaderCell struct {
	objc.Object
}

var nsTableHeaderCell = objc.Get("NSTableHeaderCell")

func NewNSTableHeaderCell(text string) NSTableHeaderCell {
	return NSTableHeaderCell{Object: nsTableHeaderCell.Alloc().Send("initTextCell:", core.String(text))}
}
