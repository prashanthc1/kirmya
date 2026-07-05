"use client";

/**
 * /verify-email — destination for the link in the verification email
 * (APP_URL/verify-email?token=<raw>, built by the backend mailer). On mount it
 * posts the token to POST /auth/verify-email and reports the outcome. If the
 * token is missing, invalid, or expired it offers a resend.
 *
 * useSearchParams() requires a Suspense boundary in the App Router, so the token
 * reader is isolated in VerifyEmailInner and the route export wraps it.
 */
import { Suspense, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { api, ApiError } from "@/lib/api/client";

type Status = "verifying" | "success" | "error" | "missing";

function VerifyEmailInner() {
  const params = useSearchParams();
  const token = params.get("token");

  const [status, setStatus] = useState<Status>(token ? "verifying" : "missing");
  const [message, setMessage] = useState<string | null>(null);

  // Resend state.
  const [email, setEmail] = useState("");
  const [resending, setResending] = useState(false);
  const [resent, setResent] = useState(false);

  // Guard against the double-invocation of effects in React 18 StrictMode so we
  // don't consume the single-use token twice (the second call would 404).
  const started = useRef(false);

  useEffect(() => {
    if (!token || started.current) return;
    started.current = true;
    (async () => {
      try {
        await api.post("/auth/verify-email", { token });
        setStatus("success");
      } catch (err) {
        setStatus("error");
        setMessage(
          err instanceof ApiError
            ? "This verification link is invalid or has expired."
            : "Something went wrong verifying your email. Please try again.",
        );
      }
    })();
  }, [token]);

  async function handleResend(e: React.FormEvent) {
    e.preventDefault();
    if (!email.trim()) return;
    setResending(true);
    try {
      // Always succeeds server-side (no account enumeration); show a neutral
      // confirmation regardless.
      await api.post("/auth/resend-verification", { email: email.trim() });
    } catch {
      /* ignore — we still show the neutral confirmation */
    } finally {
      setResending(false);
      setResent(true);
    }
  }

  if (status === "verifying") {
    return (
      <Card>
        <Heading>Verifying your email…</Heading>
        <p style={bodyStyle}>Hang tight while we confirm your link.</p>
      </Card>
    );
  }

  if (status === "success") {
    return (
      <Card>
        <Heading>Email verified ✓</Heading>
        <p style={bodyStyle}>
          Your email is confirmed. You can sign in and get started.
        </p>
        <Link href="/sign-in" style={primaryLinkStyle}>
          Go to sign in
        </Link>
      </Card>
    );
  }

  // "missing" or "error": offer a resend.
  return (
    <Card>
      <Heading>{status === "missing" ? "Missing verification token" : "Link expired"}</Heading>
      <p style={bodyStyle}>
        {message ??
          "This page needs a valid verification link. Enter your email and we'll send a fresh one."}
      </p>

      {resent ? (
        <p style={{ ...bodyStyle, color: "#2B7A4B", marginTop: "18px" }}>
          If that email has a pending account, a new verification link is on its way.
        </p>
      ) : (
        <form onSubmit={handleResend} noValidate>
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
          <button type="submit" disabled={resending} style={buttonStyle(resending)}>
            {resending ? "Sending…" : "Resend verification email"}
          </button>
        </form>
      )}

      <div style={footerStyle}>
        <Link href="/sign-in" style={{ color: "#C2683C", fontWeight: 600 }}>
          Back to sign in
        </Link>
      </div>
    </Card>
  );
}

export default function VerifyEmailPage() {
  return (
    <main style={pageStyle}>
      <Suspense
        fallback={
          <Card>
            <Heading>Loading…</Heading>
          </Card>
        }
      >
        <VerifyEmailInner />
      </Suspense>
    </main>
  );
}

// --- shared presentational bits ---------------------------------------------

function Card({ children }: { children: React.ReactNode }) {
  return <div style={cardStyle}>{children}</div>;
}

function Heading({ children }: { children: React.ReactNode }) {
  return <h1 style={headingStyle}>{children}</h1>;
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
  fontFamily: "'Public Sans',sans-serif",
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
    marginTop: "22px",
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
