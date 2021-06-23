package statusBar

import (
	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
	"github.com/whimthen/temp/macdriver-gui/widgets/alert"
	"testing"
)

func TestNewStatusBarApp(t *testing.T) {
	app := NewStatusBarApp("测试", cocoa.NSVariableStatusItemLength)
	app.AddMenuItem("item1", func(_ objc.Object) {
		nsAlert := alert.NewNSAlert()
		nsAlert.SetAlertStyle(alert.Informational)
		nsAlert.SetShowsHelp(true)
		nsAlert.Set("helpAnchor:", core.String("www.baidu.com"))
		nsAlert.SetMessageText("Alert message")
		nsAlert.SetInformativeText("Detailed description of nsAlert message")
		nsAlert.AddButtonWithTitle("Default")
		nsAlert.AddButtonWithTitle("Alternative")
		nsAlert.AddButtonWithTitle("Other")
		nsAlert.Show()
	})
	app.AddItemSeparator()
	app.AddTerminateItem()

	app.Run()
}
