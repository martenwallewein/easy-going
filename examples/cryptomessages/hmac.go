package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)

var (
	// HMAC
	cryptKey []byte
)

func generateSymmetricKeys() error {
	if err := generateSymmetricKey(senderCryptKeyName); err != nil {
		return err
	}
	return nil
}

func generateSymmetricKey(name string) error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return fmt.Errorf("Error generating symmetric key: %v", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.sym", name), key, 0600)
	if err != nil {
		return fmt.Errorf("Error saving symmetric key: %v", err)
	}
	return nil
}

func loadSymmetricKeys() error {
	var err error
	cryptKey, err = loadSymmetricKey(senderCryptKeyName)
	if err != nil {
		return err
	}
	return nil
}

func loadSymmetricKey(name string) ([]byte, error) {

	key, err := ioutil.ReadFile(fmt.Sprintf("%s.sym", name))
	if err != nil {
		return nil, fmt.Errorf("Error reading symmetric key: %v", err)
	}
	return key, nil

}

func sendMessageHMAC(addr, message string) error {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()

	// Encrypt data with AES-GCM
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return fmt.Errorf("Error creating AES cipher: %v", err)
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("Error generating nonce: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("Error creating GCM: %v", err)
	}
	encryptedMessage := aesgcm.Seal(nil, nonce, []byte(message), nil)
	encryptedMessage = append(nonce, encryptedMessage...)

	// Sign data with HMAC-SHA256
	mac := hmac.New(sha256.New, cryptKey)
	_, err = mac.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("Error signing message: %v", err)
	}
	signature := mac.Sum(nil)

	// Combine the encrypted message and signature
	signedAndEncryptedMessage := append(encryptedMessage, signature...)
	_, err = conn.Write(signedAndEncryptedMessage)
	if err != nil {
		return fmt.Errorf("Error sending signed and encrypted message: %v", err)
	}

	return nil
}

func receiveMessageHMAC(addr string) (string, error) {
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
	encryptedMessage := buffer[:n-32] // HMAC size
	signature := buffer[n-32:]

	nonce, ciphertext := encryptedMessage[:12], encryptedMessage[12:]
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return "", fmt.Errorf("Error creating AES cipher for decryption: %v", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("Error creating GCM for decryption: %v", err)
	}
	decryptedMessage, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("Error decrypting data: %v", err)
	}

	mac2 := hmac.New(sha256.New, cryptKey)
	_, err = mac2.Write(decryptedMessage)
	if err != nil {
		return "", fmt.Errorf("Error signing message for verification: %v", err)
	}
	expectedSignature := mac2.Sum(nil)
	if !hmac.Equal(signature, expectedSignature) {
		return "", fmt.Errorf("HMAC signature verified.")
	}

	return string(decryptedMessage), nil
}
