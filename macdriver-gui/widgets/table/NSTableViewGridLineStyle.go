package table

import "github.com/progrium/macdriver/core"

type NSTableViewGridLineStyle core.NSUInteger

const (
	NSTableViewGridNone                     NSTableViewGridLineStyle = 0
	NSTableViewSolidVerticalGridLineMask                             = 1 << 0
	NSTableViewSolidHorizontalGridLineMask                           = 1 << 1
	NSTableViewDashedHorizontalGridLineMask                          = 1 << 3
)
