package auth

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

type Cobra struct{}

func (c *Cobra) CreateOrChooseSSHUser(sshUser *SSHUser) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := CreateOrChooseSSHUser(sshUser)
		if err != nil {
			color.Red("Create or choose ssh user error: %s", err)
			os.Exit(128)
		}
	}
}

func (c *Cobra) AddUserCmd(f func(su *SSHUser) error) *cobra.Command {
	su := SSHUser{}
	au := &cobra.Command{
		Use: "add",
		RunE: func(cmd *cobra.Command, args []string) error {
			err, choose := AddUser(&su)
			if err != nil {
				color.Red("🌡 %+v", err)
				return err
			}

			if choose {
				return f(&su)
			}
			return nil
		},
	}

	au.Flags().StringVarP(&su.Host, "host", "H", "", "The server host")
	au.Flags().StringVarP(&su.Port, "port", "p", "", "The server port")
	au.Flags().StringVarP(&su.Username, "user", "u", "", "The ssh server user")
	au.Flags().StringVarP(&su.Password, "password", "P", "", "The server password")

	return au
}

func (c *Cobra) ListUserCmd() *cobra.Command {
	au := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			userNames := getUserNames()
			for i, username := range userNames {
				color.Green("%d. %s", i+1, username)
			}
		},
	}
	return au
}
