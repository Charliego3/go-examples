package main

import (
	"github.com/lxn/walk"
	"github.com/whimthen/temp/walk/sources/declarative"
)

func NewHomeWindow(assignTo *walk.MainWindow) error {
	return declarative.MainWindow{
		AssignTo: &assignTo,
		Title:    WindowTitle,
		Layout:   declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{Text: "HomeWindow"},
		},
	}.Create()
}
