package main

import (
	"fmt"
	"os/exec"

	"github.com/caseymrm/menuet"
)

const terminalKey = "com.tools.terminal"

type sshItem struct {
	Alias   string `json:"alias"`
	Command string `json:"command"`
}

func startCommand(command string) {
	c := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s"`, command))
	if err := c.Run(); err != nil {
		notify("Oops!", "open terminal error: %v", err)
	}
}

func openTerminal() menuet.MenuItem {
	var cmds []sshItem
	menuet.Defaults().Unmarshal(terminalKey, &cmds)

	var items []menuet.MenuItem
	for _, cmd := range cmds {
		cmd := cmd
		items = append(items, menuet.MenuItem{
			Text: cmd.Alias,
			Clicked: func() {
				startCommand(cmd.Command)
			},
		})
	}

	var filter = func(check func(k, v string) bool) bool {
		for _, cmd := range cmds {
			if check(cmd.Alias, cmd.Command) {
				return true
			}
		}
		return false
	}

	items = append(items, []menuet.MenuItem{
		{Type: menuet.Separator},
		{
			Text:       "Add new command",
			FontWeight: menuet.WeightBold,
			Clicked: func() {
				ret := menuet.App().Alert(menuet.Alert{
					MessageText: "Please enter the alias and command",
					Inputs:      []string{"Alias", "Command"},
					Buttons:     []string{"OK", "Cancel"},
				})

				if ret.Button != 0 {
					return
				}

				if ret.Inputs[0] == "" || ret.Inputs[1] == "" {
					warning("输入不正确, 请重新操作")
					return
				}

				if filter(func(k, v string) bool {
					return ret.Inputs[0] == k || ret.Inputs[1] == v
				}) {
					warning("该Command已经存在: %s", ret.Inputs[1])
					return
				}

				cmds = append(cmds, sshItem{
					Alias:   ret.Inputs[0],
					Command: ret.Inputs[1],
				})
				menuet.Defaults().Marshal(terminalKey, cmds)
				startCommand(ret.Inputs[1])
			},
		},
	}...)

	return menuet.MenuItem{
		Text:    "Terminal command",
		Clicked: nil,
		Children: func() []menuet.MenuItem {
			return items
		},
	}
}
