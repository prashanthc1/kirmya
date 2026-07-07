"use client";

import React from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import ProfileWorkspace from "@/components/profile/ProfileWorkspace";
import AuthGuard from "@/components/shared/AuthGuard";

export default function ProfilePage() {
  return (
    <AuthGuard>
      <div className="min-h-screen bg-background text-foreground flex flex-col relative overflow-x-hidden">
        {/* Premium Ambient Background Glows */}
        <div className="absolute top-0 left-1/4 w-[500px] h-[500px] bg-blue-500/5 rounded-full blur-[100px] pointer-events-none -z-10 dark:bg-blue-500/10" />
        <div className="absolute top-1/3 right-10 w-[400px] h-[400px] bg-purple-500/5 rounded-full blur-[80px] pointer-events-none -z-10 dark:bg-purple-500/10" />
        <div className="absolute bottom-10 left-10 w-[300px] h-[300px] bg-emerald-500/5 rounded-full blur-[80px] pointer-events-none -z-10 dark:bg-emerald-500/5" />

        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]}
        />

        <main className="flex-grow">
          <ProfileWorkspace />
        </main>

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
