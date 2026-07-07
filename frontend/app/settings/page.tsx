"use client";

import { useState } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";

type TabType =
  | "account"
  | "roles"
  | "appearance"
  | "preferences"
  | "privacy"
  | "notifications"
  | "security"
  | "platform"
  | "actions";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabType>("account");

  // Mock Form States
  const [fullName, setFullName] = useState("Marcus Hale");
  const [headline, setHeadline] = useState("Operations Director");
  const [email, setEmail] = useState("marcus.hale@email.com");
  const [phone, setPhone] = useState("(303) 555-0148");
  const [location, setLocation] = useState("Denver, CO");
  const [pronouns, setPronouns] = useState("he / him");

  // Notifications State
  const [emailDigests, setEmailDigests] = useState(true);
  const [newJobsAlerts, setNewJobsAlerts] = useState(true);
  const [messagesAlerts, setMessagesAlerts] = useState(true);

  // Privacy State
  const [profileVis, setProfileVis] = useState("public");
  const [showSalary, setShowSalary] = useState(false);

  // Appearance State
  const [theme, setTheme] = useState("light");

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    alert("Changes saved successfully!");
  };

  const navItems = [
    { id: "account" as TabType, label: "Account", icon: "◍" },
    { id: "roles" as TabType, label: "My roles", icon: "◈" },
    { id: "appearance" as TabType, label: "Appearance", icon: "◐" },
    { id: "preferences" as TabType, label: "Job preferences", icon: "✦" },
    { id: "privacy" as TabType, label: "Privacy & visibility", icon: "⊘" },
    { id: "notifications" as TabType, label: "Notifications", icon: "◔" },
    { id: "security" as TabType, label: "Password & security", icon: "⚿" },
    { id: "platform" as TabType, label: "Platform Info", icon: "▤" },
    { id: "actions" as TabType, label: "Account actions", icon: "⚠" },
  ];

  return (
    <AuthGuard>
      <div
        style={{
          background: "#FBF7F2",
          fontFamily: "'Public Sans', sans-serif",
          color: "#2B2620",
          minHeight: "100vh",
          overflowX: "hidden",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <SiteNav
          breadcrumb={[{ label: "Home", href: "/" }, { label: "Settings" }]}
        />

        <section style={{ maxWidth: "1180px", margin: "0 auto", width: "100%", padding: "clamp(32px,4vw,48px) 40px clamp(16px,2vw,24px)" }}>
          <div style={{ fontSize: "13px", fontWeight: 700, letterSpacing: "0.12em", textTransform: "uppercase", color: "#C2683C", marginBottom: "10px" }}>
            Settings
          </div>
          <h1 style={{ fontWeight: 800, fontSize: "clamp(30px,4vw,44px)", lineHeight: 1.02, letterSpacing: "-0.025em", margin: 0 }}>
            Manage your account
          </h1>
        </section>

        <section style={{ maxWidth: "1180px", margin: "0 auto", width: "100%", padding: "0 40px clamp(56px,6vw,90px)", display: "grid", gridTemplateColumns: "248px 1fr", gap: "28px", alignItems: "start", flex: 1 }}>
          
          {/* Navigation Sidebar */}
          <aside style={{ position: "sticky", top: "96px", background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "14px", display: "flex", flexDirection: "column", gap: "3px" }}>
            {navItems.map((item) => {
              const isActive = activeTab === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => setActiveTab(item.id)}
                  style={{
                    width: "100%",
                    textAlign: "left",
                    cursor: "pointer",
                    fontFamily: "'Public Sans', sans-serif",
                    border: "none",
                    borderRadius: "11px",
                    padding: "12px 14px",
                    display: "flex",
                    alignItems: "center",
                    gap: "12px",
                    fontSize: "15px",
                    fontWeight: isActive ? 600 : 500,
                    color: isActive ? "#2B2620" : "#5B554C",
                    background: isActive ? "#EFE7DC" : "transparent",
                    transition: "all 0.15s ease",
                  }}
                >
                  <span style={{ flex: "none", width: "20px", textAlign: "center" }}>{item.icon}</span>
                  <span>{item.label}</span>
                </button>
              );
            })}
          </aside>

          {/* Active Tab Panel */}
          <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "clamp(24px,3vw,32px)" }}>
            
            {activeTab === "account" && (
              <form onSubmit={handleSave} style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Account</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Your basic details and public profile identification.</p>
                </div>

                <div style={{ display: "flex", gap: "18px", alignItems: "center", flexWrap: "wrap" }}>
                  <img
                    src="/assets/avatar-marcus.svg"
                    alt={fullName}
                    style={{ width: "76px", height: "76px", borderRadius: "18px", objectFit: "cover", background: "#F3E7DC" }}
                  />
                  <div style={{ display: "flex", gap: "10px" }}>
                    <button type="button" style={{ border: "1px solid #E2D9CC", background: "#fff", color: "#2B2620", fontSize: "14px", fontWeight: 600, padding: "11px 20px", borderRadius: "100px", cursor: "pointer" }}>
                      Change photo
                    </button>
                    <button type="button" style={{ border: "none", background: "transparent", color: "#A8472A", fontSize: "14px", fontWeight: 600, padding: "11px 6px", cursor: "pointer" }}>
                      Remove
                    </button>
                  </div>
                </div>

                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))", gap: "16px" }}>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Full name</label>
                    <input value={fullName} onChange={(e) => setFullName(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Headline</label>
                    <input value={headline} onChange={(e) => setHeadline(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                  <div>
                    <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", gap: "8px", marginBottom: "7px" }}>
                      <label style={{ fontSize: "13px", fontWeight: 600, color: "#8A8175" }}>Email</label>
                      <span style={{ display: "inline-flex", alignItems: "center", gap: "5px", fontSize: "12px", fontWeight: 600, color: "#4F7C6A", background: "rgba(79,124,106,0.12)", padding: "3px 9px", borderRadius: "100px" }}>
                        ✓ Verified
                      </span>
                    </div>
                    <input value={email} onChange={(e) => setEmail(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                  <div>
                    <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", gap: "8px", marginBottom: "7px" }}>
                      <label style={{ fontSize: "13px", fontWeight: 600, color: "#8A8175" }}>Phone</label>
                      <span style={{ display: "inline-flex", alignItems: "center", gap: "5px", fontSize: "12px", fontWeight: 600, color: "#A8472A", background: "rgba(168,71,42,0.10)", padding: "3px 9px", borderRadius: "100px" }}>
                        ! Unverified
                      </span>
                    </div>
                    <div style={{ display: "flex", gap: "8px" }}>
                      <input value={phone} onChange={(e) => setPhone(e.target.value)} style={{ flex: 1, minWidth: 0, border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                      <button type="button" style={{ flex: "none", border: "1px solid #C2683C", background: "transparent", color: "#C2683C", fontSize: "14px", fontWeight: 600, padding: "0 16px", borderRadius: "10px", cursor: "pointer" }}>Verify</button>
                    </div>
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Location</label>
                    <input value={location} onChange={(e) => setLocation(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Pronouns</label>
                    <input value={pronouns} onChange={(e) => setPronouns(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                </div>

                <div style={{ display: "flex", gap: "12px", marginTop: "12px" }}>
                  <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", fontSize: "15px", fontWeight: 600, padding: "13px 28px", borderRadius: "100px", cursor: "pointer" }}>Save changes</button>
                </div>
              </form>
            )}

            {activeTab === "roles" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>My roles</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Select and switch your primary platform user roles.</p>
                <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "16px", fontWeight: 600 }}>Candidate Profile</h4>
                      <p style={{ margin: 0, fontSize: "14px", color: "#8A8175" }}>Apply for roles, utilize AI career coach, build resume profiles.</p>
                    </div>
                    <span style={{ background: "#4F7C6A", color: "#fff", fontSize: "12px", fontWeight: 600, padding: "6px 12px", borderRadius: "100px" }}>Active</span>
                  </div>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "16px", fontWeight: 600 }}>Recruiter Panel</h4>
                      <p style={{ margin: 0, fontSize: "14px", color: "#8A8175" }}>Sourced profiles, post jobs, and run candidates evaluations.</p>
                    </div>
                    <button style={{ border: "1px solid #C2683C", background: "transparent", color: "#C2683C", fontSize: "13px", fontWeight: 600, padding: "6px 12px", borderRadius: "100px", cursor: "pointer" }}>Activate</button>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "appearance" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Appearance</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Customize how Kirmya looks on your device.</p>
                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))", gap: "16px" }}>
                  <div
                    onClick={() => setTheme("light")}
                    style={{
                      border: theme === "light" ? "2px solid #C2683C" : "1px solid #EFE7DC",
                      borderRadius: "16px",
                      padding: "20px",
                      background: "#fff",
                      cursor: "pointer",
                      textAlign: "center",
                    }}
                  >
                    <div style={{ fontSize: "24px", marginBottom: "8px" }}>☀</div>
                    <div style={{ fontWeight: 600 }}>Light Mode</div>
                  </div>
                  <div
                    onClick={() => setTheme("dark")}
                    style={{
                      border: theme === "dark" ? "2px solid #C2683C" : "1px solid #EFE7DC",
                      borderRadius: "16px",
                      padding: "20px",
                      background: "#2B2620",
                      color: "#fff",
                      cursor: "pointer",
                      textAlign: "center",
                    }}
                  >
                    <div style={{ fontSize: "24px", marginBottom: "8px" }}>☾</div>
                    <div style={{ fontWeight: 600 }}>Dark Mode (System)</div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "preferences" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Job preferences</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Set your desired career preferences for better recommendation matching.</p>
                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <div>
                    <label style={{ display: "block", fontSize: "14px", fontWeight: 600, color: "#2B2620", marginBottom: "8px" }}>Desired Roles</label>
                    <input placeholder="e.g. Operations Director, Chief of Staff" style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }} />
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "14px", fontWeight: 600, color: "#2B2620", marginBottom: "8px" }}>Preferred Work Mode</label>
                    <div style={{ display: "flex", gap: "10px" }}>
                      {["Remote", "Hybrid", "Onsite"].map((mode) => (
                        <button key={mode} type="button" style={{ border: "1px solid #E2D9CC", background: "#fff", padding: "10px 20px", borderRadius: "100px", cursor: "pointer", fontSize: "14px", fontWeight: 600 }}>{mode}</button>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "privacy" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Privacy &amp; visibility</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Control who can view your profile and details.</p>
                <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                  <div>
                    <label style={{ display: "block", fontSize: "14px", fontWeight: 600, color: "#2B2620", marginBottom: "8px" }}>Profile Visibility</label>
                    <select value={profileVis} onChange={(e) => setProfileVis(e.target.value)} style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", background: "#FCFAF7", outline: "none" }}>
                      <option value="public">Public (Everyone can view)</option>
                      <option value="recruiters">Recruiters Only</option>
                      <option value="private">Private (Only you)</option>
                    </select>
                  </div>
                  <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                    <input type="checkbox" id="showSalary" checked={showSalary} onChange={(e) => setShowSalary(e.target.checked)} style={{ width: "18px", height: "18px", cursor: "pointer" }} />
                    <label htmlFor="showSalary" style={{ fontSize: "15px", cursor: "pointer", fontWeight: 500 }}>Show desired salary range to recruiters</label>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "notifications" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Notifications</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Decide how and when you want to receive alerts from Kirmya.</p>
                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", borderBottom: "1px solid #EFE7DC", paddingBottom: "16px" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>Weekly Digest</h4>
                      <p style={{ margin: 0, fontSize: "13px", color: "#8A8175" }}>Receive a summary of matching jobs, views, and networking invites.</p>
                    </div>
                    <input type="checkbox" checked={emailDigests} onChange={(e) => setEmailDigests(e.target.checked)} style={{ width: "36px", height: "20px", cursor: "pointer" }} />
                  </div>
                  <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", borderBottom: "1px solid #EFE7DC", paddingBottom: "16px" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>New Jobs Matching Criteria</h4>
                      <p style={{ margin: 0, fontSize: "13px", color: "#8A8175" }}>Instant notification when a job matching your preferences is posted.</p>
                    </div>
                    <input type="checkbox" checked={newJobsAlerts} onChange={(e) => setNewJobsAlerts(e.target.checked)} style={{ width: "36px", height: "20px", cursor: "pointer" }} />
                  </div>
                  <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>Message &amp; Chat Alerts</h4>
                      <p style={{ margin: 0, fontSize: "13px", color: "#8A8175" }}>Notifications for new inbox messages or career coach answers.</p>
                    </div>
                    <input type="checkbox" checked={messagesAlerts} onChange={(e) => setMessagesAlerts(e.target.checked)} style={{ width: "36px", height: "20px", cursor: "pointer" }} />
                  </div>
                </div>
              </div>
            )}

            {activeTab === "security" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Password &amp; security</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Keep your account safe by managing authentication settings.</p>
                <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "18px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "16px", fontWeight: 600 }}>Two-Factor Authentication (2FA)</h4>
                      <p style={{ margin: 0, fontSize: "14px", color: "#8A8175" }}>Enforce TOTP validation code during sign-ins.</p>
                    </div>
                    <button style={{ border: "1px solid #C2683C", background: "transparent", color: "#C2683C", fontSize: "13px", fontWeight: 600, padding: "8px 16px", borderRadius: "100px", cursor: "pointer" }}>Configure</button>
                  </div>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "18px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "16px", fontWeight: 600 }}>Change Password</h4>
                      <p style={{ margin: 0, fontSize: "14px", color: "#8A8175" }}>Update your password regularly to secure access.</p>
                    </div>
                    <button style={{ border: "1px solid #E2D9CC", background: "#fff", color: "#2B2620", fontSize: "13px", fontWeight: 600, padding: "8px 16px", borderRadius: "100px", cursor: "pointer" }}>Change</button>
                  </div>
                </div>
              </div>
            )}

            {/* Platform Information Section (Free Tier Requirement) */}
            {activeTab === "platform" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Platform Information</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Current Kirmya launch details, open-source compliance, and legal policies.</p>
                </div>

                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(260px, 1fr))", gap: "16px" }}>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7" }}>
                    <h4 style={{ margin: "0 0 6px", fontSize: "14px", fontWeight: 700, color: "#8A8175", textTransform: "uppercase", letterSpacing: "0.05em" }}>Version &amp; Build</h4>
                    <p style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#2B2620" }}>Kirmya Core v2.4.0</p>
                    <p style={{ margin: "4px 0 0", fontSize: "13px", color: "#8A8175" }}>Released July 2026. Built with Go, Gorilla Mux, and React.</p>
                  </div>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7" }}>
                    <h4 style={{ margin: "0 0 6px", fontSize: "14px", fontWeight: 700, color: "#8A8175", textTransform: "uppercase", letterSpacing: "0.05em" }}>Licensing &amp; Core</h4>
                    <p style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#2B2620" }}>100% Free &amp; Open Source</p>
                    <p style={{ margin: "4px 0 0", fontSize: "13px", color: "#8A8175" }}>Licensed under the MIT License.</p>
                  </div>
                </div>

                <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
                  <h3 style={{ margin: "0 0 12px 0", fontSize: "17px", fontWeight: 700 }}>Release Notes — What's New</h3>
                  <div style={{ display: "flex", flexDirection: "column", gap: "8px", color: "#5B554C", fontSize: "14px", lineHeight: "1.5" }}>
                    <p style={{ margin: 0 }}><strong>🚀 v2.4.0 (Current Release):</strong> Fully refactored 15-Section profile builder aggregate structure with transactional draft-snapshots, NATS JetStream search indexing, and real-time Claude LLM streaming support.</p>
                    <p style={{ margin: 0 }}><strong>🔒 v2.3.0:</strong> Implemented time-step TOTP replay prevention, secure JWT token headers, and strict same-site cookies validation.</p>
                  </div>
                </div>

                <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
                  <h3 style={{ margin: "0 0 12px 0", fontSize: "17px", fontWeight: 700 }}>Legal &amp; Compliance Policies</h3>
                  <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>
                    <a href="/legal/terms" style={{ textDecoration: "none", color: "#C2683C", fontSize: "14px", fontWeight: 600 }}>Terms of Service</a>
                    <span style={{ color: "#E2D9CC" }}>|</span>
                    <a href="/legal/privacy" style={{ textDecoration: "none", color: "#C2683C", fontSize: "14px", fontWeight: 600 }}>Privacy Policy</a>
                    <span style={{ color: "#E2D9CC" }}>|</span>
                    <a href="/legal/cookies" style={{ textDecoration: "none", color: "#C2683C", fontSize: "14px", fontWeight: 600 }}>Cookie Policy</a>
                    <span style={{ color: "#E2D9CC" }}>|</span>
                    <a href="/legal/licenses" style={{ textDecoration: "none", color: "#C2683C", fontSize: "14px", fontWeight: 600 }}>Licenses &amp; Attributions</a>
                  </div>
                </div>
              </div>
            )}

            {activeTab === "actions" && (
              <div>
                <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px", color: "#A8472A" }}>Account actions</h2>
                <p style={{ fontSize: "15px", color: "#8A8175", margin: "0 0 24px" }}>Danger zone actions to delete or export your personal details.</p>
                <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>
                  <button type="button" style={{ border: "1px solid #A8472A", background: "transparent", color: "#A8472A", fontSize: "14px", fontWeight: 600, padding: "12px 24px", borderRadius: "100px", cursor: "pointer" }}>Delete account</button>
                  <button type="button" style={{ border: "1px solid #E2D9CC", background: "#fff", color: "#2B2620", fontSize: "14px", fontWeight: 600, padding: "12px 24px", borderRadius: "100px", cursor: "pointer" }}>Export my data (JSON)</button>
                </div>
              </div>
            )}

          </div>

        </section>

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
