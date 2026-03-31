package rsa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCreatesKeyFiles(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	if err := Generate(1024); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	privatePath := filepath.Join(tmp, "private_key.pem")
	publicPath := filepath.Join(tmp, "public_key.pem")
	if _, err := os.Stat(privatePath); err != nil {
		t.Fatalf("private key not created: %v", err)
	}
	if _, err := os.Stat(publicPath); err != nil {
		t.Fatalf("public key not created: %v", err)
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()
	if err := Generate(1024); err != nil {
		t.Fatal(err)
	}
	plainText := []byte("yes this is me,do you miss me?")

	b, err := Encrypt(plainText, "./public_key.pem")
	if err != nil {
		t.Fatal(err)
	}

	bs, err := Decrypt(b, "./private_key.pem")
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != string(plainText) {
		t.Fatalf("Decrypt() got = %q, want %q", string(bs), string(plainText))
	}
}

func TestEncryptInvalidPublicKeyPath(t *testing.T) {
	_, err := Encrypt([]byte("hello"), "./not_found_public.pem")
	if err == nil {
		t.Fatalf("expected error for invalid public key path")
	}
}

func TestDecryptInvalidPrivateKeyPath(t *testing.T) {
	_, err := Decrypt([]byte("hello"), "./not_found_private.pem")
	if err == nil {
		t.Fatalf("expected error for invalid private key path")
	}
}
