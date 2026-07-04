package mailer

import "fmt"

// Brand tokens for transactional emails (mirrors the Kirmya brand identity
// system: Deep Navy primary, Ember accent, "Your Career Comeback").
const (
	brandNavy    = "#1A3A6B"
	brandEmber   = "#E07B3A"
	brandParch   = "#F4F1EC"
	brandTagline = "Your Career Comeback"
)

// verificationBody returns the plain-text and HTML bodies for the email
// verification message. The link expires in 24 hours (see verificationTTL).
func verificationBody(link string) (text, html string) {
	text = fmt.Sprintf(`Welcome to Kirmya — %s.

Confirm your email to activate your account:

%s

This link expires in 24 hours. If you didn't create a Kirmya account, you can safely ignore this email.

— The Kirmya team`, brandTagline, link)

	html = emailShell(
		"Confirm your email",
		"You're one step from your comeback. Confirm your email address to activate your Kirmya account.",
		"Verify email",
		link,
		"This link expires in 24 hours. If you didn't create a Kirmya account, you can safely ignore this email.",
	)
	return text, html
}

// resetBody returns the plain-text and HTML bodies for the password-reset
// message. The link expires in 1 hour (see resetTTL).
func resetBody(link string) (text, html string) {
	text = fmt.Sprintf(`Reset your Kirmya password.

We received a request to reset your password. Choose a new one here:

%s

This link expires in 1 hour. If you didn't request a reset, you can safely ignore this email — your password won't change.

— The Kirmya team`, link)

	html = emailShell(
		"Reset your password",
		"We received a request to reset your Kirmya password. Choose a new one using the button below.",
		"Reset password",
		link,
		"This link expires in 1 hour. If you didn't request a reset, you can safely ignore this email — your password won't change.",
	)
	return text, html
}

// emailShell renders a minimal, table-based responsive HTML email with the
// Kirmya wordmark, a single call-to-action button, and a footnote.
func emailShell(heading, intro, cta, link, note string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:%[3]s;font-family:-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background:%[3]s;padding:32px 12px;">
    <tr><td align="center">
      <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="max-width:480px;background:#ffffff;border-radius:16px;overflow:hidden;border:1px solid #E3DDD2;">
        <tr><td style="padding:28px 32px 8px 32px;">
          <span style="font-family:Georgia,'Times New Roman',serif;font-size:26px;font-weight:700;letter-spacing:-0.5px;color:%[1]s;">Kir<span style="color:%[2]s;">mya</span></span>
          <div style="font-size:10px;letter-spacing:2px;text-transform:uppercase;color:#6B7689;margin-top:6px;">%[7]s</div>
        </td></tr>
        <tr><td style="padding:16px 32px 0 32px;">
          <h1 style="font-size:20px;color:%[1]s;margin:0 0 12px 0;">%[4]s</h1>
          <p style="font-size:15px;line-height:1.6;color:#33404F;margin:0 0 24px 0;">%[5]s</p>
          <table role="presentation" cellpadding="0" cellspacing="0"><tr><td style="border-radius:10px;background:%[1]s;">
            <a href="%[8]s" style="display:inline-block;padding:13px 26px;font-size:15px;font-weight:600;color:#ffffff;text-decoration:none;border-radius:10px;">%[6]s</a>
          </td></tr></table>
          <p style="font-size:12px;line-height:1.6;color:#6B7689;margin:24px 0 4px 0;">Or paste this link into your browser:</p>
          <p style="font-size:12px;word-break:break-all;margin:0 0 24px 0;"><a href="%[8]s" style="color:%[2]s;">%[8]s</a></p>
        </td></tr>
        <tr><td style="padding:0 32px 28px 32px;border-top:1px solid #E3DDD2;">
          <p style="font-size:12px;line-height:1.6;color:#6B7689;margin:18px 0 0 0;">%[9]s</p>
        </td></tr>
      </table>
      <p style="font-size:11px;color:#9AA1AC;margin:18px 0 0 0;">© Kirmya · %[7]s</p>
    </td></tr>
  </table>
</body></html>`, brandNavy, brandEmber, brandParch, heading, intro, cta, brandTagline, link, note)
}
