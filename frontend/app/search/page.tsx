"use client";

import { Suspense, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface SearchHit {
  type: string;
  ref_id: string;
  title: string;
  subtitle: string;
  url: string;
  score: number;
}

interface SearchResponse {
  results: SearchHit[];
  engine: string;
}

function SearchResults() {
  const params = useSearchParams();
  const query = params.get("q") ?? "";

  const [results, setResults] = useState<SearchHit[]>([]);
  const [engine, setEngine] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    setLoading(true);
    (async () => {
      try {
        const data = await api.get<SearchResponse>(
          `/search?q=${encodeURIComponent(query)}`,
        );
        if (active) {
          setResults(data?.results ?? []);
          setEngine(data?.engine ?? "");
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Search is unavailable right now. Please try again.",
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [query]);

  return (
    <main
      style={{
        flex: 1,
        width: "100%",
        maxWidth: "900px",
        margin: "0 auto",
        padding: "clamp(32px,4vw,52px) 24px 72px",
      }}
    >
      <div style={eyebrowStyle}>Search</div>
      <h1 style={headingStyle}>
        {query ? <>Results for &ldquo;{query}&rdquo;</> : "Search Kirmya"}
      </h1>

      <p style={{ fontSize: "14px", color: "#8A8175", margin: "6px 0 0" }}>
        {loading
          ? "Searching…"
          : `${results.length} result${results.length === 1 ? "" : "s"}`}
        {engine ? ` · Served by ${engine}` : ""}
      </p>

      {error && (
        <div role="alert" style={alertStyle}>
          {error}
        </div>
      )}

      {!loading && !error && results.length === 0 && query && (
        <p style={{ color: "#5B554C", marginTop: "28px" }}>
          Nothing matched &ldquo;{query}&rdquo;. Try a different term.
        </p>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: "12px", marginTop: "28px" }}>
        {results.map((hit) => (
          <Link key={`${hit.type}-${hit.ref_id}`} href={hit.url || "#"} style={cardLinkStyle}>
            <div style={{ display: "flex", alignItems: "center", gap: "14px" }}>
              <span style={typeBadgeStyle}>{hit.type}</span>
              <div style={{ minWidth: 0 }}>
                <div style={titleStyle}>{hit.title}</div>
                {hit.subtitle && <div style={subtitleStyle}>{hit.subtitle}</div>}
              </div>
            </div>
          </Link>
        ))}
      </div>
    </main>
  );
}

export default function SearchPage() {
  return (
    <div style={pageStyle}>
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Search" }]} />
      <Suspense
        fallback={
          <main style={{ flex: 1, maxWidth: "900px", margin: "0 auto", padding: "52px 24px" }}>
            <p style={{ color: "#8A8175" }}>Searching…</p>
          </main>
        }
      >
        <SearchResults />
      </Suspense>
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
  marginBottom: "12px",
};

const headingStyle: React.CSSProperties = {
  fontFamily: "'Bricolage Grotesque',sans-serif",
  fontWeight: 800,
  fontSize: "clamp(28px,4vw,40px)",
  lineHeight: 1.05,
  letterSpacing: "-0.02em",
  margin: 0,
};

const cardLinkStyle: React.CSSProperties = {
  display: "block",
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "14px",
  padding: "16px 18px",
  textDecoration: "none",
  color: "inherit",
};

const typeBadgeStyle: React.CSSProperties = {
  flex: "none",
  fontSize: "11px",
  fontWeight: 700,
  letterSpacing: "0.06em",
  textTransform: "uppercase",
  color: "#5B554C",
  background: "#F3ECE2",
  padding: "5px 10px",
  borderRadius: "8px",
};

const titleStyle: React.CSSProperties = {
  fontWeight: 600,
  fontSize: "16px",
  color: "#2B2620",
};

const subtitleStyle: React.CSSProperties = {
  fontSize: "13px",
  color: "#8A8175",
  marginTop: "2px",
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
