package main

import (
	"encoding/json"
	"fmt"
	"github.com/gen2brain/dlgs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DogConfig struct {
	GitPath string
	EnvPath string
}

const (
	dogConfigFile = ".config.dog.json"
	envConfigName = ".environment.json"
)

var (
	dogConfigPath = ""
)

func init() {
	dogConfigPath = getDogConfigPath()
}

func getDogConfig() (config DogConfig, err error) {
	// 检查$HOME/.config.dog.json文件是否存在
	err = checkPath(dogConfigPath, false)
	// .config.dog.json 不存在
	var content []byte
	var dogConfig DogConfig
	if err != nil {
		dogConfig.EnvPath = askEnvConfigPath()
		gitUrl, b, _ := dlgs.Entry("Git URL", "What is you config properties git URL?", "")
		fmt.Printf("%v, %v, %v", gitUrl, b, err)
		dogConfig.GitPath = gitUrl
		log.Println("EnvConfigDirPath:", dogConfig.EnvPath, "GitPath:", dogConfig.GitPath)
		content, err = json.Marshal(dogConfig)
		if err != nil {
			log.Println(err)
			return
		}
		err = ioutil.WriteFile(dogConfigPath, content, 0666)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		content, err = ioutil.ReadFile(dogConfigPath)
		if err != nil {
			return
		}

		err = json.Unmarshal(content, &dogConfig)
		if err != nil {
			log.Println(err)
			return
		}
	}

	envConfigFilePath := filepath.Join(dogConfig.EnvPath, envConfigName)
	// 检查$envConfigDirPath/.environment.json文件是否存在
	err = checkPath(envConfigFilePath, false)
	if err != nil {
		askAndCloneEnvConfigFromGit(envConfigFilePath, dogConfig.EnvPath, dogConfig.GitPath)
	}

	bytes, err := ioutil.ReadFile(envConfigFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	var a map[string]interface{}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return
	}

	log.Printf("%+v\n", a)

	return DogConfig{}, nil
}

func askAndCloneEnvConfigFromGit(envConfigFilePath, envConfigPath, gitPath string) string {
	defer func() {
		err := checkPath(envConfigFilePath, false)
		if err != nil {
			message := "\n‼️  Config File not found: %s\n\n‼️  Confirm that the git url: %s\n\nif you want to change git url can be use `config_dog --git 'your git url'`\n"
			_, _ = dlgs.Error("Sync environment error", fmt.Sprintf(message, envConfigFilePath, gitPath))
			os.Exit(0)
		}
	}()

	// check .git is exists
	dotGit := filepath.Join(envConfigPath, ".git")
	err := checkPath(dotGit, false)
	if err == nil {
		runCmd(fmt.Sprintf("cd %s && git pull --rebase origin master", envConfigPath))
		return ""
	}

	dir := gitPath[strings.LastIndex(gitPath, "/")+1 : strings.LastIndex(gitPath, ".")]
	localGitPath := filepath.Join(envConfigPath, dir)
	runCmd(fmt.Sprintf("cd %s && git clone %s", envConfigPath, gitPath))
	log.Println("LocalGitPath:", localGitPath, "EncConfigPath:", envConfigPath)
	runCmd(fmt.Sprintf("mv -f %s/{*,.[^.]*} %s && rm -rf %s", localGitPath, envConfigPath, localGitPath)) // mv dot(hidden) and other files
	return ""
}

func askEnvConfigPath() string {
	message := "What is environment config dir path?"
	envConfigPath := ask(message, true)

	err := checkPath(envConfigPath, false)
	if err != nil {
		askEnvConfigPath()
	}
	return envConfigPath
}

func ask(message string, isDir bool) string {
	file, b, err := dlgs.File(message, "", isDir)
	if err != nil || !b {
		_, _ = dlgs.Warning("Path is empty", "You not chose any path, now exit.")
		os.Exit(0)
	}
	return file
}

func getDogConfigPath() string {
	if dogConfigPath != "" {
		return dogConfigPath
	}

	homeDir, _ := os.UserHomeDir()
	dogConfigPath = filepath.Join(homeDir, dogConfigFile)
	return dogConfigPath
}

func checkPath(configPath string, isCreate bool) error {
	stat, err := os.Stat(configPath)
	if err != nil {
		if isCreate {
			createConfigFile(configPath)
		}
		return err
	}

	if stat.IsDir() && isCreate {
		createConfigFile(configPath)
	}
	return nil
}

func createConfigFile(configPath string) {
	_, err := os.Create(configPath)
	if err != nil {
		log.Printf("创建%s失败: %+v\n", configPath, err)
	}
}
