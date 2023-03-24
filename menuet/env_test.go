package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestOpen(t *testing.T) {
	// cmd := exec.Command("neovide")
	// cmd.Dir = "/Users/charlie/dev/java/zb_main"
	// cmd.Start()

	os.Chdir("/Users/charlie/dev/java/zb_main/")
	path, _ := exec.LookPath("neovide")
	fmt.Println(path)
	// exec.Command(path).Start()
}

func TestTerminal(t *testing.T) {
	c := exec.Command("osascript", "-e", `tell application "Terminal" to do script "ll -a"`)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
