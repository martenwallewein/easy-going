package main

import (
	"fmt"
	"log"
	"os"
)

const (
	senderCryptKeyName   = "sender_crypt"
	receiverCryptKeyName = "recevier_crypt"
)

func printUsage() {
	fmt.Println("Usage: ./cryptomessages init|send|receive [args] ")
	fmt.Println("       ./cryptomessages init [rsa|hmac] - Initialize key(pair)s depending on the mode")
	fmt.Println("       ./cryptomessages send [rsa|hmac] address message - Send the given message via udp to the address")
	fmt.Println("       ./cryptomessages receive [rsa|hmac] address - Receive a message via udp listening on the address")
}

func main() {
	// You may want to use a package to parse cmd args and opts, e.g. Go's flag package or cobra
	// For this simple case, this one is enough
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}
	command := os.Args[1]
	if command == "receive" && len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}
	if command == "send" && len(os.Args) < 5 {
		printUsage()
		os.Exit(1)
	}
	mode := os.Args[2]

	switch command {
	case "init":
		err := initKeys(mode)
		if err != nil {
			log.Fatal(err)
		}
		break
	case "send":
		err := sendMessage(mode, os.Args[3], os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent message to ", os.Args[3])
		break
	case "receive":
		message, err := receiveMessage(mode, os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Received message ", message)
		break
	}

}

func sendMessage(mode, addr, message string) error {
	if mode == "rsa" {
		if err := loadRSAKeyPairs(); err != nil {
			return err
		}

		return sendMessageRSA(addr, message)
	}

	if mode == "hmac" {
		if err := loadSymmetricKeys(); err != nil {
			return err
		}

		return sendMessageHMAC(addr, message)
	}

	return nil
}

func receiveMessage(mode, addr string) (string, error) {
	if mode == "rsa" {
		if err := loadRSAKeyPairs(); err != nil {
			return "", err
		}

		return receiveMessageRSA(addr)
	}

	if mode == "hmac" {
		if err := loadSymmetricKeys(); err != nil {
			return "", err
		}

		return receiveMessageHMAC(addr)
	}

	return "", nil
}

func initKeys(mode string) error {
	switch mode {
	case "rsa":
		return generateRSAKeyPairs()
	case "hmac":
		return generateSymmetricKeys()
	default:
		return fmt.Errorf("Invalid mode. Please use 'rsa' or 'hmac'")
	}
}
