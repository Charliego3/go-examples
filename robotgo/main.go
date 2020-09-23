package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/whimthen/kits/logger"
	"time"
)

func main() {
	//for {
	//	x, y := robotgo.GetMousePos()
	//	fmt.Println("pos:", x, y)
	//	time.Sleep(time.Millisecond * 500)
	//}

	width, height := robotgo.GetScreenSize()
	logger.Debug("ScreenSize: width = %d, height = %d", width, height)
	for {
		robotgo.MoveMouseSmooth(0, 0)
		robotgo.MoveMouseSmooth(width, height)
		robotgo.MoveMouseSmooth(width, 0)
		robotgo.MoveMouseSmooth(0, height)
		time.Sleep(time.Second * 5)
	}
}
