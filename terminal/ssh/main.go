package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SSHTerminal struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
}

var isExit bool

func main() {
	Connect("niuchaolong", "Zwykj@niu2019", "172.16.100.128:2222", "130", "cd /home/appl/vip && tailf tomcat/logs/catalina.out")
}

func Connect(user, pwd, addr, env, command string) {
	//publicKey, e := ssh.ParsePublicKey(make([]byte, 11))
	//if e != nil {
	//	fmt.Println(e)
	//	return
	//}
	//hostKeyCallback := ssh.FixedHostKey(publicKey)
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = New(client, env, command)
	if err != nil {
		fmt.Println(err)
	}

}

func (t *SSHTerminal) updateTerminalSize() {
	go func() {
		// SIGWINCH is sent to the process when the window size of the terminal has
		// changed.
		sigwinchCh := make(chan os.Signal, 1)
		signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		for {
			select {
			// The client updated the size of the local PTY. This change needs to occur
			// on the server side PTY as well.
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)
				if err != nil {
					fmt.Printf("Unable to send window-change reqest: %s.", err)
					continue
				}

				// Terminal size has not changed, don't do anything.
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				err = t.Session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					fmt.Printf("Change window size error: %s.", err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight
			}
		}
	}()

}

func (t *SSHTerminal) interactiveSession(env, command string) error {

	defer func() {
		if t.exitMsg == "" {
			fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, err = t.Session.StderrPipe()

	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				//if buf[0] == 3 {
				//	isExit = true
				//	t.Session.Close()
				//	return
				//}

				if buf[0] == 26 {
					isExit = true
					t.Session.Close()
					return
				}
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	//t.Session.Start("1")
	//if err := t.Session.Run("1"); err != nil {
	//	log.Fatal(err)
	//}
	err = t.Session.Shell()
	if err != nil {
		return err
	}
	t.stdin.Write([]byte("\r\n"))
	t.stdin.Write([]byte(env + "\r\n"))
	time.Sleep(time.Millisecond * 500)
	t.stdin.Write([]byte(command))
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func New(client *ssh.Client, env, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer func() {
		if !isExit {
			if err := session.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	s := SSHTerminal{
		Session: session,
	}

	return s.interactiveSession(env, command)
}