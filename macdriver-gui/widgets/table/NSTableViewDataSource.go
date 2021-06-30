package table

import (
	"fmt"
	"github.com/progrium/macdriver/objc"
)

type NSTableViewDataSource struct {
	objc.Object
}

var defaultDataSource objc.Object

func init() {
	class := objc.NewClass("DefaultNSTableViewDataSource", "NSObject")
	class.AddMethod("tableView:objectValueForTableColumn:row:", func(dataSource, table, column objc.Object, row int) objc.Object {
		return NewNSCell(fmt.Sprintf("Row-%v, Column-%v", row, column.Uint()))
	})
	class.AddMethod("numberOfRowsInTableView:", func(object objc.Object) int {
		return 50
	})
	objc.RegisterClass(class)
	defaultDataSource = objc.Get("DefaultNSTableViewDataSource").Alloc().Init()
}

func NewNSTableViewDataSource() NSTableViewDataSource {
	return NSTableViewDataSource{Object: defaultDataSource}
}
