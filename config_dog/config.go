package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
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
	err = checkConfig(dogConfigPath, false)
	// .config.dog.json 不存在
	var content []byte
	var dogConfig DogConfig
	if err != nil {
		dogConfig.EnvPath = askEnvConfigPath(false)
		dogConfig.GitPath = ask("What is you config properties git URL?")
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
	err = checkConfig(envConfigFilePath, false)
	if err != nil {
		askAndCloneEnvConfigFromGit(dogConfig.EnvPath, dogConfig.GitPath)
	}

	bytes, err := ioutil.ReadFile(envConfigFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	if len(bytes) == 0 {
		askAndCloneEnvConfigFromGit(dogConfig.EnvPath, dogConfig.GitPath)
	}

	var a map[string]interface{}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return
	}

	log.Printf("%+v\n", a)

	return DogConfig{}, nil
}

func askAndCloneEnvConfigFromGit(envConfigPath, gitPath string) string {
	// check .git is exists
	dotGit := filepath.Join(envConfigPath, ".git")
	err := checkConfig(dotGit, false)
	if err == nil {
		runCmd(fmt.Sprintf("cd %s && git pull --rebase origin master", envConfigPath))
		return ""
	}

	dir := gitPath[strings.LastIndex(gitPath, "/")+1 : strings.LastIndex(gitPath, ".")]
	localGitPath := filepath.Join(envConfigPath, dir)
	runCmd(fmt.Sprintf("cd %s && git clone %s", envConfigPath, gitPath))
	log.Println("LocalGitPath:", localGitPath, "EncConfigPath:", envConfigPath)
	runCmd(fmt.Sprintf("mv -f %s/{*,.[^.]*,..?*} %s && rm -rf %s", localGitPath, envConfigPath, localGitPath)) // mv dot(hidden) and other files
	return ""
}

func askEnvConfigPath(repeat bool) string {
	message := "What is environment config file path?"
	if repeat {
		message = "The path does not exist, please re-enter"
	}
	envConfigPath := ask(message)

	err := checkConfig(envConfigPath, false)
	if err != nil {
		askEnvConfigPath(true)
	}
	return envConfigPath
}

func ask(message string) string {
	question := []*survey.Question{{
		Prompt:   &survey.Input{Message: color.MagentaString(message)},
		Validate: survey.Required,
	}}
	var result string
	_ = survey.Ask(question, &result)
	return result
}

func getDogConfigPath() string {
	if dogConfigPath != "" {
		return dogConfigPath
	}

	homeDir, _ := os.UserHomeDir()
	dogConfigPath = filepath.Join(homeDir, dogConfigFile)
	return dogConfigPath
}

func checkConfig(configPath string, isCreate bool) error {
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
