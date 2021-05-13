package main

import (
	"encoding/json"
	"github.com/kataras/golog"
	"io"
	"os"
	"path/filepath"
)

var cachePath = ".zszj.cache"

func init() {
	home , _ := os.UserHomeDir()
	cachePath = filepath.Join(home, cachePath)
}

type st struct {
	Name string
	Num  int
}

var sts []st

func UpdateCache(name string, num int) {
	file, err := os.OpenFile(cachePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		golog.Error("Can't open the cache file")
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		golog.Errorf("Read cache file error: %v", err)
		return
	}

	err = json.Unmarshal(fileBytes, &sts)
	if err != nil {
		golog.Errorf("Unmarshal cache error: %v", err)
		return
	}

	for _, st := range sts {
		if st.Name != name {
			continue
		}

		if st.Num < num {
			st.Num = num
			err = WriteSt(file, name, num)
			if err != nil {
				golog.Errorf("Write cache file error: %v", err)
			}
		}
	}
}

func WriteSt(file *os.File, name string, num int) error {
	bytes, err := json.Marshal(st{
		Name: name,
		Num:  num,
	})
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}