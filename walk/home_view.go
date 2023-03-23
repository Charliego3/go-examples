//build:+windows

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
)

func NewHomeWindow() *walk.MainWindow {
	var window *walk.MainWindow
	err := MainWindow{
		AssignTo: &window,
		Title:    WindowTitle,
		MinSize:  Size{Width: windowWidth, Height: windowHeight},
		Size:     Size{Width: windowWidth, Height: windowHeight},
		Layout:   Grid{MarginsZero: true, Columns: 2},
		Children: []Widget{
			// HSplitter{
			// 	HandleWidth: 1,
			// 	Children: []Widget{
			// 		CreateOrderView(),
			// 		CreateMenuView(),
			// 	},
			// },

			CreateOrderView(),
			CreateMenuView(),
		},
	}.Create()

	if err != nil {
		walk.MsgBox(window, "打开程序错误", err.Error(), walk.MsgBoxIconError)
		os.Exit(1)
	}

	Center(window, windowWidth, windowHeight)
	return window
}

func CreateOrderView() Composite {
	return Composite{
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			PushButton{
				Text: "test",
				OnClicked: func() {
				},
			},
		},
	}
}

func CreateMenuView() Composite {
	// var tab *walk.TabPage
	// err := TabPage{
	// 	AssignTo: &tab,
	// 	Name:     "name",
	// 	Title:    "title",
	// 	Children: []Widget{
	// 		Label{Text: "label"},
	// 	},
	// }.Create(nil)
	// if err != nil {
	// 	return Composite{}
	// }

	return Composite{
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			TabWidget{
				// ContentMarginsZero: true,
				ContentMargins: Margins{Left: 400},
				MinSize:        Size{Width: 500, Height: 600},
				Pages: []TabPage{
					TabPage{
						Title:   "T1",
						Content: Label{Text: "TabPage1 Content"},
					},
					TabPage{
						Title:   "T2",
						Content: Label{Text: "TabPage2 Content"},
					},
				},
			},
		},
	}
}

func CreateChildren() []Widget {
	return []Widget{
		Composite{
			Layout: Grid{Columns: 1},
			Children: []Widget{
				Label{
					Text: "Left",
				},
			},
		},
		Composite{
			Layout: Grid{Columns: 1},
			Children: []Widget{
				PushButton{
					ColumnSpan: 2,
					Text:       "test",
					OnClicked: func() {
					},
				},
			},
		},
	}
}
