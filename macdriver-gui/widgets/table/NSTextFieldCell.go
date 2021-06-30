package table

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTextFieldCell struct {
	objc.Object
}

var nsTextFieldCell = objc.Get("NSTextFieldCell")

func NewNSTextFieldCell(text string) NSTextFieldCell {
	field := nsTextFieldCell.Alloc().Send("initTextCell:", core.String(text))
	field.Set("bezelStyle:", core.NSUInteger(1))
	return NSTextFieldCell{Object: field}
}
