// Package secrets centralizes how production secrets reach the process
// environment. The rest of the codebase reads secrets with os.Getenv (e.g.
// JWT_SECRET, MFA_ENC_KEY, DATABASE_URL); this package is the single seam that
// decides *where those values come from* before anything reads them, so we can
// integrate a real secret manager without touching every call site.
//
// Backends (SECRETS_BACKEND):
//
//	env  (default) — use the ambient environment / .env as-is. No-op.
//	file           — load a secret bundle from SECRETS_FILE (default
//	                 /run/secrets/kirmya.env) and export each key into the
//	                 environment. This is the integration point for every major
//	                 secret manager that delivers secrets as a mounted file:
//	                 Kubernetes Secrets, the AWS Secrets Manager / CSI driver,
//	                 HashiCorp Vault Agent, the External Secrets Operator, and
//	                 Docker secrets all surface secrets as a file or directory.
//
// The file may be either a JSON object ({"JWT_SECRET":"…"}) or dotenv-style
// KEY=VALUE lines. Loading happens once, at startup, before OpenDatabase and
// before the identity module reads JWT_SECRET.
package secrets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Logger is the minimal logging surface Load needs (satisfied by *log.Logger).
type Logger interface{ Printf(format string, v ...any) }

const defaultSecretsFile = "/run/secrets/kirmya.env"

// Load reads SECRETS_BACKEND and applies the corresponding backend, exporting
// any loaded values into the process environment. It is safe to call once at
// startup. With the default backend it does nothing and returns nil.
//
// SECRETS_OVERRIDE controls precedence when a key already exists in the
// environment: "true" (default for the file backend — the secret store is the
// source of truth) overwrites; "false" only fills in unset keys.
func Load(logger Logger) error {
	backend := strings.ToLower(strings.TrimSpace(os.Getenv("SECRETS_BACKEND")))
	switch backend {
	case "", "env":
		return nil
	case "file":
		path := os.Getenv("SECRETS_FILE")
		if path == "" {
			path = defaultSecretsFile
		}
		return loadFile(logger, path, overrideEnabled())
	default:
		return fmt.Errorf("secrets: unknown SECRETS_BACKEND %q (want env|file)", backend)
	}
}

func overrideEnabled() bool {
	// Default true: choosing the file backend means the secret store wins.
	return os.Getenv("SECRETS_OVERRIDE") != "false"
}

func loadFile(logger Logger, path string, override bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("secrets: read %s: %w", path, err)
	}
	values, err := parse(data)
	if err != nil {
		return fmt.Errorf("secrets: parse %s: %w", path, err)
	}

	applied := 0
	for k, v := range values {
		if !override {
			if _, exists := os.LookupEnv(k); exists {
				continue
			}
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("secrets: setenv %s: %w", k, err)
		}
		applied++
	}
	if logger != nil {
		// Never log values — only the count and the source.
		logger.Printf("secrets: loaded %d value(s) from %s (override=%v)", applied, filepath.Clean(path), override)
	}
	return nil
}

// parse accepts either a JSON object of string values or dotenv-style
// KEY=VALUE lines (with optional surrounding quotes and # comments).
func parse(data []byte) (map[string]string, error) {
	trimmed := strings.TrimSpace(string(data))
	if strings.HasPrefix(trimmed, "{") {
		var obj map[string]string
		if err := json.Unmarshal([]byte(trimmed), &obj); err != nil {
			return nil, fmt.Errorf("invalid JSON secret bundle: %w", err)
		}
		return obj, nil
	}
	return parseDotenv(trimmed), nil
}

func parseDotenv(s string) map[string]string {
	out := map[string]string{}
	for _, raw := range strings.Split(s, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		// Strip matching surrounding quotes.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if key != "" {
			out[key] = val
		}
	}
	return out
}

// Require fails fast when any of the named secrets is empty. Intended for use in
// production startup so a missing secret aborts the boot with a clear message
// rather than surfacing as a confusing downstream failure.
func Require(names ...string) error {
	var missing []string
	for _, n := range names {
		if strings.TrimSpace(os.Getenv(n)) == "" {
			missing = append(missing, n)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("secrets: required secret(s) not set: %s", strings.Join(missing, ", "))
	}
	return nil
}
