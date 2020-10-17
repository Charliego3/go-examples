package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	nginxConfig = `server {
    listen                     80;
    server_name                %s;
    location / {
        proxy_pass_header      Server;
        proxy_set_header       Host $http_host;
        proxy_set_header       X-Real-IP $remote_addr;
        proxy_pass             %s;
    }
}`
)

func CompleteNginx(s *Sync, node, ip, port string) {
	nginxPath := checkNginxInstall()
	if nginxPath != "" {
		spinner.Restart()
		confPath := getNginxConfigPath(nginxPath)
		confPath = filepath.Dir(confPath)
		configurationNginx(s, confPath, node, ip, port)
		startNginx(confPath)
	}
}

func startNginx(confPath string) {
	_, err := os.Stat(filepath.Join(confPath, "nginx.pid"))
	if err != nil {
		color.Cyan("🍺 Will be start nginx: May be need enter user password")
		err := exec.Command("bash", "-c", "sudo nginx").Run()
		if err != nil {
			pe(err)
			return
		}
	} else {
		color.Cyan("🍺 Will be reload nginx conf: May be need enter user password")
		_ = exec.Command("bash", "-c", "sudo nginx -s reload").Run()
	}
}

func configurationNginx(s *Sync, confPath, node, ip, port string) {
	serversDir := filepath.Join(confPath, "servers")
	stat, err := os.Stat(serversDir)
	if err != nil {
		pe(err)
		return
	}

	if !stat.IsDir() {
		err := os.Mkdir(serversDir, 0644)
		if err != nil {
			pe(err)
			return
		}
	}

	i := ip[strings.LastIndex(ip, ".")+1:]
	confName := fmt.Sprintf("zb_%s_%s.conf", node, strings.ReplaceAll(ip, ".", "_"))
	serverName := fmt.Sprintf("tt%s2.100-%s.net", node, i)
	proxyPass := fmt.Sprintf("http://127.0.0.1:%s/", port)
	confFile := filepath.Join(serversDir, confName)
	err = ioutil.WriteFile(confFile, []byte(fmt.Sprintf(nginxConfig, serverName, proxyPass)), 0644)
	if err != nil {
		pe(err)
		return
	}

	configurationInclude(s, confPath)
	spinner.Stop()
	println(color.GreenString("🍺 Configuration nginx conf"), color.YellowString("(`%s`)", confFile), color.BlueString("server_name: `%s`, proxy_pass: `%s`", serverName, proxyPass))
}

func configurationInclude(s *Sync, confPath string) {
	var re = regexp.MustCompile(fmt.Sprintf(`(?m).*include\s+%s/servers/\*;`, confPath))
	nginxConf := filepath.Join(confPath, "nginx.conf")
	content, err := ioutil.ReadFile(nginxConf)
	if err != nil {
		pe(err)
		return
	}

	includes := re.FindAllStringSubmatch(string(content), -1)
	var hasInclude bool
	if includes != nil && len(includes) > 0 {
		for _, include := range includes {
			if !strings.HasPrefix(strings.TrimSpace(include[0]), "#") {
				hasInclude = true
				break
			}
		}
	}

	if !hasInclude {
		strings.Index(string(content), "http {")
		var re = regexp.MustCompile(`(?m).*http\s*{`)
		subIncludes := re.FindAllStringSubmatch(string(content), -1)
		subIncludeIndex := re.FindAllStringSubmatchIndex(string(content), -1)
		var index int
		if subIncludes != nil && len(subIncludes) > 0 {
			for i, subInclude := range subIncludes {
				if !strings.HasPrefix(strings.TrimSpace(subInclude[0]), "#") {
					index = i
					break
				}
			}
		}
		includeIndex := subIncludeIndex[index][1]
		newContent := string(content[:includeIndex])
		newContent += fmt.Sprintf("\n\tinclude\t%s/servers/*;\n", confPath)
		newContent += string(content[includeIndex:])
		err := ioutil.WriteFile(nginxConf, []byte(newContent), 0644)
		if err != nil {
			pe(err)
		}
	}
}

func getNginxConfigPath(nginx string) string {
	command := exec.Command("bash", "-c", nginx+" -V")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		pe(err)
		return ""
	}
	version := stderr.String()
	pathIndex := strings.Index(version, "--conf-path=")
	configPath := version[pathIndex+12 : pathIndex+strings.Index(version[pathIndex:], " ")]
	return configPath
}

func checkNginxInstall() string {
	spinner.Restart()
	path, err := checkExists("nginx")
	spinner.Stop()
	if err != nil {
		isInstall := false
		prompt := &survey.Confirm{
			Message: "Do you want to install nginx to start the project?",
			Default: true,
		}
		err = survey.AskOne(prompt, &isInstall)
		if err != nil {
			spinner.Stop()
			return ""
		}

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
		err = survey.AskOne(prompt, &isInstall)
		if err != nil {
			spinner.Stop()
			return ""
		}

		if isInstall {
			installBrew()
			path, _ = checkExists("brew")
		}
	}
	return path
}

func installNginxFromCode() {
	color.Cyan("🍺 About to install nginx from source")
	const nginxDomain = "https://nginx.org"
	downloadUrl := getNginxDownloadUrl(nginxDomain)
	if downloadUrl == "" {
		color.Red("🌡 Can't find nginx download url....")
		return
	}
	spinner.Restart()
	nginxDownloadUrl := nginxDomain + downloadUrl
	resp, err := http.Get(nginxDownloadUrl)
	if err != nil {
		pe(err)
		return
	}
	defer resp.Body.Close()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		pe(err)
		return
	}

	nginxDir := filepath.Join(homeDir, "nginx")
	stat, err := os.Stat(nginxDir)
	if err != nil {
		err := os.Mkdir(nginxDir, 0777)
		if err != nil {
			pe(err)
			return
		}
	} else if !stat.IsDir() {
		err := os.Mkdir(nginxDir, 0777)
		if err != nil {
			pe(err)
			return
		}
	}

	nginxFileName := filepath.Join(nginxDir, filepath.Base(downloadUrl))
	file, err := os.Create(nginxFileName)
	if err != nil {
		pe(err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		pe(err)
		return
	}

	spinner.Stop()
	nginxTarName := filepath.Base(file.Name())
	commandExec(fmt.Sprintf("cd %s && tar -zxvf %s", nginxDir, nginxTarName))
	unTarDir := strings.ReplaceAll(nginxTarName, ".tar.gz", "")
	commandExec(fmt.Sprintf("cd %s && ./configure", filepath.Join(nginxDir, unTarDir)))
}

func getNginxDownloadUrl(nginxDomain string) string {
	spinner.Restart()
	nginxDownloadPageUrl := nginxDomain + "/en/download.html"
	resp, err := http.Get(nginxDownloadPageUrl)
	if err != nil {
		pe(err)
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		pe(err)
		return ""
	}

	var nginxDownUri string
	doc.Find("#content").Each(func(i int, s *goquery.Selection) {
		a1 := s.Find("table").Eq(1).Find("a").Eq(1)
		href, exists := a1.Attr("href")
		if exists {
			nginxDownUri = href
		}
	})
	spinner.Stop()
	return nginxDownUri
}

func installNginxFromBrew(brew string) {
	color.Cyan("🍺 About to install nginx from HomeBrew")
	commandExec(brew + " install nginx")
}

func installBrew() {
	spinner.Restart()
	command := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh")

	output, err := command.Output()
	if err != nil {
		pe(err)
		return
	}

	fileName := "./.brew.sh"
	err = ioutil.WriteFile(fileName, output, 0644)
	if err != nil {
		pe(err)
		return
	}
	defer os.Remove(fileName)

	err = os.Chmod(fileName, 0777)
	if err != nil {
		pe(err)
		return
	}

	spinner.Stop()
	commandExec(fileName)
}

func checkExists(file string) (path string, err error) {
	path, err = exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return path, nil
}
