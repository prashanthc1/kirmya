package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// WriteJSON writes the payload verbatim (no envelope) — used where handlers
// shape their own body.
func TestWriteJSONWritesRawPayload(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteJSON(rec, http.StatusOK, map[string]string{"message": "ok"})

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected JSON content-type, got %q", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["message"] != "ok" {
		t.Fatalf("expected raw payload, got %v", body)
	}
}

// WriteSuccess wraps the payload in the standard success envelope.
func TestWriteSuccessWrapsPayload(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteSuccess(rec, http.StatusOK, map[string]string{"id": "42"})

	var response SuccessResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !response.Success {
		t.Fatal("expected success=true")
	}
	if response.Data == nil {
		t.Fatal("expected data payload")
	}
	if response.Meta.Timestamp == "" {
		t.Fatal("expected response timestamp")
	}
}

// WriteError (via the typed helpers) emits the standard error envelope.
func TestWriteErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteUnauthorizedError(rec, "missing token")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}

	var response ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Success {
		t.Fatal("expected success=false")
	}
	if response.Error == nil {
		t.Fatal("expected error payload")
	}
	if response.Error.Code != "unauthorized" {
		t.Fatalf("expected unauthorized code, got %q", response.Error.Code)
	}
	if response.Error.Message != "missing token" {
		t.Fatalf("expected message %q, got %q", "missing token", response.Error.Message)
	}
}
