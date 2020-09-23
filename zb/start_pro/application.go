package main

import (
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/temp/zb/auth"
)

func main() {

	var sshUser auth.SSHUser

	root := cobra.Command{
		Use:        "sync_zb_resource",
		Aliases:    []string{"szr"},
		SuggestFor: []string{"sync", "resource", "zb"},
		Short:      "sync the resource from test environment and module",
		Long:       "sync the resource from test environment and module, you can choose the synchronize the specified file",
		Example:    "szr",
		PreRun: func(cmd *cobra.Command, args []string) {
			auth.CreateOrChooseSSHUser(&sshUser)
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger.Error("RUN......")
			logger.Debug("SSHUser: %+v", sshUser)
		},
	}

	//root.AddCommand()

	if err := root.Execute(); err != nil {
		logger.Fatal("%+v", err)
	}
}
