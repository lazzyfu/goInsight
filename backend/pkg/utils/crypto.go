package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/lazzyfu/goinsight/internal/global"
)

// 获取 AES KEY（必须 32 字节）
func getAesKey() ([]byte, error) {
	key := global.App.Config.App.SECRET_KEY

	if len(key) != 32 {
		return nil, fmt.Errorf("SECRET_KEY must be exactly 32 bytes, got %d", len(key))
	}

	return []byte(key), nil
}

// AES-256-GCM 加密，输出 Base64
func Encrypt(plaintext string) (string, error) {
	key, err := getAesKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 输出格式：nonce + encrypted + tag
	cipherText := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AES-256-GCM 解密，输入 Base64
func Decrypt(cipherBase64 string) (string, error) {
	key, err := getAesKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
