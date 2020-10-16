package main

import (
	"archive/tar"
	"compress/gzip"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
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

var sshUser auth.SSHUser

func main() {
	var cobraAuth auth.Cobra

	root := cobra.Command{
		Use:        "sync_zb_resource",
		Aliases:    []string{"szr"},
		SuggestFor: []string{"sync", "resource", "zb"},
		Short:      "sync the resource from test environment and module",
		Long:       "sync the resource from test environment and module, you can choose the synchronize the specified file",
		Example:    "szr",
		PreRun:     cobraAuth.CreateOrChooseSSHUser(&sshUser),
		Run:        run,
	}

	root.AddCommand(cobraAuth.AddUserCmd(addUserFunc), cobraAuth.ListUserCmd())

	if err := root.Execute(); err != nil {
		color.Red("🌡 %+v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	sync := &Sync{}
	err := sync.connect(sshUser)
	if err != nil {
		color.Red(err.Error())
		return
	}
	defer sync.close()

	err = sync.shell()
	if err != nil {
		color.Red(err.Error())
		return
	}

	dashboards, servers := sync.dashboard()

	s.Stop()

	if len(dashboards) == 0 {
		color.Red("\n\tNo servers from %s", sshUser.Username)
		return
	}

	prompt := &survey.Select{
		Message: "Select a jump server node:",
		Options: dashboards,
	}

	var node string
	_ = survey.AskOne(prompt, &node)

	s.Restart()
	_, _ = sync.getPrompt(servers, node)
	sync.cd("cd /home/appl")
	modules := sync.getModules()
	s.Stop()

	if len(modules) == 0 {
		color.Red("\n\tNo modules from %s", servers[node])
		return
	}

	prompt = &survey.Select{
		Message: "Select a module for sync:",
		Options: modules,
	}

	var r string
	_ = survey.AskOne(prompt, &r)

	s.Restart()
	sync.cd("cd " + r)
	sync.cd("cd conf")
	configs := sync.getConfigs()
	s.Stop()
	if len(configs) > 0 {
		prompt := &survey.MultiSelect{
			Message: "Select files sync to current dir:",
			Options: configs,
		}

		var configFiles []string
		_ = survey.AskOne(prompt, &configFiles)

		if !doSync(s, sync, configFiles, servers[node]) {
			return
		}

		isConfigurationNginx := false
		nginxPrompt := &survey.Confirm{
			Message: "Do you want to configure nginx to access the page?",
		}
		_ = survey.AskOne(nginxPrompt, &isConfigurationNginx)

		if isConfigurationNginx {
			remotePort := sync.getRemotePort(r)
			CompleteNginx(s, sync, r, servers[node], remotePort)
		}
	}
}

func doSync(s *spinner.Spinner, sync *Sync, configs []string, ip string) bool {
	if len(configs) <= 0 {
		color.Red("\n\t🤒🤒You have not choose files to sync!!!\n\n")
		return false
	}

	// zip local files
	prompt := &survey.Confirm{
		Message: "Do you want to pack local files back up?",
		Help:    "Pack the local configuration file and place it in the `resources` directory",
		Default: true,
	}

	isZip := false
	_ = survey.AskOne(prompt, &isZip)

	path := filepath.Join("./src", "main", "resources")

	isGoing := true
	if isZip {
		err := packageFiles(s, path)
		if err != nil {
			prompt := &survey.Confirm{
				Message: "压缩tar包失败,是否继续?",
				Default: true,
				Help:    "Reason: " + err.Error(),
			}
			_ = survey.AskOne(prompt, &isGoing)
		} else {
			color.Green("🍺 Compress old files to %s", filepath.Join(path, "resources.tar.gz"))
		}
	}

	if !isGoing {
		return false
	}

	// copy remote file change local
	for _, config := range configs {
		s.Restart()
		content, err := sync.getContent(config)
		if err != nil {
			s.Stop()
			color.Red("Get remote content error, file: %s", config)
			continue
		}

		content = strings.ReplaceAll(content, ip, "")
		content = strings.ReplaceAll(content, "\r\n\r\n", "\r\n")
		f := filepath.Join(path, config)
		err = ioutil.WriteFile(f, []byte(content), 0644)
		if err != nil {
			s.Stop()
			color.Red("\t❗️❗️❗️File: %s write error: %s", config, err)
			return false
		}
		s.Stop()
		color.Green("🍺 %s file is synced", f)
	}
	return true
}

func packageFiles(s *spinner.Spinner, path string) error {
	s.Restart()
	defer s.Stop()
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
	color.Red("🌡 %+v", err)
}
