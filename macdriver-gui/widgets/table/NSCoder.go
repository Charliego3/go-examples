package table

import "github.com/progrium/macdriver/objc"

type NSCoder struct {
	objc.Object
}

var nsCoder = objc.Get("NSCoder")

func NewNSCoder() NSCoder {
	nsCoder.Class().AddMethod("decodeValueOfObjType:at:size:", func(object objc.Object) {

	})
	return NSCoder{nsCoder.Alloc().Init()}
}
