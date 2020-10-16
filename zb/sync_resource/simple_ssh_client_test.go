package main

import (
	"bytes"
	c "github.com/helloyi/go-sshclient"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := c.DialWithPasswd("127.0.0.1:22", "nzlong", "niu")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	//bs, err := client.Cmd("cmd1").Cmd("cmd2").Cmd("cmd3").Output()
	//if err != nil {
	//	t.Fatal(err)
	//}

	//t.Logf("%s", bs)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	if err := client.Cmd("which brew").SetStdio(&stdout, &stderr).Run(); err != nil {
		log.Fatal(err)
	}

	// get it
	t.Log("Output:", stdout.String())
	t.Log("Stderr:", stderr.String())
}
