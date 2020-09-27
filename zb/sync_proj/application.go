package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/temp/zb/auth"
	"log"
)

var sshUser auth.SSHUser

func main() {
	var cobraAuth auth.Cobra

	root := cobra.Command{
		Use:        "sync_zb_resource",
		Aliases:    []string{"szr"},
		SuggestFor: []string{"sync", "resource", "zb"},
		Short:      "sync the resource from test environment and module",
		Long:       "sync the resource from test environment and module, you can choose the synchronize the specified file",
		Example:    "szr",
		PreRun:     cobraAuth.CreateOrChooseSSHUser(&sshUser),
		Run:        run,
	}

	root.AddCommand(cobraAuth.AddUserCmd(addUserFunc))

	if err := root.Execute(); err != nil {
		logger.Fatal("%+v", err)
	}


	//
	//
	//exec, err = sync.exec("cd /home/appl/vip/conf")
	//if err != nil {
	//	log.Fatal("执行命令", err)
	//}
	//
	//log.Println("cd /home/appl/vip/conf", exec)
	//
	//exec, err = sync.exec("ls")
	//if err != nil {
	//	log.Fatal("执行命令", err)
	//}
	//
	//log.Println("ls", exec)
	//
	//exec, err = sync.exec("cat whiteList.json")
	//if err != nil {
	//	log.Fatal("执行命令", err)
	//}
	//
	//log.Println("cat whiteList.json", exec)

}

func run(cmd *cobra.Command, args []string) {
	sync := &Sync{}
	err := sync.connect(sshUser)
	if err != nil {
		log.Fatal(err)
	}

	err = sync.shell()
	if err != nil {
		log.Fatal(err)
	}

	dashboards, servers := sync.dashboard()

	if dashboards != nil && len(dashboards) > 0 {
		prompt := &survey.Select{
			Message: "Select a server for sync",
			Options: dashboards,
		}

		var r string
		_ = survey.AskOne(prompt, &r)

		sync.cd(servers[r])
		sync.cd("cd /home/appl")
		modules := sync.getModules()
		log.Printf("Modules: %+v", modules)

		prompt = &survey.Select{
			Message: "Select a Module for sync",
			Options: modules,
		}

		_ = survey.AskOne(prompt, &r)

		sync.cd("cd " + r)
		sync.cd("cd conf")

		exec, err := sync.exec("ls")
		if err != nil {
			log.Fatal(err)
		}

		log.Println(exec)
	}
	
	sync.close()
}
