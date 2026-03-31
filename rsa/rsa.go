package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func Generate(bits int) error {
	if bits <= 0 {
		return errors.New("bits must be greater than 0")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("generate rsa key: %w", err)
	}

	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		return fmt.Errorf("create private key file: %w", err)
	}
	defer func() { _ = privateKeyFile.Close() }()

	privateKeyBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509PrivateKey}
	if err := pem.Encode(privateKeyFile, &privateKeyBlock); err != nil {
		return fmt.Errorf("encode private key pem: %w", err)
	}

	x509PublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}
	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		return fmt.Errorf("create public key file: %w", err)
	}
	defer func() { _ = publicKeyFile.Close() }()

	publicKeyBlock := pem.Block{Type: "PUBLIC KEY", Bytes: x509PublicKey}
	if err := pem.Encode(publicKeyFile, &publicKeyBlock); err != nil {
		return fmt.Errorf("encode public key pem: %w", err)
	}

	return nil
}

func Encrypt(plainText []byte, publicKeyPath string) ([]byte, error) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	publicKeyDecodeBlock, _ := pem.Decode(publicKeyBytes)
	if publicKeyDecodeBlock == nil {
		return nil, errors.New("decode public key pem failed")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyDecodeBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not rsa")
	}

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, fmt.Errorf("encrypt plain text: %w", err)
	}

	return cipherText, nil
}

func Decrypt(cipherText []byte, privateKeyPath string) ([]byte, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, errors.New("decode private key pem failed")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return nil, fmt.Errorf("decrypt cipher text: %w", err)
	}

	return plainText, nil
}
