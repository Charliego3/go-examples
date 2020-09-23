package main

import (
	"golang.org/x/crypto/ssh"
	"log"
)

func main() {
	sshConfig := &ssh.ClientConfig{
		User: "niuchaolong",
		Auth: []ssh.AuthMethod{
			ssh.Password("Zwykj@niu2019"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "172.16.100.128:2222", sshConfig)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	//session.SendRequest("\n", true, nil)
	session.Run("1")
}
