"use client";

import { useState } from "react";
import Link from "next/link";
import {
  Sparkles,
  Mail,
  ArrowLeft,
  Loader2,
  ArrowRight,
  ShieldCheck,
} from "lucide-react";
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
    <div className="min-h-screen bg-background text-foreground flex items-center justify-center p-6 relative overflow-hidden">
      {/* Glow Orbs */}
      <div className="absolute top-[-10%] left-[-10%] w-[350px] h-[350px] rounded-full bg-blue-500/5 blur-[100px] pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[350px] h-[350px] rounded-full bg-indigo-500/5 blur-[100px] pointer-events-none" />

      <div className="w-full max-w-sm">
        {/* Brand Header */}
        <div className="text-center mb-8">
          <Link
            href="/"
            className="text-xl font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
          >
            Kirmya
          </Link>
        </div>

        {sent ? (
          <div className="bg-card border border-border/80 p-8 rounded-3xl text-center space-y-6 shadow-lg shadow-black/5">
            <div className="h-12 w-12 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary mx-auto">
              <Mail className="h-6 w-6" />
            </div>
            <div className="space-y-2">
              <h1 className="text-xl font-bold">Check your inbox</h1>
              <p className="text-xs text-muted-foreground leading-relaxed">
                If an account exists for{" "}
                <strong className="text-foreground">{email.trim()}</strong>,
                we&apos;ve sent a password-reset link. It expires in 1 hour.
              </p>
            </div>
            <Link
              href="/sign-in"
              className="w-full py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold block shadow-sm"
            >
              Back to Sign In
            </Link>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            className="bg-card border border-border/80 p-8 rounded-3xl space-y-5 shadow-lg shadow-black/5"
            noValidate
          >
            <div className="space-y-1.5 text-center sm:text-left">
              <h1 className="text-xl font-extrabold tracking-tight">
                Reset password
              </h1>
              <p className="text-xs text-muted-foreground">
                Enter your email address and we&apos;ll send you a link to reset
                your password.
              </p>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-muted-foreground block">
                Email Address
              </label>
              <div className="relative">
                <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                <input
                  type="email"
                  placeholder="name@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading || !email.trim()}
              className="w-full py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-sm font-bold shadow-sm flex items-center justify-center gap-1.5"
            >
              {loading && <Loader2 className="h-4 w-4 animate-spin shrink-0" />}
              Send Reset Link
              <ArrowRight className="h-4 w-4" />
            </button>

            <div className="text-center text-xs text-muted-foreground pt-4 border-t border-border/40">
              Remembered it?{" "}
              <Link
                href="/sign-in"
                className="font-bold text-primary hover:underline"
              >
                Sign in
              </Link>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
