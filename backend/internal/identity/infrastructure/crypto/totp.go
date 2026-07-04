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

	"github.com/pquerna/otp/totp"
)

// devMFAKey is the last-resort key used ONLY outside production. Shipping it as a
// usable production path would encrypt all TOTP secrets under a public constant.
const devMFAKey = "kirmya-dev-mfa-key"

// TOTPService implements domain.TOTP. TOTP secrets are encrypted at rest with
// AES-256-GCM using a key derived from MFA_ENC_KEY (or JWT_SECRET as a dev
// fallback).
type TOTPService struct {
	issuer string
	key    [32]byte
}

func NewTOTPService(issuer string) *TOTPService {
	src := os.Getenv("MFA_ENC_KEY")
	if src == "" {
		src = os.Getenv("JWT_SECRET")
	}
	if src == "" {
		// Fail closed in production: never encrypt MFA secrets under a public key.
		if os.Getenv("APP_ENV") == "production" {
			log.Fatalf("MFA_ENC_KEY (or JWT_SECRET) is required in production to encrypt TOTP secrets")
		}
		log.Println("WARNING: MFA_ENC_KEY/JWT_SECRET unset; using an insecure dev key for TOTP encryption.")
		src = devMFAKey
	}
	return &TOTPService{issuer: issuer, key: sha256.Sum256([]byte(src))}
}

// Generate creates a new TOTP secret, returning the encrypted secret to store
// and the otpauth:// URL (with plaintext secret) for QR enrollment.
func (s *TOTPService) Generate(accountEmail string) (secretEnc, otpauthURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: s.issuer, AccountName: accountEmail})
	if err != nil {
		return "", "", err
	}
	enc, err := s.encrypt(key.Secret())
	if err != nil {
		return "", "", err
	}
	return enc, key.URL(), nil
}

// Validate checks a 6-digit code against the encrypted secret.
func (s *TOTPService) Validate(secretEnc, code string) bool {
	secret, err := s.decrypt(secretEnc)
	if err != nil {
		return false
	}
	return totp.Validate(code, secret)
}

func (s *TOTPService) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key[:])
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

func (s *TOTPService) decrypt(encoded string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.key[:])
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
