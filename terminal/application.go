package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"github.com/whimthen/temp/terminal/ssh"
	"os"
	"path/filepath"
)

const (
	sshAddress  = "TERMINAL_SSH_ADDRESS"
	sshUser     = "TERMINAL_SSH_USER"
	sshPassword = "TERMINAL_SSH_PASSWORD"
)

var (
	logLine     int
	command     string
	logKeyword  string
	environment string
	module      string
	user        string
	pwd         string
	addr        string
	a int
	b int
	c int
	grep bool
)

func main() {
	root := &cobra.Command{
		Use:     "terminal",
		Aliases: []string{"jump"},
		Run:     runCommand,
	}

	bindFlags(root)
	root.Flags().StringVarP(&command, "cmd", "c", "", "Shell Command")
	root.AddCommand(loggerCommand())

	if err := root.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func runCommand(cmd *cobra.Command, args []string) {
	checkEnvAndModule(cmd, false)

	if !filepath.IsAbs(module) {
		module = filepath.Join("/home/appl/", module)
	}

	var tempC string
	if command != "" {
		tempC = " && " + command
	}

	command = fmt.Sprintf("cd %s%s\n", module, tempC)
	ssh.Connect(getEnv(sshUser), getEnv(sshPassword), getEnv(sshAddress), environment, command)
}

func loggerCommand() *cobra.Command {
	lc := &cobra.Command{
		Use:     "logger",
		Aliases: []string{"log"},
		Run:     showLogger,
	}

	bindFlags(lc)
	lc.Flags().StringVarP(&command, "cmd", "c", "", "Show logger with terminal")
	lc.Flags().StringVarP(&logKeyword, "keyword", "k", "", "Use keyword filtering")
	lc.Flags().IntVarP(&logLine, "line", "l", 0, "Print log and output continuously")
	lc.Flags().IntVarP(&a, "after", "A",  0, "Except for the column that conforms to the template style, the content after the row is displayed.")
	lc.Flags().IntVarP(&b, "before", "B",  0, "Except for the line that matches the style, the content before that line is displayed.")
	lc.Flags().IntVarP(&c, "context", "C",  0, "In addition to the line that meets the style, the content before and after that line is displayed.")
	lc.Flags().BoolVarP(&grep, "grep", "g", false, "Only grep the log file")
	return lc
}

func showLogger(cmd *cobra.Command, args []string) {
	checkEnvAndModule(cmd, true)

	if !filepath.IsAbs(module) {
		module = filepath.Join("/home/appl/", module)
	}

	if command == "" {
		command = "tailf tomcat/logs/catalina.out"
	}

	if logKeyword != "" {
		command = fmt.Sprintf("tailf tomcat/logs/catalina.out | grep '%s'", logKeyword)
		if a != 0 {
			command = fmt.Sprintf("grep -A %d '%s' tomcat/logs/catalina.out", a, logKeyword)
		} else if b != 0 {
			command = fmt.Sprintf("grep -B %d '%s' tomcat/logs/catalina.out", b, logKeyword)
		} else if c != 0 {
			command = fmt.Sprintf("grep -C %d '%s' tomcat/logs/catalina.out", c, logKeyword)
		}

		if grep {
			command = fmt.Sprintf("grep '%s' tomcat/logs/catalina.out", logKeyword)
		}
	} else if logLine != 0 {
		command = fmt.Sprintf("tail -%df tomcat/logs/catalina.out", logLine)
	}

	command = fmt.Sprintf("cd %s && %s\n", module, command)

	ssh.Connect(getEnv(sshUser), getEnv(sshPassword), getEnv(sshAddress), environment, command)
}

func checkEnvAndModule(cmd *cobra.Command, checkModule bool) {
	if environment == "" {
		logger.Error("No environment specified, You can use `--env string or -e string` to specify a specific test environment\n")
		_ = cmd.Help()
		os.Exit(1)
	}

	if checkModule && module == "" {
		logger.Error("No module name or path specified, You can use `--module string or -m string`\n")
		_ = cmd.Help()
		os.Exit(1)
	}
}

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&environment, "env", "e", "", "Part of IP, host name, remarks")
	cmd.Flags().StringVarP(&module, "module", "m", "", "Module name")
	cmd.Flags().StringVarP(&user, "user", "u", "", "SSH user (default: `$TERMINAL_SSH_USER`)")
	cmd.Flags().StringVarP(&pwd, "pwd", "p", "", "SSH password (default: `$TERMINAL_SSH_PASSWORD`)")
	cmd.Flags().StringVarP(&addr, "addr", "a", "", "SSH address (default: `$TERMINAL_SSH_ADDRESS`)")
}

func getEnv(key string) string {
	var commandLineP string
	switch key {
	case sshUser:
		if user != "" {
			return user
		}
		commandLineP = "user"
	case sshPassword:
		if pwd != "" {
			return pwd
		}
		commandLineP = "pwd"
	case sshAddress:
		if addr != "" {
			return addr
		}
		commandLineP = "addr"
	}
	env := os.Getenv(key)
	if env == "" {
		logger.Fatal("Please configure in environment variables `%s` from file `~/.bash_profile or ~/.zshrc` or carry this parameter in the command line: `--%s string or -%s string`", key, commandLineP, commandLineP[:1])
	}
	return env
}