package sshClient

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"sync"
)

// Firebase struct
type SSH struct {
	sync.Mutex
	key     []byte
	user    string
	server  string
	config  *ssh.ClientConfig
	hostKey ssh.PublicKey
	err     error
}

// cluster-key "./cluster-key"
func InitSSH(user string, server, clusterKey string) *SSH {
	s := &SSH{}
	s.user = user
	s.server = server

	s.key, s.err = ioutil.ReadFile(clusterKey)
	if s.err != nil {
		log.Fatalf("Unable to read private key: %v", s.err)
	}

	signer, err := ssh.ParsePrivateKey(s.key)
	if err != nil {
		log.Fatalf("Unable to parse private key: %v", err)
	}

	s.config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			// ssh.Password("yourpassword"),
		},
		HostKeyCallback: ssh.FixedHostKey(s.hostKey),
	}

	// Ignore key validation
	s.config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return s
}

func (s *SSH) Exec(command string, results chan string) {

	client, err := ssh.Dial("tcp", s.server, s.config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	results <- b.String()
}

// Simple Example
//func main() {
//
//	results := make(chan string, 0)
//
//	user := "root"
//	server := "0.0.0.0:9822"
//	command := "systemctl start docker && /gopath/bin/kind create cluster --config /root/kind.yaml"
//
//	go Exec(user, server, command, results)
//
//	fmt.Println("Waiting... <-results")
//	fmt.Println(<-results)
//	fmt.Println("Done building")
//
//	cmd2 := "/usr/local/bin/kubectl get po --all-namespaces"
//	go Exec(user, server, cmd2, results)
//	fmt.Println("running: /usr/local/bin/kubectl get po --all-namespaces")
//
//	fmt.Println(<-results)
//	fmt.Println("done kubectl")
//
//	close(results)
//}
