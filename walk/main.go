package main

import (
	"strconv"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func main() {
	loginWindow, err := NewLoginWindow()
	if err != nil {
		return
	}

	r := loginWindow.Run()

	if r != 1 {
		return
	}

	walk.MsgBox(loginWindow, "fanhuizhi", strconv.Itoa(r), walk.MsgBoxIconInformation)

	var window *walk.MainWindow

	declarative.MainWindow{
		AssignTo: &window,
		Title:    WindowTitle,
		MinSize:  declarative.Size{Width: 800, Height: 600},
		Size:     declarative.Size{Width: 800, Height: 600},
		Layout:   declarative.Grid{MarginsZero: true, Columns: 2},
		Children: []declarative.Widget{
			declarative.ImageView{
				Background: declarative.SolidColorBrush{Color: walk.RGB(255, 255, 255)},
				Image:      "resources/logo.jpg",
				Mode:       declarative.ImageViewModeCenter,
				MinSize:    declarative.Size{Width: 300},
			},

			declarative.Label{Text: "Home Window"},
		},
	}.Run()

	// win.SetWindowPos(loginWindow.Handle(), win.HWND_DESKTOP, -1, -1, -1, -1, win.SWP_NOMOVE|win.SWP_NOREPOSITION|win.SWP_NOSIZE)
	// loginWindow.Run()
}
