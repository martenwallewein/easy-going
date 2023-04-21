package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

type RsaKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var (
	// RSA
	senderCryptKeyPair   *RsaKeyPair
	receiverCryptKeyPair *RsaKeyPair
)

func sendMessageRSA(addr, message string) error {

	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()

	// Encrypt the message using the recipient's public key
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, receiverCryptKeyPair.PublicKey, []byte(message))
	if err != nil {
		return fmt.Errorf("Error encrypting message: %v", err)
	}

	// Sign the message using the sender's private key
	hash := sha256.New()
	hash.Write([]byte(message))
	hashedMessage := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, senderCryptKeyPair.PrivateKey, crypto.SHA256, hashedMessage)
	if err != nil {
		return fmt.Errorf("Error signing message: %v", err)
	}

	// Combine the encrypted message and signature
	signedAndEncryptedMessage := append(encryptedMessage, signature...)
	_, err = conn.Write(signedAndEncryptedMessage)
	if err != nil {
		return fmt.Errorf("Error sending signed and encrypted message: %v", err)
	}

	return nil
}

func receiveMessageRSA(addr string) (string, error) {
	// Listen for incoming UDP messages
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return "", fmt.Errorf("Error listening for UDP messages: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return "", fmt.Errorf("Error reading from UDP connection: %v", err)
	}
	buffer = buffer[:n]

	// Split the received message into encrypted message and signature
	encryptedMessage := buffer[:n-256] // Assuming 2048-bit RSA key; adjust size for different key lengths
	signature := buffer[n-256:]
	// Decrypt the message using the recipient's private key
	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, receiverCryptKeyPair.PrivateKey, encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("Error decrypting message: %v", err)
	}

	// Verify the signature using the sender's public key
	hash := sha256.New()
	hash.Write(decryptedMessage)
	hashedMessage := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(senderCryptKeyPair.PublicKey, crypto.SHA256, hashedMessage, signature)
	if err != nil {
		return "", fmt.Errorf("Error verifying message signature: %v", err)
	}

	return string(decryptedMessage), nil
}

func generateRSAKeyPair(name string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("Error generating RSA key pair: %v", err)
	}
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = ioutil.WriteFile(fmt.Sprintf("%s.pem", name), privateKeyPEM, 0600)
	if err != nil {
		return fmt.Errorf("Error saving private key: %v", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("Error marshaling public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = ioutil.WriteFile(fmt.Sprintf("%s.pub.pem", name), publicKeyPEM, 0644)
	if err != nil {
		return fmt.Errorf("Error saving public key: %v", err)
	}

	return nil
}

func generateRSAKeyPairs() error {
	if err := generateRSAKeyPair(senderCryptKeyName); err != nil {
		return err
	}
	if err := generateRSAKeyPair(receiverCryptKeyName); err != nil {
		return err
	}
	return nil
}

func loadRSAKeyPair(name string) (*RsaKeyPair, error) {
	privateKeyPEM, err := ioutil.ReadFile(fmt.Sprintf("%s.pem", name))
	if err != nil {
		log.Fatalf("Error reading private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing private key: %v", err)
	}

	publicKeyPEM, err := ioutil.ReadFile(fmt.Sprintf("%s.pub.pem", name))
	if err != nil {
		return nil, fmt.Errorf("Error reading public key: %v", err)
	}

	block, _ = pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %v", err)
	}

	publicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to convert public key")
	}
	return &RsaKeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}

func loadRSAKeyPairs() error {
	var err error
	senderCryptKeyPair, err = loadRSAKeyPair(senderCryptKeyName)
	if err != nil {
		return err
	}
	receiverCryptKeyPair, err = loadRSAKeyPair(receiverCryptKeyName)
	if err != nil {
		return err
	}
	return nil
}
