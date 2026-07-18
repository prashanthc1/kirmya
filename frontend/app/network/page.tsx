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
      <div className="min-h-screen bg-background text-foreground flex flex-col">
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Network Center" }]} />

        {/* Hero Section */}
        <section className="max-w-4xl mx-auto px-6 py-12 md:py-20 text-center space-y-6">
          <div className="inline-block text-xs font-bold tracking-widest uppercase text-primary bg-primary/10 border border-primary/20 px-4 py-2 rounded-full">
            Professional Networking
          </div>
          <h1 className="text-4xl md:text-6xl font-black tracking-tight leading-none text-foreground max-w-3xl mx-auto">
            Connect with verified practitioners.
          </h1>
          <p className="text-base md:text-lg text-muted-foreground max-w-xl mx-auto">
            Grow your professional network, share insider job openings, schedule mock interviews, and trade referrals without premium restrictions.
          </p>
        </section>

        {/* Navigation Tabs bar */}
        <section className="max-w-7xl mx-auto w-full px-6 border-b border-border">
          <div className="flex gap-6 overflow-x-auto scrollbar-none">
            {[
              { id: "discover", label: `🔍 Discover People` },
              { id: "connections", label: `👥 My Connections (${connections.length})` },
              { id: "requests", label: `📥 Connection Requests (${incomingRequests.length})` },
            ].map((t) => (
              <button
                key={t.id}
                onClick={() => setActiveTab(t.id as any)}
                className={`py-4 px-1 text-sm font-semibold border-b-2 cursor-pointer transition-all ${
                  activeTab === t.id
                    ? "border-primary text-primary"
                    : "border-transparent text-muted-foreground hover:text-foreground"
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>
        </section>

        {/* Main Content Area */}
        <section className="max-w-7xl mx-auto w-full px-6 py-8 flex-grow">
          
          {loading ? (
            <div className="flex justify-center items-center py-20">
              <CircularProgress className="text-primary" />
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
                  <form onSubmit={handleSearch} className="bg-card border border-border/80 rounded-3xl p-6 flex flex-col gap-4 mb-8 shadow-sm">
                    <div className="flex gap-4 flex-wrap items-center">
                      <div className="flex-grow min-w-[260px] flex items-center gap-3 border border-border bg-secondary/10 rounded-2xl px-4 py-2.5">
                        <Search size={16} className="text-muted-foreground" />
                        <input
                          type="text"
                          placeholder="Search professionals by name, title, company, or skills..."
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          className="border-none outline-none bg-transparent text-sm text-foreground w-full placeholder:text-muted-foreground"
                        />
                      </div>
                      <button
                        type="button"
                        onClick={() => setShowFilters(!showFilters)}
                        className="flex items-center gap-2 border border-border hover:bg-secondary/40 text-foreground px-4 py-2.5 rounded-2xl text-xs font-bold cursor-pointer transition-all"
                      >
                        <SlidersHorizontal size={14} /> Filters
                      </button>
                      <button
                        type="submit"
                        className="border-none bg-primary hover:bg-primary/95 text-primary-foreground px-6 py-2.5 rounded-2xl text-xs font-bold cursor-pointer transition-all shadow-sm"
                      >
                        Search
                      </button>
                    </div>

                    {showFilters && (
                      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 border-t border-border/40 pt-4">
                        <div>
                          <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">Location</label>
                          <input
                            type="text"
                            placeholder="e.g. San Francisco, CA"
                            value={selectedLocation}
                            onChange={(e) => setSelectedLocation(e.target.value)}
                            className="w-full border border-border rounded-xl px-3 py-2 text-sm bg-secondary/15 text-foreground placeholder:text-muted-foreground outline-none focus:ring-1 focus:ring-primary"
                          />
                        </div>
                        <div>
                          <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">Category</label>
                          <select
                            value={selectedCategory}
                            onChange={(e) => setSelectedCategory(e.target.value)}
                            className="w-full border border-border rounded-xl px-3 py-2 text-sm bg-secondary/15 text-foreground outline-none focus:ring-1 focus:ring-primary"
                          >
                            <option value="">All Categories</option>
                            <option value="tech">Technology</option>
                            <option value="recruiter">Recruiters</option>
                            <option value="mentor">Mentors</option>
                          </select>
                        </div>
                        <div>
                          <label className="block text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-1.5">Sort By</label>
                          <select
                            value={sortBy}
                            onChange={(e) => setSortBy(e.target.value)}
                            className="w-full border border-border rounded-xl px-3 py-2 text-sm bg-secondary/15 text-foreground outline-none focus:ring-1 focus:ring-primary"
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
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                    {professionals.map((prof) => (
                      <div
                        key={prof.id}
                        className="bg-card border border-border/80 rounded-3xl p-6 flex flex-col justify-between gap-4 shadow-sm hover:border-border transition-all"
                      >
                        <div>
                          <div className="flex gap-4 items-center mb-2">
                            <div className="w-12 h-12 rounded-full overflow-hidden bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-base select-none uppercase">
                              {prof.full_name.charAt(0)}
                            </div>
                            <div className="min-w-0">
                              <h3 className="text-sm font-bold text-foreground truncate">{prof.full_name}</h3>
                              <p className="text-xs text-muted-foreground truncate">{prof.headline}</p>
                            </div>
                          </div>
                          <div className="text-xs text-muted-foreground flex items-center gap-1">
                            <span>📍 {prof.location}</span>
                          </div>
                        </div>

                        <div className="border-t border-border/40 pt-4 flex gap-2 justify-end items-center">
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
