package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ErrResendNotConfigured is returned by NewResendMailer when RESEND_API_KEY is
// unset, so module.go can fall back to SMTP or the log-only mailer.
var ErrResendNotConfigured = errors.New("resend not configured (RESEND_API_KEY unset)")

// resendEndpoint is the Resend transactional email API. Delivery happens over
// HTTPS (port 443), so it works on hosts that block outbound SMTP (e.g. Railway
// below the Pro plan).
const resendEndpoint = "https://api.resend.com/emails"

// ResendMailer delivers the verification and password-reset emails via the
// Resend HTTP API. It implements domain.Mailer. The reset/verify tokens are
// embedded only in the message body sent to the account owner — never logged.
type ResendMailer struct {
	apiKey string
	from   string // RFC 5322 From, e.g. `Kirmya <no-reply@kirmya.com>`
	appURL string
	client *http.Client
}

// NewResendMailer builds a ResendMailer from the environment. Returns
// ErrResendNotConfigured when RESEND_API_KEY is empty.
//
//	RESEND_API_KEY  Resend API key (required to enable)
//	RESEND_FROM     From address; defaults to "<SMTP_FROM_NAME> <SMTP_FROM>",
//	                else "Kirmya <onboarding@resend.dev>" (Resend's test sender).
//	                For real delivery this MUST be an address on a domain you
//	                have verified in Resend.
//	APP_URL         base URL used to build the verify/reset links
func NewResendMailer() (*ResendMailer, error) {
	apiKey := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	if apiKey == "" {
		return nil, ErrResendNotConfigured
	}

	from := strings.TrimSpace(os.Getenv("RESEND_FROM"))
	if from == "" {
		addr := strings.TrimSpace(os.Getenv("SMTP_FROM"))
		name := strings.TrimSpace(os.Getenv("SMTP_FROM_NAME"))
		if name == "" {
			name = "Kirmya"
		}
		if addr != "" {
			from = fmt.Sprintf("%s <%s>", name, addr)
		} else {
			// Resend's shared test sender. Only delivers to the Resend account
			// owner until a custom domain is verified.
			from = name + " <onboarding@resend.dev>"
		}
	}

	appURL := strings.TrimSpace(os.Getenv("APP_URL"))
	if appURL == "" {
		appURL = "http://localhost:3000"
	}

	return &ResendMailer{
		apiKey: apiKey,
		from:   from,
		appURL: strings.TrimRight(appURL, "/"),
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// SendVerificationEmail sends the email-verification message with a link to
// APP_URL/verify-email?token=<raw>.
func (m *ResendMailer) SendVerificationEmail(ctx context.Context, email, rawToken string) error {
	link := m.appURL + "/verify-email?token=" + url.QueryEscape(rawToken)
	text, html := verificationBody(link)
	return m.send(ctx, email, "Verify your email - Kirmya", text, html)
}

// SendPasswordResetEmail sends the password-reset message with a link to
// APP_URL/reset-password?token=<raw>.
func (m *ResendMailer) SendPasswordResetEmail(ctx context.Context, email, rawToken string) error {
	link := m.appURL + "/reset-password?token=" + url.QueryEscape(rawToken)
	text, html := resetBody(link)
	return m.send(ctx, email, "Reset your password - Kirmya", text, html)
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html"`
}

// send POSTs the message to the Resend API and returns an error on any non-2xx
// response (with the API's error body included for diagnosis).
func (m *ResendMailer) send(ctx context.Context, to, subject, text, html string) error {
	if strings.ContainsAny(to, "\r\n") || strings.ContainsAny(subject, "\r\n") {
		return errors.New("resend: illegal CRLF in recipient or subject")
	}

	payload, err := json.Marshal(resendRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Text:    text,
		HTML:    html,
	})
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend send: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("resend send: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
