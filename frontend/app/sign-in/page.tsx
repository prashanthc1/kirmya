"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { 
  Sparkles, 
  Mail, 
  Lock, 
  ShieldCheck, 
  Loader2, 
  ArrowRight,
  UserCheck
} from "lucide-react";
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
    <div className="min-h-screen bg-background text-foreground flex flex-col md:flex-row">
      {/* Left Branding Panel (Desktop) */}
      <div className="hidden md:flex flex-1 bg-slate-900 dark:bg-zinc-950 p-12 flex-col justify-between relative overflow-hidden text-white border-r border-border/10">
        {/* Glow */}
        <div className="absolute bottom-[-100px] right-[-100px] w-96 h-96 rounded-full bg-blue-500/10 blur-[120px] pointer-events-none" />
        
        <Link href="/" className="text-xl font-bold tracking-tight">
          Kirmya
        </Link>

        <div className="space-y-6 max-w-md relative z-10">
          <div className="text-primary text-4xl font-extrabold">&ldquo;</div>
          <p className="text-2xl font-semibold leading-relaxed tracking-tight">
            One account. I&apos;m a job seeker, a mentor, and occasionally the one hiring. Kirmya handles all of it seamlessly.
          </p>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-xs">
              PN
            </div>
            <div>
              <p className="text-xs font-bold text-slate-100">Priya Nair</p>
              <p className="text-[10px] text-slate-400">Career Coach &bull; 3 roles active</p>
            </div>
          </div>
        </div>

        <p className="text-[10px] text-slate-500">&copy; {new Date().getFullYear()} Kirmya. Built for your comeback.</p>
      </div>

      {/* Right Login Form Panel */}
      <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background relative overflow-hidden">
        {/* Mobile Logo */}
        <div className="md:hidden absolute top-8 left-8">
          <Link href="/" className="text-lg font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
            Kirmya
          </Link>
        </div>

        <div className="w-full max-w-sm space-y-6">
          <div className="space-y-1.5 text-center md:text-left">
            <h1 className="text-2xl font-extrabold tracking-tight">Welcome back</h1>
            <p className="text-xs text-muted-foreground">Sign in to resume your active career search pipeline.</p>
          </div>

          {error && (
            <div className={`p-4 rounded-2xl text-xs font-medium border ${
              mfaRequired 
                ? "bg-amber-500/10 border-amber-500/20 text-amber-600 dark:text-amber-400" 
                : "bg-destructive/10 border-destructive/20 text-destructive"
            }`}>
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {!mfaRequired ? (
              <>
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-muted-foreground block">Email Address</label>
                  <div className="relative">
                    <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                    <input
                      type="email"
                      placeholder="name@company.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      required
                      className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-semibold text-muted-foreground block">Password</label>
                    <Link href="/forgot-password" className="text-[10px] font-bold text-primary hover:underline">
                      Forgot Password?
                    </Link>
                  </div>
                  <div className="relative">
                    <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                    <input
                      type="password"
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      required
                      className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                    />
                  </div>
                </div>
              </>
            ) : (
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground block">6-digit MFA Code</label>
                <div className="relative">
                  <ShieldCheck className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                  <input
                    type="text"
                    placeholder="000 000"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    required
                    maxLength={6}
                    className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                  />
                </div>
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-sm font-bold shadow-sm flex items-center justify-center gap-1.5"
            >
              {loading && <Loader2 className="h-4 w-4 animate-spin shrink-0" />}
              {mfaRequired ? "Verify Code" : "Sign In"}
              <ArrowRight className="h-4 w-4" />
            </button>
          </form>

          <div className="text-center text-xs text-muted-foreground pt-4 border-t border-border/40">
            Don&apos;t have an account?{" "}
            <Link href="/sign-up" className="font-bold text-primary hover:underline">
              Create one
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
