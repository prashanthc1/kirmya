"use client";

/**
 * /forgot-password — collects an email and posts to POST /auth/forgot-password.
 * The backend always responds 200 (no account enumeration), so we always show
 * the same neutral confirmation whether or not the address has an account.
 */
import { useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api/client";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!email.trim()) return;
    setLoading(true);
    try {
      await api.post("/auth/forgot-password", { email: email.trim() });
    } catch {
      /* ignore — the confirmation is intentionally identical on failure */
    } finally {
      setLoading(false);
      setSent(true);
    }
  }

  return (
    <main style={pageStyle}>
      {sent ? (
        <div style={cardStyle}>
          <h1 style={headingStyle}>Check your inbox</h1>
          <p style={bodyStyle}>
            If an account exists for <strong>{email.trim()}</strong>, we&apos;ve
            sent a password-reset link. It expires in 1 hour.
          </p>
          <Link href="/sign-in" style={primaryLinkStyle}>
            Back to sign in
          </Link>
        </div>
      ) : (
        <form style={cardStyle} onSubmit={handleSubmit} noValidate>
          <h1 style={headingStyle}>Reset your password</h1>
          <p style={{ ...bodyStyle, margin: "0 0 6px" }}>
            Enter your email and we&apos;ll send you a link to set a new password.
          </p>

          <label style={labelStyle} htmlFor="email">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@email.com"
            style={inputStyle}
          />

          <button type="submit" disabled={loading} style={buttonStyle(loading)}>
            {loading ? "Sending…" : "Send reset link"}
          </button>

          <div style={footerStyle}>
            Remembered it?{" "}
            <Link href="/sign-in" style={{ color: "#C2683C", fontWeight: 600 }}>
              Sign in
            </Link>
          </div>
        </form>
      )}
    </main>
  );
}

const pageStyle: React.CSSProperties = {
  display: "flex",
  minHeight: "100vh",
  alignItems: "center",
  justifyContent: "center",
  padding: "40px 20px",
  fontFamily: "'Public Sans',sans-serif",
  background: "#FBF7F2",
  color: "#2B2620",
};

const cardStyle: React.CSSProperties = {
  width: "100%",
  maxWidth: "440px",
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "20px",
  padding: "clamp(28px,4vw,40px)",
};

const headingStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 800,
  fontSize: "28px",
  letterSpacing: "-0.02em",
  margin: "0 0 10px",
};

const bodyStyle: React.CSSProperties = {
  fontSize: "15px",
  color: "#5B554C",
  lineHeight: 1.6,
  margin: 0,
};

const labelStyle: React.CSSProperties = {
  display: "block",
  fontSize: "13px",
  fontWeight: 600,
  color: "#8A8175",
  margin: "20px 0 7px",
};

const inputStyle: React.CSSProperties = {
  width: "100%",
  border: "1px solid #E2D9CC",
  borderRadius: "10px",
  padding: "13px 14px",
  fontSize: "15px",
  color: "#2B2620",
  outline: "none",
  background: "#FCFAF7",
  boxSizing: "border-box",
};

const primaryLinkStyle: React.CSSProperties = {
  display: "inline-block",
  marginTop: "20px",
  background: "#C2683C",
  color: "#fff",
  fontSize: "15px",
  fontWeight: 600,
  padding: "13px 24px",
  borderRadius: "100px",
};

const footerStyle: React.CSSProperties = {
  textAlign: "center",
  marginTop: "22px",
  fontSize: "15px",
  color: "#5B554C",
};

function buttonStyle(loading: boolean): React.CSSProperties {
  return {
    width: "100%",
    marginTop: "26px",
    border: "none",
    background: loading ? "#D89870" : "#C2683C",
    color: "#fff",
    fontFamily: "'Public Sans',sans-serif",
    fontSize: "16px",
    fontWeight: 600,
    padding: "15px",
    borderRadius: "100px",
    cursor: loading ? "default" : "pointer",
  };
}
