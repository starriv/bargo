package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
)

// Encryptor 加密解密类
type Encryptor struct {
	Key []byte
	Iv  []byte
}

// NewEncryptor 获得加密实例
func NewEncryptor(key []byte) *Encryptor {
	// 计算key iv
	// key为md5加密后的32位字节
	// iv为key的前16个字节
	md5er := md5.New()
	md5er.Write(key)

	newKey := []byte(hex.EncodeToString(md5er.Sum(nil)))
	iv := newKey[:aes.BlockSize]
	result := &Encryptor{
		Key: newKey,
		Iv:  iv,
	}
	return result
}

// Encrypt 加密
func (e *Encryptor) Encrypt(plaintext []byte) []byte {
	block, _ := aes.NewCipher(e.Key)
	stream := cipher.NewCFBEncrypter(block, e.Iv)
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext
}

// Decrypt 解密
func (e *Encryptor) Decrypt(ciphertext []byte) []byte {
	block, _ := aes.NewCipher(e.Key)
	stream := cipher.NewCFBDecrypter(block, e.Iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}
