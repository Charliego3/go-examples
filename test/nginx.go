package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	nginxConfig = `server {
		listen                     80;
		server_name                ttvip2.100-%s.net;
		location / {
			proxy_pass_header         Server;
			proxy_set_header          Host $http_host;
			proxy_set_header          X-Real-IP $remote_addr;
			proxy_pass                http://127.0.0.1:%s/;
		}
	}`
)

func getIncludeExists() {
	var re = regexp.MustCompile(`(?m).*include\s+/usr/local/etc/nginx/servers/\*;`)
	content, err := ioutil.ReadFile("/usr/local/etc/nginx/nginx.conf")
	if err != nil {
		color.Red(err.Error())
		return
	}

	submatch := re.FindAllStringSubmatch(string(content), -1)
	println(submatch)
}

func CompleteNginx(ip string, node string) {
	nginxPath := checkNginxInstall()
	confPath := getNginxConfigPath(nginxPath)
	configurationNginx(confPath, ip, node)
}

func configurationNginx(configPath, ip, node string) {
	base := filepath.Base(configPath)
	color.Cyan("Nginx config base path: %s", base)
	serversDir := filepath.Join(base, "servers")
	color.Cyan("Nginx config servers path: %s", serversDir)
	stat, err := os.Stat(serversDir)
	if err != nil {
		color.Red("🌡  %+v", err)
		return
	}

	if !stat.IsDir() {
		err := os.Mkdir(serversDir, 0644)
		if err != nil {
			color.Red("🌡  %+v", err)
			return
		}
	}

	confName := fmt.Sprintf("zb_%s_%s.conf", node, strings.ReplaceAll(ip, ".", "_"))
	err = ioutil.WriteFile(confName, []byte(fmt.Sprintf(nginxConfig, "", "")), 0644)
	if err != nil {
		color.Red("🌡  %+v", err)
	}
}

func getNginxConfigPath(nginx string) string {
	command := exec.Command("bash", "-c", nginx+" -V")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		color.Red("🌡  %+v", err)
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
		color.Red("🌡  %+v", err)
		return
	}
	if err = command.Wait(); err != nil {
		color.Red("🌡  %+v", err)
	}
}

func installBrew() {
	command := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh")

	output, err := command.Output()
	if err != nil {
		color.Red("🌡  %+v", err)
		return
	}

	fileName := "./.brew.sh"
	err = ioutil.WriteFile(fileName, output, 0644)
	if err != nil {
		color.Red("🌡  %+v", err)
		return
	}
	defer os.Remove(fileName)

	err = os.Chmod(fileName, 0777)
	if err != nil {
		color.Red("🌡  %+v", err)
		return
	}

	// 执行脚本文件
	command = exec.Command(fileName)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	err = command.Start()
	if err != nil {
		color.Red("🌡  %+v", err)
		return
	}
	if err = command.Wait(); err != nil {
		color.Red("🌡  %+v", err)
	}
}

func checkExists(file string) (path string, err error) {
	path, err = exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return path, nil
}
