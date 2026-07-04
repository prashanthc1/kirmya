"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, setAccessToken, ApiError } from "@/lib/api/client";
import { useAuth, type AuthUser } from "@/lib/auth/auth-context";

interface LoginResponse {
  access_token?: string;
  mfa_required?: boolean;
  user?: AuthUser;
}

export default function SignInPage() {
  const router = useRouter();
  const { setUser } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [mfaRequired, setMfaRequired] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const data = await api.post<LoginResponse>("/auth/login", {
        email: email.trim(),
        password,
        ...(code ? { code } : {}),
      });

      if (data?.mfa_required) {
        setMfaRequired(true);
        setError("Enter the 6-digit code from your authenticator app.");
        return;
      }
      if (data?.access_token) {
        setAccessToken(data.access_token);
        if (data.user) setUser(data.user);
        router.push("/dashboard");
        return;
      }
      setError("Unexpected response from the server. Please try again.");
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : "Could not sign in. Please check your connection and try again.",
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div
      style={{
        display: "flex",
        minHeight: "100vh",
        fontFamily: "'Public Sans',sans-serif",
        color: "#2B2620",
        background: "#FBF7F2",
      }}
    >
      {/* LEFT BRAND PANEL */}
      <div
        style={{
          flex: "1 1 0",
          minWidth: 0,
          background: "#2B2620",
          color: "#fff",
          padding: "clamp(36px,4vw,56px)",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          position: "relative",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            position: "absolute",
            bottom: "-160px",
            right: "-120px",
            width: "480px",
            height: "480px",
            borderRadius: "50%",
            background:
              "radial-gradient(circle, rgba(79,124,106,0.45), transparent 68%)",
            pointerEvents: "none",
          }}
        />
        <Link
          href="/"
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontSize: "26px",
            fontWeight: 800,
            letterSpacing: "-0.02em",
            color: "#fff",
            position: "relative",
          }}
        >
          Kirmya
        </Link>
        <div style={{ position: "relative", maxWidth: "440px" }}>
          <div
            style={{
              fontSize: "46px",
              lineHeight: 1,
              color: "#E7A57E",
              fontFamily: "'Bricolage Grotesque',sans-serif",
              marginBottom: "12px",
            }}
          >
            &ldquo;
          </div>
          <p
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 500,
              fontSize: "clamp(20px,2.2vw,28px)",
              lineHeight: 1.35,
              letterSpacing: "-0.01em",
              margin: "0 0 24px",
            }}
          >
            One account. I&apos;m a job seeker, a mentor, and occasionally the
            one hiring. Kirmya handles all of it.
          </p>
          <div style={{ display: "flex", alignItems: "center", gap: "14px" }}>
            <span
              style={{
                width: "46px",
                height: "46px",
                borderRadius: "50%",
                background: "#C2683C",
                flex: "none",
                display: "inline-block",
              }}
            />
            <div>
              <div style={{ fontWeight: 600, fontSize: "15px" }}>Priya Nair</div>
              <div style={{ fontSize: "13px", color: "#9C958A" }}>
                Career coach &amp; operations lead · 3 roles active
              </div>
            </div>
          </div>
        </div>
        <div
          style={{
            position: "relative",
            display: "flex",
            gap: "28px",
            fontSize: "14px",
            color: "#9C958A",
            flexWrap: "wrap",
          }}
        >
          <span>
            <strong
              style={{
                color: "#fff",
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontSize: "18px",
              }}
            >
              12.4k
            </strong>
            &nbsp;placed
          </span>
          <span>
            <strong
              style={{
                color: "#fff",
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontSize: "18px",
              }}
            >
              480
            </strong>
            &nbsp;partners
          </span>
          <span>
            <strong
              style={{
                color: "#fff",
                fontFamily: "'Bricolage Grotesque',sans-serif",
                fontSize: "18px",
              }}
            >
              100%
            </strong>
            &nbsp;free to join
          </span>
        </div>
      </div>

      {/* RIGHT FORM PANEL */}
      <div
        style={{
          flex: "1 1 0",
          minWidth: 0,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          padding: "clamp(28px,4vw,56px)",
          overflowY: "auto",
        }}
      >
        <div style={{ width: "100%", maxWidth: "440px" }}>
          <h1
            style={{
              fontFamily: "'Bricolage Grotesque',sans-serif",
              fontWeight: 800,
              fontSize: "clamp(28px,3.4vw,38px)",
              lineHeight: 1.05,
              letterSpacing: "-0.02em",
              margin: "0 0 8px",
            }}
          >
            Welcome back
          </h1>
          <p style={{ fontSize: "16px", color: "#5B554C", margin: "0 0 28px" }}>
            Sign in to your Kirmya account.
          </p>

          <form onSubmit={handleSubmit} noValidate>
            {error && (
              <div
                role="alert"
                style={{
                  background: "rgba(194,104,60,0.10)",
                  border: "1px solid rgba(194,104,60,0.35)",
                  color: "#9A4A24",
                  borderRadius: "10px",
                  padding: "11px 14px",
                  fontSize: "14px",
                  marginBottom: "16px",
                }}
              >
                {error}
              </div>
            )}

            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "16px",
                marginBottom: "22px",
              }}
            >
              <div>
                <label htmlFor="email" style={labelStyle}>
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
              </div>
              <div>
                <label htmlFor="password" style={labelStyle}>
                  Password
                </label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Your password"
                  style={inputStyle}
                />
              </div>
              {mfaRequired && (
                <div>
                  <label htmlFor="code" style={labelStyle}>
                    Authentication code
                  </label>
                  <input
                    id="code"
                    name="code"
                    inputMode="numeric"
                    autoComplete="one-time-code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="6-digit code"
                    style={inputStyle}
                  />
                </div>
              )}
            </div>

            <div style={{ textAlign: "right", marginBottom: "20px" }}>
              {/* TODO: build /forgot-password flow (backend endpoint exists). */}
              <a
                href="#"
                style={{
                  fontSize: "14px",
                  color: "#C2683C",
                  fontWeight: 600,
                  cursor: "pointer",
                }}
              >
                Forgot password ?
              </a>
            </div>

            <button
              type="submit"
              disabled={loading}
              style={{
                width: "100%",
                border: "none",
                background: loading ? "#D89870" : "#C2683C",
                color: "#fff",
                fontFamily: "'Public Sans',sans-serif",
                fontSize: "16px",
                fontWeight: 600,
                padding: "15px",
                borderRadius: "100px",
                cursor: loading ? "default" : "pointer",
                marginBottom: "14px",
              }}
            >
              {loading ? "Signing in…" : "Sign in"}
            </button>
          </form>

          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "14px",
              margin: "18px 0",
            }}
          >
            <span style={{ flex: 1, height: "1px", background: "#EFE7DC" }} />
            <span style={{ fontSize: "13px", color: "#A89C8A" }}>or</span>
            <span style={{ flex: 1, height: "1px", background: "#EFE7DC" }} />
          </div>

          <button
            type="button"
            style={{
              width: "100%",
              border: "1px solid #E2D9CC",
              background: "#fff",
              color: "#2B2620",
              fontFamily: "'Public Sans',sans-serif",
              fontSize: "15px",
              fontWeight: 600,
              padding: "14px",
              borderRadius: "100px",
              cursor: "pointer",
            }}
          >
            Continue with Google
          </button>

          <div
            style={{
              textAlign: "center",
              marginTop: "28px",
              fontSize: "15px",
              color: "#5B554C",
            }}
          >
            Don&apos;t have an account ?{" "}
            <Link
              href="/sign-up"
              style={{ color: "#C2683C", fontWeight: 600, cursor: "pointer" }}
            >
              Create one free
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

const labelStyle: React.CSSProperties = {
  display: "block",
  fontSize: "13px",
  fontWeight: 600,
  color: "#8A8175",
  marginBottom: "7px",
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
