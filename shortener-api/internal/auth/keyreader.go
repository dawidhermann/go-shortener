package auth

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type KeyManager struct {
	privateKey string
}

func NewKeyReader(privateKeyPath string) (KeyManager, error) {
	privateKey, err := readKey(privateKeyPath)
	if err != nil {
		return KeyManager{}, fmt.Errorf("failed to read private key: %w", err)
	}
	return KeyManager{
		privateKey: privateKey,
	}, nil
}

func (reader KeyManager) PrivateKeyPem() string {
	return reader.privateKey
}

func (reader KeyManager) PublicKeyPem() string {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(reader.PrivateKeyPem()))
	// fmt.Println(privateKey)
	fmt.Println(err)
	asn1Bytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	// println(err)
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var b bytes.Buffer

	// Write the public key to the public key file.
	if err := pem.Encode(&b, &publicBlock); err != nil {
		return ""
	}
	println(b.String())
	return b.String()
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
