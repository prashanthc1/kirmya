package mailer

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestSend_RejectsCRLFInjection(t *testing.T) {
	m := &SMTPMailer{host: "smtp.example.com", port: "587", from: "no-reply@kirmya.com", fromName: "Kirmya"}
	// The CRLF guard must fire before any network dial, so this returns an error
	// without attempting a connection to a non-existent host.
	err := m.send(context.Background(), "victim@example.com\r\nBcc: attacker@evil.com", "Subj", "t", "h")
	if err == nil || !strings.Contains(err.Error(), "CRLF") {
		t.Fatalf("want CRLF rejection, got %v", err)
	}
}

func TestNewSMTPMailer_NotConfigured(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	if _, err := NewSMTPMailer(); !errors.Is(err, ErrSMTPNotConfigured) {
		t.Fatalf("want ErrSMTPNotConfigured, got %v", err)
	}
}

func TestNewSMTPMailer_Defaults(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_FROM", "")
	t.Setenv("SMTP_FROM_NAME", "")
	t.Setenv("APP_URL", "https://app.kirmya.com/")

	m, err := NewSMTPMailer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.port != "587" {
		t.Errorf("default port = %q, want 587", m.port)
	}
	if m.from != "no-reply@smtp.example.com" {
		t.Errorf("default from = %q, want no-reply@smtp.example.com", m.from)
	}
	if m.fromName != "Kirmya" {
		t.Errorf("default fromName = %q, want Kirmya", m.fromName)
	}
	if m.appURL != "https://app.kirmya.com" { // trailing slash trimmed
		t.Errorf("appURL = %q, want trailing slash trimmed", m.appURL)
	}
}

func TestNewSMTPMailer_FromFallsBackToUsername(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_USERNAME", "postmaster@kirmya.com")
	t.Setenv("SMTP_FROM", "")
	m, err := NewSMTPMailer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.from != "postmaster@kirmya.com" {
		t.Errorf("from = %q, want fallback to username", m.from)
	}
}

func TestBuildMessage_MIMEStructure(t *testing.T) {
	m := &SMTPMailer{from: "no-reply@kirmya.com", fromName: "Kirmya"}
	msg := m.buildMessage("user@example.com", "Verify your email - Kirmya", "TEXT-PART", "<p>HTML-PART</p>")

	for _, want := range []string{
		"From: Kirmya <no-reply@kirmya.com>\r\n",
		"To: user@example.com\r\n",
		"Subject: Verify your email - Kirmya\r\n",
		"MIME-Version: 1.0\r\n",
		"multipart/alternative; boundary=",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"Content-Type: text/html; charset=\"utf-8\"",
		"TEXT-PART",
		"<p>HTML-PART</p>",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("message missing %q", want)
		}
	}
	// CRLF line endings are required by SMTP.
	if strings.Contains(msg, "\n") && !strings.Contains(msg, "\r\n") {
		t.Error("message must use CRLF line endings")
	}
}

func TestVerificationBody_ContainsLinkAndExpiry(t *testing.T) {
	link := "https://app.kirmya.com/verify-email?token=abc123"
	text, html := verificationBody(link)

	if !strings.Contains(text, link) {
		t.Error("text body missing verification link")
	}
	if !strings.Contains(html, link) {
		t.Error("html body missing verification link")
	}
	if !strings.Contains(text, "24 hours") {
		t.Error("text body should state the 24-hour expiry")
	}
	if !strings.Contains(html, "Verify email") {
		t.Error("html body missing CTA label")
	}
}

func TestResetBody_ContainsLinkAndExpiry(t *testing.T) {
	link := "https://app.kirmya.com/reset-password?token=xyz789"
	text, html := resetBody(link)

	if !strings.Contains(text, link) || !strings.Contains(html, link) {
		t.Error("reset bodies missing reset link")
	}
	if !strings.Contains(text, "1 hour") {
		t.Error("text body should state the 1-hour expiry")
	}
}

// raw tokens must never appear in the non-reversible token reference used for
// dev log correlation.
func TestTokenRef_DoesNotExposeToken(t *testing.T) {
	raw := "super-secret-raw-token"
	ref := tokenRef(raw)
	if strings.Contains(ref, raw) || len(ref) != 8 {
		t.Errorf("tokenRef leaked or wrong length: %q", ref)
	}
}
