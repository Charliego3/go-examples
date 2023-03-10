package table

import (
	"fmt"

	"github.com/kataras/golog"
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTableRowView struct {
	cocoa.NSView
}

var tableRowView objc.Object

func NewNSTableRowView() NSTableRowView {
	if tableRowView == nil {
		registerClass()
	}
	return NSTableRowView{NSView: cocoa.NSView{}}
}

func registerClass() {
	class := objc.NewClass("DefaultNSTableRowView", "NSTableRowView")
	class.AddMethod("viewAtColumn:", func(row, column objc.Object) objc.Object {
		golog.Errorf("Row:(%+v), Column:(%+v)", row, column)
		text := cocoa.NSView{}
		text.Set("string:", core.String(fmt.Sprintf("Row:(%+v), Column:(%+v)", row, column)))
		return text
	})
	objc.RegisterClass(class)
	tableRowView = objc.Get("DefaultNSTableRowView").Alloc().Init()
}
