package main

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type DogConfig struct {
	GitPath string
	Links   map[string]string
	isEmpty bool
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
	var envConfigPath string
	// .config.dog.json 不存在
	if err != nil {
		envConfigPath = askEnvConfigPath(false)
		log.Println("EnvConfigPath:", envConfigPath)
		err = ioutil.WriteFile(dogConfigPath, []byte(envConfigPath), 0666)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		var bytes []byte
		bytes, err = ioutil.ReadFile(dogConfigPath)
		if err != nil {
			return
		}

		envConfigPath = string(bytes)
	}

	// 检查$envConfigPath/.environment.json文件是否存在, 不存在则创建
	err = checkConfig(envConfigPath, true)
	if err != nil {
		_ = ask("What is you backup git URL?")
		progress := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		progress.Start()
		progress.Stop()
		return
	}

	bytes, err := ioutil.ReadFile(envConfigPath)
	if err != nil {
		return
	}

	var a map[string]interface{}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return
	}

	return DogConfig{
		isEmpty: true,
	}, nil
}

func askEnvConfigPath(repeat bool) string {
	message := "What is environment config file path?"
	if repeat {
		message = "The path does not exist, please re-enter"
	}
	envConfigPath := ask(message)

	envConfigPath = filepath.Join(envConfigPath, envConfigName)
	err := checkConfig(envConfigPath, false)
	if err != nil {
		askEnvConfigPath(true)
	}
	return envConfigPath
}

func ask(message string) string {
	question := []*survey.Question{{
		Prompt:   &survey.Input{Message: message},
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
