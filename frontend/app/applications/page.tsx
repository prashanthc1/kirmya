"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api } from "@/lib/api/client";
import { Briefcase, Calendar, CheckCircle2, ChevronRight, Loader2, Sparkles } from "lucide-react";

interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  salary: string;
  job_type: string;
}

interface Application {
  id: string;
  job_id: string;
  status: string;
  created_at: string;
  cover_letter: string;
}

function ApplicationsContent() {
  const [applications, setApplications] = useState<(Application & { job?: Job })[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const [jobsRes, appsRes] = await Promise.all([
          api.get<{ jobs: Job[] }>("/jobs"),
          api.get<{ applications: Application[] }>("/jobs/applications"),
        ]);

        const jobsMap = new Map<string, Job>();
        (jobsRes?.jobs ?? []).forEach((j) => jobsMap.set(j.id, j));

        const resolved = (appsRes?.applications ?? []).map((app) => ({
          ...app,
          job: jobsMap.get(app.job_id),
        }));

        setApplications(resolved);
      } catch (err: any) {
        setError(err.message || "Failed to load applications.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return "bg-amber-500/10 border-amber-500/20 text-amber-600 dark:text-amber-400";
      case "interviewing":
        return "bg-blue-500/10 border-blue-500/20 text-blue-600 dark:text-blue-400";
      case "offered":
        return "bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400";
      case "rejected":
        return "bg-rose-500/10 border-rose-500/20 text-rose-600 dark:text-rose-400";
      default:
        return "bg-secondary text-muted-foreground";
    }
  };

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/dashboard" }, { label: "My Applications" }]} />

      <main className="flex-grow max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="space-y-6">
          <div>
            <h1 className="text-2xl font-black tracking-tight">My Applications</h1>
            <p className="text-xs text-muted-foreground mt-1">
              Track the progress of your active and past job applications.
            </p>
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-20 gap-3">
              <Loader2 className="h-8 w-8 text-primary animate-spin" />
              <span className="text-xs font-semibold text-muted-foreground">Loading applications...</span>
            </div>
          ) : error ? (
            <div className="p-6 bg-destructive/10 border border-destructive/20 text-destructive rounded-2xl text-xs font-semibold">
              {error}
            </div>
          ) : applications.length === 0 ? (
            <div className="text-center py-16 border border-dashed border-border/60 rounded-3xl p-8 bg-secondary/15 space-y-4">
              <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto text-primary">
                <Briefcase className="h-6 w-6" />
              </div>
              <div className="space-y-1">
                <h3 className="text-sm font-bold">No applications found</h3>
                <p className="text-xs text-muted-foreground max-w-sm mx-auto">
                  You haven&apos;t applied to any job listings yet. Explore our jobs board and land your next role.
                </p>
              </div>
              <Link
                href="/jobs"
                className="inline-flex px-5 py-2 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold transition-all shadow-sm"
              >
                Browse Open Roles
              </Link>
            </div>
          ) : (
            <div className="space-y-3">
              {applications.map((app) => (
                <div
                  key={app.id}
                  className="bg-card border border-border/60 hover:border-border/100 p-5 rounded-2xl shadow-sm transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-4"
                >
                  <div className="flex gap-4">
                    <div className="h-12 w-12 rounded-xl bg-secondary flex items-center justify-center text-muted-foreground shrink-0 border border-border/30">
                      <Briefcase className="h-5 w-5" />
                    </div>
                    <div className="space-y-1">
                      <h3 className="text-sm font-bold hover:text-primary transition-colors">
                        {app.job ? (
                          <Link href={`/jobs/${app.job.id}`}>{app.job.title}</Link>
                        ) : (
                          "Position Unavailable"
                        )}
                      </h3>
                      <p className="text-xs font-semibold text-muted-foreground">
                        {app.job?.company || "Unknown Company"} &bull; {app.job?.location || "Remote"}
                      </p>
                      <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground/80">
                        <Calendar className="h-3.5 w-3.5" />
                        <span>Applied on {new Date(app.created_at).toLocaleDateString([], { month: "short", day: "numeric", year: "numeric" })}</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between sm:justify-end gap-3 border-t sm:border-transparent pt-3 sm:pt-0">
                    <span className={`px-2.5 py-0.5 rounded-full border text-[10px] font-extrabold uppercase tracking-wider ${getStatusColor(app.status)}`}>
                      {app.status}
                    </span>
                    {app.job && (
                      <Link
                        href={`/jobs/${app.job.id}`}
                        className="p-1.5 rounded-full hover:bg-secondary text-muted-foreground hover:text-foreground transition-all"
                        title="View job opening"
                      >
                        <ChevronRight className="h-4.5 w-4.5" />
                      </Link>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}

export default function ApplicationsPage() {
  return (
    <AuthGuard>
      <ApplicationsContent />
    </AuthGuard>
  );
}
