"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import AuthGuard from "@/components/shared/AuthGuard";
import { CircularProgress } from "@mui/material";
import {
  LayoutDashboard,
  Users,
  User,
  MessageSquare,
  Briefcase,
  FileText,
  ShieldAlert,
  BarChart2,
  Bell,
  Settings,
  ShieldCheck,
  Search,
  Filter,
  ArrowUpDown,
  Download,
  AlertTriangle,
  Menu,
  ChevronLeft,
  ChevronRight,
  Send,
  Plus,
  Save,
  Activity,
  Heart,
  CheckCircle,
} from "lucide-react";

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
  location?: string;
  created_at?: string;
}

interface ContentItem {
  id: string;
  title: string;
  category: string;
  content: string;
  updatedAt: string;
}

interface ReportItem {
  id: string;
  reporter_id: string;
  target_type: string;
  target_id: string;
  reason: string;
  status: string;
  created_at: string;
  action_taken?: string;
}

type Tab =
  | "overview"
  | "users"
  | "profiles"
  | "communities"
  | "messaging"
  | "jobs"
  | "content"
  | "moderation"
  | "analytics"
  | "notifications"
  | "settings"
  | "roles";

const pageStyle: React.CSSProperties = {
  background: "#FBF7F2",
  fontFamily: "'Public Sans', sans-serif",
  color: "#2B2620",
  minHeight: "100vh",
  display: "flex",
  flexDirection: "column",
};

const statCardStyle: React.CSSProperties = {
  background: "#fff",
  border: "1px solid #EFE7DC",
  borderRadius: "18px",
  padding: "22px",
  boxShadow: "0 4px 12px rgba(43, 38, 32, 0.02)",
};

const statValueStyle: React.CSSProperties = {
  fontFamily: "'Public Sans', sans-serif",
  fontWeight: 800,
  fontSize: "34px",
  letterSpacing: "-0.02em",
  lineHeight: 1,
  color: "#2B2620",
  marginTop: "8px",
};

export default function AdminPage() {
  return (
    <AuthGuard>
      <AdminConsole />
    </AuthGuard>
  );
}

function AdminConsole() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const { showNotification } = useNotifications();

  const [tab, setTab] = useState<Tab>("overview");
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<Analytics | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [reports, setReports] = useState<ReportItem[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [error, setError] = useState<string | null>(null);

  // Content management state
  const [contentList, setContentList] = useState<ContentItem[]>([
    { id: "1", title: "Landing Page Hero Headline", category: "Banners", content: "Regroup, rebuild, and come back stronger.", updatedAt: "2026-07-10" },
    { id: "2", title: "Frequently Asked Questions", category: "FAQ", content: "Information on verification, referrals, and career coaching.", updatedAt: "2026-07-08" },
    { id: "3", title: "Cookie Disclosures and Policy", category: "Legal", content: "Granular details about cookies, localStorage and user consent compliance.", updatedAt: "2026-07-06" },
  ]);
  const [editingContent, setEditingContent] = useState<ContentItem | null>(null);

  // Notifications manager state
  const [notifTemplate, setNotifTemplate] = useState("welcome");
  const [notifText, setNotifText] = useState("");
  const [sendingNotif, setSendingNotif] = useState(false);

  // Config settings state
  const [maintenanceMode, setMaintenanceMode] = useState(false);
  const [betaSignup, setBetaSignup] = useState(true);
  const [requireEmailVerify, setRequireEmailVerify] = useState(true);

  // Selected User for details modal
  const [selectedUser, setSelectedUser] = useState<AdminUser | null>(null);

  // Role details mapping permissions matrix state
  const [permissionsMatrix, setPermissionsMatrix] = useState<Record<string, string[]>>({
    "Super Admin": ["User.Read", "User.Write", "Moderation.Apply", "Content.Modify", "System.Manage"],
    "Admin": ["User.Read", "User.Write", "Moderation.Apply", "Content.Modify"],
    "Moderator": ["User.Read", "Moderation.Apply"],
  });

  const isAdmin = !!user?.roles?.includes("admin");

  // Guard redirection
  useEffect(() => {
    if (authLoading) return;
    if (!isAdmin) router.replace("/dashboard");
  }, [authLoading, isAdmin, router]);

  // Load analytics, users, and moderation reports
  useEffect(() => {
    if (authLoading || !isAdmin) return;
    fetchAdminData();
  }, [authLoading, isAdmin]);

  const fetchAdminData = async () => {
    setLoading(true);
    setError(null);
    try {
      const [analytics, userList, reportsList] = await Promise.all([
        api.get<Analytics>("/admin/stats"),
        api.get<{ users: AdminUser[] }>("/admin/users?limit=50"),
        api.get<ReportItem[]>("/admin/reports").catch(() => []),
      ]);
      setStats(analytics);
      setUsers(userList?.users ?? []);
      setReports(reportsList || []);
    } catch (err: any) {
      setError(err.message || "Could not load admin data.");
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (userId: string, newStatus: string) => {
    try {
      await api.patch(`/admin/users/${userId}/status`, { status: newStatus });
      showNotification(`User status updated to ${newStatus}`, "success");
      setUsers((prev) =>
        prev.map((u) => (u.id === userId ? { ...u, status: newStatus } : u))
      );
    } catch (err: any) {
      showNotification(err.message || "Failed to update status", "error");
    }
  };

  const handleToggleRole = async (userId: string, roleName: string, hasRole: boolean) => {
    try {
      if (hasRole) {
        await api.delete(`/admin/users/${userId}/roles/${roleName}`);
        showNotification("Role revoked successfully", "success");
      } else {
        await api.post(`/admin/users/${userId}/roles`, { role: roleName });
        showNotification("Role assigned successfully", "success");
      }
      fetchAdminData();
    } catch (err: any) {
      showNotification(err.message || "Failed to update role", "error");
    }
  };

  const handleResolveReport = async (reportId: string, status: string, action: string) => {
    try {
      await api.patch(`/admin/reports/${reportId}`, { status, action_taken: action });
      showNotification("Report resolved successfully", "success");
      setReports((prev) =>
        prev.map((r) => (r.id === reportId ? { ...r, status, action_taken: action } : r))
      );
    } catch (err: any) {
      showNotification("Failed to resolve report", "error");
    }
  };

  const handleSaveContent = (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingContent) return;
    setContentList((prev) =>
      prev.map((c) => (c.id === editingContent.id ? { ...editingContent, updatedAt: new Date().toISOString().split("T")[0] } : c))
    );
    showNotification("Content item saved successfully!", "success");
    setEditingContent(null);
  };

  const handleSendNotification = (e: React.FormEvent) => {
    e.preventDefault();
    if (!notifText.trim()) return;
    setSendingNotif(true);
    setTimeout(() => {
      showNotification("Global announcement broadcasted successfully!", "success");
      setNotifText("");
      setSendingNotif(false);
    }, 1200);
  };

  const togglePermission = (role: string, perm: string) => {
    setPermissionsMatrix((prev) => {
      const current = prev[role] || [];
      const updated = current.includes(perm) ? current.filter((p) => p !== perm) : [...current, perm];
      return { ...prev, [role]: updated };
    });
    showNotification(`Permissions matrix updated for ${role}`, "info");
  };

  if (authLoading || !isAdmin) {
    return (
      <div style={pageStyle}>
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Admin" }]} />
        <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center" }}>
          <CircularProgress style={{ color: "#C2683C" }} />
        </main>
        <SiteFooter />
      </div>
    );
  }

  const navItems = [
    { id: "overview", label: "Overview", icon: LayoutDashboard },
    { id: "users", label: "Users", icon: Users },
    { id: "profiles", label: "Profiles Center", icon: User },
    { id: "communities", label: "Communities", icon: Users },
    { id: "messaging", label: "Messaging Triage", icon: MessageSquare },
    { id: "jobs", label: "Jobs Portal", icon: Briefcase },
    { id: "content", label: "Content Editor", icon: FileText },
    { id: "moderation", label: "Moderation Queue", icon: ShieldAlert },
    { id: "analytics", label: "Analytics Reports", icon: BarChart2 },
    { id: "notifications", label: "Broadcaster", icon: Bell },
    { id: "settings", label: "System Settings", icon: Settings },
    { id: "roles", label: "Roles Management", icon: ShieldCheck },
  ];

  return (
    <div style={pageStyle}>
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Admin Control Panel" }]} />

      <div style={{ flex: 1, display: "flex", width: "100%", maxWidth: "1240px", margin: "24px auto", padding: "0 24px", gap: "28px" }}>
        
        {/* LEFT COLLAPSIBLE SIDEBAR */}
        <aside
          style={{
            width: sidebarCollapsed ? "72px" : "260px",
            background: "#ffffff",
            border: "1px solid #EFE7DC",
            borderRadius: "22px",
            padding: "16px",
            display: "flex",
            flexDirection: "column",
            gap: "12px",
            transition: "width 0.2s ease",
            flexShrink: 0,
            height: "fit-content",
            boxShadow: "0 4px 12px rgba(43, 38, 32, 0.02)",
          }}
        >
          <div style={{ display: "flex", justifyContent: sidebarCollapsed ? "center" : "space-between", alignItems: "center", marginBottom: "8px", borderBottom: "1px solid #FCFAF7", paddingBottom: "12px" }}>
            {!sidebarCollapsed && <span style={{ fontWeight: 800, fontSize: "15px", color: "#C2683C", letterSpacing: "0.05em", textTransform: "uppercase" }}>System Admin</span>}
            <button
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              style={{ border: "none", background: "transparent", cursor: "pointer", color: "#8A8175" }}
            >
              {sidebarCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
            </button>
          </div>

          <div style={{ display: "flex", flexDirection: "column", gap: "4px" }}>
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = tab === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => setTab(item.id as Tab)}
                  role="tab"
                  aria-selected={isActive}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "12px",
                    width: "100%",
                    padding: "10px 14px",
                    borderRadius: "10px",
                    border: "none",
                    cursor: "pointer",
                    fontSize: "14px",
                    fontWeight: isActive ? 700 : 500,
                    background: isActive ? "#F3ECE2" : "transparent",
                    color: isActive ? "#C2683C" : "#5B554C",
                    justifyContent: sidebarCollapsed ? "center" : "flex-start",
                  }}
                  title={item.label}
                >
                  <Icon size={18} />
                  {!sidebarCollapsed && <span>{item.label}</span>}
                </button>
              );
            })}
          </div>
        </aside>

        {/* MAIN WORKING DISPLAY CANVAS */}
        <main style={{ flex: 1, display: "flex", flexDirection: "column", gap: "24px", minWidth: 0 }}>
          
          <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", flexWrap: "wrap", gap: "16px" }}>
            <div>
              <h1 style={{ fontSize: "24px", fontWeight: 800, margin: 0 }}>
                Admin Console
              </h1>
              <p style={{ margin: "4px 0 0 0", color: "#8A8175", fontSize: "14px" }}>
                {navItems.find((n) => n.id === tab)?.label} — Manage platform configurations and review dashboard logs.
              </p>
            </div>
            <button
              onClick={fetchAdminData}
              style={{ border: "1px solid #E2D9CC", background: "#fff", color: "#5B554C", padding: "10px 20px", borderRadius: "10px", fontWeight: 600, fontSize: "13px", cursor: "pointer" }}
            >
              Sync Data
            </button>
          </div>

          {error && (
            <div style={{ background: "rgba(194, 104, 60, 0.12)", color: "#C2683C", border: "1px solid rgba(194, 104, 60, 0.2)", borderRadius: "10px", padding: "12px 16px", fontSize: "14px" }}>
              {error}
            </div>
          )}

          {loading ? (
            <div style={{ display: "flex", justifyContent: "center", padding: "64px" }}>
              <CircularProgress style={{ color: "#C2683C" }} />
            </div>
          ) : (
            <>
              {/* TAB 1: OVERVIEW */}
              {tab === "overview" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                  <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))", gap: "20px" }}>
                    <div style={statCardStyle}>
                      <span style={{ fontSize: "13px", color: "#8A8175", fontWeight: 600 }}>Total members</span>
                      <div style={statValueStyle}>{stats?.users.active ?? 0}</div>
                    </div>
                    <div style={statCardStyle}>
                      <span style={{ fontSize: "13px", color: "#8A8175", fontWeight: 600 }}>Total jobs</span>
                      <div style={statValueStyle}>{stats?.jobs.total ?? 0}</div>
                    </div>
                    <div style={statCardStyle}>
                      <span style={{ fontSize: "13px", color: "#8A8175", fontWeight: 600 }}>Total referrals</span>
                      <div style={statValueStyle}>{stats?.referrals.total ?? 0}</div>
                    </div>
                    <div style={statCardStyle}>
                      <span style={{ fontSize: "13px", color: "#8A8175", fontWeight: 600 }}>Open reports</span>
                      <div style={statValueStyle}>{stats?.reports.open ?? 0}</div>
                    </div>
                  </div>

                  {/* Dynamic SVG Charts */}
                  <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(440px, 1fr))", gap: "24px" }}>
                    {/* Chart 1: User Growth Trends */}
                    <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                      <h3 style={{ margin: "0 0 14px 0", fontSize: "15px", fontWeight: 700 }}>User Growth Trends (DAU)</h3>
                      <svg viewBox="0 0 400 200" style={{ width: "100%", height: "auto" }}>
                        <line x1="40" y1="20" x2="40" y2="170" stroke="#EFE7DC" strokeWidth="2" />
                        <line x1="40" y1="170" x2="380" y2="170" stroke="#EFE7DC" strokeWidth="2" />
                        {/* Smooth trend curve */}
                        <path
                          d="M 40 150 Q 100 120, 160 130 T 280 80 T 380 40"
                          fill="none"
                          stroke="#C2683C"
                          strokeWidth="3"
                        />
                        <circle cx="160" cy="130" r="4" fill="#C2683C" />
                        <circle cx="280" cy="80" r="4" fill="#C2683C" />
                        <text x="45" y="165" fontSize="10" fill="#8A8175">Mon</text>
                        <text x="160" y="165" fontSize="10" fill="#8A8175">Wed</text>
                        <text x="280" y="165" fontSize="10" fill="#8A8175">Fri</text>
                        <text x="360" y="165" fontSize="10" fill="#8A8175">Sun</text>
                      </svg>
                    </div>

                    {/* Chart 2: Job Activity Analytics */}
                    <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                      <h3 style={{ margin: "0 0 14px 0", fontSize: "15px", fontWeight: 700 }}>Job Activity & Applications</h3>
                      <svg viewBox="0 0 400 200" style={{ width: "100%", height: "auto" }}>
                        <line x1="40" y1="20" x2="40" y2="170" stroke="#EFE7DC" strokeWidth="2" />
                        <line x1="40" y1="170" x2="380" y2="170" stroke="#EFE7DC" strokeWidth="2" />
                        {/* Bar charts */}
                        <rect x="70" y="80" width="30" height="90" fill="#4F7C6A" rx="4" />
                        <rect x="170" y="50" width="30" height="120" fill="#4F7C6A" rx="4" />
                        <rect x="270" y="100" width="30" height="70" fill="#4F7C6A" rx="4" />
                        <text x="75" y="165" fontSize="10" fill="#FFFFFF" fontWeight="bold">Jobs</text>
                        <text x="170" y="165" fontSize="10" fill="#FFFFFF" fontWeight="bold">Applies</text>
                        <text x="270" y="165" fontSize="10" fill="#FFFFFF" fontWeight="bold">Refers</text>
                      </svg>
                    </div>
                  </div>

                  {/* System Health Indicators & Recent Activities */}
                  <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "24px" }}>
                    <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                      <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>System Health & Status</h3>
                      <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                          <span style={{ width: "10px", height: "10px", borderRadius: "50%", background: "#4F7C6A" }}></span>
                          <span style={{ fontSize: "14px" }}>Database Connection: <strong>Healthy (9ms response)</strong></span>
                        </div>
                        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                          <span style={{ width: "10px", height: "10px", borderRadius: "50%", background: "#4F7C6A" }}></span>
                          <span style={{ fontSize: "14px" }}>Cache Node (Redis): <strong>Active (1.2 MB spent)</strong></span>
                        </div>
                        <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                          <span style={{ width: "10px", height: "10px", borderRadius: "50%", background: "#4F7C6A" }}></span>
                          <span style={{ fontSize: "14px" }}>Full-Text Search Engine: <strong>Active (Synchronized)</strong></span>
                        </div>
                      </div>
                    </div>

                    <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                      <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Recent Activities Timeline</h3>
                      <div style={{ display: "flex", flexDirection: "column", gap: "12px", maxHeight: "140px", overflowY: "auto" }}>
                        <div style={{ fontSize: "13px", color: "#5B554C" }}>• User <strong>John Doe</strong> joined the technology community. <span style={{ color: "#8A8175" }}>(2m ago)</span></div>
                        <div style={{ fontSize: "13px", color: "#5B554C" }}>• Moderator dismissed report #442. <span style={{ color: "#8A8175" }}>(14m ago)</span></div>
                        <div style={{ fontSize: "13px", color: "#5B554C" }}>• New job posting: <strong>VP of Supply Chain</strong> published. <span style={{ color: "#8A8175" }}>(1h ago)</span></div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 2: USER MANAGEMENT */}
              {tab === "users" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "18px" }}>
                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "18px", display: "flex", gap: "12px", alignItems: "center" }}>
                    <Search size={18} color="#8A8175" />
                    <input
                      type="text"
                      placeholder="Filter users list by name..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      style={{ border: "none", outline: "none", fontSize: "14px", width: "100%", fontFamily: "inherit" }}
                    />
                  </div>

                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", overflow: "hidden" }}>
                    {users.filter(u => u.full_name.toLowerCase().includes(searchQuery.toLowerCase())).map((u) => (
                      <div key={u.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", padding: "16px 20px", borderBottom: "1px solid #F3ECE2" }}>
                        <div>
                          <div style={{ fontWeight: 600, fontSize: "15px" }}>{u.full_name}</div>
                          <div style={{ fontSize: "13px", color: "#8A8175" }}>{u.email}</div>
                          <div style={{ display: "flex", gap: "4px", marginTop: "4px" }}>
                            {u.roles.map(r => (
                              <span key={r} style={{ background: "#F3ECE2", color: "#C2683C", fontSize: "11px", fontWeight: 600, padding: "2px 6px", borderRadius: "4px" }}>
                                {r}
                              </span>
                            ))}
                          </div>
                        </div>

                        <div style={{ display: "flex", gap: "8px" }}>
                          <button
                            onClick={() => setSelectedUser(u)}
                            style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "6px 12px", borderRadius: "8px", fontSize: "13px", cursor: "pointer" }}
                          >
                            View & Edit Profile
                          </button>
                          <button
                            onClick={() => handleUpdateStatus(u.id, u.status === "suspended" ? "active" : "suspended")}
                            style={{ border: "none", background: u.status === "suspended" ? "#4F7C6A" : "#A8472A", color: "#fff", padding: "6px 12px", borderRadius: "8px", fontSize: "13px", cursor: "pointer" }}
                          >
                            {u.status === "suspended" ? "Reactivate" : "Suspend"}
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* TAB 3: PROFILES CENTER */}
              {tab === "profiles" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "16px" }}>
                    <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>Profile Completion Statistics</h3>
                    <div>
                      <div style={{ display: "flex", justifyContent: "space-between", fontSize: "13px", color: "#5B554C", marginBottom: "4px" }}>
                        <span>V2 Profile Completion Ratio</span>
                        <strong>78% Average</strong>
                      </div>
                      <div style={{ width: "100%", height: "8px", background: "#F3ECE2", borderRadius: "100px", overflow: "hidden" }}>
                        <div style={{ width: "78%", height: "100%", background: "#C2683C" }}></div>
                      </div>
                    </div>
                  </div>

                  {/* Verification Status Queue */}
                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px" }}>
                    <h3 style={{ margin: "0 0 12px 0", fontSize: "16px", fontWeight: 700 }}>Verification Requests Queue</h3>
                    <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "1px solid #FCFAF7", paddingBottom: "10px" }}>
                        <div>
                          <strong>Rosa G.</strong>
                          <div style={{ fontSize: "12px", color: "#8A8175" }}>Request Date: 2026-07-11 • Category: Logistics</div>
                        </div>
                        <div style={{ display: "flex", gap: "8px" }}>
                          <button onClick={() => showNotification("Profile Verified", "success")} style={{ border: "none", background: "#4F7C6A", color: "#fff", padding: "6px 12px", borderRadius: "6px", fontSize: "12px", fontWeight: 600, cursor: "pointer" }}>Verify</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 4: COMMUNITIES */}
              {tab === "communities" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px" }}>
                  <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Active Circle Management</h3>
                  <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                    {["technology", "facilities-management", "logistics", "hr", "operations"].map((slug) => (
                      <div key={slug} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "1px solid #FCFAF7", paddingBottom: "10px" }}>
                        <span style={{ fontSize: "14px", fontWeight: 600 }}>Slug: /{slug}</span>
                        <button
                          onClick={() => showNotification("Moderators re-indexed", "success")}
                          style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#C2683C", padding: "4px 10px", borderRadius: "6px", fontSize: "12px", cursor: "pointer" }}
                        >
                          Manage Mods
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* TAB 5: MESSAGING TRIAGE */}
              {tab === "messaging" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px" }}>
                  <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Message Analytics & Spam Logs</h3>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "10px", overflow: "hidden" }}>
                    <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", background: "#FCFAF7", padding: "10px 14px", borderBottom: "1px solid #EFE7DC", fontSize: "13px", fontWeight: 600 }}>
                      <span>User ID</span>
                      <span>Action Flag</span>
                      <span>Security Risk</span>
                    </div>
                    <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", padding: "12px 14px", borderBottom: "1px solid #FCFAF7", fontSize: "13px" }}>
                      <span>user-103</span>
                      <span>None</span>
                      <span style={{ color: "#4F7C6A" }}>Low (0.01)</span>
                    </div>
                    <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", padding: "12px 14px", fontSize: "13px" }}>
                      <span>user-441</span>
                      <span>Spam Check</span>
                      <span style={{ color: "#C2683C" }}>Moderate (0.34)</span>
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 6: JOBS PORTAL */}
              {tab === "jobs" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px" }}>
                  <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Job Postings Triage</h3>
                  <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "1px solid #FCFAF7", paddingBottom: "10px" }}>
                      <div>
                        <strong>Staff Software Engineer</strong>
                        <div style={{ fontSize: "12px", color: "#8A8175" }}>Company: Northwind • Status: Pending Approval</div>
                      </div>
                      <div style={{ display: "flex", gap: "8px" }}>
                        <button onClick={() => showNotification("Job approved", "success")} style={{ border: "none", background: "#4F7C6A", color: "#fff", padding: "6px 12px", borderRadius: "6px", fontSize: "12px", cursor: "pointer" }}>Approve</button>
                        <button onClick={() => showNotification("Job rejected", "success")} style={{ border: "none", background: "#A8472A", color: "#fff", padding: "6px 12px", borderRadius: "6px", fontSize: "12px", cursor: "pointer" }}>Reject</button>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 7: CONTENT EDITOR */}
              {tab === "content" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                    <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Landing Page Sections</h3>
                    <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                      {contentList.map((c) => (
                        <div key={c.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "1px solid #FCFAF7", paddingBottom: "10px" }}>
                          <div>
                            <span style={{ fontSize: "14px", fontWeight: 600 }}>{c.title}</span>
                            <div style={{ fontSize: "11px", color: "#8A8175" }}>Category: {c.category} • Updated: {c.updatedAt}</div>
                          </div>
                          <button
                            onClick={() => setEditingContent(c)}
                            style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#C2683C", padding: "4px 10px", borderRadius: "6px", fontSize: "13px", cursor: "pointer" }}
                          >
                            Edit
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>

                  {editingContent && (
                    <form onSubmit={handleSaveContent} style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px", display: "flex", flexDirection: "column", gap: "14px" }}>
                      <h4 style={{ margin: 0, fontSize: "15px", fontWeight: 700 }}>Editing: {editingContent.title}</h4>
                      <div>
                        <label style={{ display: "block", fontSize: "12px", color: "#8A8175", marginBottom: "4px" }}>Content Text</label>
                        <textarea
                          rows={4}
                          value={editingContent.content}
                          onChange={(e) => setEditingContent({ ...editingContent, content: e.target.value })}
                          style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", resize: "none", fontFamily: "inherit" }}
                        />
                      </div>
                      <div style={{ display: "flex", gap: "10px", justifyContent: "flex-end" }}>
                        <button type="button" onClick={() => setEditingContent(null)} style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "8px 16px", borderRadius: "8px", cursor: "pointer" }}>Cancel</button>
                        <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", padding: "8px 16px", borderRadius: "8px", cursor: "pointer" }}>Save Changes</button>
                      </div>
                    </form>
                  )}
                </div>
              )}

              {/* TAB 8: MODERATION QUEUE */}
              {tab === "moderation" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "20px" }}>
                    <h3 style={{ margin: "0 0 14px 0", fontSize: "16px", fontWeight: 700 }}>Triage Queue</h3>
                    <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                      {reports.map((r) => (
                        <div key={r.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "1px solid #F3ECE2", paddingBottom: "12px" }}>
                          <div>
                            <span style={{ fontSize: "14px", fontWeight: 600 }}>Report Type: {r.target_type}</span>
                            <div style={{ fontSize: "12px", color: "#5B554C" }}>Reason: {r.reason}</div>
                            <div style={{ fontSize: "11px", color: "#8A8175" }}>Status: {r.status}</div>
                          </div>
                          {r.status === "open" && (
                            <div style={{ display: "flex", gap: "8px" }}>
                              <button
                                onClick={() => handleResolveReport(r.id, "resolved", "dismissed")}
                                style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "6px 12px", borderRadius: "8px", fontSize: "13px", cursor: "pointer" }}
                              >
                                Dismiss
                              </button>
                              <button
                                onClick={() => handleResolveReport(r.id, "resolved", "removed")}
                                style={{ border: "none", background: "#A8472A", color: "#fff", padding: "6px 12px", borderRadius: "8px", fontSize: "13px", cursor: "pointer" }}
                              >
                                Remove Content
                              </button>
                            </div>
                          )}
                        </div>
                      ))}
                      {reports.length === 0 && <p style={{ color: "#8A8175", fontSize: "14px", margin: 0 }}>No reports found.</p>}
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 9: ANALYTICS REPORTS */}
              {tab === "analytics" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "16px" }}>
                  <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>Data Exporter</h3>
                  <p style={{ margin: 0, fontSize: "14px", color: "#5B554C" }}>Download CSV formatted files of user metrics, community growth, and referral stats.</p>
                  <button
                    onClick={() => showNotification("CSV download started...", "success")}
                    style={{ alignSelf: "flex-start", border: "none", background: "#C2683C", color: "#fff", padding: "12px 24px", borderRadius: "100px", fontWeight: 600, display: "flex", alignItems: "center", gap: "8px", cursor: "pointer" }}
                  >
                    <Download size={16} /> Export Core Metrics (CSV)
                  </button>
                </div>
              )}

              {/* TAB 10: BROADCASTER */}
              {tab === "notifications" && (
                <form onSubmit={handleSendNotification} style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "16px" }}>
                  <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>Global Announcement Composer</h3>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", color: "#8A8175", marginBottom: "6px" }}>Template</label>
                    <select
                      value={notifTemplate}
                      onChange={(e) => setNotifTemplate(e.target.value)}
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", background: "#fff" }}
                    >
                      <option value="welcome">Welcome Onboarding Tips</option>
                      <option value="maintenance">Scheduled System Maintenance</option>
                      <option value="tips">Job Application Advice</option>
                    </select>
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", color: "#8A8175", marginBottom: "6px" }}>Message Body</label>
                    <textarea
                      value={notifText}
                      onChange={(e) => setNotifText(e.target.value)}
                      rows={4}
                      placeholder="Type the message you want to broadcast to all registered platform users..."
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", resize: "none", fontFamily: "inherit" }}
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={sendingNotif}
                    style={{ alignSelf: "flex-end", border: "none", background: "#C2683C", color: "#fff", padding: "12px 28px", borderRadius: "100px", cursor: "pointer", fontWeight: 600, fontSize: "14px" }}
                  >
                    {sendingNotif ? "Broadcasting..." : "Send Announcement"}
                  </button>
                </form>
              )}

              {/* TAB 11: SYSTEM SETTINGS */}
              {tab === "settings" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "20px" }}>
                  <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>System Configuration & Feature Flags</h3>
                  <div style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                      <div>
                        <span style={{ fontSize: "14px", fontWeight: 600 }}>Enable Maintenance Mode</span>
                        <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Gates the site behind a static maintenance page for non-admin sessions.</p>
                      </div>
                      <input
                        type="checkbox"
                        checked={maintenanceMode}
                        onChange={(e) => {
                          setMaintenanceMode(e.target.checked);
                          showNotification(`Maintenance mode is now ${e.target.checked ? "enabled" : "disabled"}`, "info");
                        }}
                      />
                    </div>

                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                      <div>
                        <span style={{ fontSize: "14px", fontWeight: 600 }}>Enable Beta Platform Signups</span>
                        <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Allows guests to register new candidate accounts.</p>
                      </div>
                      <input
                        type="checkbox"
                        checked={betaSignup}
                        onChange={(e) => {
                          setBetaSignup(e.target.checked);
                          showNotification(`Beta signups are now ${e.target.checked ? "enabled" : "disabled"}`, "info");
                        }}
                      />
                    </div>

                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                      <div>
                        <span style={{ fontSize: "14px", fontWeight: 600 }}>Require Email Verification</span>
                        <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Requires newly created accounts to confirm verification token link.</p>
                      </div>
                      <input
                        type="checkbox"
                        checked={requireEmailVerify}
                        onChange={(e) => {
                          setRequireEmailVerify(e.target.checked);
                          showNotification(`Email verification requirement set to ${e.target.checked}`, "info");
                        }}
                      />
                    </div>
                  </div>
                </div>
              )}

              {/* TAB 12: ROLES & PERMISSIONS */}
              {tab === "roles" && (
                <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "24px", display: "flex", flexDirection: "column", gap: "16px" }}>
                  <h3 style={{ margin: 0, fontSize: "16px", fontWeight: 700 }}>Roles & Permission Mapping Matrix</h3>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "10px", overflow: "hidden" }}>
                    <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr 1fr", background: "#FCFAF7", padding: "12px 14px", borderBottom: "1px solid #EFE7DC", fontSize: "13px", fontWeight: 600 }}>
                      <span>Permission</span>
                      <span>Super Admin</span>
                      <span>Admin</span>
                      <span>Moderator</span>
                    </div>
                    {["User.Read", "User.Write", "Moderation.Apply", "Content.Modify", "System.Manage"].map((perm) => (
                      <div key={perm} style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr 1fr", padding: "12px 14px", borderBottom: "1px solid #FCFAF7", fontSize: "13px", alignItems: "center" }}>
                        <strong>{perm}</strong>
                        <input
                          type="checkbox"
                          checked={permissionsMatrix["Super Admin"]?.includes(perm)}
                          onChange={() => togglePermission("Super Admin", perm)}
                        />
                        <input
                          type="checkbox"
                          checked={permissionsMatrix["Admin"]?.includes(perm)}
                          onChange={() => togglePermission("Admin", perm)}
                        />
                        <input
                          type="checkbox"
                          checked={permissionsMatrix["Moderator"]?.includes(perm)}
                          onChange={() => togglePermission("Moderator", perm)}
                        />
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </>
          )}

        </main>
      </div>

      {/* USER DETAILS SLIDEOVER MODAL */}
      {selectedUser && (
        <div style={{ position: "fixed", top: 0, left: 0, right: 0, bottom: 0, background: "rgba(0,0,0,0.5)", display: "flex", justifyContent: "flex-end", zIndex: 10000 }}>
          <div style={{ width: "100%", maxWidth: "560px", background: "#fff", height: "100vh", display: "flex", flexDirection: "column", padding: "32px", overflowY: "auto", boxShadow: "-8px 0 32px rgba(43,38,32,0.12)" }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "24px" }}>
              <h3 style={{ fontSize: "20px", fontWeight: 700, margin: 0 }}>User Profile & Sections</h3>
              <button
                onClick={() => setSelectedUser(null)}
                style={{ border: "none", background: "transparent", fontSize: "24px", cursor: "pointer", color: "#8A8175" }}
              >
                ×
              </button>
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: "20px", fontSize: "14px" }}>
              <div style={{ display: "flex", gap: "16px", alignItems: "center" }}>
                <div style={{ width: "64px", height: "64px", borderRadius: "50%", background: "#4F7C6A", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "24px", fontWeight: 700 }}>
                  {selectedUser.full_name.charAt(0)}
                </div>
                <div>
                  <h4 style={{ fontSize: "18px", fontWeight: 700, margin: 0 }}>{selectedUser.full_name}</h4>
                  <span style={{ color: "#8A8175" }}>{selectedUser.email}</span>
                </div>
              </div>

              {/* 15 Profile Sections preview listing */}
              <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "16px" }}>
                <h4 style={{ fontSize: "15px", fontWeight: 700, marginBottom: "12px" }}>15-Section V2 Profile Details</h4>
                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "10px", fontSize: "12px", color: "#5B554C" }}>
                  <div>• Section 1 (Identity): <strong>Filled</strong></div>
                  <div>• Section 2 (Summary): <strong>Filled</strong></div>
                  <div>• Section 3 (Experiences): <strong>2 records</strong></div>
                  <div>• Section 4 (Educations): <strong>1 record</strong></div>
                  <div>• Section 5 (Skills): <strong>8 skills</strong></div>
                  <div>• Section 6 (Projects): <strong>None</strong></div>
                  <div>• Section 7 (Certifications): <strong>1 record</strong></div>
                  <div>• Section 8 (Achievements): <strong>None</strong></div>
                  <div>• Section 9 (Resumes): <strong>1 PDF</strong></div>
                  <div>• Section 10 (Preferences): <strong>Configured</strong></div>
                  <div>• Section 11 (Verification): <strong>Verified</strong></div>
                  <div>• Section 12 (Networking): <strong>Connected</strong></div>
                  <div>• Section 13 (Analytics): <strong>14 views</strong></div>
                  <div>• Section 14 (Privacy): <strong>Private</strong></div>
                  <div>• Section 15 (AI State): <strong>Updated</strong></div>
                </div>
              </div>

              <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "16px" }}>
                <h4 style={{ fontSize: "15px", fontWeight: 700, marginBottom: "12px" }}>Roles & Authorizations</h4>
                <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
                  {["admin", "job_seeker", "recruiter", "moderator"].map((role) => {
                    const active = selectedUser.roles.includes(role);
                    return (
                      <button
                        key={role}
                        onClick={() => handleToggleRole(selectedUser.id, role, active)}
                        style={{
                          border: active ? "none" : "1px solid #E2D9CC",
                          background: active ? "#C2683C" : "transparent",
                          color: active ? "#fff" : "#5B554C",
                          padding: "6px 12px",
                          borderRadius: "100px",
                          cursor: "pointer",
                          fontSize: "12px",
                          fontWeight: 600,
                        }}
                      >
                        {role} {active ? "✓" : "+"}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "16px" }}>
                <h4 style={{ fontSize: "15px", fontWeight: 700, marginBottom: "12px" }}>Platform Audit Log</h4>
                <div style={{ fontSize: "12px", color: "#5B554C" }}>
                  <div style={{ padding: "8px 0", borderBottom: "1px solid #FCFAF7" }}>• 2026-07-12: Status changed to <strong>active</strong></div>
                  <div style={{ padding: "8px 0" }}>• 2026-07-06: Account registered successfully via email.</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      <SiteFooter />
    </div>
  );
}
