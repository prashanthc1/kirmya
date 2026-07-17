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
import ConnectionsList from "@/components/connections/ConnectionsList";
import PendingRequestsPanel from "@/components/connections/PendingRequestsPanel";
import SuggestionsCarousel from "@/components/connections/SuggestionsCarousel";
import ConnectButton from "@/components/connections/ConnectButton";

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
                  <div className="mb-8">
                    <SuggestionsCarousel />
                  </div>
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

                        <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "14px", display: "flex", gap: "10px", justifyContent: "flex-end", alignItems: "center" }}>
                          <ConnectButton
                            targetUserId={prof.id}
                            targetUserName={prof.full_name}
                            currentConnectionStatus={prof.connectionStatus as any || "none"}
                            connectionId={prof.connectionId}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* CONNECTIONS TAB */}
              {activeTab === "connections" && (
                <ConnectionsList />
              )}

              {/* REQUESTS TAB */}
              {activeTab === "requests" && (
                <PendingRequestsPanel />
              )}
            </>
          )}

        </section>

        <SiteFooter />
      </div>
    </AuthGuard>
  );
}
