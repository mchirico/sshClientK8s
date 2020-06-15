package sshClient

import (
	"fmt"
	"testing"
)

func TestInitSSH(t *testing.T) {
	s := InitSSH("root", "0.0.0.0:9822", "../../cluster-key")
	fmt.Println(s.err)
}

func TestExec(t *testing.T) {
	results := make(chan string, 0)
	s := InitSSH("root", "0.0.0.0:9822", "../../cluster-key")
	command := "uptime"
	go s.Exec(command, results)
	fmt.Println(<-results)
	close(results)

}
