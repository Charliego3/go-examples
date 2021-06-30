package table

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSCell struct {
	objc.Object
}

var nsCell = objc.Get("NSCell")

func NewNSCell(text string) NSCell {
	cell := nsCell.Alloc().Send("initTextCell:", core.String(text))
	cell.Set("type:", core.NSUInteger(1))
	cell.Set("controlSize:", core.NSUInteger(1))
	cell.Send("calcDrawInfo:", core.Rect(0, 0, 100, 40))
	return NSCell{Object: cell}
}
