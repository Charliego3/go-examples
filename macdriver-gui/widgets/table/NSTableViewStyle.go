package table

import "github.com/progrium/macdriver/core"

type NSTableViewStyle core.NSUInteger

const (
	NSTableViewStyleAutomatic NSTableViewStyle = iota
	NSTableViewStyleFullWidth
	NSTableViewStyleInset
	NSTableViewStyleSourceList
	NSTableViewStylePlain
)
