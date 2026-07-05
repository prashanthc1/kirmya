"use client";

import { useEffect, useState } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface Referral {
  id: string;
  company: string;
  message: string;
  status: string;
  outcome: string;
  created_at: string;
}

const STATUS_COLORS: Record<string, string> = {
  pending: "#B0852E",
  accepted: "#4F7C6A",
  declined: "#9A4A24",
};

export default function ReferralsPage() {
  const [referrals, setReferrals] = useState<Referral[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await api.get<{ referrals: Referral[] }>("/referrals/outgoing");
        if (active) setReferrals(data?.referrals ?? []);
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Could not load your referrals. Please try again.",
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
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Referrals" }]} />

      <main
        style={{
          flex: 1,
          width: "100%",
          maxWidth: "860px",
          margin: "0 auto",
          padding: "clamp(32px,4vw,52px) 24px 72px",
        }}
      >
        <div style={eyebrowStyle}>Referrals</div>
        <h1 style={headingStyle}>Your requests</h1>
        <p style={leadStyle}>
          Track the referral requests you&apos;ve sent and where each one stands.
        </p>

        {error && (
          <div role="alert" style={alertStyle}>
            {error}
          </div>
        )}

        {loading && <p style={{ color: "#8A8175", marginTop: "24px" }}>Loading…</p>}

        {!loading && !error && referrals.length === 0 && (
          <p style={{ color: "#5B554C", marginTop: "24px" }}>
            You haven&apos;t requested any referrals yet.
          </p>
        )}

        <div style={{ display: "flex", flexDirection: "column", gap: "14px", marginTop: "28px" }}>
          {referrals.map((ref) => (
            <article key={ref.id} style={cardStyle}>
              <div style={{ display: "flex", justifyContent: "space-between", gap: "16px", flexWrap: "wrap" }}>
                <h2 style={companyStyle}>{ref.company || "Referral request"}</h2>
                <span
                  style={{
                    ...statusPillStyle,
                    color: STATUS_COLORS[ref.status] ?? "#5B554C",
                    background:
                      ref.status === "accepted"
                        ? "rgba(79,124,106,0.12)"
                        : "#F3ECE2",
                  }}
                >
                  {ref.outcome || ref.status}
                </span>
              </div>
              {ref.message && <p style={messageStyle}>{ref.message}</p>}
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
  color: "#4F7C6A",
  marginBottom: "14px",
};

const headingStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 800,
  fontSize: "clamp(30px,5vw,46px)",
  lineHeight: 1.03,
  letterSpacing: "-0.025em",
  margin: "0 0 12px",
};

const leadStyle: React.CSSProperties = {
  fontSize: "clamp(16px,2vw,18px)",
  lineHeight: 1.6,
  color: "#5B554C",
  maxWidth: "520px",
  margin: 0,
};

const cardStyle: React.CSSProperties = {
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "18px",
  padding: "clamp(18px,3vw,24px)",
  display: "flex",
  flexDirection: "column",
  gap: "10px",
};

const companyStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 700,
  fontSize: "19px",
  margin: 0,
  letterSpacing: "-0.01em",
};

const statusPillStyle: React.CSSProperties = {
  flex: "none",
  height: "fit-content",
  fontSize: "12px",
  fontWeight: 600,
  padding: "6px 12px",
  borderRadius: "100px",
  textTransform: "capitalize",
};

const messageStyle: React.CSSProperties = {
  fontSize: "14px",
  lineHeight: 1.55,
  color: "#6B6357",
  margin: 0,
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
