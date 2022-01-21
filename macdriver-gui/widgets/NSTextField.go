package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTextField struct {
	cocoa.NSView `objc:"GoNSTextField : NSTextField"`
}

func performKeyEquivalent(obj, event objc.Object) bool {
	modifierFlags := event.Get("modifierFlags").Int()
	character := event.Get("charactersIgnoringModifiers")
	window := obj.Get("window")
	//responder := window.Get("firstResponder")

	// 点击 esc 取消焦点
	if modifierFlags&0xffff0000 == 0 && character.String() == "\x1b" {
		window.Send("makeFirstResponder:", nil)
		return true
	}

	if modifierFlags&0xffff0000 != 1<<20 {
		return obj.SendSuper("performKeyEquivalent:", event).Bool()
	}

	editor := window.Send("fieldEditor:forObject:", true, obj)
	var rtn objc.Object
	switch character.String() {
	case "a":
		//rtn = obj.Send("selectText:")
		rtn = editor.Send("selectAll:")
	case "c":
		rtn = editor.Send("copy:")
		//rtn = obj.SendSuper("sendAction:to:", objc.Sel("copy:"), responder)
	case "v":
		rtn = editor.Send("paste:")
		//rtn = obj.SendSuper("sendAction:to:", objc.Sel("paste:"), responder)
	case "x":
		rtn = editor.Send("cut:")
		//rtn = obj.SendSuper("sendAction:to:", objc.Sel("cut:"), responder)
		//case "z":
		//editor.Set("allowsUndo:", true)
		//golog.Error("按下 Command + z")
		//undoManager := window.Get("undoManager")
		//golog.Errorf("UndoManager: %v", undoManager)
		//undoManager.Send("undo:")
		//delegate := editor.Get("delegate")
		//golog.Errorf("Delegate: %v", delegate)
		//undoManager := delegate.Send("undoManager:for:", editor)
		//golog.Errorf("UndoManager: %v", undoManager)
		//rtn = editor.Send("undo:")
		//undoManager := responder.Get("undoManager")
		//golog.Errorf("UndoManger: %v", undoManager)
		//undoManager.Send("undo:", nil)
	}
	if rtn == nil {
		return true
	}
	return rtn.Bool()
}

var textFieldObj objc.Object

func init() {
	class := objc.NewClassFromStruct(NSTextField{})
	class.AddMethod("performKeyEquivalent:", performKeyEquivalent)
	objc.RegisterClass(class)
	textFieldObj = objc.Get("GoNSTextField")
}

func NewNSTextField(frame core.NSRect) NSTextField {
	return NSTextField{NSView: cocoa.NSView{}}
}

func (t NSTextField) SetBackgroundColor(color cocoa.NSColor) {
	t.Set("backgroundColor:", &color)
}

func (t NSTextField) SetBordered(isBordered bool) {
	t.Set("bordered:", isBordered)
}

func (t NSTextField) SetStringValue(val string) {
	t.Set("stringValue:", core.String(val))
}
