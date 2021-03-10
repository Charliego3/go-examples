package main

import (
	"github.com/progrium/macdriver/bridge"
	"os"
	"time"
)

func main() {
	// start a bridge subprocess
	host, _ := bridge.NewHost(os.Stderr)
	go host.Run()

	// create a window
	window := bridge.Window{
		Title:       "My Title",
		Size:        bridge.Size{W: 480, H: 240},
		Position:    bridge.Point{X: 200, Y: 200},
		Closable:    true,
		Minimizable: false,
		Resizable:   false,
		Borderless:  false,
		AlwaysOnTop: true,
		Background:  &bridge.Color{R: 1, G: 1, B: 1, A: 0.5},
	}
	bridge.Sync(host.Peer, &window)

	// change its title
	window.Title = "My New Title"
	bridge.Sync(host.Peer, &window)

	time.Sleep(time.Hour)
	// destroy the window
	bridge.Release(host.Peer, &window)
}
