/*
 * 文件作用：敏感数据加密/解密工具
 * 负责功能：
 *   - 使用 AES-GCM 加密敏感字段
 *   - 生成/校验加密前缀
 */
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"cli-proxy/internal/config"
)

const encryptedPrefix = "enc:"

// IsEncrypted 判断是否为已加密值
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}

// EncryptString 加密字符串（已加密则直接返回）
func EncryptString(value string) (string, error) {
	if value == "" || IsEncrypted(value) {
		return value, nil
	}
	key, err := getEncryptionKey()
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

	ciphertext := gcm.Seal(nil, nonce, []byte(value), nil)
	payload := append(nonce, ciphertext...)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(payload), nil
}

// DecryptString 解密字符串（非加密值则直接返回）
func DecryptString(value string) (string, error) {
	if value == "" || !IsEncrypted(value) {
		return value, nil
	}
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, encryptedPrefix))
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

	if len(raw) < gcm.NonceSize() {
		return "", errors.New("密文长度无效")
	}

	nonce := raw[:gcm.NonceSize()]
	ciphertext := raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func getEncryptionKey() ([]byte, error) {
	secret := ""
	if config.Cfg != nil {
		secret = config.Cfg.Security.DataKey
		if secret == "" {
			secret = config.Cfg.JWT.Secret
		}
	}
	if secret == "" {
		return nil, errors.New("加密密钥为空")
	}
	sum := sha256.Sum256([]byte(secret))
	return sum[:], nil
}
