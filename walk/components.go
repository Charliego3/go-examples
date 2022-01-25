package main

import (
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func NewFont(size int, bold ...bool) declarative.Font {
	var b bool
	if len(bold) > 0 {
		b = bold[0]
	}
	return declarative.Font{
		Family:    "楷体",
		PointSize: size,
		Bold:      b,
	}
}

func GetWinScreen() (width, height int32) {
	width = win.GetSystemMetrics(win.SM_CXSCREEN)
	height = win.GetSystemMetrics(win.SM_CYSCREEN)
	return
}
