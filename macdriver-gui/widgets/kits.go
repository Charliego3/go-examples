package widgets

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

func NewStringNSArray(strs ...string) core.NSArray {
	objsInterface := make([]interface{}, len(strs))
	for i, str := range strs {
		objsInterface[i] = core.String(str)
	}
	return core.NSArray{Object: objc.Get("NSArray").Send("arrayWithObjects:", objsInterface...)}
}
