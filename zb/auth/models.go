package auth

import "github.com/spf13/cobra"

type Command struct {
	*cobra.Command
}

type Authentication interface {
	CreateOrChooseSSHUser(su *SSHUser) error
	AddUser() (*Command, error)
}