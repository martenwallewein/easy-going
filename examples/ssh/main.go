package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/kevinburke/ssh_config"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <host> <command> ", os.Args[0])
	}

	host := os.Args[1]
	command := os.Args[2]

	// Load user's SSH config file
	f, err := os.Open(os.ExpandEnv("$HOME/.ssh/config"))
	sshConfig, err := ssh_config.Decode(f)
	if err != nil {
		log.Fatalf("Failed to load SSH config: %s", err)
	}

	// Get the host configuration
	user, err := sshConfig.Get(host, "User")
	if err != nil {
		log.Fatalf("Failed to get ssh User from config: %s", err)
	}

	hostName, err := sshConfig.Get(host, "HostName")
	if err != nil {
		log.Fatalf("Failed to get ssh HostName from config: %s", err)
	}

	port, err := sshConfig.Get(host, "Port")
	if err != nil {
		log.Fatalf("Failed to get ssh Port from config: %s", err)
	}

	// Load private key for authentication
	keyPath := os.ExpandEnv("$HOME/.ssh/id_rsa")
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte("password"))
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}

	// Or without password
	//key, err := ssh.ParsePrivateKey(keyBytes)
	//if err != nil {
	//	log.Fatalf("Failed to parse private key: %s", err)
	//}

	// Configure SSH client
	sshClientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Connect to the host
	address := net.JoinHostPort(hostName, port)
	client, err := ssh.Dial("tcp", address, sshClientConfig)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", host, err)
	}
	defer client.Close()

	fmt.Printf("Connected to %s\n", host)

	// Create an SSH session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %s", err)
	}
	defer session.Close()

	// Run the command on the remote machine
	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to get stdout: %s", err)
	}
	if err := session.Start(command); err != nil {
		log.Fatalf("failed to start command: %s", err)
	}
	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatalf("failed to read stdout: %s", err)
	}
	if err := session.Wait(); err != nil {
		log.Fatalf("command failed: %s", err)
	}
	fmt.Println(string(output))
	client.Close()
}
