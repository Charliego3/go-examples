package main

import (
	"os/exec"
	"testing"
)

func TestOpen(t *testing.T) {
    cmd := exec.Command("neovide")
    cmd.Dir = "/Users/charlie/dev/java/zb_main"
    cmd.Start()
}
