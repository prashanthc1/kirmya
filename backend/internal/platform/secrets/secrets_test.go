package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_JSON(t *testing.T) {
	got, err := parse([]byte(`{"JWT_SECRET":"abc","DATABASE_URL":"postgres://x"}`))
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if got["JWT_SECRET"] != "abc" || got["DATABASE_URL"] != "postgres://x" {
		t.Fatalf("unexpected parse result: %+v", got)
	}
}

func TestParse_Dotenv(t *testing.T) {
	in := `
# a comment
export JWT_SECRET="quoted value"
MFA_ENC_KEY='single'
EMPTY=
NO_EQUALS_LINE
DATABASE_URL=postgres://u:p@h/db?sslmode=disable
`
	got, err := parse([]byte(in))
	if err != nil {
		t.Fatalf("parse dotenv: %v", err)
	}
	if got["JWT_SECRET"] != "quoted value" {
		t.Fatalf("JWT_SECRET = %q, want %q", got["JWT_SECRET"], "quoted value")
	}
	if got["MFA_ENC_KEY"] != "single" {
		t.Fatalf("MFA_ENC_KEY = %q, want single", got["MFA_ENC_KEY"])
	}
	if got["DATABASE_URL"] != "postgres://u:p@h/db?sslmode=disable" {
		t.Fatalf("DATABASE_URL = %q", got["DATABASE_URL"])
	}
	if _, ok := got["NO_EQUALS_LINE"]; ok {
		t.Fatal("line without '=' should be ignored")
	}
	if v, ok := got["EMPTY"]; !ok || v != "" {
		t.Fatalf("EMPTY should map to empty string, got ok=%v v=%q", ok, v)
	}
}

func TestLoad_FileBackend_OverrideAndFill(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.env")
	if err := os.WriteFile(path, []byte("JWT_SECRET=from-file\nNEW_KEY=created\n"), 0o600); err != nil {
		t.Fatalf("write secret file: %v", err)
	}

	t.Setenv("SECRETS_BACKEND", "file")
	t.Setenv("SECRETS_FILE", path)

	// Override mode (default): file wins over an existing env value.
	t.Setenv("JWT_SECRET", "from-env")
	t.Setenv("SECRETS_OVERRIDE", "true")
	if err := Load(nil); err != nil {
		t.Fatalf("load (override): %v", err)
	}
	if os.Getenv("JWT_SECRET") != "from-file" {
		t.Fatalf("override mode should let the file win, got %q", os.Getenv("JWT_SECRET"))
	}
	if os.Getenv("NEW_KEY") != "created" {
		t.Fatalf("expected NEW_KEY to be created, got %q", os.Getenv("NEW_KEY"))
	}

	// Fill-only mode: existing env value is preserved.
	t.Setenv("JWT_SECRET", "from-env")
	t.Setenv("SECRETS_OVERRIDE", "false")
	if err := Load(nil); err != nil {
		t.Fatalf("load (fill): %v", err)
	}
	if os.Getenv("JWT_SECRET") != "from-env" {
		t.Fatalf("fill mode should preserve env, got %q", os.Getenv("JWT_SECRET"))
	}
}

func TestLoad_DefaultBackendIsNoop(t *testing.T) {
	t.Setenv("SECRETS_BACKEND", "")
	if err := Load(nil); err != nil {
		t.Fatalf("default backend should be a no-op, got %v", err)
	}
}

func TestLoad_UnknownBackend(t *testing.T) {
	t.Setenv("SECRETS_BACKEND", "vault-direct")
	if err := Load(nil); err == nil {
		t.Fatal("expected error for unknown backend")
	}
}

func TestRequire(t *testing.T) {
	t.Setenv("PRESENT_SECRET", "x")
	t.Setenv("ABSENT_SECRET", "")
	if err := Require("PRESENT_SECRET"); err != nil {
		t.Fatalf("present secret should pass: %v", err)
	}
	if err := Require("PRESENT_SECRET", "ABSENT_SECRET"); err == nil {
		t.Fatal("expected error when a required secret is absent")
	}
}
