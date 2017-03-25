package bargo

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
	h := md5.New()
	h.Write(key)
	newKey := []byte(hex.EncodeToString(h.Sum(nil)))
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
	ciphertext := make([]byte, len(plaintext))

	stream := cipher.NewCFBEncrypter(block, e.Iv)
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext
}

// Decrypt 解密
func (e *Encryptor) Decrypt(ciphertext []byte) []byte {
	block, _ := aes.NewCipher(e.Key)

	stream := cipher.NewCFBDecrypter(block, e.Iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}