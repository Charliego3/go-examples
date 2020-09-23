package auth

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
)

var qs = []*survey.Question{
	{
		Name:     "Host",
		Prompt:   &survey.Input{Message: "What is server ssh host?"},
		Validate: survey.Required,
	},
	{
		Name:     "Port",
		Prompt:   &survey.Input{Message: "What is server ssh port?"},
		Validate: survey.Required,
	},
	{
		Name:     "User",
		Prompt:   &survey.Input{Message: "What is username?"},
		Validate: survey.Required,
	},
	{
		Name:     "Password",
		Prompt:   &survey.Password{Message: "What is the password of this user?"},
		Validate: survey.Required,
	},
}

type SSHUser struct {
	Host     string
	Port     string
	User     string
	Password string
}

func init() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("$HOME")
}

func Auth(sshUser *SSHUser) error {
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
				"host":     answers.Host,
				"port":     answers.Port,
				"username": answers.User,
				"password": answers.Password,
			}
			viper.Set(usersKey, users)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}
			*sshUser = answers
		}
	} else {
		selectedUser := ""
		separator := " - "
		var userNames []string
		for _, su := range users {
			m := su.(map[string]interface{})
			name := net.JoinHostPort(m["host"].(string), m["port"].(string))
			name += separator
			name += m["username"].(string)
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

		for _, su := range users {
			m := su.(map[string]interface{})
			if h == m["host"].(string) && p == m["port"].(string) && msg[1] == m["username"].(string) {
				*sshUser = SSHUser{
					Host:     h,
					Port:     p,
					User:     msg[1],
					Password: m["password"].(string),
				}
				return nil
			}
		}
	}
	return nil
}

func AddUser() *cobra.Command {
	au := &cobra.Command{

	}
	return au
}
