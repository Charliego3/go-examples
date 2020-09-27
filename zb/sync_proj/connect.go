package main

import (
	"bytes"
	"errors"
	"github.com/whimthen/temp/zb/auth"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Sync struct {
	session *ssh.Session
	client  *ssh.Client
	stdin   io.WriteCloser
	output  bytes.Buffer
	read    int
}

var (
	ipValidator = regexp.MustCompile("(?m)(?:(?:\\d|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(?:\\d|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])")
	colorMatch = regexp.MustCompile("\u001B\\[(\\d+;)?\\d+m")
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
		log.Println(err)
	}
	s.stdin, err = session.StdinPipe()
	if err != nil {
		return err
	}
	s.session.Stdout = &s.output

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
		s.client.Close()
	}
}

func (s *Sync) exec(cmd string) (string, error) {
	if s.session == nil {
		return "", errors.New("please connect first")
	}

	_, err := s.stdin.Write([]byte(cmd + "\r"))
	if err != nil {
		return "", err
	}

	time.Sleep(time.Millisecond * 500)
	rtn := s.output.String()[s.read+len(cmd)+2:]
	s.read = s.output.Len()
	return rtn, nil
}

func (s *Sync) dashboard() ([]string, map[string]string) {
	txt, err := s.exec("")
	if err != nil {
		log.Fatal(err)
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
			if err != nil {
				log.Fatal("matchString", err)
			}
			space := strings.TrimSpace(line)
			splits := strings.Split(space, "|")
			var ls []string
			for _, split := range splits {
				ls = append(ls, strings.TrimSpace(split))
			}
			n := strings.Join(ls, " - ")
			servers = append(servers, n)
			serverMap[n] = strings.TrimSpace(splits[2])
		}
	}

	return servers, serverMap
}

func (s *Sync) cd(dir string) {
	_, err := s.exec(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Sync) getModules() []string {
	rtn, err := s.exec("ls")
	if err != nil {
		log.Fatal(err)
	}
	var modules []string

	buffer := bytes.NewBufferString(rtn)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		line = string(colorMatch.ReplaceAll([]byte(line), []byte("")))
		splits := strings.Split(line, " ")
		for _, module := range splits {
			ext := filepath.Ext(module)
			if module != "" && ext == "" {
				module = strings.ReplaceAll(module, "\r\n", "")
				modules = append(modules, module)
			}
		}
	}
	return modules
}