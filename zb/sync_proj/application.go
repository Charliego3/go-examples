package main

import (
	"archive/tar"
	"compress/gzip"
	"github.com/AlecAivazis/survey/v2"
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

	if dashboards != nil && len(dashboards) > 0 {
		prompt := &survey.Select{
			Message: "Select a jump server node 🗂:",
			Options: dashboards,
		}

		var r string
		_ = survey.AskOne(prompt, &r)

		sync.cd(servers[r])
		sync.cd("cd /home/appl")
		modules := sync.getModules()

		prompt = &survey.Select{
			Message: "Select a module for sync 🧩:",
			Options: modules,
		}

		_ = survey.AskOne(prompt, &r)

		sync.cd("cd " + r)
		confDir := string(colorMatch.ReplaceAll([]byte(sync.cd("cd conf")), []byte("")))
		configs := sync.getConfigs()
		if len(configs) > 0 {
			prompt := &survey.MultiSelect{
				Message: "Select files sync to current dir 📜:",
				Options: configs,
			}

			var configFiles []string
			_ = survey.AskOne(prompt, &configFiles)

			syncing(sync, configFiles, servers[r], confDir)
		}
	}

	sync.close()
}

func syncing(s *Sync, configs []string, ip string, confDir string) {
	if len(configs) <= 0 {
		color.Red("\n\t🤒🤒You have not choose files to sync!!!\n\n")
		return
	}

	// zip local files
	prompt := &survey.Confirm{
		Message: "Do you want to pack local files back up 📦?",
		Help:    "Pack the local configuration file and place it in the `resources` directory",
		Default: true,
	}

	isZip := false
	_ = survey.AskOne(prompt, &isZip)

	log.Println("是否打包", isZip)
	path := filepath.Join("./src", "main", "resources")

	isGoing := true
	if isZip {
		err := tarit(path)
		if err != nil {
			prompt := &survey.Confirm{
				Message: "压缩tar包失败,是否继续?",
				Default: true,
				Help:    "Reason: " + err.Error(),
			}
			_ = survey.AskOne(prompt, &isGoing)
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

		log.Printf("File: %s, Content: %s", config, content)
		content = strings.ReplaceAll(content, ip, "")
		content = strings.ReplaceAll(content, "\r\n\r\n", "\r\n")
		f := filepath.Join(path, config)
		err = ioutil.WriteFile(f, []byte(content), 0644)
		if err != nil {
			color.Red("\t❗️❗️❗️File: %s write error: %s", config, err)
		}
	}
}

func tarit(path string) error {
	// 创建文件
	zfn := "resources.tar.gz"
	fw, err := os.Create(filepath.Join(path, zfn))
	if err != nil {
		return err
	}
	defer fw.Close()

	// 将 tar 包使用 gzip 压缩，其实添加压缩功能很简单，
	// 只需要在 fw 和 tw 之前加上一层压缩就行了，和 Linux 的管道的感觉类似
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(gw)
	// 如果需要启用 gzip 将上面代码注释，换成下面的

	defer tw.Close()

	// 下面就该开始处理数据了，这里的思路就是递归处理目录及目录下的所有文件和目录
	// 这里可以自己写个递归来处理，不过 Golang 提供了 filepath.Walk 函数，可以很方便的做这个事情
	// 直接将这个函数的处理结果返回就行，需要传给它一个源文件或目录，它就可以自己去处理
	// 我们就只需要去实现我们自己的 打包逻辑即可，不需要再去路径相关的事情
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
		// 这里需要处理下 hdr 中的 Name，因为默认文件的名字是不带路径的，
		// 打包之后所有文件就会堆在一起，这样就破坏了原本的目录结果
		// 例如： 将原本 hdr.Name 的 syslog 替换程 log/syslog
		// 这个其实也很简单，回调函数的 fileName 字段给我们返回来的就是完整路径的 log/syslog
		// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
		// 熟悉 Linux 的都知道为什么要去掉这个
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

		// 记录下过程，这个可以不记录，这个看需要，这样可以看到打包的过程
		return nil
	})

	return err
}
