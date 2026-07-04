package crypto

import "testing"

func TestArgon2Hash_VerifyRoundTrip(t *testing.T) {
	h := NewArgon2Hasher()
	encoded, err := h.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if encoded == "" {
		t.Fatal("expected non-empty hash")
	}

	ok, err := h.Verify("correct horse battery staple", encoded)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatal("expected password to verify")
	}

	bad, err := h.Verify("wrong password", encoded)
	if err != nil {
		t.Fatalf("verify wrong: %v", err)
	}
	if bad {
		t.Fatal("expected wrong password to fail")
	}
}

func TestArgon2Hash_DistinctSalts(t *testing.T) {
	h := NewArgon2Hasher()
	a, _ := h.Hash("same")
	b, _ := h.Hash("same")
	if a == b {
		t.Fatal("expected different hashes due to random salt")
	}
}
