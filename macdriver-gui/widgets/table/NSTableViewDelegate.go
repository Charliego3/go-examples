package table

import (
	"fmt"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTableViewDelegate struct {
	objc.Object
}

var nsTableViewDelegate objc.Object

func NewNSTableViewDelegate() NSTableViewDelegate {
	if nsTableViewDelegate == nil {
		lazyLoadDefaultDelegate()
	}
	return NSTableViewDelegate{Object: nsTableViewDelegate}
}

func lazyLoadDefaultDelegate() {
	class := objc.NewClass("DefaultNSTableViewDelegate", "NSObject")
	class.AddMethod("tableView:viewForTableColumn:row:", func(dataSource, table, column objc.Object, row int) objc.Object {
		text := cocoa.NSView{}
		identifier := column.Get("identifier").String()
		if identifier == "Column1" {
			//return NewNSCell(fmt.Sprintf("Delegate: Row-%d, Column1", row))
			text.Set("string:", core.String(fmt.Sprintf("Delegate: Row-%d, Column1", row)))
		} else if identifier == "Column2" {
			//return NewNSCell(fmt.Sprintf("Delegate: Row-%d, Column2", row))
			text.Set("string:", core.String(fmt.Sprintf("Delegate: Row-%d, Column2", row)))
		} else if identifier == "Number" {
			//return NewNSCell(fmt.Sprintf("Delegate: Row-%d, Number", row))
			text.Set("string:", core.String(fmt.Sprintf("Delegate: Row-%d, Number", row)))
		} else {
			//return NewNSCell(fmt.Sprintf("Delegate: Row-%d, Nobody's Here!", row))
			text.Set("string:", core.String(fmt.Sprintf("Delegate: Row-%d, Nobody's Here!", row)))
		}
		return text
	})
	//class.AddMethod("tableView:rowViewForRow:", func(delegate, table objc.Object, row int) objc.Object {
	//	golog.Errorf("Delegate:%v, Table:%v, Row:%v", delegate, table, row)
	//	return NewNSTableRowView().NSView
	//})
	objc.RegisterClass(class)
	nsTableViewDelegate = objc.Get("DefaultNSTableViewDelegate").Alloc().Init()
}
