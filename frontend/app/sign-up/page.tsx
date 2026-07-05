"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { api, setAccessToken, ApiError } from "@/lib/api/client";
import { useAuth, type AuthUser } from "@/lib/auth/auth-context";

type Role = "job_seeker" | "referrer" | "mentor";

const ROLES: { value: Role; label: string }[] = [
  { value: "job_seeker", label: "Job seeker" },
  { value: "referrer", label: "Referrer" },
  { value: "mentor", label: "Mentor" },
];

interface RegisterResponse {
  access_token?: string;
  verification_required?: boolean;
  user?: AuthUser;
}

export default function SignUpPage() {
  const router = useRouter();
  const { setUser } = useAuth();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("job_seeker");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [verifyEmail, setVerifyEmail] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!fullName.trim()) {
      setError("Please enter your name.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }

    setLoading(true);
    try {
      const data = await api.post<RegisterResponse>("/auth/register", {
        full_name: fullName.trim(),
        email: email.trim(),
        password,
        role,
      });

      if (data?.verification_required) {
        // Account created but email verification is gated on — show a notice
        // instead of routing into the app.
        setVerifyEmail(email.trim());
        return;
      }

      if (data?.access_token) {
        // Auto-login: keep the token, populate auth state, head into the app.
        setAccessToken(data.access_token);
        if (data.user) setUser(data.user);
        router.push("/dashboard");
        return;
      }

      // Created without a session and without the verification flag — send the
      // user to sign in.
      router.push("/sign-in");
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : "Something went wrong creating your account. Please try again.",
      );
    } finally {
      setLoading(false);
    }
  }

  if (verifyEmail) {
    return (
      <main style={pageStyle}>
        <div style={cardStyle}>
          <h1 style={headingStyle}>Check your inbox</h1>
          <p style={{ fontSize: "15px", color: "#5B554C", lineHeight: 1.6 }}>
            We sent a verification link to <strong>{verifyEmail}</strong>. Confirm
            your email to finish setting up your account, then sign in.
          </p>
          <Link href="/sign-in" style={primaryLinkStyle}>
            Go to sign in
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main style={pageStyle}>
      <form style={cardStyle} onSubmit={handleSubmit} noValidate>
        <div style={{ marginBottom: "8px" }}>
          <div
            style={{
              fontFamily: "'Public Sans',sans-serif",
              fontSize: "22px",
              fontWeight: 800,
              letterSpacing: "-0.02em",
              color: "#2B2620",
            }}
          >
            Kirmya
          </div>
        </div>
        <h1 style={headingStyle}>Create your account</h1>
        <p style={{ fontSize: "15px", color: "#5B554C", margin: "0 0 22px" }}>
          Free to join. Built for your comeback.
        </p>

        {error && (
          <div role="alert" style={errorStyle}>
            {error}
          </div>
        )}

        <label style={labelStyle} htmlFor="name">
          Full name
        </label>
        <input
          id="name"
          name="name"
          autoComplete="name"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
          placeholder="Jordan Rivera"
          style={inputStyle}
        />

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

        <label style={labelStyle} htmlFor="password">
          Password
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

        <label style={labelStyle} htmlFor="role">
          I&apos;m joining as
        </label>
        <select
          id="role"
          name="role"
          value={role}
          onChange={(e) => setRole(e.target.value as Role)}
          style={{ ...inputStyle, cursor: "pointer" }}
        >
          {ROLES.map((r) => (
            <option key={r.value} value={r.value}>
              {r.label}
            </option>
          ))}
        </select>

        <button type="submit" disabled={loading} style={buttonStyle(loading)}>
          {loading ? "Creating account…" : "Create one free"}
        </button>

        <div
          style={{
            textAlign: "center",
            marginTop: "22px",
            fontSize: "15px",
            color: "#5B554C",
          }}
        >
          Already have an account?{" "}
          <Link href="/sign-in" style={{ color: "#C2683C", fontWeight: 600 }}>
            Sign in
          </Link>
        </div>
      </form>
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
  fontFamily: "'Public Sans',sans-serif",
  fontWeight: 800,
  fontSize: "28px",
  letterSpacing: "-0.02em",
  margin: "0 0 6px",
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
  marginBottom: "4px",
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
