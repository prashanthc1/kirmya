package mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

// ErrSMTPNotConfigured is returned by NewSMTPMailer when SMTP_HOST is unset, so
// the caller (module.go) can fall back to the log-only mailer in development.
var ErrSMTPNotConfigured = errors.New("smtp not configured (SMTP_HOST unset)")

// SMTPMailer delivers transactional auth emails (verification, password reset)
// over SMTP. It implements domain.Mailer. Connection security is chosen by port:
// 465 uses implicit TLS, anything else attempts STARTTLS and refuses to send
// credentials over a plaintext link. Tokens are embedded only in the message
// body sent to the account owner — never logged.
type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string // envelope + From: address
	fromName string // display name in the From: header
	appURL   string
	timeout  time.Duration
}

// NewSMTPMailer builds an SMTPMailer from the environment. Returns
// ErrSMTPNotConfigured when SMTP_HOST is empty.
//
//	SMTP_HOST      smtp host (required to enable)
//	SMTP_PORT      smtp port (default 587)
//	SMTP_USERNAME  auth username (optional; if set, PLAIN auth is used)
//	SMTP_PASSWORD  auth password
//	SMTP_FROM      From: address (default: SMTP_USERNAME, else no-reply@<host-domain>)
//	SMTP_FROM_NAME From: display name (default "Kirmya")
func NewSMTPMailer() (*SMTPMailer, error) {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		return nil, ErrSMTPNotConfigured
	}
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USERNAME")
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
	if from == "" {
		if username != "" {
			from = username
		} else {
			from = "no-reply@" + host
		}
	}
	fromName := strings.TrimSpace(os.Getenv("SMTP_FROM_NAME"))
	if fromName == "" {
		fromName = "Kirmya"
	}
	appURL := strings.TrimSpace(os.Getenv("APP_URL"))
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: os.Getenv("SMTP_PASSWORD"),
		from:     from,
		fromName: fromName,
		appURL:   strings.TrimRight(appURL, "/"),
		timeout:  10 * time.Second,
	}, nil
}

// SendVerificationEmail composes and sends the email-verification message with a
// link to APP_URL/verify-email?token=<raw>.
func (m *SMTPMailer) SendVerificationEmail(ctx context.Context, email, rawToken string) error {
	link := m.appURL + "/verify-email?token=" + url.QueryEscape(rawToken)
	text, html := verificationBody(link)
	return m.send(ctx, email, "Verify your email - Kirmya", text, html)
}

// SendPasswordResetEmail composes and sends the password-reset message with a
// link to APP_URL/reset-password?token=<raw>.
func (m *SMTPMailer) SendPasswordResetEmail(ctx context.Context, email, rawToken string) error {
	link := m.appURL + "/reset-password?token=" + url.QueryEscape(rawToken)
	text, html := resetBody(link)
	return m.send(ctx, email, "Reset your password - Kirmya", text, html)
}

// send dials the SMTP server (implicit TLS on :465, otherwise STARTTLS),
// authenticates if a username is configured, and delivers a multipart message.
func (m *SMTPMailer) send(ctx context.Context, to, subject, text, html string) error {
	// Defend against SMTP header injection: a CRLF in the recipient (or subject)
	// could smuggle extra headers/recipients. Reject rather than sanitise.
	if strings.ContainsAny(to, "\r\n") || strings.ContainsAny(subject, "\r\n") {
		return errors.New("smtp: illegal CRLF in recipient or subject")
	}

	addr := net.JoinHostPort(m.host, m.port)
	d := net.Dialer{Timeout: m.timeout}

	var conn net.Conn
	var err error
	if m.port == "465" {
		conn, err = tls.DialWithDialer(&d, "tcp", addr, &tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12})
	} else {
		conn, err = d.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	// Bound the whole exchange so a slow/hung server can't pin a request
	// goroutine indefinitely.
	_ = conn.SetDeadline(time.Now().Add(m.timeout))

	c, err := smtp.NewClient(conn, m.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer func() { _ = c.Quit() }()

	// Upgrade to TLS via STARTTLS when not already on an implicit-TLS port.
	if m.port != "465" {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}); err != nil {
				return fmt.Errorf("smtp starttls: %w", err)
			}
		} else if m.username != "" {
			// Refuse to send credentials over an un-encrypted connection.
			return errors.New("smtp: server does not support STARTTLS; refusing to send credentials in plaintext")
		}
	}

	if m.username != "" {
		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := c.Mail(m.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write([]byte(m.buildMessage(to, subject, text, html))); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}
	return nil
}

// buildMessage assembles a multipart/alternative MIME message (plain text +
// HTML) with the standard auth headers.
func (m *SMTPMailer) buildMessage(to, subject, text, html string) string {
	boundary := "kirmya-boundary-7f3a9c1e"
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s <%s>\r\n", m.fromName, m.from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(&b, "Content-Type: multipart/alternative; boundary=%q\r\n", boundary)
	b.WriteString("\r\n")

	fmt.Fprintf(&b, "--%s\r\n", boundary)
	b.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n")
	b.WriteString(text)
	b.WriteString("\r\n\r\n")

	fmt.Fprintf(&b, "--%s\r\n", boundary)
	b.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n\r\n")
	b.WriteString(html)
	b.WriteString("\r\n\r\n")

	fmt.Fprintf(&b, "--%s--\r\n", boundary)
	return b.String()
}
