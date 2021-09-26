package main

import (
	"github.com/kataras/golog"
	"github.com/spf13/cobra"
)

func main() {
	golog.SetLevel("debug")
	rootCmd := &cobra.Command{
		Use:     "downloader [flags] url",
		Example: "  downloader https://github.com/release/xxx.dmg",
		Version: "v0.0.1",
		RunE:    cmdRun,
	}
	err := rootCmd.Execute()
	if err != nil {
		golog.Fatalf("The command execute catch an error: %+v", err)
	}
}

func cmdRun(cmd *cobra.Command, args []string) error {
	golog.Debug("the cmd running....")
	return nil
}
