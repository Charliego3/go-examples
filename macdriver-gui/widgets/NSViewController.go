package widgets

import "github.com/progrium/macdriver/objc"

type NSViewController struct {
	objc.Object `objc:"SomeViewController : NSViewController"`
}

func (c NSViewController) ViewDidLoad() {
	println("NSViewController ViewDidLoad has executed....")
}

func NewNSViewController() objc.Object {
	class := objc.NewClassFromStruct(NSViewController{})
	class.AddMethod("viewDidLoad", (*NSViewController).ViewDidLoad)
	objc.RegisterClass(class)

	object := objc.Get("SomeViewController").Send("alloc").Send("init")
	return object
}
