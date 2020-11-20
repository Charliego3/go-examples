package main

import (
	"github.com/fatih/color"
	"github.com/sqweek/dialog"
	"os"
	"os/exec"
)

func main() {
	//directory, err := dialog.Directory().Title("Load images").Browse()
	//color.Yellow("You chose directory is: %s, err: %+v", directory, err)
	_, err := getDogConfig()
	if err != nil {

	}

	yesNo := dialog.Message("%s", "Do you want to add deamon?").Title("Are you sure?").YesNo()
	color.Red("YesNo: %v", yesNo)
}

func runCmd(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	_ = command.Run()
}
