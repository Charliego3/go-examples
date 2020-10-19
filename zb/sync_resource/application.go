package main

import (
	"archive/tar"
	"compress/gzip"
	"github.com/AlecAivazis/survey/v2"
	spinner2 "github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/whimthen/temp/zb/auth"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	sshUser   *auth.SSHUser
	cobraAuth auth.Cobra
	spinner   *spinner2.Spinner
)

func main() {
	root := cobra.Command{
		Use:        "sync_resource",
		SuggestFor: []string{"sync", "resource", "zb"},
		Short:      "sync the resource from test environment and module",
		Long:       "sync the resource from test environment and module, you can choose the synchronize the specified file",
		Example:    "sync_resource",
		Run:        run,
	}

	root.AddCommand(cobraAuth.AddUserCmd(addUserFunc), cobraAuth.ListUserCmd())

	if err := root.Execute(); err != nil {
		color.Red("🌡  %+v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	path := filepath.Join("./src", "main", "resources")
	_, err := os.Stat(path)
	if err != nil {
		color.Red("🌡  %+v", "this directory is not the correct Maven structure")
		return
	}

	sshUser = &auth.SSHUser{}
	cobraAuth.CreateOrChooseSSHUser(sshUser)(nil, nil)

	if sshUser == nil {
		return
	}

	if sshUser.Host == "" || sshUser.Port == "" || sshUser.Username == "" || sshUser.Password == "" {
		color.Red("🌡  No choose user. Exit.")
		return
	}

	spinner = spinner2.New(spinner2.CharSets[11], 100*time.Millisecond)
	spinner.Start()
	sync := &Sync{}
	err = sync.connect(*sshUser)
	if err != nil {
		pe(err)
		return
	}
	defer sync.close()

	err = sync.shell()
	if err != nil {
		pe(err)
		return
	}

	dashboards, servers := sync.dashboard()

	spinner.Stop()

	if len(dashboards) == 0 {
		color.Red("🌡  No servers from %s", sshUser.Username)
		return
	}

	prompt := &survey.Select{
		Message: "Select a jump server node:",
		Options: dashboards,
	}

	var node string
	err = survey.AskOne(prompt, &node)
	if err != nil {
		spinner.Stop()
		return
	}

	spinner.Restart()
	_, err = sync.getPrompt(servers, node)
	if err != nil {
		pe(err)
		return
	}
	sync.cd("cd /home/appl")
	modules := sync.getModules()
	spinner.Stop()

	if len(modules) == 0 {
		spinner.Stop()
		color.Red("🌡  No modules from %s", servers[node])
		return
	}

	prompt = &survey.Select{
		Message: "Select a module for sync:",
		Options: modules,
	}

	var module string
	err = survey.AskOne(prompt, &module)
	if err != nil {
		spinner.Stop()
		return
	}

	spinner.Restart()
	sync.cd("cd " + module)
	sync.cd("cd conf")
	configs := sync.getConfigs()
	spinner.Stop()
	if len(configs) > 0 {
		prompt := &survey.MultiSelect{
			Message: "Select files sync to current dir:",
			Options: configs,
		}

		var configFiles []string
		err = survey.AskOne(prompt, &configFiles)
		if err != nil {
			spinner.Stop()
			return
		}

		if !doSync(sync, module, configFiles, servers[node]) {
			return
		}

		isConfigurationNginx := false
		nginxPrompt := &survey.Confirm{
			Message: "Do you want to configure nginx to access the page?",
		}
		err = survey.AskOne(nginxPrompt, &isConfigurationNginx)
		if err != nil {
			spinner.Stop()
			return
		}

		if isConfigurationNginx {
			remotePort := sync.getRemotePort(module)
			CompleteNginx(sync, module, servers[node], remotePort)
		}
	}
}

func doSync(sync *Sync, module string, configs []string, ip string) bool {
	if len(configs) <= 0 {
		color.Red("🌡  You have not choose files to sync!!!\n\n")
		return false
	}

	// zip local files
	prompt := &survey.Confirm{
		Message: "Do you want to pack local files back up?",
		Help:    "Pack the local configuration file and place it in the `resources` directory",
		Default: true,
	}

	isZip := false
	err := survey.AskOne(prompt, &isZip)
	if err != nil {
		spinner.Stop()
		return false
	}

	path := filepath.Join("./src", "main", "resources")

	isGoing := true
	if isZip {
		err := packageFiles(path)
		if err != nil {
			prompt := &survey.Confirm{
				Message: "压缩tar包失败,是否继续?",
				Default: true,
				Help:    "Reason: " + err.Error(),
			}
			err = survey.AskOne(prompt, &isGoing)
			if err != nil {
				spinner.Stop()
				return false
			}
		} else {
			color.Green("🍺 Compress old files to %s", filepath.Join(path, "resources.tar.gz"))
		}
	}

	if !isGoing {
		return false
	}

	// copy remote file change local
	for _, config := range configs {
		spinner.Restart()
		content, err := sync.getContent(config)
		if err != nil {
			spinner.Stop()
			color.Red("Get remote content error, file: %s", config)
			continue
		}

		pwd, err := os.Getwd()
		if err != nil {
			pe(err)
			return false
		}
		content = strings.ReplaceAll(content, filepath.Join("/home/appl/", module), pwd)
		content = strings.ReplaceAll(content, "127.0.0.1", ip)
		content = strings.ReplaceAll(content, "\r\n\r\n", "\r\n")
		f := filepath.Join(path, config)
		err = ioutil.WriteFile(f, []byte(content), 0644)
		if err != nil {
			spinner.Stop()
			color.Red("\t❗️❗️❗️File: %s write error: %s", config, err)
			return false
		}
		spinner.Stop()
		color.Green("🍺 %s file is synced", f)
	}
	return true
}

func packageFiles(path string) error {
	spinner.Restart()
	defer spinner.Stop()
	// 创建文件
	zfn := "resources.tar.gz"
	fw, err := os.Create(filepath.Join(path, zfn))
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(gw)
	defer tw.Close()

	err = filepath.Walk(path, func(fileName string, fi os.FileInfo, err error) error {
		// 因为这个闭包会返回个 error ，所以先要处理一下这个
		if err != nil {
			return err
		}

		base := filepath.Base(fileName)
		if base == zfn {
			return nil
		}

		if strings.HasPrefix(base, ".") {
			return nil
		}

		if fi.IsDir() && strings.HasSuffix(path, fileName) {
			return nil
		}

		var link string
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		// 这里就不需要我们自己再 os.Stat 了，它已经做好了，我们直接使用 fi 即可
		hdr, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		hdr.Name = strings.TrimPrefix(fileName, path)

		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer fr.Close()

		// copy 文件数据到 tw
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func commandExec(cmd string) bool {
	command := exec.Command("bash", "-c", cmd)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	err := command.Run()
	if err != nil {
		pe(err)
		return false
	}
	return true
}

func pe(err error) {
	spinner.Stop()
	color.Red("🌡  %+v", err)
}
