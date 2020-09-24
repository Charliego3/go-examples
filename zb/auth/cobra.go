package auth

import (
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
)

type Cobra struct {}

func (c *Cobra) CreateOrChooseSSHUser(sshUser *SSHUser) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := CreateOrChooseSSHUser(sshUser)
		if err != nil {
			logger.Fatal("Create or choose ssh user error: %s", err)
		}
	}
}

func (c *Cobra) AddUser() (*Command, error) {
	return nil, nil
}

func (c *Cobra) AddUserForCobra() *cobra.Command {
	su := SSHUser{}
	au := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			err, choose := AddUser(&su)
			if err != nil {
				logger.Error("Create new user fail: %s", err)
			}

			logger.Debug("Cobra Add User: %+v, Choose: %v", su,choose)
		},
	}

	au.Flags().StringVarP(&su.Host, "host", "H", "", "The server host")
	au.Flags().StringVarP(&su.Port, "port", "p", "", "The server port")
	au.Flags().StringVarP(&su.Username, "user", "u", "", "The ssh server user")
	au.Flags().StringVarP(&su.Password, "password", "P", "", "The server password")

	return au
}