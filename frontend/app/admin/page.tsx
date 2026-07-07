"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import AuthGuard from "@/components/shared/AuthGuard";

interface Analytics {
  users: { total: number; active: number; new_7d: number };
  jobs: { total: number; applications: number };
  referrals: { total: number; accepted: number; hired: number };
  communities: { total: number; posts: number };
  reports: { open: number };
}

interface AdminUser {
  id: string;
  email: string;
  full_name: string;
  status: string;
  roles: string[];
}

type Tab = "overview" | "users";

function AdminConsole() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();

  const [tab, setTab] = useState<Tab>("overview");
  const [stats, setStats] = useState<Analytics | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [error, setError] = useState<string | null>(null);

  const isAdmin = !!user?.roles?.includes("admin");

  // Guard: once the session has settled, bounce non-admins to the dashboard.
  useEffect(() => {
    if (authLoading) return;
    if (!isAdmin) router.replace("/dashboard");
  }, [authLoading, isAdmin, router]);

  // Load analytics + users only for confirmed admins.
  useEffect(() => {
    if (authLoading || !isAdmin) return;
    let active = true;
    (async () => {
      try {
        const [analytics, userList] = await Promise.all([
          api.get<Analytics>("/admin/stats"),
          api.get<{ users: AdminUser[] }>("/admin/users?limit=50"),
        ]);
        if (active) {
          setStats(analytics);
          setUsers(userList?.users ?? []);
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Could not load admin data. Please try again.",
          );
        }
      }
    })();
    return () => {
      active = false;
    };
  }, [authLoading, isAdmin]);

  // While the guard decides, or for non-admins mid-redirect, render nothing heavy.
  if (authLoading || !isAdmin) {
    return (
      <div style={pageStyle}>
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Admin" }]}
        />
        <main
          style={{
            flex: 1,
            maxWidth: "1100px",
            margin: "0 auto",
            padding: "52px 24px",
          }}
        >
          <p style={{ color: "#8A8175" }}>Checking access…</p>
        </main>
        <SiteFooter />
      </div>
    );
  }

  const statCards: { label: string; value: number | undefined }[] = [
    { label: "Total members", value: stats?.users.total },
    { label: "Total jobs", value: stats?.jobs.total },
    { label: "Total referrals", value: stats?.referrals.total },
    { label: "Total communities", value: stats?.communities.total },
  ];

  return (
    <div style={pageStyle}>
      <SiteNav
        breadcrumb={[{ label: "Home", href: "/" }, { label: "Admin" }]}
      />

      <main
        style={{
          flex: 1,
          width: "100%",
          maxWidth: "1100px",
          margin: "0 auto",
          padding: "clamp(28px,4vw,44px) 24px 72px",
        }}
      >
        <div style={eyebrowStyle}>Platform</div>
        <h1 style={headingStyle}>Admin Console</h1>

        <div
          role="tablist"
          style={{
            display: "flex",
            gap: "6px",
            background: "#F3ECE2",
            borderRadius: "100px",
            padding: "4px",
            width: "fit-content",
            margin: "22px 0 28px",
          }}
        >
          <button
            type="button"
            role="tab"
            aria-selected={tab === "overview"}
            onClick={() => setTab("overview")}
            style={tabStyle(tab === "overview")}
          >
            Overview
          </button>
          <button
            type="button"
            role="tab"
            aria-selected={tab === "users"}
            onClick={() => setTab("users")}
            style={tabStyle(tab === "users")}
          >
            Users
          </button>
        </div>

        {error && (
          <div role="alert" style={alertStyle}>
            {error}
          </div>
        )}

        {tab === "overview" && (
          <section
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(auto-fit,minmax(212px,1fr))",
              gap: "16px",
            }}
          >
            {statCards.map((card) => (
              <div key={card.label} style={statCardStyle}>
                <div
                  style={{
                    fontSize: "13.5px",
                    color: "#8A8175",
                    marginBottom: "12px",
                  }}
                >
                  {card.label}
                </div>
                <div style={statValueStyle}>{card.value ?? "—"}</div>
              </div>
            ))}
          </section>
        )}

        {tab === "users" && (
          <section
            style={{
              background: "#fff",
              border: "1px solid #EFE7DC",
              borderRadius: "18px",
              overflow: "hidden",
            }}
          >
            {users.length === 0 ? (
              <p style={{ padding: "20px", color: "#8A8175", margin: 0 }}>
                No users found.
              </p>
            ) : (
              <ul style={{ listStyle: "none", margin: 0, padding: 0 }}>
                {users.map((u) => (
                  <li key={u.id} style={userRowStyle}>
                    <div style={{ minWidth: 0 }}>
                      <div style={{ fontWeight: 600, fontSize: "15px" }}>
                        {u.full_name}
                      </div>
                      <div style={{ fontSize: "13px", color: "#8A8175" }}>
                        {u.email}
                      </div>
                    </div>
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        flexWrap: "wrap",
                      }}
                    >
                      {u.roles.map((r) => (
                        <span key={r} style={roleTagStyle}>
                          {r}
                        </span>
                      ))}
                      <span style={{ fontSize: "12px", color: "#5B554C" }}>
                        {u.status}
                      </span>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </section>
        )}
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
  marginBottom: "12px",
};

const headingStyle: React.CSSProperties = {
  fontFamily: "'Public Sans',sans-serif",
  fontWeight: 800,
  fontSize: "clamp(28px,4vw,42px)",
  lineHeight: 1.03,
  letterSpacing: "-0.02em",
  margin: 0,
};

function tabStyle(active: boolean): React.CSSProperties {
  return {
    border: "none",
    background: active ? "#fff" : "transparent",
    color: active ? "#2B2620" : "#8A8175",
    fontFamily: "'Public Sans',sans-serif",
    fontSize: "14px",
    fontWeight: 600,
    padding: "9px 20px",
    borderRadius: "100px",
    cursor: "pointer",
    boxShadow: active ? "0 1px 3px rgba(0,0,0,0.10)" : "none",
  };
}

const statCardStyle: React.CSSProperties = {
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "18px",
  padding: "22px",
};

const statValueStyle: React.CSSProperties = {
  fontFamily: "'Public Sans',sans-serif",
  fontWeight: 800,
  fontSize: "34px",
  letterSpacing: "-0.02em",
  lineHeight: 1,
  color: "#2B2620",
};

const userRowStyle: React.CSSProperties = {
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
  gap: "16px",
  flexWrap: "wrap",
  padding: "14px 18px",
  borderBottom: "1px solid #F3ECE2",
};

const roleTagStyle: React.CSSProperties = {
  fontSize: "11px",
  fontWeight: 700,
  letterSpacing: "0.04em",
  textTransform: "uppercase",
  color: "#5B554C",
  background: "#F3ECE2",
  padding: "4px 9px",
  borderRadius: "7px",
};

const alertStyle: React.CSSProperties = {
  marginBottom: "20px",
  background: "rgba(194,104,60,0.10)",
  border: "1px solid rgba(194,104,60,0.35)",
  color: "#9A4A24",
  borderRadius: "10px",
  padding: "12px 14px",
  fontSize: "14px",
};

export default function AdminPage() {
  return (
    <AuthGuard>
      <AdminConsole />
    </AuthGuard>
  );
}
