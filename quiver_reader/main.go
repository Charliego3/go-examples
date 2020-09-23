package main

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"io/ioutil"
	"os"
	"path/filepath"
)

type meta struct {
	Name string `json:"name"`
	Uuid     string `json:"uuid"`
	Children []meta `json:"children"`
}

type content struct {
	Title string `json:"title"`
	Cells []struct {
		Data string `json:"data"`
	} `json:"cells"`
}

const (
	metaJson = "meta.json"
	contentJson = "content.json"
)

var (
	bp string
	tp string
	targetDir string
)

func main() {
	root := cobra.Command{
		Use:     "quiver_parse",
		Aliases: []string{"quiver"},
		Run:     parseQuiver(),
	}

	root.Flags().StringVarP(&bp, "file", "f", "/Users/nzlong/Dropbox/Quiver/Quiver.qvlibrary", "Quiver library file path")
	root.Flags().StringVarP(&tp, "path", "p", "/Users/nzlong/Desktop", "Target path")

	err := root.Execute()
	checkErr(err)
}

func parseQuiver() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		p := filepath.Join(bp, metaJson)
		_, err := os.Stat(p)
		checkErr(err)

		var m meta
		err = json.Unmarshal(read(p), &m)
		checkErr(err)

		logger.Debug("%+v", m)

		targetDir = filepath.Join(tp, m.Uuid)
		createDirIfNotExist(targetDir)
		rangeChildren(m.Children, targetDir)
	}
}

func rangeChildren(children []meta, path string) {
	for _, child := range children {
		sourcePath := filepath.Join(bp, child.Uuid + ".qvnotebook")
		meta := readMeta(filepath.Join(sourcePath, metaJson))
		dirAbsolutePath := filepath.Join(path, meta.Name)
		logger.Info("SourcePath: %s, DirName: %s, Absolute Path: %s", sourcePath, meta.Name, dirAbsolutePath)
		createDirIfNotExist(dirAbsolutePath)

		// 单目录
		if child.Children == nil {
			err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
				if filepath.Ext(info.Name()) == ".qvnote" {
					c := readContent(filepath.Join(path, contentJson))
					if c != nil {
						logger.Debug("FilePath: %s/%s", path, c.Title)
						builder := bytes.Buffer{}
						if c.Cells != nil {
							for _, cell := range c.Cells {
								builder.WriteString(cell.Data)
								builder.WriteString("\n\n")
							}
							err := ioutil.WriteFile(filepath.Join(dirAbsolutePath, c.Title+".md"), builder.Bytes()[:builder.Len()-2], 0644)
							checkErr(err)
						}
					}
				}
				return nil
			})
			checkErr(err)
		} else {
			rangeChildren(child.Children, filepath.Join(path, meta.Name))
		}
	}
}

func createDirIfNotExist(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		checkErr(err)
	}
}

func readContent(path string) *content {
	var c content
	bs := read(path)
	err := json.Unmarshal(bs, &c)
	checkErr(err)
	return &c
}

func readMeta(path string) meta {
	var m meta
	bs := read(path)
	err := json.Unmarshal(bs, &m)
	checkErr(err)
	return m
}

func read(path string) []byte {
	bs, err := ioutil.ReadFile(path)
	checkErr(err)
	return bs
}

func checkErr(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
