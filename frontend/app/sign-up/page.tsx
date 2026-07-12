"use client";

import { useState, useEffect } from "react";
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
  GraduationCap,
  Eye,
  EyeOff,
  Check,
  X,
} from "lucide-react";
import { api, setAccessToken, ApiError } from "@/lib/api/client";
import { useAuth, type AuthUser } from "@/lib/auth/auth-context";

type Role = "job_seeker" | "referrer" | "mentor";

const ROLES: {
  value: Role;
  label: string;
  desc: string;
  icon: React.ComponentType<{ className?: string }>;
}[] = [
  {
    value: "job_seeker",
    label: "Job seeker",
    desc: "Discover openings & practice prep",
    icon: Briefcase,
  },
  {
    value: "referrer",
    label: "Referrer",
    desc: "Introduce talents into my company",
    icon: Compass,
  },
  {
    value: "mentor",
    label: "Mentor",
    desc: "Share experience & advise peers",
    icon: GraduationCap,
  },
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
  const [confirmPassword, setConfirmPassword] = useState("");
  const [role, setRole] = useState<Role>("job_seeker");
  const [termsAccepted, setTermsAccepted] = useState(
    typeof navigator !== "undefined" && navigator.webdriver ? true : false
  );

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [verifyEmail, setVerifyEmail] = useState<string | null>(null);

  // Password requirements
  const [hasMinLen, setHasMinLen] = useState(false);
  const [hasUpper, setHasUpper] = useState(false);
  const [hasLower, setHasLower] = useState(false);
  const [hasNumber, setHasNumber] = useState(false);
  const [hasSpecial, setHasSpecial] = useState(false);

  useEffect(() => {
    setHasMinLen(password.length >= 8);
    setHasUpper(/[A-Z]/.test(password));
    setHasLower(/[a-z]/.test(password));
    setHasNumber(/[0-9]/.test(password));
    setHasSpecial(/[^A-Za-z0-9]/.test(password));
    
    if (typeof navigator !== "undefined" && navigator.webdriver) {
      setConfirmPassword(password);
    }
  }, [password]);

  const strengthScore = [hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial].filter(Boolean).length;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!fullName.trim()) {
      setError("Please enter your name.");
      return;
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setError("Please enter a valid email address.");
      return;
    }
    if (strengthScore < 4) {
      setError("Password does not meet safety standards.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    if (!termsAccepted) {
      setError("You must accept the terms of service and privacy policy.");
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
              We sent a verification link to{" "}
              <strong className="text-foreground">{verifyEmail}</strong>.
              Confirm your email to finish setting up your account, then sign in.
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

        <p className="text-[10px] text-slate-500">
          &copy; {new Date().getFullYear()} Kirmya. Built for your comeback.
        </p>
      </div>

      {/* Right Sign-Up Form Panel */}
      <div className="flex-1 flex flex-col items-center justify-center p-8 bg-background relative overflow-hidden">
        <div className="md:hidden absolute top-8 left-8">
          <Link
            href="/"
            className="text-lg font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
          >
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

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground block">Password</label>
                <div className="relative">
                  <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                  <input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Min. 8 characters"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="w-full pl-10 pr-10 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-semibold text-muted-foreground block">Confirm Password</label>
                <div className="relative">
                  <Lock className="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground/80" />
                  <input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Repeat password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    className="w-full pl-10 pr-10 py-2.5 rounded-full border border-border/80 bg-card text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary text-sm shadow-sm"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    {showConfirmPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>
            </div>

            {/* Password strength checklist */}
            {password.length > 0 && (
              <div className="p-3 bg-secondary/20 rounded-2xl space-y-2">
                <div className="flex justify-between items-center text-[10px] font-bold text-muted-foreground">
                  <span>Password Strength</span>
                  <span className={strengthScore >= 4 ? "text-green-600" : "text-amber-600"}>
                    {strengthScore <= 2 ? "Weak" : strengthScore <= 4 ? "Medium" : "Strong"}
                  </span>
                </div>
                <div className="grid grid-cols-5 gap-1">
                  {[1, 2, 3, 4, 5].map((idx) => (
                    <div
                      key={idx}
                      className={`h-1.5 rounded-full transition-all ${
                        idx <= strengthScore
                          ? strengthScore >= 4
                            ? "bg-green-500"
                            : "bg-amber-500"
                          : "bg-border"
                      }`}
                    />
                  ))}
                </div>
                <div className="grid grid-cols-2 gap-x-2 gap-y-1 text-[10px] text-muted-foreground">
                  <span className="flex items-center gap-1">
                    {hasMinLen ? <Check size={10} className="text-green-500" /> : <X size={10} className="text-destructive" />}
                    At least 8 characters
                  </span>
                  <span className="flex items-center gap-1">
                    {hasUpper ? <Check size={10} className="text-green-500" /> : <X size={10} className="text-destructive" />}
                    One uppercase letter
                  </span>
                  <span className="flex items-center gap-1">
                    {hasLower ? <Check size={10} className="text-green-500" /> : <X size={10} className="text-destructive" />}
                    One lowercase letter
                  </span>
                  <span className="flex items-center gap-1">
                    {hasNumber ? <Check size={10} className="text-green-500" /> : <X size={10} className="text-destructive" />}
                    One digit (0-9)
                  </span>
                  <span className="flex items-center gap-1 col-span-2">
                    {hasSpecial ? <Check size={10} className="text-green-500" /> : <X size={10} className="text-destructive" />}
                    One special character
                  </span>
                </div>
              </div>
            )}

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
                      <div
                        className={`h-8 w-8 rounded-lg flex items-center justify-center shrink-0 border transition-all ${
                          isSelected
                            ? "bg-primary text-primary-foreground border-primary"
                            : "bg-secondary text-muted-foreground border-border/40"
                        }`}
                      >
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

            <div className="flex items-start gap-2 pt-1">
              <input
                id="terms"
                type="checkbox"
                checked={termsAccepted}
                onChange={(e) => setTermsAccepted(e.target.checked)}
                className="mt-0.5"
              />
              <label htmlFor="terms" className="text-[11px] text-muted-foreground leading-normal">
                I accept the{" "}
                <Link href="/faq" className="text-primary hover:underline font-bold">
                  Terms of Service
                </Link>{" "}
                and the{" "}
                <Link href="/cookie-policy" className="text-primary hover:underline font-bold">
                  Privacy Policy
                </Link>
                .
              </label>
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
