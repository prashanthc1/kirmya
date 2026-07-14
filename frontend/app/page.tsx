"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";
import { motion, AnimatePresence } from "framer-motion";
import {
  Search,
  Sparkles,
  Briefcase,
  Compass,
  GraduationCap,
  Users,
  FileText,
  TrendingUp,
  DollarSign,
  Check,
  ChevronDown,
  ArrowRight,
  ShieldCheck,
  Building,
  UploadCloud,
  Cpu,
  BrainCircuit,
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

const COMPANIES = [
  "Stripe",
  "Linear",
  "Vercel",
  "Airbnb",
  "Notion",
  "Framer",
  "Figma",
];

const HERO_FEATURES = [
  {
    icon: FileText,
    title: "AI Resume Optimization",
    description:
      "Scan your resume against real-time ATS parser guidelines and insert missing industry terms.",
  },
  {
    icon: Compass,
    title: "Employee Referral Networks",
    description:
      "Connect with verified employees at top companies and secure warm internal referrals.",
  },
  {
    icon: GraduationCap,
    title: "Skills Gap Training",
    description:
      "Visually trace the missing technical skills you need to qualify for high-tier open positions.",
  },
];

const PREVIEW_JOBS = [
  {
    company: "Linear",
    logo: "L",
    title: "Senior Product Designer",
    location: "Remote (US)",
    salary: "$160k – $210k",
    match: 98,
    skills: ["Figma", "Design Systems", "Prototyping"],
    type: "Full-time",
  },
  {
    company: "Vercel",
    logo: "V",
    title: "Staff Frontend Engineer",
    location: "Remote (Global)",
    salary: "$180k – $240k",
    match: 94,
    skills: ["React", "Next.js", "Tailwind CSS"],
    type: "Full-time",
  },
  {
    company: "Stripe",
    logo: "S",
    title: "Product Engineer (Payments)",
    location: "San Francisco, CA",
    salary: "$170k – $220k",
    match: 89,
    skills: ["Ruby", "React", "System Design"],
    type: "Hybrid",
  },
];



const FAQS = [
  {
    q: "Is Kirmya really free to use?",
    a: "Yes. During career transitions, the last thing you need is another bill. The core tools—resume parsing, AI coaching suggestions, community access, and mentorship sessions—are free. We focus on one metric: the speed at which you land your next interview.",
  },
  {
    q: "How does the referral system prevent spam?",
    a: "We do not allow cold, automated messaging. Referral requests require structured introductions, target specific open roles, and are matched based on mutual interest, protecting referrers from spam while keeping candidate signal high.",
  },
  {
    q: "What makes Kirmya different from LinkedIn?",
    a: "LinkedIn optimizes for screen time, influencers, and vanity metrics. Kirmya is a recovery workspace. There is no public content feed. We provide tools to directly prepare you for interviews and connect with internal champions.",
  },
];

export default function HomePage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [activeFaq, setActiveFaq] = useState<number | null>(null);
  const [resumeScore, setResumeScore] = useState<number | null>(null);
  const [scanning, setScanning] = useState(false);

  useEffect(() => {
    if (!loading && user) {
      router.replace("/dashboard");
    }
  }, [user, loading, router]);

  const handleMockScan = () => {
    if (scanning) return;
    setScanning(true);
    setResumeScore(null);
    setTimeout(() => {
      setScanning(false);
      setResumeScore(87);
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col relative overflow-hidden">
      {/* Background Orbs */}
      <div className="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] rounded-full bg-blue-500/10 blur-[120px] pointer-events-none" />
      <div className="absolute top-[20%] right-[-10%] w-[600px] h-[600px] rounded-full bg-indigo-500/10 blur-[150px] pointer-events-none" />

      <SiteNav />

      {/* Hero Section */}
      <section className="relative pt-20 pb-16 md:pt-32 md:pb-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center relative z-10">
          {/* Tagline */}
          <motion.div
            initial={{ opacity: 0, y: 15 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20 mb-6"
          >
            <Sparkles className="h-3 w-3" />
            Built for the moment between jobs
          </motion.div>

          {/* Heading */}
          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-4xl sm:text-6xl md:text-7xl font-extrabold tracking-tight mb-6 leading-[1.1] max-w-5xl mx-auto"
          >
            You didn&apos;t lose your career. <br />
            <span className="bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 bg-clip-text text-transparent">
              You just lost that one job.
            </span>
          </motion.h1>

          {/* Subheading */}
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-lg md:text-xl text-muted-foreground max-w-3xl mx-auto mb-10 leading-relaxed"
          >
            The gap on your resume isn&apos;t a red flag—it&apos;s a chapter.
            Kirmya is the AI-powered operating system where top-tier
            professionals regroup, optimize materials, swap referrals, and land
            interviews.
          </motion.p>

          {/* Call to Actions */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16"
          >
            <Link
              href="/sign-up"
              className="w-full sm:w-auto px-8 py-3.5 rounded-full text-sm font-semibold bg-primary text-primary-foreground hover:bg-primary/95 transition-all shadow-lg shadow-blue-500/10 flex items-center justify-center gap-2 group"
            >
              Start Free Today
              <ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
            </Link>
            <Link
              href="/jobs"
              className="w-full sm:w-auto px-8 py-3.5 rounded-full text-sm font-semibold border border-border hover:bg-secondary transition-all flex items-center justify-center"
            >
              Browse Open Roles
            </Link>
          </motion.div>

          {/* Trusted Companies */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.8, delay: 0.4 }}
            className="border-t border-border/40 pt-8"
          >
            <p className="text-xs font-semibold text-muted-foreground uppercase tracking-widest mb-6">
              Connect with professionals at companies like
            </p>
            <div className="flex flex-wrap justify-center items-center gap-8 md:gap-12 opacity-60">
              {COMPANIES.map((company, idx) => (
                <span
                  key={idx}
                  className="text-lg font-bold tracking-tight text-foreground select-none"
                >
                  {company}
                </span>
              ))}
            </div>
          </motion.div>
        </div>
      </section>

      {/* Feature grid */}
      <section className="py-20 bg-secondary/35 border-y border-border/40 relative">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-3xl font-extrabold tracking-tight mb-4">
              Everything you need for a faster comeback
            </h2>
            <p className="text-muted-foreground text-base">
              Kirmya replaces traditional, passive job searching with active,
              AI-assisted tools that prioritize candidate placement speed.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {HERO_FEATURES.map((feature, idx) => {
              const Icon = feature.icon;
              return (
                <div
                  key={idx}
                  className="bg-card border border-border/60 p-8 rounded-3xl space-y-4 hover:shadow-lg hover:border-border transition-all duration-300"
                >
                  <div className="h-10 w-10 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary">
                    <Icon className="h-5 w-5" />
                  </div>
                  <h3 className="text-lg font-bold">{feature.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    {feature.description}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      {/* AI Sandbox Interactive Showcase */}
      <section className="py-20">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Left Column: Copy */}
            <div className="space-y-6">
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20">
                <BrainCircuit className="h-3.5 w-3.5" />
                AI Career Co-Pilot
              </div>
              <h2 className="text-3xl sm:text-4xl font-extrabold tracking-tight leading-tight">
                Benchmark your profile with automated intelligence
              </h2>
              <p className="text-muted-foreground leading-relaxed">
                Our parsing model reads your resume just like a corporate
                Applicant Tracking System (ATS), scores it against live job
                parameters, and offers structural corrections instantly.
              </p>

              <ul className="space-y-3">
                {[
                  "Visual match score against real open positions",
                  "Missing keywords highlighted in red for instant additions",
                  "Direct recommendations for skill gaps and online courses",
                ].map((item, idx) => (
                  <li key={idx} className="flex items-start gap-2.5 text-sm">
                    <div className="mt-0.5 h-4 w-4 rounded-full bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-500">
                      <Check className="h-3 w-3" />
                    </div>
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* Right Column: Interactive Sandbox */}
            <div className="bg-card border border-border p-6 rounded-3xl shadow-xl relative overflow-hidden">
              <div className="flex items-center justify-between border-b border-border/40 pb-4 mb-6">
                <div className="flex items-center gap-2">
                  <Cpu className="h-4 w-4 text-primary" />
                  <span className="text-sm font-bold">
                    Resume ATS Simulator
                  </span>
                </div>
                <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              </div>

              <div className="space-y-4">
                <div
                  className="border border-dashed border-border/80 p-8 rounded-2xl text-center space-y-3 bg-secondary/25 hover:bg-secondary/40 transition-colors cursor-pointer"
                  onClick={handleMockScan}
                >
                  <UploadCloud className="h-8 w-8 text-muted-foreground mx-auto" />
                  <div>
                    <p className="text-sm font-bold">
                      Click here to upload your resume
                    </p>
                    <p className="text-xs text-muted-foreground">
                      PDF or Word files, max 5MB
                    </p>
                  </div>
                </div>

                {scanning && (
                  <div className="space-y-2">
                    <div className="h-1 bg-secondary rounded-full overflow-hidden">
                      <motion.div
                        className="h-full bg-primary"
                        initial={{ width: 0 }}
                        animate={{ width: "100%" }}
                        transition={{ duration: 2 }}
                      />
                    </div>
                    <p className="text-xs text-center text-muted-foreground">
                      Scanning resume keywords against Live ATS databases...
                    </p>
                  </div>
                )}

                {resumeScore !== null && (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="p-4 bg-secondary/50 rounded-2xl space-y-3"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-semibold">
                        Resume Strength Index
                      </span>
                      <span className="text-sm font-extrabold text-emerald-500">
                        {resumeScore}% Match
                      </span>
                    </div>
                    <div className="space-y-1">
                      <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">
                        Suggested Enhancements
                      </p>
                      <p className="text-xs">
                        Add{" "}
                        <strong className="text-primary font-semibold">
                          Next.js App Router
                        </strong>{" "}
                        and{" "}
                        <strong className="text-primary font-semibold">
                          CI/CD Pipeline
                        </strong>{" "}
                        to match targeted React roles.
                      </p>
                    </div>
                  </motion.div>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Jobs preview */}
      <section className="py-20 bg-secondary/15 border-y border-border/40">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col sm:flex-row sm:items-end justify-between mb-12 gap-4">
            <div>
              <h2 className="text-3xl font-extrabold tracking-tight">
                Featured Live Roles
              </h2>
              <p className="text-muted-foreground text-sm mt-2">
                Roles optimized for our community members featuring warm
                internal referrals.
              </p>
            </div>
            <Link
              href="/jobs"
              className="text-sm font-bold text-primary flex items-center gap-1 hover:underline"
            >
              Search all open positions
              <ArrowRight className="h-4 w-4" />
            </Link>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {PREVIEW_JOBS.map((job, idx) => (
              <div
                key={idx}
                className="bg-card border border-border/60 p-6 rounded-3xl space-y-6 flex flex-col justify-between hover:shadow-lg hover:border-border transition-all"
              >
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2.5">
                      <div className="h-8 w-8 rounded-lg bg-secondary flex items-center justify-center text-sm font-bold border border-border/40 select-none">
                        {job.logo}
                      </div>
                      <span className="text-xs font-semibold text-muted-foreground">
                        {job.company}
                      </span>
                    </div>
                    <div className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-extrabold text-emerald-500 flex items-center gap-0.5">
                      <Sparkles className="h-2.5 w-2.5" />
                      {job.match}% AI Match
                    </div>
                  </div>

                  <h3 className="text-base font-bold truncate">{job.title}</h3>

                  <div className="flex flex-wrap gap-1.5">
                    {job.skills.map((skill, sIdx) => (
                      <span
                        key={sIdx}
                        className="text-[10px] px-2 py-0.5 bg-secondary text-muted-foreground rounded-full border border-border/40"
                      >
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>

                <div className="border-t border-border/40 pt-4 flex items-center justify-between text-xs">
                  <span className="text-muted-foreground">{job.location}</span>
                  <span className="font-bold">{job.salary}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>



      {/* FAQ Accordions */}
      <section className="py-20 bg-secondary/15 border-t border-border/40">
        <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-extrabold tracking-tight text-center mb-12">
            Frequently Asked Questions
          </h2>

          <div className="space-y-4">
            {FAQS.map((faq, idx) => {
              const isOpen = activeFaq === idx;
              return (
                <div
                  key={idx}
                  className="border border-border/80 bg-card rounded-2xl overflow-hidden transition-all duration-300"
                >
                  <button
                    onClick={() => setActiveFaq(isOpen ? null : idx)}
                    className="w-full flex items-center justify-between p-5 text-left font-bold text-sm focus:outline-none"
                  >
                    <span>{faq.q}</span>
                    <ChevronDown
                      className={`h-4 w-4 text-muted-foreground transition-transform duration-300 ${isOpen ? "rotate-180" : ""}`}
                    />
                  </button>

                  <AnimatePresence initial={false}>
                    {isOpen && (
                      <motion.div
                        initial={{ height: 0, opacity: 0 }}
                        animate={{ height: "auto", opacity: 1 }}
                        exit={{ height: 0, opacity: 0 }}
                        transition={{ duration: 0.2 }}
                      >
                        <div className="px-5 pb-5 pt-1 text-sm text-muted-foreground leading-relaxed border-t border-border/20">
                          {faq.a}
                        </div>
                      </motion.div>
                    )}
                  </AnimatePresence>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      <SiteFooter />
    </div>
  );
}
