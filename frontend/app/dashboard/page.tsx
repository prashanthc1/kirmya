"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";
import AuthGuard from "@/components/shared/AuthGuard";
import {
  Briefcase,
  Bookmark,
  Users,
  CheckCircle,
  MapPin,
  Edit,
  Eye,
  MessageSquare,
  UserPlus,
  UserCheck,
  X,
  TrendingUp,
  Sparkles,
  ArrowRight,
  Send,
  Loader2,
  Clock,
  FileText,
  GraduationCap,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api } from "@/lib/api/client";
import { profileClient, Profile } from "@/lib/api/profile";
import { connectionsClient, Connection } from "@/lib/api/connections";
import SuggestionsCarousel from "@/components/connections/SuggestionsCarousel";

interface DashboardSummary {
  unread_notifications: number;
  job_seeker: {
    applications: number;
    saved_jobs: number;
    outgoing_referrals: number;
  };
  recruiter: {
    posted_jobs: number;
    total_applicants: number;
    incoming_referrals: number;
  };
  mentor: {
    upcoming_sessions: number;
    pending_requests: number;
    completed_sessions: number;
  };
}

function DashboardContent() {
  const router = useRouter();
  const { user } = useAuth();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [incomingRequests, setIncomingRequests] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Quick Coach Chat state
  const [coachMsg, setCoachMsg] = useState("");
  const [coachReplies, setCoachReplies] = useState<
    Array<{ sender: "user" | "coach"; text: string }>
  >([
    {
      sender: "coach",
      text: "Hello! I noticed you have some active applications. How can I help you prepare today?",
    },
  ]);
  const [sendingCoach, setSendingCoach] = useState(false);

  const loadData = async () => {
    try {
      const p = await profileClient.getMe();
      setProfile(p);

      const s = await api.get<DashboardSummary>("/me/dashboard");
      setSummary(s);

      const conns = await connectionsClient.getConnections(1, 100);
      setConnections(conns || []);

      const reqs = await connectionsClient.getPendingRequests("incoming");
      setIncomingRequests(reqs || []);
    } catch (err) {
      setError("Failed to load dashboard data.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAccept = async (id: string) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      await connectionsClient.acceptConnection(id);
      await loadData();
    } catch (e) {
      console.error(e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async (id: string) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      await connectionsClient.declineConnection(id);
      await loadData();
    } catch (e) {
      console.error(e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleSendCoach = () => {
    if (!coachMsg.trim() || sendingCoach) return;
    const userText = coachMsg;
    setCoachReplies((prev) => [...prev, { sender: "user", text: userText }]);
    setCoachMsg("");
    setSendingCoach(true);

    setTimeout(() => {
      setSendingCoach(false);
      setCoachReplies((prev) => [
        ...prev,
        {
          sender: "coach",
          text: `Based on your profile, I recommend updating your work experiences with metrics. Would you like me to rewrite your headline?`,
        },
      ]);
    }, 1500);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]}
        />
        <div className="flex-grow flex flex-col items-center justify-center py-20 gap-3">
          <Loader2 className="h-8 w-8 text-primary animate-spin" />
          <span className="text-sm font-semibold text-muted-foreground">
            Analyzing career metrics...
          </span>
        </div>
        <SiteFooter />
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]}
        />
        <main className="flex-grow max-w-lg mx-auto w-full px-4 py-20">
          <div className="p-6 bg-destructive/10 border border-destructive/20 rounded-3xl text-destructive text-center space-y-3">
            <p className="text-sm font-bold">
              {error || "Could not load dashboard."}
            </p>
            <button
              onClick={loadData}
              className="px-4 py-2 bg-destructive text-destructive-foreground rounded-full text-xs font-bold"
            >
              Try Again
            </button>
          </div>
        </main>
        <SiteFooter />
      </div>
    );
  }

  const firstName = user?.full_name?.split(" ")[0] || "Professional";

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]}
      />

      <main className="flex-grow max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        {/* Welcome Header */}
        <div className="space-y-1">
          <span className="text-xs font-bold uppercase tracking-widest text-primary">
            Workspace Overview
          </span>
          <h1 className="text-3xl font-extrabold tracking-tight">
            Good morning, {firstName}.
          </h1>
          <p className="text-sm text-muted-foreground">
            Here&apos;s a quick snapshot of your professional recovery network
            activity.
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {[
            {
              label: "Active Applications",
              value: summary?.job_seeker?.applications ?? 0,
              icon: Briefcase,
              color: "text-blue-500",
            },
            {
              label: "Saved Openings",
              value: summary?.job_seeker?.saved_jobs ?? 0,
              icon: Bookmark,
              color: "text-indigo-500",
            },
            {
              label: "Connections",
              value: connections.length,
              icon: Users,
              color: "text-emerald-500",
            },
            {
              label: "Profile Completeness",
              value: `${profile.profile_completeness_score}%`,
              icon: CheckCircle,
              color: "text-amber-500",
            },
          ].map((stat, idx) => {
            const Icon = stat.icon;
            return (
              <div
                key={idx}
                className="bg-card border border-border/60 p-5 rounded-2xl space-y-3 shadow-sm hover:border-border transition-all"
              >
                <div className="flex items-center justify-between">
                  <span className="text-xs font-semibold text-muted-foreground">
                    {stat.label}
                  </span>
                  <Icon className={`h-4.5 w-4.5 ${stat.color}`} />
                </div>
                <p className="text-2xl font-black tracking-tight">
                  {stat.value}
                </p>
              </div>
            );
          })}
        </div>

        {/* Quick Access Navigation Cards */}
        <div className="space-y-3">
          <h3 className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Quick Actions</h3>
          <div className="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-8 gap-3">
            {[
              { label: "My Applications", href: "/applications", icon: CheckCircle, color: "text-blue-500 bg-blue-500/10" },
              { label: "Saved Jobs", href: "/jobs/saved", icon: Bookmark, color: "text-indigo-500 bg-indigo-500/10" },
              { label: "Resume Score", href: "/resume", icon: FileText, color: "text-emerald-500 bg-emerald-500/10" },
              { label: "Recommended Jobs", href: "/jobs", icon: Briefcase, color: "text-amber-500 bg-amber-500/10" },
              { label: "Communities", href: "/communities", icon: Users, color: "text-violet-500 bg-violet-500/10" },
              { label: "Messages", href: "/messages", icon: MessageSquare, color: "text-pink-500 bg-pink-500/10" },
              { label: "Upcoming Interviews", href: "/jobs", icon: Clock, color: "text-cyan-500 bg-cyan-500/10" },
              { label: "Mentors", href: "/mentorship", icon: GraduationCap, color: "text-rose-500 bg-rose-500/10" },
            ].map((card, idx) => {
              const Icon = card.icon;
              return (
                <Link
                  key={idx}
                  href={card.href}
                  className="flex flex-col items-center justify-center text-center p-4 bg-card border border-border/60 hover:border-primary/40 rounded-2xl shadow-sm hover:shadow transition-all group"
                >
                  <div className={`p-2.5 rounded-xl ${card.color} mb-2.5 group-hover:scale-110 transition-transform`}>
                    <Icon className="h-5 w-5" />
                  </div>
                  <span className="text-[11px] font-bold text-foreground line-clamp-2 leading-snug">{card.label}</span>
                </Link>
              );
            })}
          </div>
        </div>

        {/* Main Content Areas */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
          {/* Left Column (Span 2): Active Modules */}
          <div className="lg:col-span-2 space-y-6">
            {/* Profile snapshot card */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
              <div className="flex flex-col sm:flex-row items-center sm:items-start gap-4">
                <div className="h-16 w-16 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-xl uppercase select-none">
                  {user?.full_name?.charAt(0) || "P"}
                </div>

                <div className="space-y-1 text-center sm:text-left flex-grow">
                  <div className="flex flex-col sm:flex-row items-center gap-2">
                    <h2 className="text-lg font-bold">{user?.full_name}</h2>
                    {profile.career_status && (
                      <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-extrabold text-emerald-500 uppercase tracking-wider">
                        {profile.career_status.replace("_", " ")}
                      </span>
                    )}
                  </div>
                  <p className="text-sm font-semibold text-muted-foreground">
                    {profile.headline || "Add a career headline"}
                  </p>
                  <div className="flex items-center justify-center sm:justify-start gap-1 text-xs text-muted-foreground">
                    <MapPin className="h-3.5 w-3.5" />
                    <span>
                      {profile.location || "Add location"} &bull;{" "}
                      {profile.open_to_remote ? "Open to Remote" : "On-site"}
                    </span>
                  </div>
                </div>

                <div className="flex gap-2 shrink-0">
                  <Link
                    href="/profile/edit"
                    className="px-4 py-2 rounded-full border border-border hover:bg-secondary text-xs font-bold flex items-center gap-1"
                  >
                    <Edit className="h-3.5 w-3.5" />
                    Edit Profile
                  </Link>
                  <Link
                    href="/profile"
                    className="px-4 py-2 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold flex items-center gap-1 shadow-sm"
                  >
                    <Eye className="h-3.5 w-3.5" />
                    Public View
                  </Link>
                </div>
              </div>
            </div>

            {/* Recommendations Widget */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm">
              <SuggestionsCarousel />
            </div>

            {/* Pending Requests Widget */}
            {incomingRequests.length > 0 && (
              <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
                <h3 className="text-base font-bold flex items-center gap-2">
                  <UserPlus className="h-4.5 w-4.5 text-primary" />
                  Connection Requests ({incomingRequests.length})
                </h3>

                <div className="space-y-3">
                  {incomingRequests.map((req) => (
                    <div
                      key={req.id}
                      className="p-4 bg-secondary/15 border border-border/40 rounded-2xl flex items-center justify-between gap-4"
                    >
                      <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-xs select-none">
                          {req.user.name?.charAt(0) || "M"}
                        </div>
                        <div>
                          <Link
                            href={`/profile/${req.user.id}`}
                            className="text-sm font-bold hover:underline"
                          >
                            {req.user.name}
                          </Link>
                          <p className="text-xs text-muted-foreground line-clamp-1">
                            {req.user.headline || "Professional"}
                          </p>
                        </div>
                      </div>

                      <div className="flex items-center gap-1.5">
                        <button
                          disabled={actionLoading}
                          onClick={() => handleAccept(req.id)}
                          className="p-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-500 hover:bg-emerald-500/20 transition-all"
                          title="Accept request"
                        >
                          <UserCheck className="h-4.5 w-4.5" />
                        </button>
                        <button
                          disabled={actionLoading}
                          onClick={() => handleReject(req.id)}
                          className="p-1.5 rounded-full bg-destructive/10 border border-destructive/20 text-destructive hover:bg-destructive/20 transition-all"
                          title="Ignore request"
                        >
                          <X className="h-4.5 w-4.5" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Connections Widget */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
              <h3 className="text-base font-bold flex items-center gap-2">
                <Users className="h-4.5 w-4.5 text-primary" />
                Network Connections ({connections.length})
              </h3>

              {connections.length === 0 ? (
                <div className="text-center py-8 border border-dashed border-border/60 rounded-2xl p-6 bg-secondary/15 space-y-3">
                  <p className="text-xs text-muted-foreground">
                    Get connected with peers, sponsors, and mentors to unlock
                    messaging.
                  </p>
                  <Link
                    href="/search?type=user"
                    className="inline-flex px-4 py-2 rounded-full border border-border hover:bg-secondary text-xs font-bold"
                  >
                    Find People
                  </Link>
                </div>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {connections.map((c) => {
                    const connName = c.user.name;
                    const connHeadline = c.user.headline;
                    const connID = c.user.id;
                    return (
                      <div
                        key={c.id}
                        className="p-4 bg-secondary/15 border border-border/40 rounded-2xl flex items-center justify-between gap-3"
                      >
                        <div className="flex items-center gap-2.5 min-w-0">
                          <div className="h-8 w-8 rounded-full bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-500 font-extrabold text-xs select-none">
                            {connName?.charAt(0) || "M"}
                          </div>
                          <div className="min-w-0">
                            <Link
                              href={`/profile/${connID}`}
                              className="text-xs font-bold hover:underline block truncate"
                            >
                              {connName}
                            </Link>
                            <span className="text-[10px] text-muted-foreground block truncate">
                              {connHeadline || "Professional"}
                            </span>
                          </div>
                        </div>
                        <Link
                          href="/inbox"
                          className="p-1.5 rounded-full hover:bg-secondary text-muted-foreground hover:text-foreground transition-all"
                        >
                          <MessageSquare className="h-4.5 w-4.5" />
                        </Link>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          </div>

          {/* Right Column (Span 1): Sidebar widgets */}
          <div className="space-y-6">
            {/* Quick AI Coach chatbot widget */}
            <div className="bg-card border border-border/80 rounded-3xl shadow-sm p-6 space-y-4">
              <div className="flex items-center justify-between pb-3 border-b border-border/40">
                <div className="flex items-center gap-2">
                  <Sparkles className="h-4.5 w-4.5 text-primary" />
                  <h3 className="text-sm font-bold text-foreground">
                    AI Career Coach
                  </h3>
                </div>
                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
              </div>

              {/* Chat replies */}
              <div className="space-y-3 max-h-[220px] overflow-y-auto pr-1">
                {coachReplies.map((reply, idx) => (
                  <div
                    key={idx}
                    className={`p-3 rounded-2xl text-xs leading-relaxed max-w-[90%] ${
                      reply.sender === "coach"
                        ? "bg-secondary/40 text-muted-foreground mr-auto rounded-tl-none"
                        : "bg-primary text-primary-foreground ml-auto rounded-tr-none"
                    }`}
                  >
                    {reply.text}
                  </div>
                ))}
                {sendingCoach && (
                  <div className="p-3 bg-secondary/40 rounded-2xl rounded-tl-none text-xs leading-relaxed mr-auto max-w-[90%] flex items-center gap-1 text-muted-foreground">
                    <Loader2 className="h-3 w-3 animate-spin text-primary" />
                    Analyzing profile updates...
                  </div>
                )}
              </div>

              {/* Input box */}
              <div className="relative pt-2 border-t border-border/40">
                <input
                  type="text"
                  placeholder="Ask Coach a question..."
                  value={coachMsg}
                  onChange={(e) => setCoachMsg(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSendCoach()}
                  className="w-full pl-3 pr-10 py-2 rounded-full border border-border/60 bg-secondary/15 placeholder:text-muted-foreground text-xs focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                />
                <button
                  onClick={handleSendCoach}
                  className="absolute right-2.5 top-[18px] p-1 rounded-full text-muted-foreground hover:text-foreground"
                >
                  <Send className="h-3.5 w-3.5" />
                </button>
              </div>
            </div>

            {/* Recommended Positions list */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
              <h3 className="text-sm font-bold flex items-center gap-1.5">
                <TrendingUp className="h-4.5 w-4.5 text-primary" />
                Recommended Roles
              </h3>

              <div className="space-y-4">
                {[
                  {
                    title: "VP, Supply Chain Operations",
                    company: "Atlas Co",
                    location: "Bangalore / Remote",
                  },
                  {
                    title: "Director of Logistics & Operations",
                    company: "Vertex Corp",
                    location: "Mumbai / Hybrid",
                  },
                ].map((rec, idx) => (
                  <div key={idx} className="space-y-1 block group">
                    <Link
                      href="/jobs"
                      className="text-xs font-bold group-hover:text-primary transition-colors block"
                    >
                      {rec.title}
                    </Link>
                    <div className="flex justify-between items-center text-[10px] text-muted-foreground">
                      <span>{rec.company}</span>
                      <span>{rec.location}</span>
                    </div>
                  </div>
                ))}
              </div>

              <Link
                href="/jobs"
                className="text-xs font-bold text-primary flex items-center gap-1 hover:underline pt-2 border-t border-border/40"
              >
                Explore all jobs
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}

export default function DashboardPage() {
  return (
    <AuthGuard>
      <DashboardContent />
    </AuthGuard>
  );
}
