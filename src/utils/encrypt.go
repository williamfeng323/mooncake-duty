package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt encrypt the string with pre-configured key
func Encrypt(str string) ([]byte, error) {
	key := []byte(GetConf().EncryptKey)
	text := []byte(str)
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

// Decrypt dencrypt the string with pre-configured key
func Decrypt(str string) ([]byte, error) {
	key := []byte(GetConf().EncryptKey)
	text := []byte(str)
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := text[:nonceSize], text[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
