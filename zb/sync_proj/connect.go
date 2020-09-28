package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/whimthen/temp/zb/auth"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Sync struct {
	session *ssh.Session
	client  *ssh.Client
	stdin   io.WriteCloser
	stdout  io.Reader
	output  bytes.Buffer
	read    int
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
		log.Println(err)
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

	//time.Sleep(time.Millisecond * 500)
	//rtn := s.output.String()[s.read+len(cmd)+2:]
	//s.read = s.output.Len()


	//_, err = s.stdin.Write([]byte("\r"))
	//time.Sleep(time.Millisecond * 500)
	//var bs []byte
	//read, err := s.stdout.Read(bs)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//log.Println("Read:", read)
	//
	//end := string(colorMatch.ReplaceAll([]byte(s.cd("")), []byte("")))
	//c := make(chan string)
	//
	//go func() {
	//	var (
	//		buf [65 * 1024]byte
	//		t   int
	//	)
	//	waitingString := ""
	//	for {
	//		n, err := s.stdout.Read(buf[t:]) //this reads the ssh terminal
	//		if err != nil && err != io.EOF{
	//			fmt.Println(err)
	//			break
	//		}
	//		if err == io.EOF || n == 0 {
	//			c <- string(buf[:t])
	//			t = 0
	//			break
	//		}
	//		t += n
	//		waitingString += string(buf[:n])
	//	}
	//
	//	//for {
	//	//	//buf := make([]byte, 1000)
	//	//
	//	//	n, err := stdout.Read(buf[t:])
	//	//	if err != nil {
	//	//		fmt.Println(err.Error())
	//	//		close(c)
	//	//		return
	//	//	}
	//	//	t += n
	//	//	result := string(buf[:t])
	//	//	if strings.HasSuffix(result, end) {
	//	//		c <- string(buf[:t])
	//	//		t = 0
	//	//	}
	//	//}
	//}()
	//
	//str := <-c
	//log.Println(str)

	//var (
	//	buf [65 * 1024]byte
	//	t   int
	//)
	//waitingString := ""
	//for {
	//	n, err := s.stdout.Read(buf[t:]) //this reads the ssh terminal
	//	if err != nil && err != io.EOF{
	//		fmt.Println(err)
	//		break
	//	}
	//	if err == io.EOF {
	//		//c <- string(buf[:t])
	//		t = 0
	//		break
	//	}
	//	t += n
	//	waitingString += string(buf[:n])
	//	if n < t {
	//		break
	//	}
	//}

	//buf := make([]byte, 1000)
	//n, err := s.stdout.Read(buf) //this reads the ssh terminal
	//waitingString := ""
	//if err == nil {
	//	println(fmt.Sprintf("%s", buf[:n]))
	//	waitingString = string(buf[:n])
	//}
	//for err == nil {
	//	// this loop will not end!!
	//	n, err = s.stdout.Read(buf)
	//	waitingString += string(buf[:n])
	//	println(fmt.Sprintf("%s", buf[:n]))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, s.stdout); err != nil {
		log.Fatalf("reading failed: %s", err)
	}

	log.Println(buf.String())

	return buf.String(), nil
}

func (s *Sync) getStdout(cmd, end string) (string, error) {
	if s.session == nil {
		return "", errors.New("please connect first")
	}
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		wg.Add(1)
		_, err := s.stdin.Write([]byte(cmd + "\n"))
		if err != nil {
			//return "", err
			log.Panic(err)
		}
		wg.Wait()
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		stdout, err := s.session.StdoutPipe()
		if err != nil {
			log.Panic(err)
		}
		for {
			n, err := stdout.Read(buf[t:])
			if err != nil {
				fmt.Println(err.Error())
				close(out)
				return
			}
			t += n
			result := string(buf[:t])
			if strings.HasSuffix(result, end) {
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()

	return <-out, nil
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

func (s *Sync) cd(dir string) string {
	dir, err := s.exec(dir)
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func (s *Sync) ls() string {
	rtn, err := s.exec("ls")
	if err != nil {
		log.Fatal(err)
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

func (s *Sync) getContent(fn string, confDir string) (string, error) {
	content, err := s.exec("cat " + fn)
	//content, err := s.getStdout("cat "+fn, confDir)
	if err != nil {
		return "", err
	}

	lastIndex := strings.LastIndex(content, "\r\n")
	if lastIndex <= 0 {
		lastIndex = len(content)
	}
	return content[:lastIndex], nil
}
