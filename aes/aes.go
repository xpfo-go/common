package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

var pwdKey = []byte("default aes key.")

func InitAesKey(key []byte) error {
	size := len(key)
	if size != 16 && size != 24 && size != 32 {
		return errors.New("key size must be 16,24,32")
	}
	pwdKey = key
	return nil
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// Encrypt 加密
func Encrypt(data string) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	dataByte := []byte(data)
	encryptBytes := pkcs7Padding(dataByte, blockSize)
	//初始化加密数据接收切片
	cryptBytes := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, pwdKey[:blockSize])
	//执行加密
	blockMode.CryptBlocks(cryptBytes, encryptBytes)
	return cryptBytes, nil
}

// Decrypt 解密
func Decrypt(data []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(pwdKey)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, pwdKey[:blockSize])
	//初始化解密数据接收切片
	cryptBytes := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(cryptBytes, data)
	//去除填充
	cryptBytes, err = pkcs7UnPadding(cryptBytes)
	if err != nil {
		return nil, err
	}
	return cryptBytes, nil
}

// EncryptAndBase64 Aes加密 后 base64
func EncryptAndBase64(data string) (string, error) {
	res, err := Encrypt(data)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(res), nil
}

// DecryptAndBase64 base64解码后 Aes 解密
func DecryptAndBase64(data string) ([]byte, error) {
	dataByte, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return Decrypt(dataByte)
}
