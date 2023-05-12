package rsa

import "testing"

func TestRSAGenerate(t *testing.T) {
	Generate(1024) // 密钥随机取源长度
}

func TestRsaEncrypt(t *testing.T) {
	plainText := []byte("yes this is me,do you miss me?")

	b := Encrypt(plainText, "./public_key.pem")
	t.Log(string(b))
}

func TestRsaDecrypt(t *testing.T) {
	plainText := []byte("yes this is me,do you miss me?")

	b := Encrypt(plainText, "./public_key.pem")

	bs := Decrypt(b, "./private_key.pem")
	t.Log(len(bs))
	t.Log(string(bs))
}
