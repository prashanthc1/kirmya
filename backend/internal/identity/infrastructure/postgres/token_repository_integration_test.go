//go:build integration

package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"workspace-app/internal/identity/domain"
	"workspace-app/internal/platform/tx"
	"workspace-app/internal/testsupport"
)

var errForcedRollback = errors.New("forced rollback")

func verificationExpiry() time.Time { return time.Now().Add(24 * time.Hour) }

// TestVerificationRepository_StoreEmailToken_JoinsAmbientTx is the direct
// regression guard for the registration foreign-key failure
// (email_verification_tokens_user_id_fkey). Register creates the user and
// stores its verification token inside a single transaction. If the
// VerificationRepository executes on the pooled connection instead of the
// active *sql.Tx, the just-inserted (uncommitted) user row is invisible to it
// and the token insert is rejected by the FK. This test creates a user and
// stores its token within the same RunInTx and asserts both commit atomically;
// it fails with the FK violation if VerificationRepository is not tx-aware.
func TestVerificationRepository_StoreEmailToken_JoinsAmbientTx(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	users := NewUserRepository(db)
	verif := NewVerificationRepository(db)
	mgr := tx.NewTxManager(db)
	ctx := context.Background()

	u := &domain.User{Email: "reg@cb.test", PasswordHash: "hash", FullName: "Reg User", Status: domain.StatusActive}

	err := mgr.RunInTx(ctx, func(ctx context.Context) error {
		if err := users.Create(ctx, u); err != nil {
			return err
		}
		// Same transaction: the user row exists but is not yet committed. The
		// token write must run on the same *sql.Tx to satisfy the FK.
		return verif.StoreEmailToken(ctx, u.ID, "token-hash", verificationExpiry())
	})
	if err != nil {
		t.Fatalf("register-in-tx should commit user + verification token atomically: %v", err)
	}

	// The token is consumable after commit and maps back to the created user.
	gotUserID, err := verif.ConsumeEmailToken(ctx, "token-hash")
	if err != nil {
		t.Fatalf("consume email token: %v", err)
	}
	if gotUserID != u.ID {
		t.Fatalf("token mapped to user %q, want %q", gotUserID, u.ID)
	}
}

// TestVerificationRepository_RollbackDropsToken confirms the token write
// participates in the transaction on the failure path too: if the enclosing
// transaction rolls back, no token (and no user) is left behind.
func TestVerificationRepository_RollbackDropsToken(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	users := NewUserRepository(db)
	verif := NewVerificationRepository(db)
	mgr := tx.NewTxManager(db)
	ctx := context.Background()

	u := &domain.User{Email: "rollback@cb.test", PasswordHash: "hash", FullName: "Rollback User", Status: domain.StatusActive}

	sentinel := errForcedRollback
	err := mgr.RunInTx(ctx, func(ctx context.Context) error {
		if err := users.Create(ctx, u); err != nil {
			return err
		}
		if err := verif.StoreEmailToken(ctx, u.ID, "doomed-hash", verificationExpiry()); err != nil {
			return err
		}
		return sentinel // force rollback
	})
	if err != sentinel {
		t.Fatalf("expected forced rollback error, got %v", err)
	}

	// Nothing committed: the token is gone.
	if _, err := verif.ConsumeEmailToken(ctx, "doomed-hash"); err != domain.ErrTokenNotFound {
		t.Fatalf("expected ErrTokenNotFound after rollback, got %v", err)
	}
}
