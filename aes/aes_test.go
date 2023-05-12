package aes

import "testing"

func TestDecrypt(t *testing.T) {
	if err := InitAesKey([]byte("ase key ase key ase key ase key ")); err != nil {
		t.Error(err)
		return
	}
	bs, err := Encrypt("yes this is me,do you miss me?")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bs))
	ret, err := Decrypt(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(ret))
}

func TestDecryptByAes(t *testing.T) {
	bs, err := EncryptAndBase64("yes this is me,do you miss me?")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bs))
	ret, err := DecryptAndBase64(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(ret))
}
