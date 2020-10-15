package main

import (
	"bytes"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func CompleteNginx() {
	nginxPath := checkNginxInstall()
	confPath := getNginxConfigPath(nginxPath)
	configurationNginx(confPath)
}

func configurationNginx(configPath string) {

}

func getNginxConfigPath(nginx string) string {
	command := exec.Command("bash", "-c", nginx+" -V")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		color.Red(err.Error())
		return ""
	}
	version := stderr.String()
	pathIndex := strings.Index(version, "--conf-path=")
	configPath := version[pathIndex+12 : pathIndex+strings.Index(version[pathIndex:], " ")]
	return configPath
}

func checkNginxInstall() string {
	path, err := checkExists("nginx")
	if err != nil {
		isInstall := false
		prompt := &survey.Confirm{
			Message: "Do you want to install nginx to start the project?",
			Default: true,
		}
		_ = survey.AskOne(prompt, &isInstall)

		if isInstall {
			brewPath := checkBrewInstall()
			if brewPath == "" {
				installNginxFromCode()
			} else {
				installNginxFromBrew(brewPath)
			}
		}
	}
	return path
}

func checkBrewInstall() string {
	path, err := checkExists("brew")
	if err != nil {
		isInstall := false
		prompt := &survey.Confirm{
			Message: "Do you want to install brew to install nginx?",
			Default: true,
		}
		_ = survey.AskOne(prompt, &isInstall)

		if isInstall {
			installBrew()
			path, _ = checkExists("brew")
		}
	}
	return path
}

func installNginxFromCode() {

}

func installNginxFromBrew(brew string) {
	command := exec.Command("bash", "-c", "brew install nginx")
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Start()
	if err != nil {
		color.Red(err.Error())
		return
	}
	if err = command.Wait(); err != nil {
		color.Red(err.Error())
	}
}

func installBrew() {
	command := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh")

	output, err := command.Output()
	if err != nil {
		color.Red(err.Error())
		return
	}

	fileName := "./.brew.sh"
	err = ioutil.WriteFile(fileName, output, 0644)
	if err != nil {
		color.Red(err.Error())
		return
	}
	defer os.Remove(fileName)

	err = os.Chmod(fileName, 0777)
	if err != nil {
		color.Red(err.Error())
		return
	}

	// 执行脚本文件
	command = exec.Command(fileName)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	err = command.Start()
	if err != nil {
		color.Red(err.Error())
		return
	}
	if err = command.Wait(); err != nil {
		color.Red(err.Error())
	}
}

func checkExists(file string) (path string, err error) {
	path, err = exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return path, nil
}
