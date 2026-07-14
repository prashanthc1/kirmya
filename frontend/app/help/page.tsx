"use client";

import React, { useState } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { GraduationCap, Mail, MessageSquare, Shield, HelpCircle, ArrowRight } from "lucide-react";
import Link from "next/link";

export default function HelpCenterPage() {
  const [activeFaq, setActiveFaq] = useState<number | null>(null);

  const FAQS = [
    {
      q: "Is Kirmya completely free to use?",
      a: "Yes. During career transitions, the last thing you need is another bill. The core tools—resume parsing, AI coaching suggestions, community access, and mentorship sessions—are free. We focus on one metric: the speed at which you land your next interview.",
    },
    {
      q: "How does the referral system protect my privacy?",
      a: "We do not allow cold, automated messaging. Referral requests require structured introductions, target specific open roles, and are matched based on mutual interest, protecting referrers from spam while keeping candidate signal high.",
    },
    {
      q: "What makes Kirmya different from LinkedIn?",
      a: "LinkedIn optimizes for screen time, influencers, and vanity metrics. Kirmya is a recovery workspace. There is no public content feed. We provide tools to directly prepare you for interviews and connect with internal champions.",
    },
  ];

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Help Center" }]} />

      <main className="flex-grow max-w-4xl mx-auto px-4 sm:px-6 py-12 w-full space-y-12">
        {/* Hero */}
        <div className="text-center space-y-3">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20">
            <GraduationCap className="h-3.5 w-3.5" />
            Support & Resources
          </div>
          <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight">Help Center</h1>
          <p className="text-sm text-muted-foreground max-w-xl mx-auto">
            Find answers to common questions about Kirmya, get help with your account, or contact support.
          </p>
        </div>

        {/* Action blocks */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Link
            href="/contact"
            className="p-5 bg-card border border-border/60 hover:border-primary/40 rounded-2xl shadow-sm transition-all group flex flex-col gap-3"
          >
            <div className="h-9 w-9 rounded-xl bg-blue-500/10 flex items-center justify-center text-blue-500">
              <Mail className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-xs font-bold">Contact Support</h3>
              <p className="text-[10px] text-muted-foreground mt-1">Get in touch with our team for direct help.</p>
            </div>
            <span className="text-[10px] font-semibold text-primary mt-auto flex items-center gap-1">
              Send Message <ArrowRight className="h-3 w-3 group-hover:translate-x-0.5 transition-transform" />
            </span>
          </Link>

          <Link
            href="/privacy"
            className="p-5 bg-card border border-border/60 hover:border-primary/40 rounded-2xl shadow-sm transition-all group flex flex-col gap-3"
          >
            <div className="h-9 w-9 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-500">
              <Shield className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-xs font-bold">Privacy & Security</h3>
              <p className="text-[10px] text-muted-foreground mt-1">Learn how we protect and manage your data.</p>
            </div>
            <span className="text-[10px] font-semibold text-primary mt-auto flex items-center gap-1">
              Read Policy <ArrowRight className="h-3 w-3 group-hover:translate-x-0.5 transition-transform" />
            </span>
          </Link>

          <Link
            href="/settings"
            className="p-5 bg-card border border-border/60 hover:border-primary/40 rounded-2xl shadow-sm transition-all group flex flex-col gap-3"
          >
            <div className="h-9 w-9 rounded-xl bg-violet-500/10 flex items-center justify-center text-violet-500">
              <HelpCircle className="h-5 w-5" />
            </div>
            <div>
              <h3 className="text-xs font-bold">Account Preferences</h3>
              <p className="text-[10px] text-muted-foreground mt-1">Manage notifications, account deactivation, etc.</p>
            </div>
            <span className="text-[10px] font-semibold text-primary mt-auto flex items-center gap-1">
              Go to Settings <ArrowRight className="h-3 w-3 group-hover:translate-x-0.5 transition-transform" />
            </span>
          </Link>
        </div>

        {/* FAQs */}
        <div className="space-y-4">
          <h2 className="text-base font-bold text-center">Frequently Asked Questions</h2>
          <div className="bg-card border border-border/60 rounded-2xl divide-y divide-border/40 overflow-hidden shadow-sm">
            {FAQS.map((faq, idx) => (
              <div key={idx} className="p-4 space-y-2">
                <button
                  onClick={() => setActiveFaq(activeFaq === idx ? null : idx)}
                  className="w-full flex items-center justify-between text-left text-xs font-bold text-foreground focus:outline-none cursor-pointer"
                >
                  <span>{faq.q}</span>
                  <span className="text-muted-foreground/60">{activeFaq === idx ? "−" : "+"}</span>
                </button>
                {activeFaq === idx && (
                  <p className="text-[11px] leading-relaxed text-muted-foreground/90 pt-1">{faq.a}</p>
                )}
              </div>
            ))}
          </div>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
