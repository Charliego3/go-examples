package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type Dog struct {
	watcher *fsnotify.Watcher
}

func NewDog() *Dog {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {

	}
	return &Dog{
		watcher: watcher,
	}
}

func (d *Dog) Watch(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		// dir is not exists
		return
	}

	if stat.IsDir() {
		d.watchDir(path)
	} else {
		d.watchFile(path)
	}
}

func (d *Dog) watchDir(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = d.watcher.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}

		return nil
	})

	if err != nil {
		log.Println("WalkWatch error", err)
	}
}

func (d *Dog) watchFile(file string) {
	abs, _ := filepath.Abs(file)
	_ = d.watcher.Add(abs)
}

func (d *Dog) EventHandler() {
	for {
		select {
		case ev := <-d.watcher.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("创建文件 : ", ev.Name)
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						d.watcher.Add(ev.Name)
						fmt.Println("添加监控 : ", ev.Name)
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					//fmt.Println("写入文件 : ", ev.Name)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("删除文件 : ", ev.Name)
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						d.watcher.Remove(ev.Name)
						fmt.Println("删除监控 : ", ev.Name)
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					fmt.Println("重命名文件 : ", ev.Name)
					d.watcher.Remove(ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-d.watcher.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}
}
