package application

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

const encryptionPrefix = "enc:"

// encryptAESGCM encrypts plaintext using AES-GCM with a key derived from secretKey.
// The output is prefixed with "enc:" and base64 encoded.
func encryptAESGCM(plaintext string, secretKey string) (string, error) {
	if secretKey == "" {
		return plaintext, nil
	}

	// Derive 32-byte key using SHA-256
	key := sha256.Sum256([]byte(secretKey))

	block, err := aes.NewCipher(key[:])
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

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptionPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptAESGCM decrypts a base64 encoded string using AES-GCM.
// If the string does not start with "enc:", it is returned as-is (for backward compatibility).
func decryptAESGCM(ciphertextStr string, secretKey string) (string, error) {
	if secretKey == "" || !strings.HasPrefix(ciphertextStr, encryptionPrefix) {
		return ciphertextStr, nil
	}

	rawCiphertext := strings.TrimPrefix(ciphertextStr, encryptionPrefix)
	data, err := base64.StdEncoding.DecodeString(rawCiphertext)
	if err != nil {
		// Fallback to original string if not valid base64
		return ciphertextStr, nil
	}

	key := sha256.Sum256([]byte(secretKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return ciphertextStr, nil // Not enough bytes for nonce
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Return original string if decryption fails (safeguard for legacy/corrupt)
		return ciphertextStr, nil
	}

	return string(plaintext), nil
}
