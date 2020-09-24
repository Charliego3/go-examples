package auth

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	configName = ".sync_zb_resource"
	configType = "json"
	usersKey   = "users"

	host     = "host"
	port     = "port"
	username = "username"
	password = "password"
)

var (
	hostQuestion = &survey.Question{
		Name:     "Host",
		Prompt:   &survey.Input{Message: "What is ssh server host?"},
		Validate: survey.Required,
	}

	portQuestion = &survey.Question{
		Name:     "Port",
		Prompt:   &survey.Input{Message: "What is ssh server port?"},
		Validate: survey.Required,
	}

	userQuestion = &survey.Question{
		Name:     "Username",
		Prompt:   &survey.Input{Message: "What is username?"},
		Validate: survey.Required,
	}

	passwordQuestion = &survey.Question{
		Name:     "Password",
		Prompt:   &survey.Password{Message: "What is the password of this user?"},
		Validate: survey.Required,
	}

	qs = []*survey.Question{
		hostQuestion,
		portQuestion,
		userQuestion,
		passwordQuestion,
	}


	bred = color.New(color.FgRed, color.Bold)
	prompt = color.New(color.FgBlue, color.Bold).Sprint("\n==>")
)

type SSHUser struct {
	Host     string
	Port     string
	Username string
	Password string
}

func init() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("$HOME")
}

func initConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			dir, _ := os.UserHomeDir()
			err := viper.WriteConfigAs(filepath.Join(dir, configName+"."+configType))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func CreateOrChooseSSHUser(su *SSHUser) error {
	if err := initConfig(); err != nil {
		return err
	}
	users := viper.GetStringMap(usersKey)
	if len(users) == 0 {
		createUser := false
		prompt := &survey.Confirm{
			Message: color.GreenString("Do you want to create a new user?"),
		}
		_ = survey.AskOne(prompt, &createUser)

		if createUser {
			answers := SSHUser{}

			err := survey.Ask(qs, &answers)
			if err != nil {
				return err
			}

			users["0"] = map[string]string{
				host:     answers.Host,
				port:     answers.Port,
				username: answers.Username,
				password: answers.Password,
			}
			viper.Set(usersKey, users)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}
			*su = answers
		}
	} else {
		_ = selectUser(users, su)
	}
	return nil
}

func selectUser(users map[string]interface{}, su *SSHUser) error {
	selectedUser := ""
	separator := " - "
	var userNames []string
	for _, u := range users {
		m := u.(map[string]interface{})
		name := net.JoinHostPort(m[host].(string), m[port].(string))
		name += separator
		name += m[username].(string)
		userNames = append(userNames, name)
	}
	prompt := &survey.Select{
		Message: "Choose a user:",
		Options: userNames,
	}
	err := survey.AskOne(prompt, &selectedUser)
	if err != nil {
		return err
	}

	msg := strings.Split(selectedUser, separator)
	h, p, err := net.SplitHostPort(msg[0])
	if err != nil {
		return err
	}

	return loopUsers(h, p, msg[1], func(users map[string]interface{}, m map[string]interface{}) error {
		*su = SSHUser{
			Host:     h,
			Port:     p,
			Username: msg[1],
			Password: m[password].(string),
		}
		return nil
	})
}

func AddUser(su *SSHUser) (err error, choose bool) {
	if err := initConfig(); err != nil {
		return err, false
	}
	var hq []*survey.Question
	if su.Host == "" {
		hq = append(hq, hostQuestion)
	}
	if su.Port == "" {
		hq = append(hq, portQuestion)
	}
	if su.Username == "" {
		hq = append(hq, userQuestion)
	}
	if su.Password == "" {
		hq = append(hq, passwordQuestion)
	}
	err = survey.Ask(hq, su)
	if err != nil {
		return err, false
	}

	err = loopUsers(su.Host, su.Port, su.Username, func(users, m map[string]interface{}) error {
		println(prompt, bred.Sprint("This user for server is already exist."))
		update := false
		err := survey.AskOne(&survey.Confirm{
			Message: color.RedString(fmt.Sprintf("Do you want to update user(%s - %s)?", net.JoinHostPort(su.Host, su.Port), su.Username)),
		}, &update)
		if err != nil {
			return err
		}

		if update {
			m[password] = su.Password
			viper.Set(usersKey, users)
			err := viper.WriteConfig()
			if err != nil {
				println(prompt, bred.Sprintf("Update this user fail: %s", err))
				return err
			}

			connect := false
			err = survey.AskOne(&survey.Confirm{
				Message: color.RedString("Do you want to sync server prop?"),
			}, &connect)
			if err != nil {
				return err
			}
			if connect {
				_ = selectUser(users, su)
				choose = true
			}
		}
		return nil
	})
	if err != nil {
		return err, false
	}
	return
}

func loopUsers(h, p, u string, f func(map[string]interface{}, map[string]interface{}) error) error {
	users := viper.GetStringMap(usersKey)
	if len(users) > 0 {
		for _, user := range users {
			m := user.(map[string]interface{})
			if h == m[host].(string) && p == m[port].(string) && u == m[username].(string) {
				return f(users, m)
			}
		}
	}
	return nil
}
