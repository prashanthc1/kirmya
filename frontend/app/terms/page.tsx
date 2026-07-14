"use client";

import React from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { FileText } from "lucide-react";

export default function TermsOfServicePage() {
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Terms of Service" }]} />

      <main className="flex-grow max-w-3xl mx-auto px-4 sm:px-6 py-12 w-full space-y-6">
        <div className="flex items-center gap-3 border-b border-border/40 pb-4">
          <div className="h-10 w-10 bg-primary/10 border border-primary/20 rounded-xl flex items-center justify-center text-primary">
            <FileText className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-xl font-black tracking-tight">Terms of Service</h1>
            <p className="text-[10px] text-muted-foreground">Last updated: July 13, 2026</p>
          </div>
        </div>

        <div className="space-y-4 text-xs leading-relaxed text-muted-foreground/90">
          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">1. Platform Scope & Free Access</h2>
            <p>
              Kirmya is a completely free career transition recovery platform. We do not charge subscriptions, payment plans, or processing fees. All features are accessible to all registered users without charge.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">2. Code of Conduct</h2>
            <p>
              Users must interact respectfully. Spamming, automated mass outreach, unconstructive communication, or malicious resume stuffing will result in permanent account suspension.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-sm font-bold text-foreground">3. Disclaimer of Warranties</h2>
            <p>
              Kirmya tools and advice (including AI coach insights and mentor sessions) are provided &quot;as is&quot; without guarantees of specific employment or interview outcomes. We serve as a transition support workspace.
            </p>
          </section>
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}
