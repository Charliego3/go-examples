package type_alias

import (
	"github.com/progrium/macdriver/core"
)

type NSUserInterfaceItemIdentifier struct {
	core.NSString
}

func NewNSUserInterfaceItemIdentifier(val string) NSUserInterfaceItemIdentifier {
	return NSUserInterfaceItemIdentifier{NSString: core.String(val)}
}
