"use client";

import React, { useState } from "react";
import Link from "next/link";
import { 
  User, 
  MapPin, 
  Clock, 
  CheckCircle, 
  FileText, 
  Building2, 
  ArrowLeft, 
  Calendar,
  Layers,
  Sparkles,
  TrendingUp,
  Award,
  Globe2,
  FileCheck
} from "lucide-react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

const OUTCOMES = [
  { metric: "−18%", label: "Logistics cost per unit held over 3 years" },
  { metric: "99.2%", label: "On-time delivery kept through peak season supply" },
  { metric: "14 → 6", label: "Distribution centers consolidated without loss" }
];

const EXPERIENCES = [
  {
    title: "Operations Director",
    company: "Cascade Freight",
    duration: "2014 – 2025 (11 years)",
    desc: "Ran a 120-person operations organization across 14 distribution centers. Led the network consolidation that reduced overall cost-per-unit by 18% while lifting delivery consistency to 99.2%."
  },
  {
    title: "Senior Operations Manager",
    company: "Northstar Distribution",
    duration: "2007 – 2014 (7 years)",
    desc: "Scaled regional fulfillment from 2 to 6 facilities through a high-growth period, establishing the standard S&OP pipeline still in use today."
  },
  {
    title: "Operations Manager",
    company: "Meridian Logistics",
    duration: "2003 – 2007 (4 years)",
    desc: "Began on the warehouse floor and worked up to managing a regional operations hub of 60 active fulfillment staff."
  }
];

const SKILLS = [
  "Network Strategy",
  "S&OP Planning",
  "P&L Ownership",
  "Carrier Negotiation",
  "Cost Reduction",
  "Crisis Operations",
  "Fulfillment Logistics"
];

export default function ProfilePage() {
  const [activeRole, setActiveRole] = useState("Job Seeker");

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]} />

      <main className="flex-grow max-w-5xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 space-y-6">
        
        {/* Core Identity Panel */}
        <div className="bg-card border border-border/80 p-6 md:p-8 rounded-3xl shadow-sm space-y-6">
          <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
            <div className="h-24 w-24 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-3xl select-none shrink-0">
              MH
            </div>

            <div className="flex-grow space-y-3 text-center md:text-left">
              <div className="flex flex-col md:flex-row items-center gap-2 justify-center md:justify-start">
                <h1 className="text-2xl md:text-3xl font-extrabold tracking-tight">Marcus Hale</h1>
                <span className="px-2.5 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-extrabold text-emerald-500 uppercase tracking-widest flex items-center gap-1">
                  <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                  Open to work
                </span>
              </div>

              <p className="text-base font-semibold text-muted-foreground">Operations Director &bull; Supply Chain &amp; Logistics</p>

              <div className="flex flex-wrap items-center justify-center md:justify-start gap-2 text-xs">
                <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40 flex items-center gap-1">
                  <MapPin className="h-3.5 w-3.5" />
                  Denver, CO &bull; Remote-ready
                </span>
                <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40 flex items-center gap-1">
                  <Clock className="h-3.5 w-3.5" />
                  22 Years Experience
                </span>
                <span className="px-3 py-1 bg-secondary text-muted-foreground rounded-full border border-border/40 flex items-center gap-1">
                  <CheckCircle className="h-3.5 w-3.5 text-emerald-500" />
                  References Verified
                </span>
              </div>
            </div>
          </div>

          {/* Role selection tab */}
          <div className="flex items-center gap-3 border-t border-border/40 pt-4 flex-wrap">
            <span className="text-xs font-bold text-muted-foreground uppercase tracking-widest">Active as</span>
            <div className="flex bg-secondary p-1 rounded-full border border-border/40 gap-1">
              {["Job Seeker", "Recruiter", "Mentor"].map((role) => (
                <button
                  key={role}
                  onClick={() => setActiveRole(role)}
                  className={`px-4 py-1.5 rounded-full text-xs font-semibold transition-all ${
                    activeRole === role
                      ? "bg-background text-foreground shadow-sm"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {role}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Dynamic Detail grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 items-start">
          {/* Main profile contents */}
          <div className="md:col-span-2 space-y-8">
            {/* About / Summary */}
            <div className="space-y-3">
              <h2 className="text-lg font-bold text-foreground">Summary</h2>
              <p className="text-sm leading-relaxed text-muted-foreground">
                Operations leader with 22 years steadying complex supply chains through growth, restructuring, and two downturns. I&apos;ve owned cost bases north of $400M, rebuilt distribution networks under pressure, and kept service levels high when budgets were not. I lead calmly, hire well, and make the unglamorous call when it&apos;s the right one.
              </p>
            </div>

            {/* Proven outcomes */}
            <div className="space-y-4">
              <h2 className="text-lg font-bold text-foreground">Proven Outcomes</h2>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                {OUTCOMES.map((item, idx) => (
                  <div key={idx} className="bg-card border border-border/80 p-5 rounded-2xl space-y-2">
                    <span className="text-2xl font-black tracking-tight text-primary">{item.metric}</span>
                    <p className="text-xs text-muted-foreground leading-normal">{item.label}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Experience Timeline */}
            <div className="space-y-4">
              <h2 className="text-lg font-bold text-foreground">Work Experience</h2>
              <div className="space-y-6">
                {EXPERIENCES.map((exp, idx) => (
                  <div key={idx} className="flex gap-4 group">
                    <div className="flex flex-col items-center shrink-0">
                      <div className="h-7 w-7 rounded-full bg-secondary border border-border flex items-center justify-center text-[10px] font-bold text-muted-foreground group-hover:border-primary group-hover:text-primary transition-all">
                        {idx + 1}
                      </div>
                      {idx < EXPERIENCES.length - 1 && (
                        <div className="w-[1.5px] bg-border/80 flex-grow my-1.5" />
                      )}
                    </div>

                    <div className="space-y-2 pb-2">
                      <div>
                        <h3 className="text-sm font-bold text-foreground">{exp.title}</h3>
                        <div className="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
                          <span className="font-semibold text-foreground">{exp.company}</span>
                          <span>&bull;</span>
                          <span>{exp.duration}</span>
                        </div>
                      </div>
                      <p className="text-xs text-muted-foreground leading-relaxed">{exp.desc}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Hiring snapshot */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-4">
              <h3 className="text-sm font-bold flex items-center gap-1.5">
                <Layers className="h-4.5 w-4.5 text-primary" />
                Hiring Snapshot
              </h3>
              
              <div className="space-y-3 text-xs border-b border-border/40 pb-4">
                {[
                  { label: "Target Level", val: "Director / VP" },
                  { label: "Availability", val: "Immediate", color: "text-emerald-500 font-bold" },
                  { label: "Work Preference", val: "Remote / Hybrid" },
                  { label: "Target Comp", val: "$180,000+" }
                ].map((row, idx) => (
                  <div key={idx} className="flex justify-between items-center">
                    <span className="text-muted-foreground">{row.label}</span>
                    <span className={`font-semibold ${row.color || "text-foreground"}`}>{row.val}</span>
                  </div>
                ))}
              </div>

              {/* Skills list */}
              <div className="space-y-2">
                <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest block">Skills &amp; Domain Expertise</span>
                <div className="flex flex-wrap gap-1.5">
                  {SKILLS.map((skill) => (
                    <span key={skill} className="px-2.5 py-1 bg-secondary text-muted-foreground border border-border/40 rounded-full text-[10px] font-semibold">
                      {skill}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            {/* References verified */}
            <div className="bg-card border border-border/80 p-6 rounded-3xl shadow-sm space-y-3">
              <div className="flex items-center gap-2">
                <FileCheck className="h-5 w-5 text-emerald-500" />
                <h3 className="text-sm font-bold">References Screened</h3>
              </div>
              <p className="text-xs text-muted-foreground leading-relaxed">
                Three professional references (including former VP and direct reports) have been verified by Kirmya staff. Full logs are available on request to recruiters.
              </p>
            </div>
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
