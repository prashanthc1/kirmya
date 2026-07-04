package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"workspace-app/internal/identity/domain"
)

// --- in-memory UserRepository ---

type fakeUsers struct {
	mu    sync.Mutex
	seq   int
	byID  map[string]*domain.User
	roles map[string][]string
}

func newFakeUsers() *fakeUsers {
	return &fakeUsers{byID: map[string]*domain.User{}, roles: map[string][]string{}}
}

func (f *fakeUsers) Create(_ context.Context, u *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, e := range f.byID {
		if equalFold(e.Email, u.Email) {
			return domain.ErrEmailTaken
		}
	}
	f.seq++
	u.ID = fmt.Sprintf("user-%d", f.seq)
	u.Version = 1
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	clone := *u
	f.byID[u.ID] = &clone
	return nil
}

func (f *fakeUsers) GetByID(_ context.Context, id string) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	out := *u
	out.Roles = append([]string(nil), f.roles[id]...)
	return &out, nil
}

func (f *fakeUsers) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.byID {
		if equalFold(u.Email, email) {
			out := *u
			out.Roles = append([]string(nil), f.roles[u.ID]...)
			return &out, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (f *fakeUsers) Update(_ context.Context, u *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur, ok := f.byID[u.ID]
	if !ok {
		return domain.ErrUserNotFound
	}
	if cur.Version != u.Version {
		return domain.ErrOptimisticLock
	}
	cur.FullName = u.FullName
	cur.Status = u.Status
	cur.EmailVerified = u.EmailVerified
	cur.MFAEnabled = u.MFAEnabled
	cur.Version++
	u.Version = cur.Version
	return nil
}

func (f *fakeUsers) AssignRole(_ context.Context, userID, role string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.roles[userID] = append(f.roles[userID], role)
	return nil
}

func (f *fakeUsers) RemoveRole(_ context.Context, userID, role string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	kept := []string{}
	for _, r := range f.roles[userID] {
		if r != role {
			kept = append(kept, r)
		}
	}
	f.roles[userID] = kept
	return nil
}

func (f *fakeUsers) GetRoles(_ context.Context, userID string) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.roles[userID]...), nil
}

func (f *fakeUsers) SetEmailVerified(_ context.Context, userID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[userID]; ok {
		u.EmailVerified = true
	}
	return nil
}

func (f *fakeUsers) SetPasswordHash(_ context.Context, userID, hash string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[userID]; ok {
		u.PasswordHash = hash
	}
	return nil
}

func (f *fakeUsers) SetMFAEnabled(_ context.Context, userID string, enabled bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.byID[userID]; ok {
		u.MFAEnabled = enabled
	}
	return nil
}

func (f *fakeUsers) UpdateLastLogin(_ context.Context, userID string) error { return nil }

func (f *fakeUsers) Search(_ context.Context, query string, limit int) ([]domain.DirectoryEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := []domain.DirectoryEntry{}
	for _, u := range f.byID {
		if equalFold(u.Email, query) || u.FullName == query {
			out = append(out, domain.DirectoryEntry{ID: u.ID, FullName: u.FullName, Email: u.Email})
		}
	}
	return out, nil
}

func (f *fakeUsers) GetDirectory(_ context.Context, id string) (*domain.DirectoryEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	u, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return &domain.DirectoryEntry{ID: u.ID, FullName: u.FullName, Email: u.Email}, nil
}

// --- in-memory RefreshTokenRepository ---

type fakeRefresh struct {
	mu     sync.Mutex
	seq    int
	byHash map[string]*domain.RefreshToken
	byID   map[string]*domain.RefreshToken
}

func newFakeRefresh() *fakeRefresh {
	return &fakeRefresh{byHash: map[string]*domain.RefreshToken{}, byID: map[string]*domain.RefreshToken{}}
}

func (f *fakeRefresh) Store(_ context.Context, t *domain.RefreshToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.seq++
	t.ID = fmt.Sprintf("rt-%d", f.seq)
	clone := *t
	f.byHash[t.TokenHash] = &clone
	f.byID[t.ID] = &clone
	return nil
}

func (f *fakeRefresh) FindByHash(_ context.Context, hash string) (*domain.RefreshToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, ok := f.byHash[hash]
	if !ok {
		return nil, domain.ErrTokenNotFound
	}
	out := *t
	return &out, nil
}

func (f *fakeRefresh) MarkReplaced(_ context.Context, id, replacedBy string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.byID[id]; ok {
		now := time.Now()
		t.ReplacedBy = &replacedBy
		t.RevokedAt = &now
	}
	return nil
}

func (f *fakeRefresh) Revoke(_ context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.byID[id]; ok {
		now := time.Now()
		t.RevokedAt = &now
	}
	return nil
}

func (f *fakeRefresh) RevokeFamily(_ context.Context, familyID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for _, t := range f.byID {
		if t.FamilyID == familyID {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (f *fakeRefresh) RevokeAllForUser(_ context.Context, userID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	for _, t := range f.byID {
		if t.UserID == userID {
			t.RevokedAt = &now
		}
	}
	return nil
}

// --- in-memory VerificationRepository ---

type fakeVerif struct {
	mu       sync.Mutex
	email    map[string]string // hash -> userID
	password map[string]string
}

func newFakeVerif() *fakeVerif {
	return &fakeVerif{email: map[string]string{}, password: map[string]string{}}
}

func (f *fakeVerif) StoreEmailToken(_ context.Context, userID, hash string, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.email[hash] = userID
	return nil
}

func (f *fakeVerif) ConsumeEmailToken(_ context.Context, hash string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	id, ok := f.email[hash]
	if !ok {
		return "", domain.ErrTokenNotFound
	}
	delete(f.email, hash)
	return id, nil
}

func (f *fakeVerif) StorePasswordToken(_ context.Context, userID, hash string, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.password[hash] = userID
	return nil
}

func (f *fakeVerif) ConsumePasswordToken(_ context.Context, hash string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	id, ok := f.password[hash]
	if !ok {
		return "", domain.ErrTokenNotFound
	}
	delete(f.password, hash)
	return id, nil
}

// --- trivial fakes ---

type noopAudit struct{}

func (noopAudit) Record(context.Context, string, string, string, string, map[string]any, string) error {
	return nil
}

type noopEvents struct{}

func (noopEvents) Publish(context.Context, string, string, map[string]any) error { return nil }

type fakeMailer struct {
	lastVerifyRaw string
	lastResetRaw  string
	verifyErr     error // when set, SendVerificationEmail fails (e.g. prod log-mailer)
}

func newFakeMailer() *fakeMailer { return &fakeMailer{} }

func (m *fakeMailer) SendVerificationEmail(_ context.Context, _, rawToken string) error {
	m.lastVerifyRaw = rawToken
	return m.verifyErr
}

func (m *fakeMailer) SendPasswordResetEmail(_ context.Context, _, rawToken string) error {
	m.lastResetRaw = rawToken
	return nil
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
