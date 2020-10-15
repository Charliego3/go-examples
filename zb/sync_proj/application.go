package main

import (
	"archive/tar"
	"compress/gzip"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/temp/zb/auth"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	root.AddCommand(cobraAuth.AddUserCmd(addUserFunc))

	if err := root.Execute(); err != nil {
		logger.Fatal("%+v", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	sync := &Sync{}
	err := sync.connect(sshUser)
	if err != nil {
		log.Fatal(err)
	}

	err = sync.shell()
	if err != nil {
		log.Fatal(err)
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

	var r string
	_ = survey.AskOne(prompt, &r)

	s.Restart()
	_, _ = sync.getPrompt(servers, r)
	sync.cd("cd /home/appl")
	modules := sync.getModules()
	s.Stop()

	if len(modules) == 0 {
		color.Red("\n\tNo modules from %s", servers[r])
		return
	}

	prompt = &survey.Select{
		Message: "Select a module for sync:",
		Options: modules,
	}

	_ = survey.AskOne(prompt, &r)

	s.Restart()
	sync.cd("cd " + r)
	confDir := string(colorMatch.ReplaceAll([]byte(sync.cd("cd conf")), []byte("")))
	configs := sync.getConfigs()
	s.Stop()
	if len(configs) > 0 {
		prompt := &survey.MultiSelect{
			Message: "Select files sync to current dir:",
			Options: configs,
		}

		var configFiles []string
		_ = survey.AskOne(prompt, &configFiles)

		doSync(sync, configFiles, servers[r], confDir)

		isConfigurationNginx := false
		nginxPrompt := &survey.Confirm{
			Message: "Do you want to configure nginx to access the page?",
		}
		_ = survey.AskOne(nginxPrompt, &isConfigurationNginx)

		if isConfigurationNginx {
			CompleteNginx()
		}
	}

	sync.close()
}

func doSync(s *Sync, configs []string, ip string, confDir string) {
	if len(configs) <= 0 {
		color.Red("\n\t🤒🤒You have not choose files to sync!!!\n\n")
		return
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
		err := packageFiles(path)
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
		return
	}

	// copy remote file change local
	for _, config := range configs {
		content, err := s.getContent(config, confDir)
		if err != nil {
			color.Red("Get remote content error, file: %s", config)
			continue
		}

		content = strings.ReplaceAll(content, ip, "")
		content = strings.ReplaceAll(content, "\r\n\r\n", "\r\n")
		f := filepath.Join(path, config)
		color.Green("🍺 %s file is synced", f)
		err = ioutil.WriteFile(f, []byte(content), 0644)
		if err != nil {
			color.Red("\t❗️❗️❗️File: %s write error: %s", config, err)
		}
	}
}

func packageFiles(path string) error {
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
