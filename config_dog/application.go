package main

import (
	"os"
	"os/exec"
)

func main() {
	_, err := getDogConfig()
	if err != nil {

	}
}

func runCmd(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	_ = command.Run()
}
