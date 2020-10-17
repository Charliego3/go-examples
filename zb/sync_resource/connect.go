package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/whimthen/temp/zb/auth"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Sync struct {
	session *ssh.Session
	client  *ssh.Client
	stdin   io.WriteCloser
	stdout  io.Reader
	output  bytes.Buffer
	read    int
	prompt  string
}

var (
	ipValidator = regexp.MustCompile("(?m)(?:(?:\\d|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(?:\\d|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])")
	colorMatch  = regexp.MustCompile("\u001B\\[(\\d+;)?\\d+m")
)

func (s *Sync) connect(su auth.SSHUser) error {
	config := &ssh.ClientConfig{
		User:            su.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(su.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(su.Host, su.Port), config)
	if err != nil {
		return err
	}
	s.client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	s.session = session
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}
	// Request pseudo terminal
	if err := session.RequestPty(termType, 40, 200, ssh.TerminalModes{}); err != nil {
		color.Red(err.Error())
		return err
	}
	s.stdin, err = session.StdinPipe()
	if err != nil {
		return err
	}
	//s.session.Stdout = &s.output
	s.stdout, _ = s.session.StdoutPipe()

	return nil
}

func (s *Sync) shell() error {
	if s.session == nil {
		return errors.New("please connect first")
	}

	return s.session.Shell()
}

func (s *Sync) close() {
	if s.session != nil {
		_ = s.session.Close()
	}

	if s.client != nil {
		_ = s.client.Close()
	}
}

func (s *Sync) exec(cmd string) (string, error) {
	return s.getStdout(cmd, "")
}

func (s *Sync) getRemotePort(node string) string {
	spinner.Restart()
	serverPath := fmt.Sprintf("/home/appl/%s/tomcat/conf/server.xml", node)
	content, err := s.getContent(serverPath)
	if err != nil {
		spinner.Stop()
		color.Red(err.Error())
		return "8080"
	}

	type Server struct {
		Service struct {
			Connector []struct {
				Port string `xml:"port,attr"`
			} `xml:"Connector"`
		} `xml:"Service"`
	}

	server := &Server{}
	_ = xml.Unmarshal([]byte(content), server)
	if len(server.Service.Connector) < 1 {
		return "8080"
	}

	return server.Service.Connector[0].Port
}

func (s *Sync) getPrompt(servers map[string]string, ip string) (string, error) {
	return s.getStdout(servers[ip], ip)
}

func (s *Sync) getStdout(cmd, ip string) (string, error) {
	if s.session == nil {
		return "", errors.New("please connect first")
	}

	if strings.HasPrefix(cmd, "cd") {
		tempPrompt := strings.ReplaceAll(cmd, "cd ", "")
		if strings.HasPrefix(s.prompt[1:], "/") {
			tempPrompt = s.prompt[1:len(s.prompt)-1] + "/" + tempPrompt
		}
		s.prompt = ":" + tempPrompt + "$"
	}

	_, err := s.stdin.Write([]byte(cmd + "\r"))
	if err != nil {
		return "", err
	}

	waitingString := ""
	for {
		var buf [65 * 1024]byte
		n, err := s.stdout.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			break
		}
		waitingString += string(buf[:n])
		line := string(colorMatch.ReplaceAll([]byte(waitingString), []byte("")))
		if ip != "" && s.prompt == "" {
			if strings.HasSuffix(strings.TrimSpace(line), ":~$") {
				reader := bufio.NewReader(bytes.NewBufferString(line))
				var line []byte
				for {
					lineBytes, _, err := reader.ReadLine()
					if err == io.EOF {
						break
					}
					line = lineBytes
				}
				s.prompt = string(line)
				return s.prompt, nil
			}
		}
		line = strings.TrimSpace(line)
		if (strings.Count(line, "Opt>") == 3) || (s.prompt != "" && strings.HasSuffix(line, s.prompt)) {
			break
		}
	}
	return waitingString, nil
}

func (s *Sync) dashboard() ([]string, map[string]string) {
	txt, err := s.exec("")
	if err != nil {
		color.Red(err.Error())
		return nil, nil
	}

	var servers []string
	serverMap := make(map[string]string, 0)

	buffer := bytes.NewBufferString(txt)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		if ipValidator.MatchString(line) {
			space := strings.TrimSpace(line)
			splits := strings.Split(space, "|")
			var ls []string
			ls = append(ls, strings.TrimSpace(splits[2]))
			ls = append(ls, strings.TrimSpace(splits[1]))
			ls = append(ls, strings.TrimSpace(splits[3]))
			n := strings.Join(ls, " - ")
			servers = append(servers, n)
			serverMap[n] = ls[0]
		}
	}

	return servers, serverMap
}

func (s *Sync) cd(dir string) string {
	dir, err := s.exec(dir)
	if err != nil {
		color.Red(err.Error())
		return ""
	}
	return dir
}

func (s *Sync) ls() string {
	rtn, err := s.exec("ls")
	if err != nil {
		color.Red(err.Error())
		return ""
	}
	return rtn
}

func (s *Sync) getModules() []string {
	rtn := s.ls()
	var modules []string
	s.clearColor(rtn, func(node string) {
		ext := filepath.Ext(node)
		if node != "" && ext == "" {
			node = strings.ReplaceAll(node, "\r\n", "")
			modules = append(modules, node)
		}
	})
	return modules
}

func (s *Sync) clearColor(rtn string, f func(node string)) {
	buffer := bytes.NewBufferString(rtn)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		line = string(colorMatch.ReplaceAll([]byte(line), []byte("")))
		splits := strings.Split(line, " ")
		for _, module := range splits {
			f(module)
		}
	}
}

func (s *Sync) getConfigs() []string {
	configString := s.ls()

	var configs []string
	s.clearColor(configString, func(node string) {
		ext := filepath.Ext(node)
		if node != "" && ext != "" {
			node = strings.ReplaceAll(node, "\r\n", "")
			configs = append(configs, node)
		}
	})
	return configs
}

func (s *Sync) getContent(filePath string) (string, error) {
	cmd := "cat " + filePath
	content, err := s.exec(cmd)
	if err != nil {
		return "", err
	}

	lastIndex := strings.LastIndex(content, "\r\n")
	if lastIndex <= 0 {
		lastIndex = len(content)
	}
	content = strings.ReplaceAll(content[len(cmd):lastIndex], "\r\n", "\n")
	return content[1:], nil
}
