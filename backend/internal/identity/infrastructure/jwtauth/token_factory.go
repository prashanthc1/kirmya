// Package jwtauth implements domain.TokenFactory: short-lived HS256 JWT access
// tokens plus opaque (refresh / verification / reset) token generation.
package jwtauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the access-token payload.
type Claims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

// Factory implements domain.TokenFactory.
type Factory struct {
	secret    []byte
	accessTTL time.Duration
}

var (
	genOnce sync.Once
	genKey  []byte
)

// NewFactory builds a token factory. JWT_SECRET configures signing; if unset a
// random per-process key is used (dev only). JWT_ACCESS_TTL (seconds) overrides
// the 15-minute default.
func NewFactory() *Factory {
	return &Factory{secret: signingKey(), accessTTL: accessTTL()}
}

// minSecretLen is the minimum acceptable JWT_SECRET length for HS256 (256-bit).
const minSecretLen = 32

// maxAccessTTL caps the configurable access-token lifetime. Because access
// tokens are stateless and cannot be revoked before expiry (M3), an over-long
// TTL widens the window in which a suspended user / revoked role keeps acting.
const maxAccessTTL = time.Hour

func isProd() bool { return os.Getenv("APP_ENV") == "production" }

func signingKey() []byte {
	if s := strings.TrimSpace(os.Getenv("JWT_SECRET")); s != "" {
		// In production a short/weak secret is fatal — HS256 security depends on it.
		if isProd() && len(s) < minSecretLen {
			log.Fatalf("JWT_SECRET must be at least %d bytes in production (got %d)", minSecretLen, len(s))
		}
		if len(s) < minSecretLen {
			log.Printf("WARNING: JWT_SECRET is shorter than %d bytes; use a stronger secret.", minSecretLen)
		}
		return []byte(s)
	}
	// No secret set: fail closed in production; allow an ephemeral dev key otherwise.
	if isProd() {
		log.Fatalf("JWT_SECRET is required in production")
	}
	genOnce.Do(func() {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			log.Fatalf("could not generate JWT secret: %v", err)
		}
		genKey = []byte(hex.EncodeToString(b))
		log.Println("WARNING: JWT_SECRET unset; using a random per-process key. Set JWT_SECRET in production.")
	})
	return genKey
}

func accessTTL() time.Duration {
	if v := os.Getenv("JWT_ACCESS_TTL"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			ttl := time.Duration(secs) * time.Second
			if ttl > maxAccessTTL {
				log.Printf("WARNING: JWT_ACCESS_TTL=%s exceeds the %s cap; clamping (M3 revocation-lag).", ttl, maxAccessTTL)
				return maxAccessTTL
			}
			return ttl
		}
	}
	return 15 * time.Minute
}

// IssueAccessToken signs a JWT for the user.
func (f *Factory) IssueAccessToken(userID, email string, roles []string) (string, int, error) {
	now := time.Now()
	claims := Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(f.accessTTL)),
			Issuer:    "kirmya",
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(f.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int(f.accessTTL.Seconds()), nil
}

// Parse validates a JWT and returns its claims.
func (f *Factory) Parse(token string) (*Claims, error) {
	var claims Claims
	parsed, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return f.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return &claims, nil
}

// GenerateOpaqueToken returns a random URL-safe token and its SHA-256 hash.
// Only the hash is persisted.
func (f *Factory) GenerateOpaqueToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(b)
	return raw, f.HashOpaqueToken(raw), nil
}

// HashOpaqueToken returns the SHA-256 hex hash used for lookups.
func (f *Factory) HashOpaqueToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
