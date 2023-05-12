package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"runtime"
)

func Generate(bits int) {
	//get current path
	_, currentPath, _, _ := runtime.Caller(0)
	currentPath = filepath.Dir(currentPath)

	//----------------------------------------------private key

	// GenerateKey generates an RSA keypair of the given bit size using the
	// random source random (for example, crypto/rand.Reader).
	privateKey, _ := rsa.GenerateKey(rand.Reader, bits)

	//serialize private key to ASN.1 der by x509.MarshalPKCS1PrivateKey
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

	//encode x509 to pem and save to file
	//1. create private file
	privateKeyFile, err := os.Create(currentPath + "/private_key.pem")
	if err != nil {
		panic(err)
	}
	defer func() { _ = privateKeyFile.Close() }()
	//2. new a pem block struct object
	privateKeyBlock := pem.Block{
		Type:    "RSA Private Key",
		Headers: nil,
		Bytes:   x509PrivateKey,
	}
	//3. save to file
	_ = pem.Encode(privateKeyFile, &privateKeyBlock)

	//----------------------------------------------public key

	//get public key
	publicKey := privateKey.PublicKey
	//serialize public key to ASN.1 der by x509.MarshalPKCS1PublicKey
	x509PublicKey, _ := x509.MarshalPKIXPublicKey(&publicKey)

	//encode x509 to pem and save to file
	//1. create public key file
	publicKeyFile, err := os.Create(currentPath + "/public_key.pem")
	if err != nil {
		panic(err)
	}
	defer func() { _ = publicKeyFile.Close() }()

	//2. new a pem block struct object
	publicKeyBlock := pem.Block{
		Type:    "RSA Public Key",
		Headers: nil,
		Bytes:   x509PublicKey,
	}

	//3. save to file
	_ = pem.Encode(publicKeyFile, &publicKeyBlock)
}

// Encrypt publicKeyPath 传错了会panic
func Encrypt(plainText []byte, publicKeyPath string) []byte {
	publicKeyFile, _ := os.Open(publicKeyPath)

	defer func() { _ = publicKeyFile.Close() }()

	publicKeyFileInfo, _ := publicKeyFile.Stat()

	//1. make size
	buf := make([]byte, publicKeyFileInfo.Size())
	//2. read file to buf
	_, _ = publicKeyFile.Read(buf)
	//3. decode pem
	publicKeyDecodeBlock, _ := pem.Decode(buf)
	//4. x509 decode
	publicKeyInterface, _ := x509.ParsePKIXPublicKey(publicKeyDecodeBlock.Bytes)

	//assert
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//encrypt plainText
	cipherText, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)

	return cipherText
}

// Decrypt privateKeyPath 传得不对会panic 解密出问题会返回 空[]byte
func Decrypt(cipherText []byte, privateKeyPath string) []byte {

	privateKeyFile, _ := os.Open(privateKeyPath)

	defer func() { _ = privateKeyFile.Close() }()

	privateKeyInfo, _ := privateKeyFile.Stat()
	buf := make([]byte, privateKeyInfo.Size())
	_, _ = privateKeyFile.Read(buf)
	//pem decode
	privateKeyBlock, _ := pem.Decode(buf)
	//X509 decode
	privateKey, _ := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)

	return plainText
}
