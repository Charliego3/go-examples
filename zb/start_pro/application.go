package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/temp/zb/auth"
	"log"
	"os/exec"
)

func main() {

	var sshUser auth.SSHUser

	cobraAuth := auth.Cobra{}

	root := cobra.Command{
		Use:        "sync_zb_resource",
		Aliases:    []string{"szr"},
		SuggestFor: []string{"sync", "resource", "zb"},
		Short:      "sync the resource from test environment and module",
		Long:       "sync the resource from test environment and module, you can choose the synchronize the specified file",
		Example:    "szr",
		PreRun:     cobraAuth.CreateOrChooseSSHUser(&sshUser),
		Run: func(cmd *cobra.Command, args []string) {
			logger.Error("RUN......")
			logger.Debug("SSHUser: %+v", sshUser)
		},
	}

	path, err := exec.LookPath("aaaa")
	if err != nil {
		log.Fatal("installing brew is in your future")
	}
	fmt.Printf("fortune is available at %s\n", path)

	root.AddCommand(cobraAuth.AddUserForCobra())

	if err := root.Execute(); err != nil {
		logger.Fatal("%+v", err)
	}
}
