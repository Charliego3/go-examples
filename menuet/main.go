package main

import (
	"log"
	"time"

	"github.com/caseymrm/menuet"
)

func helloClock() {
	for {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "Hello World " + time.Now().Format(":05"),
		})
		time.Sleep(time.Second)
	}
}

func main() {
	//go helloClock()
	//menuet.App().RunApplication()

	go func() {
		alert := menuet.Alert{
			MessageText:     "This is a Alert",
			InformativeText: "This is InformativeText",
			Buttons:         []string{"Btn1", "Btn2"},
			Inputs:          []string{"Input1", "Input2"},
		}
		clicked := menuet.App().Alert(alert)

		log.Printf("AlertCliecked: %#v", clicked)
	}()

	menuet.App().RunApplication()
}
