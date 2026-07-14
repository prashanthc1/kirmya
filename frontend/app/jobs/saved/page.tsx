"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api } from "@/lib/api/client";
import { Bookmark, Briefcase, ChevronRight, Loader2, MapPin } from "lucide-react";

interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  salary: string;
  job_type: string;
}

function SavedJobsContent() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<{ jobs: Job[] }>("/jobs/saved");
        setJobs(data?.jobs ?? []);
      } catch (err: any) {
        setError(err.message || "Failed to load saved jobs.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <SiteNav breadcrumb={[{ label: "Home", href: "/dashboard" }, { label: "Saved Jobs" }]} />

      <main className="flex-grow max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        <div className="space-y-6">
          <div>
            <h1 className="text-2xl font-black tracking-tight">Saved Jobs</h1>
            <p className="text-xs text-muted-foreground mt-1">
              Review and apply to jobs you have saved for later.
            </p>
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-20 gap-3">
              <Loader2 className="h-8 w-8 text-primary animate-spin" />
              <span className="text-xs font-semibold text-muted-foreground">Loading saved jobs...</span>
            </div>
          ) : error ? (
            <div className="p-6 bg-destructive/10 border border-destructive/20 text-destructive rounded-2xl text-xs font-semibold">
              {error}
            </div>
          ) : jobs.length === 0 ? (
            <div className="text-center py-16 border border-dashed border-border/60 rounded-3xl p-8 bg-secondary/15 space-y-4">
              <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center mx-auto text-primary">
                <Bookmark className="h-6 w-6" />
              </div>
              <div className="space-y-1">
                <h3 className="text-sm font-bold">No saved jobs</h3>
                <p className="text-xs text-muted-foreground max-w-sm mx-auto">
                  You haven&apos;t saved any job opportunities yet. Keep browsing the jobs board and bookmark roles you like.
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
              {jobs.map((job) => (
                <div
                  key={job.id}
                  className="bg-card border border-border/60 hover:border-border p-5 rounded-2xl shadow-sm transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-4"
                >
                  <div className="flex gap-4">
                    <div className="h-12 w-12 rounded-xl bg-secondary flex items-center justify-center text-muted-foreground shrink-0 border border-border/30">
                      <Briefcase className="h-5 w-5" />
                    </div>
                    <div className="space-y-1">
                      <h3 className="text-sm font-bold hover:text-primary transition-colors">
                        <Link href={`/jobs/${job.id}`}>{job.title}</Link>
                      </h3>
                      <p className="text-xs font-semibold text-muted-foreground">
                        {job.company} &bull; {job.location}
                      </p>
                      <div className="flex items-center gap-1.5 text-[10px] text-muted-foreground/85">
                        <span className="font-semibold text-foreground">{job.salary}</span>
                        <span>&bull;</span>
                        <span className="capitalize">{job.job_type.replace("_", " ")}</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between sm:justify-end gap-3 border-t sm:border-transparent pt-3 sm:pt-0">
                    <Link
                      href={`/jobs/${job.id}`}
                      className="px-4 py-1.5 rounded-full bg-primary hover:bg-primary/95 text-primary-foreground text-xs font-bold shadow-sm"
                    >
                      Apply Now
                    </Link>
                    <Link
                      href={`/jobs/${job.id}`}
                      className="p-1.5 rounded-full hover:bg-secondary text-muted-foreground hover:text-foreground transition-all"
                      title="View job details"
                    >
                      <ChevronRight className="h-4.5 w-4.5" />
                    </Link>
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

export default function SavedJobsPage() {
  return (
    <AuthGuard>
      <SavedJobsContent />
    </AuthGuard>
  );
}
