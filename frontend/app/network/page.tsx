"use client";

import React, { useState, useEffect } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import AuthGuard from "@/components/shared/AuthGuard";
import { api } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import { CircularProgress } from "@mui/material";
import { Search, Filter, SlidersHorizontal, Check, UserMinus, Ban, Send, MessageSquare } from "lucide-react";

interface ConnectionUser {
  id: string;
  requester_id: string;
  receiver_id: string;
  status: string;
  origin: string;
  created_at: string;
  requester_name?: string;
  requester_headline?: string;
  requester_photo_url?: string;
  receiver_name?: string;
  receiver_headline?: string;
  receiver_photo_url?: string;
}

interface Professional {
  id: string;
  full_name: string;
  email: string;
  headline: string;
  photo_url?: string;
  location?: string;
  skills?: string[];
  connectionStatus?: string;
  connectionId?: string;
}

export default function NetworkPage({ initialTab }: { initialTab?: "discover" | "connections" | "requests" }) {
  const { user } = useAuth();
  const { showNotification } = useNotifications();

  const [activeTab, setActiveTab] = useState<"discover" | "connections" | "requests">(initialTab || "discover");
  const [loading, setLoading] = useState(true);

  // Lists
  const [connections, setConnections] = useState<ConnectionUser[]>([]);
  const [incomingRequests, setIncomingRequests] = useState<ConnectionUser[]>([]);
  const [professionals, setProfessionals] = useState<Professional[]>([]);

  // Search & Filter state
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedLocation, setSelectedLocation] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("");
  const [showFilters, setShowFilters] = useState(false);
  const [sortBy, setSortBy] = useState("relevance");

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      // 1. Fetch connections
      const connData = await api.get<ConnectionUser[]>("/network/connections");
      setConnections(connData || []);

      // 2. Fetch requests
      const reqData = await api.get<ConnectionUser[]>("/network/requests/incoming");
      setIncomingRequests(reqData || []);

      // 3. Fetch initial professionals discovery list
      const hits = await api.get<{ results: any[] }>("/search?q=a&type=user");
      const mapped: Professional[] = (hits.results || [])
        .filter((h) => h.id !== user?.id)
        .map((h) => ({
          id: h.id,
          full_name: h.title,
          email: h.description,
          headline: h.snippet || "Experienced Professional",
          location: "San Francisco, CA",
          skills: ["Engineering", "Strategy"],
        }));
      setProfessionals(mapped);
    } catch (err: any) {
      showNotification(err.message || "Failed to load networking directory", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) {
      loadData();
      return;
    }
    setLoading(true);
    try {
      const hits = await api.get<{ results: any[] }>(`/search?q=${encodeURIComponent(searchQuery)}&type=user`);
      const mapped: Professional[] = (hits.results || [])
        .filter((h) => h.id !== user?.id)
        .map((h) => ({
          id: h.id,
          full_name: h.title,
          email: h.description,
          headline: h.snippet || "Experienced Professional",
          location: "San Francisco, CA",
          skills: ["Engineering", "Strategy"],
        }));
      setProfessionals(mapped);
    } catch (err: any) {
      showNotification("Search failed", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleConnect = async (receiverId: string) => {
    try {
      await api.post("/network/requests", {
        receiver_id: receiverId,
        origin: "manual_request",
      });
      showNotification("Connection request sent!", "success");
      // Update professional local status
      setProfessionals((prev) =>
        prev.map((p) => (p.id === receiverId ? { ...p, connectionStatus: "pending" } : p))
      );
    } catch (err: any) {
      showNotification(err.message || "Failed to send connection request", "error");
    }
  };

  const handleAccept = async (requestId: string) => {
    try {
      await api.put(`/network/requests/${requestId}/accept`, {});
      showNotification("Connection request accepted!", "success");
      // Refresh Lists
      loadData();
    } catch (err: any) {
      showNotification("Failed to accept request", "error");
    }
  };

  const handleReject = async (requestId: string) => {
    try {
      await api.put(`/network/requests/${requestId}/reject`, {});
      showNotification("Request declined.", "success");
      // Refresh Lists
      loadData();
    } catch (err: any) {
      showNotification("Failed to decline request", "error");
    }
  };

  const handleUnconnect = async (targetUserId: string) => {
    if (!window.confirm("Are you sure you want to remove this connection?")) return;
    try {
      await api.delete(`/network/connections/${targetUserId}`);
      showNotification("Connection removed.", "success");
      loadData();
    } catch (err: any) {
      showNotification("Failed to remove connection", "error");
    }
  };

  const handleBlock = async (targetUserId: string) => {
    if (!window.confirm("Are you sure you want to block this professional?")) return;
    try {
      await api.post("/network/block", {
        blocked_id: targetUserId,
      });
      showNotification("Professional blocked.", "success");
      loadData();
    } catch (err: any) {
      showNotification("Failed to block professional", "error");
    }
  };

  const handleStartChat = async (targetUserId: string, targetName: string) => {
    try {
      const conv = await api.post<any>("/conversations", {
        participant_ids: [targetUserId],
        title: targetName,
      });
      window.location.href = `/inbox?convId=${conv.id}`;
    } catch (err: any) {
      showNotification("Failed to start conversation", "error");
    }
  };

  return (
    <AuthGuard>
      <div
        style={{
          background: "#FBF7F2",
          fontFamily: "'Public Sans', sans-serif",
          color: "#2B2620",
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Network Center" }]} />

        {/* Hero Section */}
        <section style={{ maxWidth: "920px", margin: "0 auto", padding: "clamp(56px,7vw,100px) 40px clamp(40px,5vw,56px)", textAlign: "center" }}>
          <div style={{ display: "inline-block", fontSize: "13px", fontWeight: 700, letterSpacing: "0.1em", textTransform: "uppercase", color: "#C2683C", background: "rgba(194,104,60,0.12)", padding: "8px 16px", borderRadius: "100px", marginBottom: "26px" }}>
            Professional Networking
          </div>
          <h1 style={{ fontWeight: 800, fontSize: "clamp(40px,6.5vw,72px)", lineHeight: 1.02, letterSpacing: "-0.025em", margin: "0 auto 22px", maxWidth: "760px" }}>
            Connect with verified practitioners.
          </h1>
          <p style={{ fontSize: "clamp(17px,2vw,20px)", lineHeight: 1.6, color: "#5B554C", maxWidth: "600px", margin: "0 auto 34px" }}>
            Grow your professional network, share insider job openings, schedule mock interviews, and trade referrals without premium restrictions.
          </p>
        </section>

        {/* Navigation Tabs bar */}
        <section style={{ maxWidth: "1240px", margin: "0 auto", width: "100%", padding: "0 40px", borderBottom: "1px solid #EFE7DC" }}>
          <div style={{ display: "flex", gap: "24px" }}>
            {[
              { id: "discover", label: `🔍 Discover People` },
              { id: "connections", label: `👥 My Connections (${connections.length})` },
              { id: "requests", label: `📥 Connection Requests (${incomingRequests.length})` },
            ].map((t) => (
              <button
                key={t.id}
                onClick={() => setActiveTab(t.id as any)}
                style={{
                  border: "none",
                  background: "transparent",
                  padding: "16px 8px",
                  fontSize: "15px",
                  fontWeight: activeTab === t.id ? 700 : 500,
                  color: activeTab === t.id ? "#C2683C" : "#5B554C",
                  borderBottom: activeTab === t.id ? "3px solid #C2683C" : "none",
                  cursor: "pointer",
                }}
              >
                {t.label}
              </button>
            ))}
          </div>
        </section>

        {/* Main Content Area */}
        <section style={{ maxWidth: "1240px", margin: "32px auto", width: "100%", padding: "0 40px clamp(48px,6vw,72px)", flex: 1 }}>
          
          {loading ? (
            <div style={{ display: "flex", justifyContent: "center", alignItems: "center", padding: "64px" }}>
              <CircularProgress style={{ color: "#C2683C" }} />
            </div>
          ) : (
            <>
              {/* DISCOVER TAB */}
              {activeTab === "discover" && (
                <div>
                  {/* Search Bar & Filters */}
                  <form onSubmit={handleSearch} style={{ background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "18px 24px", display: "flex", flexDirection: "column", gap: "16px", marginBottom: "32px", boxShadow: "0 4px 12px rgba(43, 38, 32, 0.02)" }}>
                    <div style={{ display: "flex", gap: "16px", flexWrap: "wrap", alignItems: "center" }}>
                      <div style={{ flex: 1, minWidth: "260px", display: "flex", alignItems: "center", gap: "10px", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px 14px", background: "#FCFAF7" }}>
                        <Search size={18} color="#8A8175" />
                        <input
                          type="text"
                          placeholder="Search professionals by name, title, company, or skills..."
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          style={{ border: "none", outline: "none", background: "transparent", fontSize: "15px", color: "#2B2620", width: "100%", fontFamily: "inherit" }}
                        />
                      </div>
                      <button
                        type="button"
                        onClick={() => setShowFilters(!showFilters)}
                        style={{ display: "flex", alignItems: "center", gap: "8px", border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "12px 20px", borderRadius: "10px", fontWeight: 600, cursor: "pointer", fontSize: "14px" }}
                      >
                        <SlidersHorizontal size={16} /> Filters
                      </button>
                      <button
                        type="submit"
                        style={{ border: "none", background: "#C2683C", color: "#fff", padding: "12px 28px", borderRadius: "10px", fontWeight: 600, cursor: "pointer", fontSize: "14px" }}
                      >
                        Search
                      </button>
                    </div>

                    {showFilters && (
                      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))", gap: "16px", borderTop: "1px solid #F6EFE6", paddingTop: "16px" }}>
                        <div>
                          <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Location</label>
                          <input
                            type="text"
                            placeholder="e.g. San Francisco, CA"
                            value={selectedLocation}
                            onChange={(e) => setSelectedLocation(e.target.value)}
                            style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", background: "#FCFAF7", outline: "none" }}
                          />
                        </div>
                        <div>
                          <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Category</label>
                          <select
                            value={selectedCategory}
                            onChange={(e) => setSelectedCategory(e.target.value)}
                            style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", background: "#FCFAF7", outline: "none" }}
                          >
                            <option value="">All Categories</option>
                            <option value="tech">Technology</option>
                            <option value="recruiter">Recruiters</option>
                            <option value="mentor">Mentors</option>
                          </select>
                        </div>
                        <div>
                          <label style={{ display: "block", fontSize: "12px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Sort By</label>
                          <select
                            value={sortBy}
                            onChange={(e) => setSortBy(e.target.value)}
                            style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "8px", padding: "10px", fontSize: "14px", background: "#FCFAF7", outline: "none" }}
                          >
                            <option value="relevance">Relevance</option>
                            <option value="connections">Most Connected</option>
                            <option value="recent">Recently Joined</option>
                          </select>
                        </div>
                      </div>
                    )}
                  </form>

                  {/* Grid layout for discovery */}
                  <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))", gap: "24px" }}>
                    {professionals.map((prof) => (
                      <div
                        key={prof.id}
                        style={{
                          background: "#ffffff",
                          border: "1px solid #EFE7DC",
                          borderRadius: "20px",
                          padding: "24px",
                          display: "flex",
                          flexDirection: "column",
                          justifyContent: "space-between",
                          gap: "18px",
                          boxShadow: "0 4px 12px rgba(43, 38, 32, 0.03)",
                        }}
                      >
                        <div>
                          <div style={{ display: "flex", gap: "14px", alignItems: "center", marginBottom: "12px" }}>
                            <div style={{ width: "52px", height: "52px", borderRadius: "50%", background: "#4F7C6A", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "18px", fontWeight: 700 }}>
                              {prof.full_name.charAt(0)}
                            </div>
                            <div>
                              <h3 style={{ fontSize: "16px", fontWeight: 700, margin: 0, color: "#2B2620" }}>{prof.full_name}</h3>
                              <p style={{ fontSize: "13px", color: "#8A8175", margin: 0 }}>{prof.headline}</p>
                            </div>
                          </div>
                          <div style={{ fontSize: "13px", color: "#5B554C" }}>📍 {prof.location}</div>
                        </div>

                        <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "14px", display: "flex", gap: "10px" }}>
                          <button
                            onClick={() => handleConnect(prof.id)}
                            disabled={prof.connectionStatus === "pending"}
                            style={{
                              flex: 1,
                              border: "none",
                              background: prof.connectionStatus === "pending" ? "rgba(43,38,32,0.08)" : "#C2683C",
                              color: prof.connectionStatus === "pending" ? "#8A8175" : "#fff",
                              padding: "10px 16px",
                              borderRadius: "100px",
                              fontWeight: 600,
                              fontSize: "13px",
                              cursor: prof.connectionStatus === "pending" ? "not-allowed" : "pointer",
                            }}
                          >
                            {prof.connectionStatus === "pending" ? "Pending..." : "Connect"}
                          </button>
                          <button
                            onClick={() => handleStartChat(prof.id, prof.full_name)}
                            style={{
                              border: "1px solid #E2D9CC",
                              background: "transparent",
                              color: "#C2683C",
                              padding: "10px 16px",
                              borderRadius: "100px",
                              fontWeight: 600,
                              fontSize: "13px",
                              cursor: "pointer",
                            }}
                          >
                            Message
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* CONNECTIONS TAB */}
              {activeTab === "connections" && (
                <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))", gap: "24px" }}>
                  {connections.map((c) => {
                    const isReq = c.requester_id === user?.id;
                    const partnerName = isReq ? c.receiver_name : c.requester_name;
                    const partnerHeadline = isReq ? c.receiver_headline : c.requester_headline;
                    const partnerId = isReq ? c.receiver_id : c.requester_id;

                    return (
                      <div
                        key={c.id}
                        style={{
                          background: "#ffffff",
                          border: "1px solid #EFE7DC",
                          borderRadius: "20px",
                          padding: "24px",
                          display: "flex",
                          flexDirection: "column",
                          justifyContent: "space-between",
                          gap: "18px",
                          boxShadow: "0 4px 12px rgba(43, 38, 32, 0.03)",
                        }}
                      >
                        <div>
                          <div style={{ display: "flex", gap: "14px", alignItems: "center", marginBottom: "12px" }}>
                            <div style={{ width: "52px", height: "52px", borderRadius: "50%", background: "#C2683C", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "18px", fontWeight: 700 }}>
                              {partnerName?.charAt(0) || "P"}
                            </div>
                            <div>
                              <h3 style={{ fontSize: "16px", fontWeight: 700, margin: 0, color: "#2B2620" }}>{partnerName}</h3>
                              <p style={{ fontSize: "13px", color: "#8A8175", margin: 0 }}>{partnerHeadline}</p>
                            </div>
                          </div>
                        </div>

                        <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "14px", display: "flex", gap: "10px" }}>
                          <button
                            onClick={() => handleStartChat(partnerId, partnerName || "Partner")}
                            style={{
                              flex: 1,
                              border: "none",
                              background: "#C2683C",
                              color: "#fff",
                              padding: "10px 16px",
                              borderRadius: "100px",
                              fontWeight: 600,
                              fontSize: "13px",
                              cursor: "pointer",
                            }}
                          >
                            Chat
                          </button>
                          <button
                            onClick={() => handleUnconnect(partnerId)}
                            style={{
                              border: "1px solid #E2D9CC",
                              background: "transparent",
                              color: "#8A8175",
                              padding: "10px 12px",
                              borderRadius: "100px",
                              cursor: "pointer",
                            }}
                            title="Remove Connection"
                          >
                            <UserMinus size={16} />
                          </button>
                          <button
                            onClick={() => handleBlock(partnerId)}
                            style={{
                              border: "1px solid #E2D9CC",
                              background: "transparent",
                              color: "#A8472A",
                              padding: "10px 12px",
                              borderRadius: "100px",
                              cursor: "pointer",
                            }}
                            title="Block User"
                          >
                            <Ban size={16} />
                          </button>
                        </div>
                      </div>
                    );
                  })}
                  {connections.length === 0 && (
                    <div style={{ gridColumn: "1/-1", textAlign: "center", padding: "64px 24px", color: "#8A8175" }}>
                      <span style={{ fontSize: "40px" }}>👥</span>
                      <h3 style={{ margin: "16px 0 6px", fontSize: "18px", fontWeight: 700, color: "#2B2620" }}>No connections yet</h3>
                      <p style={{ margin: 0, fontSize: "14px" }}>Start discovering other professionals and build your network.</p>
                    </div>
                  )}
                </div>
              )}

              {/* REQUESTS TAB */}
              {activeTab === "requests" && (
                <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
                  {incomingRequests.map((req) => (
                    <div
                      key={req.id}
                      style={{
                        background: "#ffffff",
                        border: "1px solid #EFE7DC",
                        borderRadius: "18px",
                        padding: "20px 24px",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "space-between",
                        flexWrap: "wrap",
                        gap: "16px",
                        boxShadow: "0 4px 12px rgba(43, 38, 32, 0.02)",
                      }}
                    >
                      <div style={{ display: "flex", gap: "14px", alignItems: "center" }}>
                        <div style={{ width: "48px", height: "48px", borderRadius: "50%", background: "#4F7C6A", color: "#fff", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "16px", fontWeight: 700 }}>
                          {req.requester_name?.charAt(0) || "R"}
                        </div>
                        <div>
                          <h3 style={{ fontSize: "16px", fontWeight: 700, margin: 0, color: "#2B2620" }}>{req.requester_name}</h3>
                          <p style={{ fontSize: "13px", color: "#8A8175", margin: 0 }}>{req.requester_headline}</p>
                        </div>
                      </div>

                      <div style={{ display: "flex", gap: "10px" }}>
                        <button
                          onClick={() => handleAccept(req.id)}
                          style={{
                            border: "none",
                            background: "#C2683C",
                            color: "#fff",
                            padding: "10px 20px",
                            borderRadius: "100px",
                            fontWeight: 600,
                            fontSize: "13px",
                            cursor: "pointer",
                          }}
                        >
                          Accept
                        </button>
                        <button
                          onClick={() => handleReject(req.id)}
                          style={{
                            border: "1px solid #E2D9CC",
                            background: "transparent",
                            color: "#5B554C",
                            padding: "10px 20px",
                            borderRadius: "100px",
                            fontWeight: 600,
                            fontSize: "13px",
                            cursor: "pointer",
                          }}
                        >
                          Decline
                        </button>
                      </div>
                    </div>
                  ))}
                  {incomingRequests.length === 0 && (
                    <div style={{ textAlign: "center", padding: "64px 24px", color: "#8A8175" }}>
                      <span style={{ fontSize: "40px" }}>📥</span>
                      <h3 style={{ margin: "16px 0 6px", fontSize: "18px", fontWeight: 700, color: "#2B2620" }}>No pending requests</h3>
                      <p style={{ margin: 0, fontSize: "14px" }}>Incoming connection invites will appear here.</p>
                    </div>
                  )}
                </div>
              )}
            </>
          )}

        </section>

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
