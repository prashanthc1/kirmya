"use client";

import { useState, useEffect } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api } from "@/lib/api/client";
import { useNotifications } from "@/components/shared/Notifications";
import { CircularProgress } from "@mui/material";

type TabType =
  | "general"
  | "profile-pref"
  | "account"
  | "security"
  | "privacy"
  | "notifications"
  | "ai"
  | "job-pref"
  | "learning"
  | "connected"
  | "cookies"
  | "accessibility"
  | "platform";

interface NavSection {
  id: TabType;
  label: string;
  desc: string;
  icon: string;
}

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabType>("general");
  const [searchQuery, setSearchQuery] = useState("");
  const [loading, setLoading] = useState(true);
  const [rateLimitWarning, setRateLimitWarning] = useState<string | null>(null);
  const { showNotification } = useNotifications();

  // Settings State
  const [settings, setSettings] = useState<any>(null);
  const [profileSettings, setProfileSettings] = useState<any>(null);
  const [securityActivity, setSecurityActivity] = useState<any>(null);
  const [cookieConsent, setCookieConsent] = useState<any>(null);

  // Forms states
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [usernameInput, setUsernameInput] = useState("");
  const [customURLInput, setCustomURLInput] = useState("");
  const [deactivateCheck, setDeactivateCheck] = useState(false);
  const [showDeactivateModal, setShowDeactivateModal] = useState(false);

  const [newGoal, setNewGoal] = useState("");

  const navSections: NavSection[] = [
    { id: "general", label: "General Settings", desc: "Language, timezone, theme preferences", icon: "⚙" },
    { id: "profile-pref", label: "Profile Settings", desc: "Custom URL, visibility, career status", icon: "◍" },
    { id: "account", label: "Account Settings", desc: "Username, email, account actions", icon: "👤" },
    { id: "security", label: "Security Settings", desc: "Change password, active sessions", icon: "🔒" },
    { id: "privacy", label: "Privacy Settings", desc: "Search discoverability, messaging control", icon: "⊘" },
    { id: "notifications", label: "Notification Settings", desc: "Job matches, mentions, digests", icon: "◔" },
    { id: "ai", label: "AI Preferences", desc: "AI assistant and recommendation options", icon: "🤖" },
    { id: "job-pref", label: "Job Preferences", desc: "Desired roles, work mode, salary expectations", icon: "💼" },
    { id: "learning", label: "Learning Preferences", desc: "Skill goals, study reminders", icon: "📚" },
    { id: "connected", label: "Connected Accounts", desc: "Google, LinkedIn, GitHub, Microsoft, Apple", icon: "🔗" },
    { id: "cookies", label: "Privacy & Cookies", desc: "Cookie consents, tracking options", icon: "🍪" },
    { id: "accessibility", label: "Accessibility Settings", desc: "Font size, high contrast, reduced motion", icon: "👁" },
    { id: "platform", label: "Platform Information", desc: "Version, license, release notes", icon: "▤" },
  ];

  const filteredSections = navSections.filter(
    (sec) =>
      sec.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
      sec.desc.toLowerCase().includes(searchQuery.toLowerCase())
  );

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [sData, pData, secData, cookieData] = await Promise.all([
        api.get<any>("/settings"),
        api.get<any>("/profile/settings"),
        api.get<any>("/security/activity"),
        api.get<any>("/privacy/cookies"),
      ]);
      setSettings(sData);
      setProfileSettings(pData);
      setSecurityActivity(secData);
      setCookieConsent(cookieData);

      setUsernameInput(pData.username || "");
      setCustomURLInput(pData.custom_url || "");

      // Apply initial accessibility settings
      if (sData.accessibility) {
        applyAccessibilityStyles(sData.accessibility);
      }
    } catch (err: any) {
      if (err.status === 429) {
        setRateLimitWarning("Too many requests. Please wait before saving any preferences.");
      }
      showNotification(err.message || "Failed to load settings data", "error");
    } finally {
      setLoading(false);
    }
  };

  const applyAccessibilityStyles = (acc: any) => {
    if (typeof document === "undefined") return;
    let styleTag = document.getElementById("accessibility-styles");
    if (!styleTag) {
      styleTag = document.createElement("style");
      styleTag.id = "accessibility-styles";
      document.head.appendChild(styleTag);
    }

    let css = "";
    if (acc.font_size === "small") {
      css += "body, html { font-size: 14px !important; }";
    } else if (acc.font_size === "large") {
      css += "body, html { font-size: 18px !important; }";
    } else if (acc.font_size === "extra-large") {
      css += "body, html { font-size: 20px !important; }";
    } else {
      css += "body, html { font-size: 16px !important; }";
    }

    if (acc.high_contrast) {
      css += `
        body, html {
          filter: contrast(1.25) !important;
          background: #FFFFFF !important;
          color: #000000 !important;
        }
        button, input, select, textarea {
          border: 2px solid #000000 !important;
          color: #000000 !important;
          background: #FFFFFF !important;
        }
      `;
    }

    if (acc.reduced_motion) {
      css += `
        *, *::before, *::after {
          animation-delay: -1ms !important;
          animation-duration: 1ms !important;
          animation-iteration-count: 1 !important;
          background-attachment: initial !important;
          scroll-behavior: auto !important;
          transition-duration: 0s !important;
          transition-delay: 0s !important;
        }
      `;
    }
    styleTag.innerHTML = css;
  };

  const saveSettingsSegment = async (segment: string, payload: any) => {
    try {
      const updated = await api.patch<any>("/settings", { [segment]: payload });
      setSettings(updated);
      showNotification("Settings updated successfully", "success");
      if (segment === "accessibility") {
        applyAccessibilityStyles(payload);
      }
    } catch (err: any) {
      if (err.status === 429) {
        showNotification("Rate limit reached. Please wait before retrying.", "warning");
      } else {
        showNotification(err.message || "Failed to save settings segment", "error");
      }
    }
  };

  const handleProfileSettingsSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (usernameInput.trim().length > 0 && !/^[a-zA-Z0-9_-]+$/.test(usernameInput)) {
      showNotification("Username must contain only alphanumeric characters, dashes, or underscores", "error");
      return;
    }
    if (customURLInput.trim().length > 0 && !/^[a-zA-Z0-9_-]+$/.test(customURLInput)) {
      showNotification("Custom URL must contain only alphanumeric characters, dashes, or underscores", "error");
      return;
    }

    try {
      await api.patch("/profile/settings", {
        username: usernameInput,
        custom_url: customURLInput,
        profile_visibility: profileSettings.profile_visibility,
        field_visibility: profileSettings.field_visibility,
        open_to_work: profileSettings.open_to_work,
        referral_eligible: profileSettings.referral_eligible,
        willing_to_mentor: profileSettings.willing_to_mentor,
      });
      showNotification("Profile settings saved", "success");
    } catch (err: any) {
      showNotification(err.message || "Failed to update profile settings", "error");
    }
  };

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      showNotification("New password confirmation does not match", "error");
      return;
    }
    if (newPassword.length < 8) {
      showNotification("New password must be at least 8 characters long", "error");
      return;
    }

    try {
      await api.post("/security/password/change", {
        current_password: currentPassword,
        new_password: newPassword,
      });
      showNotification("Password changed successfully", "success");
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err: any) {
      showNotification(err.message || "Failed to change password", "error");
    }
  };

  const handleRevokeSession = async (sessionId: string) => {
    try {
      await api.post("/security/logout-device", { session_id: sessionId });
      setSecurityActivity((prev: any) => ({
        ...prev,
        sessions: prev.sessions.filter((s: any) => s.id !== sessionId),
      }));
      showNotification("Session revoked", "success");
    } catch (err: any) {
      showNotification("Failed to revoke session", "error");
    }
  };

  const handleDisconnectConnectedAccount = async (provider: string) => {
    try {
      await api.delete(`/security/connected-account/${provider}`);
      setSecurityActivity((prev: any) => ({
        ...prev,
        connected_accounts: prev.connected_accounts.filter((a: any) => a.provider !== provider),
      }));
      showNotification(`Disconnected ${provider} account`, "success");
    } catch (err: any) {
      try {
        await api.post("/security/logout-device", { provider });
      } catch (inner: any) {
        showNotification(`Failed to disconnect ${provider}`, "error");
      }
    }
  };

  const handleSimulateOAuthConnect = (provider: string) => {
    showNotification(`Connecting to ${provider}... OAuth simulation complete!`, "success");
    setSecurityActivity((prev: any) => {
      const exists = prev.connected_accounts.some((a: any) => a.provider === provider);
      if (exists) return prev;
      return {
        ...prev,
        connected_accounts: [
          ...prev.connected_accounts,
          {
            id: Math.random().toString(),
            provider,
            provider_uid: `${provider}_user_123`,
            created_at: new Date().toISOString(),
          },
        ],
      };
    });
  };

  const handleCookieConsentSave = async (functional: boolean, analytics: boolean, ai: boolean) => {
    try {
      const cc = await api.patch<any>("/privacy/cookies", {
        functional,
        analytics,
        ai_personalization: ai,
      });
      setCookieConsent(cc);
      showNotification("Cookie preferences saved successfully", "success");
    } catch (err: any) {
      showNotification("Failed to save cookie consent", "error");
    }
  };

  const handleDeactivate = async () => {
    if (!deactivateCheck) {
      showNotification("Please check the confirmation box", "error");
      return;
    }
    try {
      await api.delete("/users/me");
      showNotification("Your account has been deactivated. Logging out...", "success");
      setTimeout(() => {
        window.location.href = "/auth/login";
      }, 2000);
    } catch (err: any) {
      showNotification(err.message || "Failed to deactivate account", "error");
    }
  };

  const handleSimulateExport = () => {
    const dataStr = JSON.stringify({ settings, profileSettings, cookieConsent }, null, 2);
    const blob = new Blob([dataStr], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `kirmya_user_export_${new Date().toISOString().slice(0, 10)}.json`;
    link.click();
    showNotification("Data exported successfully", "success");
  };

  if (loading) {
    return (
      <div style={{ background: "#FBF7F2", minHeight: "100vh", display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center" }}>
        <CircularProgress style={{ color: "#C2683C" }} />
        <p style={{ marginTop: "16px", color: "#5B554C", fontFamily: "'Public Sans', sans-serif" }}>Loading Settings Center...</p>
      </div>
    );
  }

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
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Settings" }]} />

        <section style={{ maxWidth: "1180px", margin: "0 auto", width: "100%", padding: "clamp(32px,4vw,48px) 40px clamp(16px,2vw,24px)" }}>
          <div style={{ fontSize: "13px", fontWeight: 700, letterSpacing: "0.12em", textTransform: "uppercase", color: "#C2683C", marginBottom: "10px" }}>
            Settings Control Room
          </div>
          <h1 style={{ fontWeight: 800, fontSize: "clamp(30px,4vw,44px)", lineHeight: 1.02, letterSpacing: "-0.025em", margin: 0 }}>
            Account &amp; Platform Settings
          </h1>
          {rateLimitWarning && (
            <div style={{ marginTop: "16px", padding: "12px 18px", background: "rgba(168,71,42,0.08)", color: "#A8472A", border: "1px solid rgba(168,71,42,0.15)", borderRadius: "10px", fontSize: "14px" }}>
              ⚠️ {rateLimitWarning}
            </div>
          )}
        </section>

        <section style={{ maxWidth: "1180px", margin: "0 auto", width: "100%", padding: "0 40px clamp(56px,6vw,90px)", display: "grid", gridTemplateColumns: "300px 1fr", gap: "28px", alignItems: "start", flex: 1 }}>
          
          {/* Navigation Sidebar with search */}
          <aside style={{ position: "sticky", top: "96px", background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "18px", padding: "14px", display: "flex", flexDirection: "column", gap: "10px" }}>
            <div style={{ position: "relative" }}>
              <input
                type="text"
                placeholder="Search settings..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{
                  width: "100%",
                  border: "1px solid #E2D9CC",
                  borderRadius: "10px",
                  padding: "10px 12px",
                  fontSize: "14px",
                  outline: "none",
                  background: "#FCFAF7",
                  fontFamily: "'Public Sans', sans-serif",
                }}
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery("")}
                  style={{ position: "absolute", right: "12px", top: "50%", transform: "translateY(-50%)", border: "none", background: "transparent", cursor: "pointer", color: "#8A8175" }}
                >
                  ×
                </button>
              )}
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: "3px" }}>
              {filteredSections.map((item) => {
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
                      flexDirection: "column",
                      gap: "2px",
                      background: isActive ? "#EFE7DC" : "transparent",
                      transition: "all 0.15s ease",
                    }}
                  >
                    <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                      <span style={{ fontSize: "16px", color: isActive ? "#C2683C" : "#8A8175" }}>{item.icon}</span>
                      <span style={{ fontSize: "14px", fontWeight: isActive ? 600 : 500, color: isActive ? "#2B2620" : "#5B554C" }}>{item.label}</span>
                    </div>
                    <span style={{ fontSize: "11px", color: "#8A8175", marginLeft: "26px" }}>{item.desc}</span>
                  </button>
                );
              })}
              {filteredSections.length === 0 && (
                <p style={{ padding: "14px", color: "#8A8175", fontSize: "14px", textAlign: "center" }}>No results found</p>
              )}
            </div>
          </aside>

          {/* Active Tab Panel */}
          <div style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "clamp(24px,3vw,32px)" }}>
            
            {/* GENERAL SETTINGS */}
            {activeTab === "general" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>General Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Language, localization, and theme configuration.</p>
                </div>

                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px" }}>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>System Language</label>
                    <select
                      value={settings.general.language}
                      onChange={(e) => saveSettingsSegment("general", { ...settings.general, language: e.target.value })}
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                    >
                      <option value="en">English (US)</option>
                      <option value="fr">Français (French)</option>
                      <option value="de">Deutsch (German)</option>
                      <option value="es">Español (Spanish)</option>
                      <option value="ja">日本語 (Japanese)</option>
                    </select>
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Local Time Zone</label>
                    <select
                      value={settings.general.timezone}
                      onChange={(e) => saveSettingsSegment("general", { ...settings.general, timezone: e.target.value })}
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                    >
                      <option value="UTC">UTC</option>
                      <option value="America/New_York">Eastern Time (US/Canada)</option>
                      <option value="America/Denver">Mountain Time (US/Canada)</option>
                      <option value="Europe/London">London / GMT</option>
                      <option value="Asia/Kolkata">IST (Asia/Kolkata)</option>
                    </select>
                  </div>
                </div>

                <div>
                  <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "12px" }}>Select Interface Theme</label>
                  <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: "12px" }}>
                    {[
                      { id: "light", label: "☀ Light Mode" },
                      { id: "dark", label: "☾ Dark Mode" },
                      { id: "system", label: "◐ Follow System" },
                    ].map((t) => (
                      <button
                        key={t.id}
                        type="button"
                        onClick={() => saveSettingsSegment("general", { ...settings.general, theme: t.id })}
                        style={{
                          border: settings.general.theme === t.id ? "2px solid #C2683C" : "1px solid #E2D9CC",
                          background: "#fff",
                          padding: "16px",
                          borderRadius: "12px",
                          cursor: "pointer",
                          fontWeight: 600,
                          fontSize: "14px",
                          color: "#2B2620",
                          textAlign: "center",
                        }}
                      >
                        {t.label}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Email Digest Settings</label>
                  <select
                    value={settings.general.email_digest}
                    onChange={(e) => saveSettingsSegment("general", { ...settings.general, email_digest: e.target.value })}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                  >
                    <option value="off">Off (No digest emails)</option>
                    <option value="daily">Daily summary</option>
                    <option value="weekly">Weekly digest (Recommended)</option>
                  </select>
                </div>
              </div>
            )}

            {/* PROFILE SETTINGS */}
            {activeTab === "profile-pref" && (
              <form onSubmit={handleProfileSettingsSave} style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Profile Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Configure how your professional identity and custom handle look.</p>
                </div>

                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "16px" }}>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Unique Username</label>
                    <input
                      value={usernameInput}
                      onChange={(e) => setUsernameInput(e.target.value)}
                      placeholder="e.g. marcushale"
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                    />
                  </div>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Custom Profile URL</label>
                    <div style={{ display: "flex", alignItems: "center" }}>
                      <span style={{ padding: "12px 10px", border: "1px solid #E2D9CC", borderRight: "none", borderRadius: "10px 0 0 10px", background: "#EFE7DC", fontSize: "14px", color: "#8A8175" }}>kirmya.com/p/</span>
                      <input
                        value={customURLInput}
                        onChange={(e) => setCustomURLInput(e.target.value)}
                        placeholder="marcushale"
                        style={{ flex: 1, minWidth: 0, border: "1px solid #E2D9CC", borderRadius: "0 10px 10px 0", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                      />
                    </div>
                  </div>
                </div>

                <div>
                  <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Profile Privacy Visibility</label>
                  <select
                    value={profileSettings.profile_visibility}
                    onChange={(e) => setProfileSettings({ ...profileSettings, profile_visibility: e.target.value })}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", color: "#2B2620", outline: "none", background: "#FCFAF7" }}
                  >
                    <option value="public">Public (Visible to everyone)</option>
                    <option value="network">Connections only</option>
                    <option value="private">Private (Only you &amp; recruiters)</option>
                  </select>
                </div>

                <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
                  <h4 style={{ margin: "0 0 12px 0", fontSize: "15px", fontWeight: 600 }}>Toggles &amp; Career Flags</h4>
                  <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                    <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                      <input
                        type="checkbox"
                        checked={profileSettings.open_to_work}
                        onChange={(e) => setProfileSettings({ ...profileSettings, open_to_work: e.target.checked })}
                        style={{ width: "18px", height: "18px" }}
                      />
                      <span>Active candidate state (Open to opportunities)</span>
                    </label>
                    <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                      <input
                        type="checkbox"
                        checked={profileSettings.referral_eligible}
                        onChange={(e) => setProfileSettings({ ...profileSettings, referral_eligible: e.target.checked })}
                        style={{ width: "18px", height: "18px" }}
                      />
                      <span>Available for giving referrals</span>
                    </label>
                    <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                      <input
                        type="checkbox"
                        checked={profileSettings.willing_to_mentor}
                        onChange={(e) => setProfileSettings({ ...profileSettings, willing_to_mentor: e.target.checked })}
                        style={{ width: "18px", height: "18px" }}
                      />
                      <span>Willing to mentor peers / students</span>
                    </label>
                  </div>
                </div>

                <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", fontSize: "15px", fontWeight: 600, padding: "13px 28px", borderRadius: "100px", cursor: "pointer", width: "fit-content" }}>
                  Save Profile Settings
                </button>
              </form>
            )}

            {/* ACCOUNT SETTINGS & ACTIONS */}
            {activeTab === "account" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Account Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Control account-level configurations and perform actions.</p>
                </div>

                <div style={{ borderBottom: "1px solid #EFE7DC", paddingBottom: "24px" }}>
                  <h4 style={{ margin: "0 0 10px 0", fontSize: "15px", fontWeight: 600 }}>Connected Account Info</h4>
                  <p style={{ fontSize: "14px", color: "#5B554C" }}>You are signed in with Kirmya. Connect third-party providers under the <strong>Connected Accounts</strong> section on the left.</p>
                </div>

                <div>
                  <h4 style={{ margin: "0 0 14px 0", fontSize: "15px", fontWeight: 600, color: "#A8472A" }}>Danger Zone</h4>
                  <div style={{ display: "flex", gap: "12px", flexWrap: "wrap" }}>
                    <button
                      type="button"
                      onClick={handleSimulateExport}
                      style={{ border: "1px solid #E2D9CC", background: "#fff", color: "#2B2620", fontSize: "14px", fontWeight: 600, padding: "12px 24px", borderRadius: "100px", cursor: "pointer" }}
                    >
                      Export my data (JSON)
                    </button>
                    <button
                      type="button"
                      onClick={() => setShowDeactivateModal(true)}
                      style={{ border: "1px solid #A8472A", background: "transparent", color: "#A8472A", fontSize: "14px", fontWeight: 600, padding: "12px 24px", borderRadius: "100px", cursor: "pointer" }}
                    >
                      Deactivate account
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* PASSWORD & SECURITY */}
            {activeTab === "security" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Security Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Change password, review security alerts and revoke sessions.</p>
                </div>

                <form onSubmit={handlePasswordChange} style={{ borderBottom: "1px solid #EFE7DC", paddingBottom: "24px", display: "flex", flexDirection: "column", gap: "16px" }}>
                  <h4 style={{ margin: 0, fontSize: "16px", fontWeight: 600 }}>Update Password</h4>
                  <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))", gap: "16px" }}>
                    <div>
                      <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Current Password</label>
                      <input
                        type="password"
                        value={currentPassword}
                        onChange={(e) => setCurrentPassword(e.target.value)}
                        required
                        style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
                      />
                    </div>
                    <div>
                      <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>New Password</label>
                      <input
                        type="password"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        required
                        style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
                      />
                    </div>
                    <div>
                      <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Confirm New Password</label>
                      <input
                        type="password"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                        style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
                      />
                    </div>
                  </div>
                  <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", fontSize: "14px", fontWeight: 600, padding: "10px 20px", borderRadius: "100px", cursor: "pointer", width: "fit-content" }}>
                    Save Password
                  </button>
                </form>

                <div style={{ borderBottom: "1px solid #EFE7DC", paddingBottom: "24px" }}>
                  <h4 style={{ margin: "0 0 12px 0", fontSize: "16px", fontWeight: 600 }}>Active Logged-In Sessions</h4>
                  <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                    {securityActivity.sessions.map((s: any) => (
                      <div key={s.id} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", border: "1px solid #EFE7DC", borderRadius: "12px", padding: "14px", background: "#FCFAF7" }}>
                        <div>
                          <div style={{ fontSize: "14px", fontWeight: 600 }}>{s.user_agent}</div>
                          <div style={{ fontSize: "12px", color: "#8A8175" }}>IP: {s.ip_address} • Started: {new Date(s.created_at).toLocaleDateString()}</div>
                        </div>
                        <button
                          type="button"
                          onClick={() => handleRevokeSession(s.id)}
                          style={{ border: "none", background: "transparent", color: "#A8472A", fontSize: "13px", fontWeight: 600, cursor: "pointer" }}
                        >
                          Revoke Session
                        </button>
                      </div>
                    ))}
                    {securityActivity.sessions.length === 0 && <p style={{ color: "#8A8175", fontSize: "14px" }}>No active sessions.</p>}
                  </div>
                </div>

                <div>
                  <h4 style={{ margin: "0 0 12px 0", fontSize: "16px", fontWeight: 600 }}>Security Logs &amp; Activity History</h4>
                  <div style={{ display: "flex", flexDirection: "column", gap: "8px", maxHeight: "200px", overflowY: "auto" }}>
                    {securityActivity.history.map((h: any) => (
                      <div key={h.id} style={{ display: "flex", justifyContent: "space-between", padding: "8px 0", borderBottom: "1px solid #F3E7DC", fontSize: "13px" }}>
                        <span style={{ fontWeight: 600, color: "#2B2620" }}>{h.action}</span>
                        <span style={{ color: "#8A8175" }}>{h.ip_address} • {new Date(h.created_at).toLocaleString()}</span>
                      </div>
                    ))}
                    {securityActivity.history.length === 0 && <p style={{ color: "#8A8175", fontSize: "14px" }}>No history events recorded.</p>}
                  </div>
                </div>
              </div>
            )}

            {/* PRIVACY SETTINGS */}
            {activeTab === "privacy" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Privacy Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Configure search, messaging discoverability, and AI data settings.</p>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.privacy.show_email}
                      onChange={(e) => saveSettingsSegment("privacy", { ...settings.privacy, show_email: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Expose my email address on my public profile page</span>
                  </label>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.privacy.discoverable}
                      onChange={(e) => saveSettingsSegment("privacy", { ...settings.privacy, discoverable: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Allow search engines (Google, Bing) to index my profile</span>
                  </label>
                  <div>
                    <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Allow Direct Messages From</label>
                    <select
                      value={settings.privacy.allow_messages}
                      onChange={(e) => saveSettingsSegment("privacy", { ...settings.privacy, allow_messages: e.target.value })}
                      style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
                    >
                      <option value="everyone">Everyone</option>
                      <option value="network">Connections only</option>
                      <option value="none">No one</option>
                    </select>
                  </div>
                </div>
              </div>
            )}

            {/* NOTIFICATION SETTINGS */}
            {activeTab === "notifications" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Notification Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Configure channels and alerts for Jobs, Mentorship, Messaging, and Referrals.</p>
                </div>

                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "24px" }}>
                  <div>
                    <h4 style={{ margin: "0 0 14px 0", fontSize: "15px", fontWeight: 600 }}>Email Alerts</h4>
                    <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.email_jobs}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, email_jobs: e.target.checked })}
                        />
                        <span>Jobs matching preferences</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.email_mentorship}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, email_mentorship: e.target.checked })}
                        />
                        <span>Mentorship invitations</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.email_messages}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, email_messages: e.target.checked })}
                        />
                        <span>New inbox messages</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.email_referrals}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, email_referrals: e.target.checked })}
                        />
                        <span>Referral requests</span>
                      </label>
                    </div>
                  </div>

                  <div>
                    <h4 style={{ margin: "0 0 14px 0", fontSize: "15px", fontWeight: 600 }}>In-App Notifications</h4>
                    <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.inapp_jobs}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, inapp_jobs: e.target.checked })}
                        />
                        <span>Jobs matching preferences</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.inapp_mentorship}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, inapp_mentorship: e.target.checked })}
                        />
                        <span>Mentorship invitations</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.inapp_messages}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, inapp_messages: e.target.checked })}
                        />
                        <span>New inbox messages</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.notifications.inapp_referrals}
                          onChange={(e) => saveSettingsSegment("notifications", { ...settings.notifications, inapp_referrals: e.target.checked })}
                        />
                        <span>Referral requests</span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* AI PREFERENCES */}
            {activeTab === "ai" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>AI Preferences</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Configure AI Assistant features and suggestion automation.</p>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.ai.enable_ai_assistant}
                      onChange={(e) => saveSettingsSegment("ai", { ...settings.ai, enable_ai_assistant: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span style={{ fontWeight: 600 }}>Enable Kirmya AI Assistant globally</span>
                  </label>

                  <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
                    <h4 style={{ margin: "0 0 12px 0", fontSize: "15px", fontWeight: 600 }}>Enable AI Suggestions For:</h4>
                    <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.ai.ai_job_recommendations}
                          onChange={(e) => saveSettingsSegment("ai", { ...settings.ai, ai_job_recommendations: e.target.checked })}
                        />
                        <span>Job Matching Recommendations</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.ai.ai_resume_suggestions}
                          onChange={(e) => saveSettingsSegment("ai", { ...settings.ai, ai_resume_suggestions: e.target.checked })}
                        />
                        <span>Resume Enhancement &amp; Keywording suggestions</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.ai.ai_roadmap_suggestions}
                          onChange={(e) => saveSettingsSegment("ai", { ...settings.ai, ai_roadmap_suggestions: e.target.checked })}
                        />
                        <span>Interactive Career Roadmaps &amp; Courses</span>
                      </label>
                      <label style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                        <input
                          type="checkbox"
                          checked={settings.ai.ai_skill_gap_analysis}
                          onChange={(e) => saveSettingsSegment("ai", { ...settings.ai, ai_skill_gap_analysis: e.target.checked })}
                        />
                        <span>Skill Gap Assessments</span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* JOB PREFERENCES */}
            {activeTab === "job-pref" && (
              <form onSubmit={handleProfileSettingsSave} style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Job Preferences</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Define career expectations and visibility flags.</p>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={profileSettings.open_to_work}
                      onChange={(e) => setProfileSettings({ ...profileSettings, open_to_work: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Open to job opportunities (displays green badge to recruiters)</span>
                  </label>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={profileSettings.willing_to_mentor}
                      onChange={(e) => setProfileSettings({ ...profileSettings, willing_to_mentor: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Willing to accept mentorship connections</span>
                  </label>
                </div>

                <button type="submit" style={{ border: "none", background: "#C2683C", color: "#fff", fontSize: "15px", fontWeight: 600, padding: "13px 28px", borderRadius: "100px", cursor: "pointer", width: "fit-content" }}>
                  Save Job Preferences
                </button>
              </form>
            )}

            {/* LEARNING PREFERENCES */}
            {activeTab === "learning" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Learning Preferences</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Establish skills study goals and reminders.</p>
                </div>

                <div>
                  <h4 style={{ margin: "0 0 10px 0", fontSize: "15px", fontWeight: 600 }}>Active Goals</h4>
                  <div style={{ display: "flex", flexWrap: "wrap", gap: "8px", marginBottom: "12px" }}>
                    {settings.learning.learning_goals.map((g: string) => (
                      <span key={g} style={{ background: "#EFE7DC", color: "#2B2620", fontSize: "13px", padding: "6px 12px", borderRadius: "100px", display: "inline-flex", alignItems: "center", gap: "6px" }}>
                        {g}
                        <button
                          type="button"
                          onClick={() => saveSettingsSegment("learning", { ...settings.learning, learning_goals: settings.learning.learning_goals.filter((x: string) => x !== g) })}
                          style={{ border: "none", background: "transparent", color: "#A8472A", cursor: "pointer", padding: 0 }}
                        >
                          ×
                        </button>
                      </span>
                    ))}
                    {settings.learning.learning_goals.length === 0 && <span style={{ color: "#8A8175", fontSize: "13px" }}>No learning goals set.</span>}
                  </div>
                  <div style={{ display: "flex", gap: "10px" }}>
                    <input
                      value={newGoal}
                      onChange={(e) => setNewGoal(e.target.value)}
                      placeholder="Add learning goal..."
                      style={{ flex: 1, border: "1px solid #E2D9CC", borderRadius: "10px", padding: "8px 12px", fontSize: "14px", outline: "none" }}
                    />
                    <button
                      type="button"
                      onClick={() => {
                        if (!newGoal) return;
                        saveSettingsSegment("learning", { ...settings.learning, learning_goals: [...settings.learning.learning_goals, newGoal] });
                        setNewGoal("");
                      }}
                      style={{ border: "none", background: "#C2683C", color: "#fff", padding: "8px 16px", borderRadius: "10px", cursor: "pointer", fontSize: "14px" }}
                    >
                      Add
                    </button>
                  </div>
                </div>

                <div style={{ borderTop: "1px solid #EFE7DC", paddingTop: "20px" }}>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.learning.learning_reminders}
                      onChange={(e) => saveSettingsSegment("learning", { ...settings.learning, learning_reminders: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Receive learning &amp; skills upskilling reminders</span>
                  </label>
                </div>
              </div>
            )}

            {/* CONNECTED ACCOUNTS */}
            {activeTab === "connected" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Connected Accounts</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Link third-party platforms to simplify login and profile syncing.</p>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                  {[
                    { key: "google", label: "Google Account", color: "#4285F4" },
                    { key: "linkedin", label: "LinkedIn Account", color: "#0A66C2" },
                    { key: "github", label: "GitHub Account", color: "#24292E" },
                    { key: "microsoft", label: "Microsoft Account", color: "#00A4EF" },
                    { key: "apple", label: "Apple ID", color: "#000000" },
                  ].map((prov) => {
                    const linked = securityActivity.connected_accounts.some((a: any) => a.provider === prov.key);
                    return (
                      <div key={prov.key} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7" }}>
                        <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                          <span style={{ width: "8px", height: "8px", borderRadius: "50%", background: linked ? "#4F7C6A" : "#8A8175" }} />
                          <span style={{ fontWeight: 600, color: "#2B2620" }}>{prov.label}</span>
                        </div>
                        {linked ? (
                          <button
                            type="button"
                            onClick={() => handleDisconnectConnectedAccount(prov.key)}
                            style={{ border: "1px solid #A8472A", background: "transparent", color: "#A8472A", fontSize: "13px", fontWeight: 600, padding: "6px 12px", borderRadius: "100px", cursor: "pointer" }}
                          >
                            Disconnect
                          </button>
                        ) : (
                          <button
                            type="button"
                            onClick={() => handleSimulateOAuthConnect(prov.key)}
                            style={{ border: "1px solid #C2683C", background: "transparent", color: "#C2683C", fontSize: "13px", fontWeight: 600, padding: "6px 12px", borderRadius: "100px", cursor: "pointer" }}
                          >
                            Connect
                          </button>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            )}

            {/* PRIVACY & COOKIES CENTER */}
            {activeTab === "cookies" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Privacy &amp; Cookies</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Control tracking consents and cookies stored on your device.</p>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>Essential Cookies</h4>
                      <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Required for sign-ins, security, and CSRF protection. Cannot be disabled.</p>
                    </div>
                    <span style={{ color: "#4F7C6A", fontSize: "13px", fontWeight: 600 }}>Always Active</span>
                  </div>

                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>Functional Cookies</h4>
                      <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Saves language selection, custom themes, and landing page choice.</p>
                    </div>
                    <input
                      type="checkbox"
                      checked={cookieConsent.functional}
                      onChange={(e) => handleCookieConsentSave(e.target.checked, cookieConsent.analytics, cookieConsent.ai_personalization)}
                      style={{ width: "36px", height: "20px", cursor: "pointer" }}
                    />
                  </div>

                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>Analytics &amp; Tracking Cookies</h4>
                      <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Anonymously logs page views, job searches, and interaction metrics.</p>
                    </div>
                    <input
                      type="checkbox"
                      checked={cookieConsent.analytics}
                      onChange={(e) => handleCookieConsentSave(cookieConsent.functional, e.target.checked, cookieConsent.ai_personalization)}
                      style={{ width: "36px", height: "20px", cursor: "pointer" }}
                    />
                  </div>

                  <div style={{ border: "1px solid #EFE7DC", borderRadius: "12px", padding: "16px", background: "#FCFAF7", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div>
                      <h4 style={{ margin: "0 0 4px", fontSize: "15px", fontWeight: 600 }}>AI Personalization Cookies</h4>
                      <p style={{ margin: 0, fontSize: "12px", color: "#8A8175" }}>Saves context of coach chats and resume evaluations locally to optimize LLM calls.</p>
                    </div>
                    <input
                      type="checkbox"
                      checked={cookieConsent.ai_personalization}
                      onChange={(e) => handleCookieConsentSave(cookieConsent.functional, cookieConsent.analytics, e.target.checked)}
                      style={{ width: "36px", height: "20px", cursor: "pointer" }}
                    />
                  </div>
                </div>
              </div>
            )}

            {/* ACCESSIBILITY SETTINGS */}
            {activeTab === "accessibility" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                <div>
                  <h2 style={{ fontWeight: 700, fontSize: "22px", margin: "0 0 4px" }}>Accessibility Settings</h2>
                  <p style={{ fontSize: "15px", color: "#8A8175", margin: 0 }}>Configure interface accessibility rules for display and screen readers.</p>
                </div>

                <div>
                  <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "7px" }}>Interface Font Size</label>
                  <select
                    value={settings.accessibility.font_size}
                    onChange={(e) => saveSettingsSegment("accessibility", { ...settings.accessibility, font_size: e.target.value })}
                    style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
                  >
                    <option value="small">Small text</option>
                    <option value="medium">Medium text (Default)</option>
                    <option value="large">Large text</option>
                    <option value="extra-large">Extra large text</option>
                  </select>
                </div>

                <div style={{ display: "flex", flexDirection: "column", gap: "12px" }}>
                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.accessibility.high_contrast}
                      onChange={(e) => saveSettingsSegment("accessibility", { ...settings.accessibility, high_contrast: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>High Contrast Mode (Enhances line borders and shadows)</span>
                  </label>

                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.accessibility.reduced_motion}
                      onChange={(e) => saveSettingsSegment("accessibility", { ...settings.accessibility, reduced_motion: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Reduced Motion (Disables sliders and micro-animations)</span>
                  </label>

                  <label style={{ display: "flex", alignItems: "center", gap: "10px", cursor: "pointer" }}>
                    <input
                      type="checkbox"
                      checked={settings.accessibility.compact_mode}
                      onChange={(e) => saveSettingsSegment("accessibility", { ...settings.accessibility, compact_mode: e.target.checked })}
                      style={{ width: "18px", height: "18px" }}
                    />
                    <span>Compact layout density (Reduces padding heights)</span>
                  </label>
                </div>
              </div>
            )}

            {/* PLATFORM INFORMATION */}
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
                  <h3 style={{ margin: "0 0 12px 0", fontSize: "17px", fontWeight: 700 }}>Release Notes — What&apos;s New</h3>
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

          </div>

        </section>

        {/* DEACTIVATE ACCOUNT CONFIRMATION MODAL */}
        {showDeactivateModal && (
          <div style={{ position: "fixed", top: 0, left: 0, right: 0, bottom: 0, background: "rgba(0,0,0,0.5)", display: "flex", justifyContent: "center", alignItems: "center", zIndex: 1000, padding: "20px" }}>
            <div style={{ background: "#fff", border: "1px solid #EFE7DC", borderRadius: "18px", maxWidth: "480px", width: "100%", padding: "24px", display: "flex", flexDirection: "column", gap: "18px" }}>
              <h3 style={{ margin: 0, color: "#A8472A", fontWeight: 700, fontSize: "20px" }}>Warning: Deactivate Account</h3>
              <p style={{ margin: 0, fontSize: "14px", color: "#5B554C", lineHeight: 1.5 }}>
                Deactivating your account will hide your profile from all candidates and recruiters. Your search discoverability is suspended instantly. You will be logged out of all connected sessions. You can reactivate your account at any time by logging back in.
              </p>
              <label style={{ display: "flex", alignItems: "flex-start", gap: "10px", cursor: "pointer", fontSize: "13px", color: "#2B2620" }}>
                <input
                  type="checkbox"
                  checked={deactivateCheck}
                  onChange={(e) => setDeactivateCheck(e.target.checked)}
                  style={{ marginTop: "3px" }}
                />
                <span>I understand the consequences and wish to deactivate my account.</span>
              </label>
              <div style={{ display: "flex", justifyContent: "flex-end", gap: "12px" }}>
                <button
                  type="button"
                  onClick={() => { setShowDeactivateModal(false); setDeactivateCheck(false); }}
                  style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "10px 20px", borderRadius: "100px", cursor: "pointer", fontSize: "14px", fontWeight: 600 }}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  onClick={handleDeactivate}
                  disabled={!deactivateCheck}
                  style={{ border: "none", background: deactivateCheck ? "#A8472A" : "#E2D9CC", color: "#fff", padding: "10px 20px", borderRadius: "100px", cursor: deactivateCheck ? "pointer" : "not-allowed", fontSize: "14px", fontWeight: 600 }}
                >
                  Deactivate
                </button>
              </div>
            </div>
          </div>
        )}

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
