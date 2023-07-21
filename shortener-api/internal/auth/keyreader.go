package auth

import (
	"fmt"
	"io"
	"os"
)

type KeyManager struct {
	privateKey string
	publicKey  string
}

func NewKeyReader(privateKeyPath string, publicKeyPath string) (KeyManager, error) {
	privateKey, err := readKey(privateKeyPath)
	if err != nil {
		return KeyManager{}, fmt.Errorf("failed to read private key: %w", err)
	}
	publicKey, err := readKey(publicKeyPath)
	if err != nil {
		return KeyManager{}, fmt.Errorf("failed to read private key: %w", err)
	}
	return KeyManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (reader KeyManager) PrivateKeyPem() string {
	return reader.privateKey
}

func (reader KeyManager) PublicKeyPem() string {
	return reader.publicKey
}

func readKey(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
