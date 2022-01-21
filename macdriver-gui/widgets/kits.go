package widgets

import (
	"github.com/progrium/macdriver/core"
)

func NewStringNSArray(strs ...string) core.NSArray {
	objsInterface := make([]interface{}, len(strs))
	for i, str := range strs {
		objsInterface[i] = core.String(str)
	}
	return core.NSArray{}
}
