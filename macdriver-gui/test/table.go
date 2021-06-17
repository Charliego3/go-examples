package main

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSTableView struct {
	objc.Object
}

var NSTableView_ = objc.Get("NSTableView")

func NSTableView_Init(frame core.NSRect) NSTableView {
	return NSTableView{cocoa.NSView{Object: NSTableView_.Alloc().Send("initWithFrame:", frame)}}
}

func NSTableView_New() NSTableView {
	return NSTableView{NSTableView_.Alloc().Init()}
}
