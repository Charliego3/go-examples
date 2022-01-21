package widgets

import (
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

type NSUserNotification struct {
	objc.Object
}

var NSUserNotification_ = objc.Get("NSUserNotification")

type NSUserNotificationCenter struct {
	objc.Object
}

var NSUserNotificationCenter_ = objc.Get("NSUserNotificationCenter")

func NewNotificationCenter() NSUserNotification {
	return NSUserNotification{NSUserNotification_.Alloc().Init()}
}

func (n NSUserNotification) SetTitle(title string) {
	n.Set("title:", core.String(title))
}

func (n NSUserNotification) SetInformativeText(text string) {
	n.Set("informativeText:", core.String(text))
}

func (n NSUserNotification) Show() {
	center := NSUserNotificationCenter{NSUserNotificationCenter_.Send("defaultUserNotificationCenter")}
	center.Send("deliverNotification:", n)
	n.Release()
}
