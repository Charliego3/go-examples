package main

import (
	"encoding/json"
	"flag"
	"github.com/whimthen/kits/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	uid      string
	clientId string
	token    string
	src      string
)

func main() {
	//quiver()
	targetDir := "/Users/nzlong/Desktop/redis"

	juejin(targetDir)
}

func quiver() {
	baseDir := "/Users/nzlong/Dropbox/Quiver/Quiver.qvlibrary"

	files, err := os.Open(baseDir)
	if err != nil {
		logger.Error("Open Quiver Error: %+v", err)
		return
	}

	stat, err := files.Stat()
	if stat.IsDir() {
		logger.Info("%s is Dir", baseDir)
		dirs, err := files.Readdir(-1)
		if err != nil {
			logger.Error("Read dir error: %+v", err)
			return
		}

		for _, d := range dirs {
			logger.Info("ChDir Name: %s, IsDir: %v", d.Name(), d.IsDir())
			if d.Name() == "meta.json" {
				path := filepath.Join(baseDir, d.Name())

				bytes, err := ioutil.ReadFile(path)
				if err != nil {
					logger.Error("Read %s error: %+v", path, err)
					return
				}

				//logger.Info("%s: %s", path, bytes)

				meta := Meta{}
				_ = json.Unmarshal(bytes, &meta)


				children := meta.Children[3].Children
				for _, child := range children {

					chdj := filepath.Join(baseDir, child.UUID+".qvnotebook", "meta.json")
					file, _ := ioutil.ReadFile(chdj)

					section := QSection{}
					_ = json.Unmarshal(file, &section)

					if section.Name == "MySql" {
						logger.Info("%+v", section)
					}
				}

			}
		}
	}
}

func juejin(targetDir string) {
	flag.StringVar(&uid, "uid", "5b30bc9ae51d4558ac4890aa", "The Juejin login user uid")
	flag.StringVar(&clientId, "cid", "1595904852284", "The Juejin login user clientId")
	flag.StringVar(&token, "token", "eyJhY2Nlc3NfdG9rZW4iOiJkRkNBd0lHS1FnWGIwVnh3IiwicmVmcmVzaF90b2tlbiI6InlpMmhNejBVSXVISzFFbGwiLCJ0b2tlbl90eXBlIjoibWFjIiwiZXhwaXJlX2luIjoyNTkyMDAwfQ==", "Login user token")
	flag.StringVar(&src, "src", "web", "src, emmm.........")
	flag.Parse()

	q := &QueryBase{
		UID:      uid,
		ClientId: clientId,
		Token:    token,
		Src:      src,
	}

	logger.Debug("The Q is %+v\n", q)

	listSection := q.GetListSection("5afc2e5f6fb9a07a9b362527")
	for index, section := range listSection.D {
		content := q.GetSection(section.Id)
		mdp := filepath.Join(targetDir, strconv.Itoa(index+1) + "、 "+content.D.Title+".md")

		err := ioutil.WriteFile(mdp, []byte(content.D.Content), 0644)
		if err != nil {
			logger.Error("WriteToFile error: %+v", err)
		}
		time.Sleep(time.Second * 2)
	}
}
