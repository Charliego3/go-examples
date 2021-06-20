package widgets

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSOpenPanel struct {
	cocoa.NSView
}

var nsOpenPanel_ = objc.Get("NSOpenPanel")

func NewNSOpenPanel() NSOpenPanel {
	return NSOpenPanel{NSView: cocoa.NSView{Object: nsOpenPanel_.Alloc().Init()}}
}

func (p NSOpenPanel) SetMessage(message string) {
	p.Set("message:", core.String(message))
}

func (p NSOpenPanel) SetDirectoryURL(url string) {
	p.Set("directoryURL:", core.URL(url))
}

func (p NSOpenPanel) HomeDirectory() string {
	return `\(NSHomeDirectory())/Downloads`
}

func (p NSOpenPanel) SetAllowsMultipleSelection(allow bool) {
	p.Set("allowsMultipleSelection:", allow)
}

func (p NSOpenPanel) SetCanChooseDirectories(choose bool) {
	p.Set("canChooseDirectories:", choose)
}

func (p NSOpenPanel) SetCanCreateDirectories(create bool) {
	p.Set("canCreateDirectories:", create)
}

func (p NSOpenPanel) SetCanChooseFiles(choose bool) {
	p.Set("canChooseFiles:", choose)
}

func (p NSOpenPanel) SetAllowedFileTypes(choose ...string) {
	p.Set("canChooseFiles:", core.String(choose[0]))
}

func (p NSOpenPanel) GetPath() string {
	object := p.Get("path")
	println(object)
	return ""
}

func (p NSOpenPanel) Show() objc.Object {
	return p.Send("runModal")
}

func (p NSOpenPanel) RunModalForDirectory() objc.Object {
	return p.Send("runModalForDirectory:file:", nil, nil)
}
