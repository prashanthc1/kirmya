"use client";

import React, { useRef } from "react";
import {
  User,
  CheckCircle,
  Award,
  TrendingUp,
  ShieldCheck,
  FileCheck,
  Download,
  ExternalLink,
  HelpCircle,
  RotateCcw,
  RotateCw,
  RefreshCw,
  UploadCloud,
  AlertCircle,
} from "lucide-react";
import { ExtendedProfile } from "./types";

interface ProfileLeftSidebarProps {
  profile: ExtendedProfile;
  canUndo: boolean;
  canRedo: boolean;
  onUndo: () => void;
  onRedo: () => void;
  onReset: () => void;
  onOpenOnboarding: () => void;
  onScrollToSection: (sectionId: string) => void;
}

export default function ProfileLeftSidebar({
  profile,
  canUndo,
  canRedo,
  onUndo,
  onRedo,
  onReset,
  onOpenOnboarding,
  onScrollToSection,
}: ProfileLeftSidebarProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Helper to calculate scores based on profile fields (simulated)
  const completeness = profile.profile_completeness_score || 70;
  const atsScore = profile.analytics?.ats_score || 85;
  const recruiterScore = 90; // Recruiter readiness based on contact & status
  const verificationScore = profile.trust_score || 80;
  const resumeScore = profile.resumes?.[0]?.ats_score || 82;

  const handleAvatarClick = () => {
    fileInputRef.current?.click();
  };

  // Simulated profile photo change
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Normally uploads to server
    alert("Profile photo upload simulation triggered.");
  };

  // Calculate circular progress dash offset
  const radius = 38;
  const circumference = 2 * Math.PI * radius;
  const strokeDashoffset = circumference - (completeness / 100) * circumference;

  return (
    <aside className="w-full lg:w-80 flex-shrink-0 lg:sticky lg:top-24 space-y-6 self-start">
      {/* Sticky Top Profile Panel */}
      <div className="bg-card border border-border/80 rounded-3xl p-6 shadow-sm overflow-hidden relative">
        {/* Cover Banner Mockup */}
        <div className="absolute top-0 left-0 right-0 h-20 overflow-hidden bg-slate-900">
          {profile.cover_banner ? (
            <img
              src={profile.cover_banner}
              alt="Cover Banner"
              className="w-full h-full object-cover opacity-70"
            />
          ) : (
            <div className="w-full h-full bg-gradient-to-r from-blue-600 to-indigo-600 opacity-80" />
          )}
        </div>

        {/* User Card Avatar and Basics */}
        <div className="relative pt-10 flex flex-col items-center text-center">
          {/* Avatar frame with radial ring */}
          <div
            className="relative group cursor-pointer"
            onClick={handleAvatarClick}
          >
            <input
              type="file"
              ref={fileInputRef}
              className="hidden"
              onChange={handleFileChange}
              accept="image/*"
            />
            {/* Circular completeness tracker */}
            <svg className="w-24 h-24 transform -rotate-90">
              <circle
                cx="48"
                cy="48"
                r={radius}
                className="stroke-secondary fill-transparent"
                strokeWidth="4"
              />
              <circle
                cx="48"
                cy="48"
                r={radius}
                className="stroke-primary fill-transparent transition-all duration-500 ease-out"
                strokeWidth="4"
                strokeDasharray={circumference}
                strokeDashoffset={strokeDashoffset}
                strokeLinecap="round"
              />
            </svg>

            {/* Photo Container */}
            <div className="absolute inset-2.5 rounded-full overflow-hidden bg-secondary border border-border flex items-center justify-center">
              {profile.photo_url ? (
                <img
                  src={profile.photo_url}
                  alt={profile.headline}
                  className="w-full h-full object-cover"
                />
              ) : (
                <span className="text-xl font-bold text-primary">
                  {profile.preferred_name?.charAt(0) ||
                    profile.headline?.charAt(0) ||
                    "K"}
                </span>
              )}
            </div>

            {/* Upload Overlay */}
            <div className="absolute inset-2.5 rounded-full bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center text-white text-[9px] font-bold">
              <UploadCloud className="h-4.5 w-4.5 mb-0.5" />
              <span>Change Photo</span>
            </div>
          </div>

          <div className="mt-3 space-y-1.5 w-full">
            <h3 className="font-extrabold text-lg text-foreground tracking-tight flex items-center gap-1 justify-center">
              {profile.preferred_name || "Professional Identity"}
              {profile.trust_score && profile.trust_score > 90 && (
                <CheckCircle className="h-4 w-4 text-emerald-500 fill-emerald-500/10 shrink-0" />
              )}
            </h3>
            <p className="text-xs text-muted-foreground font-semibold px-4 line-clamp-2">
              {profile.headline || "Ready to excel"}
            </p>
            {profile.location && (
              <p className="text-[11px] text-muted-foreground/80 font-medium">
                📍 {profile.location}
              </p>
            )}
          </div>
        </div>

        {/* Completion details widget */}
        <div className="mt-6 border-t border-border/40 pt-4 flex justify-between items-center text-xs">
          <span className="text-muted-foreground font-medium">
            Profile Completeness
          </span>
          <span className="font-bold text-primary bg-primary/5 px-2 py-0.5 rounded-full">
            {completeness}%
          </span>
        </div>
      </div>

      {/* Profile Scores Dashboard */}
      <div className="bg-card border border-border/80 rounded-3xl p-5 shadow-sm space-y-4">
        <h4 className="text-xs font-bold text-muted-foreground uppercase tracking-widest flex items-center gap-1.5">
          <TrendingUp className="h-4 w-4 text-primary" />
          Kirmya Readiness Scores
        </h4>

        <div className="grid grid-cols-2 gap-2.5">
          {[
            {
              label: "Career Readiness",
              val: atsScore + 2,
              icon: Award,
              color: "text-blue-500 bg-blue-500/5 border-blue-500/20",
            },
            {
              label: "ATS Score",
              val: atsScore,
              icon: FileCheck,
              color: "text-emerald-500 bg-emerald-500/5 border-emerald-500/20",
            },
            {
              label: "Recruiter Reach",
              val: recruiterScore,
              icon: TrendingUp,
              color: "text-violet-500 bg-violet-500/5 border-violet-500/20",
            },
            {
              label: "Verification Score",
              val: verificationScore,
              icon: ShieldCheck,
              color: "text-indigo-500 bg-indigo-500/5 border-indigo-500/20",
            },
          ].map((score) => {
            const Icon = score.icon;
            return (
              <div
                key={score.label}
                className={`border p-3 rounded-2xl flex flex-col justify-between gap-2.5 ${score.color}`}
              >
                <div className="flex justify-between items-start">
                  <Icon className="h-4 w-4 opacity-80" />
                  <span className="font-black text-base leading-none">
                    {score.val}%
                  </span>
                </div>
                <span className="text-[10px] font-bold leading-tight uppercase tracking-wider text-muted-foreground/80">
                  {score.label}
                </span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-card border border-border/80 rounded-3xl p-5 shadow-sm space-y-3">
        <h4 className="text-xs font-bold text-muted-foreground uppercase tracking-widest">
          Quick Actions
        </h4>

        <div className="space-y-2">
          {/* Onboarding Trigger */}
          <button
            onClick={onOpenOnboarding}
            className="w-full text-left px-3.5 py-2.5 bg-primary/5 hover:bg-primary/10 border border-primary/15 text-primary text-xs font-bold rounded-2xl flex items-center justify-between transition-all cursor-pointer"
          >
            <span>Restart Onboarding Setup</span>
            <RefreshCw className="h-3.5 w-3.5" />
          </button>

          {/* Export PDF */}
          <button
            onClick={() => window.print()}
            className="w-full text-left px-3.5 py-2.5 hover:bg-secondary border border-border text-foreground text-xs font-bold rounded-2xl flex items-center justify-between transition-all cursor-pointer"
          >
            <span>Export Resume PDF</span>
            <Download className="h-3.5 w-3.5 text-muted-foreground" />
          </button>

          {/* History Stack Controls */}
          <div className="grid grid-cols-2 gap-2 pt-1">
            <button
              onClick={onUndo}
              disabled={!canUndo}
              className={`px-3 py-2 border rounded-xl text-[11px] font-bold flex items-center justify-center gap-1.5 transition-all ${
                canUndo
                  ? "border-border hover:bg-secondary text-foreground cursor-pointer"
                  : "border-border/40 text-muted-foreground/40 cursor-not-allowed"
              }`}
              title="Undo edit (Ctrl+Z)"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              Undo
            </button>
            <button
              onClick={onRedo}
              disabled={!canRedo}
              className={`px-3 py-2 border rounded-xl text-[11px] font-bold flex items-center justify-center gap-1.5 transition-all ${
                canRedo
                  ? "border-border hover:bg-secondary text-foreground cursor-pointer"
                  : "border-border/40 text-muted-foreground/40 cursor-not-allowed"
              }`}
              title="Redo edit (Ctrl+Y)"
            >
              <RotateCw className="h-3.5 w-3.5" />
              Redo
            </button>
          </div>

          <button
            onClick={onReset}
            className="w-full text-center py-2 text-[10px] text-destructive font-bold hover:underline cursor-pointer transition-all"
          >
            Reset Profile to Defaults
          </button>
        </div>
      </div>

      {/* AI Completeness Checklist Suggestions */}
      <div className="bg-card border border-border/80 rounded-3xl p-5 shadow-sm space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-xs font-bold text-muted-foreground uppercase tracking-widest">
            AI Profile Booster
          </h4>
          <span className="text-[10px] bg-amber-500/10 text-amber-600 font-bold px-1.5 py-0.5 rounded-md flex items-center gap-0.5">
            <AlertCircle className="h-3 w-3" />
            Urgent
          </span>
        </div>

        <div className="space-y-2.5">
          {[
            {
              label: "Connect calendar link for booking",
              score: "+5%",
              section: "identity",
            },
            {
              label: "Add Elevator Pitch to branding",
              score: "+8%",
              section: "summary",
            },
            {
              label: "Add KPIs to Northstar work experience",
              score: "+10%",
              section: "experience",
            },
            {
              label: "Complete Skills verification test",
              score: "+12%",
              section: "verification",
            },
          ].map((item, idx) => (
            <button
              key={idx}
              onClick={() => onScrollToSection(item.section)}
              className="w-full text-left p-2.5 bg-secondary/40 hover:bg-secondary/80 border border-border/50 hover:border-primary/20 rounded-xl transition-all flex justify-between items-start gap-2 text-xs group cursor-pointer"
            >
              <div className="space-y-0.5">
                <p className="font-semibold text-foreground leading-snug group-hover:text-primary transition-colors">
                  {item.label}
                </p>
                <p className="text-[10px] text-muted-foreground">
                  Scrolls to section
                </p>
              </div>
              <span className="text-[10px] font-bold text-emerald-500 bg-emerald-500/5 px-1.5 py-0.5 rounded shrink-0">
                {item.score}
              </span>
            </button>
          ))}
        </div>
      </div>
    </aside>
  );
}
