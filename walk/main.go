package main

import (
	"fmt"
	"github.com/lxn/walk"
	"os"
)

var settings *walk.IniFileSettings

func init() {
	settings = walk.NewIniFileSettings("hty_settings.ini")
	if err := settings.Load(); err != nil {
		walk.MsgBox(nil, "加载配置错误", fmt.Sprintf("%s", err.Error()), walk.MsgBoxIconError)
		os.Exit(1)
	}
}

func main() {
	// loginWindow, err := NewLoginWindow()
	// if err != nil {
	// 	return
	// }
	//
	// r := loginWindow.Run()
	//
	// if r != 1 {
	// 	return
	// }

	window := NewHomeWindow()
	window.Run()
}
