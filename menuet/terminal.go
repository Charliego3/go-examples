package main

import (
	"fmt"
	"os/exec"

	"github.com/caseymrm/menuet"
)

func startCommand(command string) {
	c := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s"`, command))
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}

func openTerminal() menuet.MenuItem {
	return menuet.MenuItem{
		Text:    "Open Remote SSH",
		Clicked: nil,
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				{
					Text: "测试环境(ttjms)",
					Clicked: func() {
						startCommand("ttjms")
					},
				},
				{
					Text: "量化环境(jms4)",
					Clicked: func() {
						startCommand("jms")
					},
				},
			}
		},
	}
}
