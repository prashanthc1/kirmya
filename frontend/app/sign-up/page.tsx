"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { 
  Sparkles, 
  Mail, 
  Lock, 
  User, 
  ChevronRight, 
  Loader2, 
  ArrowRight,
  Briefcase,
  Compass,
  GraduationCap
} from "lucide-react";
import { api, setAccessToken, ApiError } from "@/lib/api/client";
import { useAuth, type AuthUser } from "@/lib/auth/auth-context";

type Role = "job_seeker" | "referrer" | "mentor";

const ROLES: { value: Role; label: string; desc: string; icon: React.ComponentType<{ className?: string }> }[] = [
  { value: "job_seeker", label: "Job seeker", desc: "Discover openings & practice prep", icon: Briefcase },
  { value: "referrer", label: "Referrer", desc: "Introduce talents into my company", icon: Compass },
  { value: "mentor", label: "Mentor", desc: "Share experience & advise peers", icon: GraduationCap },
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
        setVerifyEmail(email.trim());
        return;
      }

      if (data?.access_token) {
        setAccessToken(data.access_token);
        if (data.user) setUser(data.user);
        router.push("/dashboard");
        return;
      }

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
      <div className="min-h-screen bg-background text-foreground flex items-center justify-center p-6">
        <div className="w-full max-w-sm bg-card border border-border/80 p-8 rounded-3xl text-center space-y-6 shadow-lg shadow-black/5">
          <div className="h-12 w-12 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary mx-auto">
            <Mail className="h-6 w-6" />
          </div>
          <div className="space-y-2">
            <h1 className="text-xl font-bold">Check your inbox</h1>
            <p className="text-xs text-muted-foreground leading-relaxed">
              We sent a verification link to <strong className="text-foreground">{verifyEmail}</strong>. Confirm your email to finish setting up your account, then sign in.
            </p>
          </div>
          <Link
            href="/sign-in"
            className="w-full py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold block shadow-sm"
          >
            Go to Sign In
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col md:flex-row">
      {/* Left Branding Panel (Desktop) */}
      <div className="hidden md:flex flex-1 bg-slate-900 dark:bg-zinc-950 p-12 flex-col justify-between relative overflow-hidden text-white border-r border-border/10">
        {/* Glow */}
        <div className="absolute bottom-[-100px] right-[-100px] w-96 h-96 rounded-full bg-indigo-500/10 blur-[120px] pointer-events-none" />
        
        <Link href="/" className="text-xl font-bold tracking-tight">
          Kirmya
        </Link>

        <div className="space-y-6 max-w-md relative z-10">
          <div className="text-primary text-4xl font-extrabold">&ldquo;</div>
          <p className="text-2xl font-semibold leading-relaxed tracking-tight">
            Designed to guide you through transition. Discover warm connections, polish interview scripts, and come back stronger.
          </p>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-xs">
              MH
            </div>
            <div>
              <p className="text-xs font-bold text-slate-100">Marcus Hale</p>
              <p className="text-[10px] text-slate-400">Operations Lead &bull; Rebuilt in 2026</p>
            </div>
          </div>
        </div>

        <p className="text-[10px] text-slate-500">&copy; {new Date().getFullYear()} Kirmya. Built for your comeback.</p>
      </div>

      {/* Right Sign-Up Form Panel */}
      <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background relative overflow-hidden">
        {/* Mobile Logo */}
        <div className="md:hidden absolute top-8 left-8">
          <Link href="/" className="text-lg font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
            Kirmya
          </Link>
        </div>

        <div className="w-full max-w-md space-y-6">
          <div className="space-y-1.5 text-center md:text-left">
            <h1 className="text-2xl font-extrabold tracking-tight">Create your account</h1>
            <p className="text-xs text-muted-foreground">Free to join. Built for your professional comeback.</p>
          </div>

          {error && (
            <div className="p-4 rounded-2xl bg-destructive/10 border border-destructive/20 text-destructive text-xs font-medium">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4" noValidate>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground block">Full Name</label>
                <div className="relative">
                  <User className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                  <input
                    id="name"
                    type="text"
                    placeholder="Jordan Rivera"
                    value={fullName}
                    onChange={(e) => setFullName(e.target.value)}
                    required
                    className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground block">Email Address</label>
                <div className="relative">
                  <Mail className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                  <input
                    id="email"
                    type="email"
                    placeholder="name@company.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                  />
                </div>
              </div>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-muted-foreground block">Password</label>
              <div className="relative">
                <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                <input
                  id="password"
                  type="password"
                  placeholder="Min. 8 characters"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  className="w-full pl-10 pr-4 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                />
              </div>
            </div>

            {/* Select Role Panel */}
            <div className="space-y-2">
              <label className="text-xs font-semibold text-muted-foreground block">Select primary account type</label>
              <div className="grid grid-cols-1 gap-2">
                {ROLES.map((item) => {
                  const Icon = item.icon;
                  const isSelected = role === item.value;
                  return (
                    <div
                      key={item.value}
                      onClick={() => setRole(item.value)}
                      className={`p-3 border rounded-2xl cursor-pointer flex items-center gap-3 transition-all ${
                        isSelected 
                          ? "bg-primary/5 border-primary shadow-sm" 
                          : "border-border bg-card hover:bg-secondary/40"
                      }`}
                    >
                      <div className={`h-8 w-8 rounded-lg flex items-center justify-center shrink-0 border transition-all ${
                        isSelected
                          ? "bg-primary text-primary-foreground border-primary"
                          : "bg-secondary text-muted-foreground border-border/40"
                      }`}>
                        <Icon className="h-4.5 w-4.5" />
                      </div>
                      <div className="text-left">
                        <span className="text-xs font-bold block">{item.label}</span>
                        <span className="text-[10px] text-muted-foreground block leading-tight">{item.desc}</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full py-2.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-sm font-bold shadow-sm flex items-center justify-center gap-1.5"
            >
              {loading && <Loader2 className="h-4 w-4 animate-spin shrink-0" />}
              Create one free
              <ArrowRight className="h-4 w-4" />
            </button>
          </form>

          <div className="text-center text-xs text-muted-foreground pt-4 border-t border-border/40">
            Already have an account?{" "}
            <Link href="/sign-in" className="font-bold text-primary hover:underline">
              Sign in
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
