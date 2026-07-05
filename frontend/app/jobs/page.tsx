"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  salary: string;
  job_type: string;
  description: string;
}

export default function JobsPage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await api.get<{ jobs: Job[] }>("/jobs");
        if (active) setJobs(data?.jobs ?? []);
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Could not load jobs. Please try again.",
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  return (
    <div style={pageStyle}>
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Jobs" }]} />

      <main
        style={{
          flex: 1,
          width: "100%",
          maxWidth: "1000px",
          margin: "0 auto",
          padding: "clamp(32px,4vw,52px) 24px 72px",
        }}
      >
        <div style={eyebrowStyle}>
          {loading ? "Loading roles…" : `${jobs.length} open roles`}
        </div>
        <h1 style={headingStyle}>Roles that want your experience</h1>
        <p style={leadStyle}>
          Hand-vetted openings from recruiters who hire for proven track records.
        </p>

        {error && (
          <div role="alert" style={alertStyle}>
            {error}
          </div>
        )}

        {!loading && !error && jobs.length === 0 && (
          <p style={{ color: "#5B554C" }}>No open roles right now. Check back soon.</p>
        )}

        <div style={{ display: "flex", flexDirection: "column", gap: "14px", marginTop: "28px" }}>
          {jobs.map((job) => (
            <article key={job.id} style={cardStyle}>
              <div style={{ display: "flex", justifyContent: "space-between", gap: "16px", flexWrap: "wrap" }}>
                <div style={{ minWidth: 0 }}>
                  <h2 style={jobTitleStyle}>{job.title}</h2>
                  <div style={companyStyle}>
                    {job.company}
                    {job.location ? ` · ${job.location}` : ""}
                  </div>
                </div>
                {job.job_type && <span style={pillStyle}>{job.job_type}</span>}
              </div>
              {job.description && <p style={descStyle}>{job.description}</p>}
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between",
                  gap: "12px",
                  flexWrap: "wrap",
                  marginTop: "6px",
                }}
              >
                {job.salary && (
                  <span style={{ fontSize: "14px", color: "#4F7C6A", fontWeight: 600 }}>
                    {job.salary}
                  </span>
                )}
                <Link href="/jobs/detail" style={applyLinkStyle}>
                  View role →
                </Link>
              </div>
            </article>
          ))}
        </div>
      </main>

      <SiteFooter />
    </div>
  );
}

const pageStyle: React.CSSProperties = {
  background: "#FBF7F2",
  fontFamily: "'Public Sans',sans-serif",
  color: "#2B2620",
  minHeight: "100vh",
  overflowX: "hidden",
  display: "flex",
  flexDirection: "column",
};

const eyebrowStyle: React.CSSProperties = {
  fontSize: "13px",
  fontWeight: 700,
  letterSpacing: "0.12em",
  textTransform: "uppercase",
  color: "#C2683C",
  marginBottom: "14px",
};

const headingStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 800,
  fontSize: "clamp(30px,5vw,48px)",
  lineHeight: 1.03,
  letterSpacing: "-0.025em",
  margin: "0 0 12px",
};

const leadStyle: React.CSSProperties = {
  fontSize: "clamp(16px,2vw,19px)",
  lineHeight: 1.6,
  color: "#5B554C",
  maxWidth: "560px",
  margin: 0,
};

const cardStyle: React.CSSProperties = {
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "18px",
  padding: "clamp(20px,3vw,26px)",
  display: "flex",
  flexDirection: "column",
  gap: "12px",
};

const jobTitleStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 700,
  fontSize: "20px",
  margin: "0 0 4px",
  letterSpacing: "-0.01em",
};

const companyStyle: React.CSSProperties = { fontSize: "15px", color: "#5B554C" };

const pillStyle: React.CSSProperties = {
  flex: "none",
  height: "fit-content",
  fontSize: "12px",
  fontWeight: 600,
  color: "#4F7C6A",
  background: "rgba(79,124,106,0.12)",
  padding: "6px 12px",
  borderRadius: "100px",
  textTransform: "capitalize",
};

const descStyle: React.CSSProperties = {
  fontSize: "14px",
  lineHeight: 1.55,
  color: "#6B6357",
  margin: 0,
};

const applyLinkStyle: React.CSSProperties = {
  fontSize: "14px",
  fontWeight: 600,
  color: "#C2683C",
};

const alertStyle: React.CSSProperties = {
  marginTop: "20px",
  background: "rgba(194,104,60,0.10)",
  border: "1px solid rgba(194,104,60,0.35)",
  color: "#9A4A24",
  borderRadius: "10px",
  padding: "12px 14px",
  fontSize: "14px",
};
