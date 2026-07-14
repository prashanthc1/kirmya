"use client";

import React, { useState, use } from "react";
import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { Award, Calendar, CheckCircle2, Clock, Globe, GraduationCap, Link2, MapPin, Sparkles } from "lucide-react";

interface Mentor {
  username: string;
  name: string;
  title: string;
  company: string;
  location: string;
  bio: string;
  rating: string;
  sessions: string;
  topics: string[];
  experience: string;
}

const MOCK_MENTORS: Record<string, Mentor> = {
  "sarah-jenkins": {
    username: "sarah-jenkins",
    name: "Sarah Jenkins",
    title: "Director of UX",
    company: "Notion",
    location: "San Francisco, CA / Remote",
    bio: "Helps designers trace system architectures, build portfolio pieces, and navigate transition phases after sudden layoffs.",
    rating: "4.9",
    sessions: "124",
    topics: ["Portfolio Critique", "Resume & Profile Review", "Design Leadership"],
    experience: "12 years",
  },
  "marcus-chen": {
    username: "marcus-chen",
    name: "Marcus Chen",
    title: "Engineering VP",
    company: "Linear",
    location: "Remote",
    bio: "Ex-Stripe, ex-Meta lead engineer focused on high-scale distributed systems, database schema layouts, and Go backend patterns.",
    rating: "5.0",
    sessions: "210",
    topics: ["Backend System Design", "Engineering Management", "Technical Mock Interviews"],
    experience: "15 years",
  },
  "asha-rao": {
    username: "asha-rao",
    name: "Asha Rao",
    title: "Principal Recruiter",
    company: "Stripe",
    location: "SF / Hybrid",
    bio: "Career transition expert. Provides insider tips on ATS resume optimization, screening benchmarks, and salary counters.",
    rating: "4.8",
    sessions: "95",
    topics: ["Salary Negotiation", "ATS Optimization", "Recruiting Strategy"],
    experience: "10 years",
  },
  "devon-webb": {
    username: "devon-webb",
    name: "Devon Webb",
    title: "Senior iOS Lead",
    company: "Vercel",
    location: "Remote (Global)",
    bio: "Passionate about mobile apps, Next.js setups, and helping engineers recover from transition gaps into top mobile roles.",
    rating: "4.9",
    sessions: "68",
    topics: ["iOS Architecture", "Vite/Next.js Deployments", "Transition Strategy"],
    experience: "8 years",
  },
};

interface PageProps {
  params: Promise<{ username: string }> | { username: string };
}

function MentorProfileContent({ username }: { username: string }) {
  const mentor = MOCK_MENTORS[username] || {
    username,
    name: username.split("-").map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(" "),
    title: "Professional Advisor",
    company: "Kirmya Network",
    location: "Remote",
    bio: "Experienced professional offering transition guidance, skill optimization, and interview confidence coaching.",
    rating: "4.9",
    sessions: "10+",
    topics: ["Career Strategy", "Resume Optimization", "Mock Interview"],
    experience: "5+ years",
  };

  const [booked, setBooked] = useState(false);
  const [selectedSlot, setSelectedSlot] = useState<string | null>(null);

  const SLOTS = ["Tomorrow, 10:00 AM", "Tomorrow, 2:00 PM", "Wednesday, 11:00 AM", "Wednesday, 4:00 PM"];

  const handleBook = () => {
    if (!selectedSlot) return;
    setBooked(true);
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Mentorship", href: "/mentorship" }, { label: mentor.name }]} />

      <main className="flex-grow max-w-4xl mx-auto px-4 sm:px-6 py-8 w-full">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          {/* Main Info Column */}
          <div className="lg:col-span-2 space-y-6">
            <div className="bg-card border border-border/60 p-6 rounded-3xl shadow-sm space-y-4">
              <div className="flex gap-4">
                <div className="h-16 w-16 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-500 font-bold text-2xl select-none">
                  {mentor.name.charAt(0)}
                </div>
                <div className="space-y-1">
                  <h1 className="text-xl font-black tracking-tight">{mentor.name}</h1>
                  <p className="text-sm font-semibold text-muted-foreground">
                    {mentor.title} at <span className="text-foreground">{mentor.company}</span>
                  </p>
                  <div className="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <MapPin className="h-3.5 w-3.5" />
                      {mentor.location}
                    </span>
                    <span className="flex items-center gap-1">
                      <GraduationCap className="h-3.5 w-3.5" />
                      {mentor.experience} experience
                    </span>
                  </div>
                </div>
              </div>

              <div className="border-t border-border/40 pt-4">
                <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-2">About Me</h3>
                <p className="text-xs leading-relaxed text-muted-foreground/90">{mentor.bio}</p>
              </div>

              <div className="border-t border-border/40 pt-4">
                <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-2.5">Focus Areas</h3>
                <div className="flex flex-wrap gap-2">
                  {mentor.topics.map((t, idx) => (
                    <span
                      key={idx}
                      className="px-2.5 py-0.5 rounded-full border border-border text-[10px] font-semibold bg-secondary/30"
                    >
                      {t}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Booking Sidebar Column */}
          <div className="space-y-6">
            <div className="bg-card border border-border/60 p-6 rounded-3xl shadow-sm space-y-4">
              <div className="flex justify-between items-center pb-3 border-b border-border/40">
                <h3 className="text-sm font-bold">Book a Session</h3>
                <span className="text-xs font-extrabold text-emerald-500 bg-emerald-500/10 border border-emerald-500/20 px-2 py-0.5 rounded-full">
                  FREE
                </span>
              </div>

              {booked ? (
                <div className="text-center py-6 space-y-3">
                  <div className="h-10 w-10 bg-emerald-500/10 border border-emerald-500/20 rounded-full flex items-center justify-center mx-auto text-emerald-500">
                    <CheckCircle2 className="h-5 w-5" />
                  </div>
                  <div className="space-y-1">
                    <h4 className="text-xs font-bold text-foreground">Session Requested!</h4>
                    <p className="text-[10px] text-muted-foreground">
                      Marcus has been notified. You will receive an email confirmation with calendar details.
                    </p>
                  </div>
                  <button
                    onClick={() => setBooked(false)}
                    className="w-full py-1.5 rounded-full border border-border text-xs font-semibold hover:bg-secondary transition-all cursor-pointer"
                  >
                    Book another slot
                  </button>
                </div>
              ) : (
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground block">
                      Select Slot
                    </label>
                    <div className="space-y-1.5">
                      {SLOTS.map((slot) => (
                        <button
                          key={slot}
                          onClick={() => setSelectedSlot(slot)}
                          className={`w-full text-left px-3.5 py-2.5 rounded-xl text-xs transition-all border ${
                            selectedSlot === slot
                              ? "bg-primary/10 border-primary text-primary font-semibold"
                              : "bg-secondary/20 border-border/40 hover:border-border text-muted-foreground hover:text-foreground"
                          }`}
                        >
                          {slot}
                        </button>
                      ))}
                    </div>
                  </div>

                  <button
                    onClick={handleBook}
                    disabled={!selectedSlot}
                    className="w-full py-2 bg-primary hover:bg-primary/95 disabled:bg-muted disabled:text-muted-foreground text-primary-foreground text-xs font-bold rounded-full transition-all shadow-sm cursor-pointer"
                  >
                    Confirm Booking
                  </button>
                </div>
              )}
            </div>
          </div>

        </div>
      </main>

      <SiteFooter />
    </div>
  );
}

export default function MentorProfilePage({ params }: PageProps) {
  // Next.js 15 uses promises for route parameters; safe lookup support for all versions:
  const resolvedParams = params && "then" in params ? use(params) : (params as { username: string });
  const username = resolvedParams?.username || "";

  return (
    <AuthGuard>
      <MentorProfileContent username={username} />
    </AuthGuard>
  );
}
