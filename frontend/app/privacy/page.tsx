"use client";

import React from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { Shield } from "lucide-react";

export default function PrivacyPolicyPage() {
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Privacy Policy" }]} />

      <main className="flex-grow max-w-3xl mx-auto px-4 sm:px-6 py-12 w-full space-y-6">
        <div className="flex items-center gap-3 border-b border-border/40 pb-4">
          <div className="h-10 w-10 bg-primary/10 border border-primary/20 rounded-xl flex items-center justify-center text-primary">
            <Shield className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-xl font-black tracking-tight">Privacy Policy</h1>
            <p className="text-[10px] text-muted-foreground">Last updated: July 13, 2026</p>
          </div>
        </div>

        <div className="space-y-4 text-xs leading-relaxed text-muted-foreground/90">
          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">1. Information We Collect</h2>
            <p>
              We collect information that you directly provide to us, including when you register an account, fill in your profile metadata, upload your resume file for optimization scanning, book a session with a network mentor, or participate in community channels. This may include your name, email, credentials, and career materials.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">2. How We Use Your Data</h2>
            <p>
              We use the collected information to power Kirmya.com transition support features, parse and suggest resume improvements, match you with verified employee referrers, schedule mentorship sessions, and connect you with relevant community circles.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">3. Cookie Preferences & Consent</h2>
            <p>
              You can review, modify, or withdraw your cookie preferences at any time under your account settings. Kirmya respects standard DNT (Do Not Track) signals and only deploys non-essential cookies (such as analytics and personalized recommendations) with your explicit consent.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">4. Contact Us</h2>
            <p>
              If you have any questions or feedback regarding this Privacy Policy, please reach out to us at privacy@kirmya.com or via our help center.
            </p>
          </section>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
