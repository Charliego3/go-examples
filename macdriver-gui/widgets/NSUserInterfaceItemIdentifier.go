package widgets

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSUserInterfaceItemIdentifier struct {
	objc.Object
}

var identifier = objc.Get("NSUserInterfaceItemIdentifier")

func NewNSUserInterfaceItemIdentifier(val string) NSUserInterfaceItemIdentifier {
	return NSUserInterfaceItemIdentifier{Object: identifier.Alloc().Send("init:", core.String(val))}
}
