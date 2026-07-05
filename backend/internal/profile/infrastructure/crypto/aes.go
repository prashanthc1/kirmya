package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"os"
)

const devProfileKey = "kirmya-dev-profile-key"

// Encryptor handles AES-256-GCM encryption/decryption for sensitive profile fields.
type Encryptor struct {
	key [32]byte
}

func NewEncryptor() *Encryptor {
	src := os.Getenv("PROFILE_ENC_KEY")
	if src == "" {
		src = os.Getenv("JWT_SECRET")
	}
	if src == "" {
		if os.Getenv("APP_ENV") == "production" {
			log.Fatalf("PROFILE_ENC_KEY (or JWT_SECRET) is required in production to encrypt profile fields")
		}
		log.Println("WARNING: PROFILE_ENC_KEY/JWT_SECRET unset; using an insecure dev key for profile field encryption.")
		src = devProfileKey
	}
	return &Encryptor{key: sha256.Sum256([]byte(src))}
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	block, err := aes.NewCipher(e.key[:])
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
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (e *Encryptor) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(e.key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
