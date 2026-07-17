//go:build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"workspace-app/internal/common"
	"workspace-app/internal/connections"
	"workspace-app/internal/platform/cache"
	"workspace-app/internal/profile/application"
	"workspace-app/internal/profile/domain"
	"workspace-app/internal/profile/infrastructure/postgres"
	"workspace-app/internal/testsupport"
)

func TestGinProfileHandlers(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	repo := postgres.NewRepository(db)
	cch := cache.NewMemory()
	svc := application.NewService(repo, nil, cch)

	ctx := context.Background()

	// Insert test users
	userA := testsupport.InsertUser(t, db, "user_a@kirmya.test", "User A")
	userB := testsupport.InsertUser(t, db, "user_b@kirmya.test", "User B")
	userC := testsupport.InsertUser(t, db, "user_c@kirmya.test", "User C")

	// Set up connection A <-> B (accepted) using real Connections Service
	connectionsRepo := connections.NewRepository(db)
	connectionsSvc := connections.NewService(db, connectionsRepo, nil)

	err := connectionsSvc.SendRequest(ctx, userA, userB, nil, nil)
	if err != nil {
		t.Fatalf("send connection request: %v", err)
	}

	conn, err := connectionsRepo.GetConnection(ctx, userA, userB)
	if err != nil {
		t.Fatalf("get connection: %v", err)
	}

	err = connectionsSvc.AcceptRequest(ctx, conn.ID, userB)
	if err != nil {
		t.Fatalf("accept connection request: %v", err)
	}

	// Setup initial profile for B using service
	pB, err := svc.Get(ctx, userB)
	if err != nil {
		t.Fatalf("get B draft: %v", err)
	}
	pB.Identity.Headline = "Eng Lead"
	pB.Identity.Email = "user_b@kirmya.test"
	pB.Identity.Phone = "+123456"
	pB.Preferences.SalaryMin = 100000
	pB.Preferences.SalaryMax = 150000
	pB.Preferences.SalaryCurrency = "USD"
	pB.Privacy.FieldVisibility["profile"] = "public"
	pB.Privacy.FieldVisibility["salary"] = "connections_only"
	pB.Privacy.FieldVisibility["experience"] = "public"

	_, err = svc.UpdateProfile(ctx, userB, pB.Version, domain.AggregateUpdate{
		Identity:    &pB.Identity,
		Preferences: &pB.Preferences,
		Privacy:     &pB.Privacy,
	})
	if err != nil {
		t.Fatalf("update B draft: %v", err)
	}

	// Publish B profile to version snapshot to support GetPublished read
	_, err = svc.Publish(ctx, userB, userB, "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("publish B: %v", err)
	}

	// Set up Gin Router
	gin.SetMode(gin.TestMode)

	// Helper to perform requests with a specific authenticated user
	performRequest := func(method, path string, body []byte, authedUserID string) *httptest.ResponseRecorder {
		r := gin.New()
		authMock := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				authHeader := req.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					uid := strings.TrimPrefix(authHeader, "Bearer ")
					req = req.WithContext(common.ContextWithUserID(req.Context(), uid))
				}
				next.ServeHTTP(w, req)
			})
		}
		RegisterGinRoutes(r, db, authMock, svc, cch)

		w := httptest.NewRecorder()
		var req *http.Request
		if body != nil {
			req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(method, path, nil)
		}
		if authedUserID != "" {
			req.Header.Set("Authorization", "Bearer "+authedUserID)
		}
		r.ServeHTTP(w, req)
		return w
	}

	// 1. Test GET /api/profile/me (decrypted own view)
	t.Run("Get own profile", func(t *testing.T) {
		w := performRequest("GET", "/api/profile/me", nil, userB)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. body: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		headline := resp["headline"].(string)
		if headline != "Eng Lead" {
			t.Fatalf("expected headline 'Eng Lead', got %s", headline)
		}
		email := resp["email"].(string)
		if email != "user_b@kirmya.test" {
			t.Fatalf("expected decrypted email 'user_b@kirmya.test', got %s", email)
		}
	})

	// 2. Test GET /api/profile/public/:userId
	t.Run("Get public profile as connection (A -> B)", func(t *testing.T) {
		w := performRequest("GET", "/api/profile/public/"+userB, nil, userA)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. body: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["is_connection"].(bool) != true {
			t.Fatal("expected is_connection to be true")
		}
		// Connection-gated sensitive details should be present
		email, ok := resp["email"].(string)
		if !ok || email != "user_b@kirmya.test" {
			t.Fatalf("expected email 'user_b@kirmya.test', got %s (ok=%t)", email, ok)
		}
		salMin, ok := resp["salary_min"].(float64)
		if !ok || salMin != 100000 {
			t.Fatalf("expected salary_min 100000, got %v", resp["salary_min"])
		}
	})

	t.Run("Get public profile as non-connection (C -> B)", func(t *testing.T) {
		w := performRequest("GET", "/api/profile/public/"+userB, nil, userC)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["is_connection"].(bool) != false {
			t.Fatal("expected is_connection to be false")
		}
		// Sensitive contact/salary details must be omitted entirely
		if _, ok := resp["email"]; ok {
			t.Fatal("sensitive field email was NOT omitted for non-connection")
		}
		if _, ok := resp["salary_min"]; ok {
			t.Fatal("sensitive field salary_min was NOT omitted for non-connection")
		}
	})

	// 3. Test PATCH endpoints
	t.Run("PATCH basic-info", func(t *testing.T) {
		body := []byte(`{"headline":"Staff Software Engineer","location":"Dubai"}`)
		w := performRequest("PATCH", "/api/profile/me/basic-info", body, userB)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d. body: %s", w.Code, w.Body.String())
		}
		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)

		if resp["headline"].(string) != "Staff Software Engineer" {
			t.Fatalf("expected new headline 'Staff Software Engineer', got %s", resp["headline"])
		}
	})

	t.Run("PATCH contact details and validation", func(t *testing.T) {
		// Valid Email
		body := []byte(`{"email":"new_b@kirmya.test","phone":"+97150"}`)
		w := performRequest("PATCH", "/api/profile/me/contact", body, userB)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// Invalid Email Format
		bodyInvalid := []byte(`{"email":"bad_email_no_at"}`)
		w2 := performRequest("PATCH", "/api/profile/me/contact", bodyInvalid, userB)
		if w2.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 validation error, got %d", w2.Code)
		}
	})
}
