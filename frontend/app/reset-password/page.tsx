"use client";

/**
 * /reset-password — destination for the link in the password-reset email
 * (APP_URL/reset-password?token=<raw>, built by the backend mailer). Reads the
 * token from the query string, collects a new password, and posts to
 * POST /auth/reset-password { token, password }. On success the user is sent to
 * sign in (the backend has revoked every existing session).
 *
 * useSearchParams() requires a Suspense boundary in the App Router, so the token
 * reader is isolated in ResetPasswordInner and the route export wraps it.
 */
import { Suspense, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { api, ApiError } from "@/lib/api/client";

function ResetPasswordInner() {
  const router = useRouter();
  const params = useSearchParams();
  const token = params.get("token");

  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!token) {
      setError("This reset link is missing its token. Request a new one.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }
    if (password !== confirm) {
      setError("Passwords don't match.");
      return;
    }

    setLoading(true);
    try {
      await api.post("/auth/reset-password", { token, password });
      setDone(true);
      // Give the confirmation a moment, then move to sign in.
      setTimeout(() => router.push("/sign-in"), 1500);
    } catch (err) {
      setError(
        err instanceof ApiError
          ? "This reset link is invalid or has expired. Request a new one."
          : "Something went wrong resetting your password. Please try again.",
      );
    } finally {
      setLoading(false);
    }
  }

  if (done) {
    return (
      <div style={cardStyle}>
        <h1 style={headingStyle}>Password updated ✓</h1>
        <p style={bodyStyle}>
          Your password has been reset and all other sessions were signed out.
          Taking you to sign in…
        </p>
        <Link href="/sign-in" style={primaryLinkStyle}>
          Go to sign in
        </Link>
      </div>
    );
  }

  return (
    <form style={cardStyle} onSubmit={handleSubmit} noValidate>
      <h1 style={headingStyle}>Choose a new password</h1>
      <p style={{ ...bodyStyle, margin: "0 0 6px" }}>
        Enter a new password for your account.
      </p>

      {error && (
        <div role="alert" style={errorStyle}>
          {error}
        </div>
      )}

      <label style={labelStyle} htmlFor="password">
        New password
      </label>
      <input
        id="password"
        name="password"
        type="password"
        autoComplete="new-password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="At least 8 characters"
        style={inputStyle}
      />

      <label style={labelStyle} htmlFor="confirm">
        Confirm new password
      </label>
      <input
        id="confirm"
        name="confirm"
        type="password"
        autoComplete="new-password"
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        placeholder="Re-enter your password"
        style={inputStyle}
      />

      <button type="submit" disabled={loading} style={buttonStyle(loading)}>
        {loading ? "Resetting…" : "Reset password"}
      </button>

      <div style={footerStyle}>
        <Link href="/sign-in" style={{ color: "#C2683C", fontWeight: 600 }}>
          Back to sign in
        </Link>
      </div>
    </form>
  );
}

export default function ResetPasswordPage() {
  return (
    <main style={pageStyle}>
      <Suspense
        fallback={
          <div style={cardStyle}>
            <h1 style={headingStyle}>Loading…</h1>
          </div>
        }
      >
        <ResetPasswordInner />
      </Suspense>
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
  margin: "16px 0 7px",
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

const errorStyle: React.CSSProperties = {
  background: "rgba(194,104,60,0.10)",
  border: "1px solid rgba(194,104,60,0.35)",
  color: "#9A4A24",
  borderRadius: "10px",
  padding: "11px 14px",
  fontSize: "14px",
  margin: "16px 0 4px",
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
