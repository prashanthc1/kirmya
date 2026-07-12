"use client";

import React, { useState, useEffect } from "react";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api } from "@/lib/api/client";
import { useAuth } from "@/lib/auth/auth-context";
import { useNotifications } from "@/components/shared/Notifications";
import { CircularProgress } from "@mui/material";

interface Community {
  id: string;
  slug: string;
  name: string;
  description: string;
  category: string;
  member_count: number;
}

const CATEGORIES = [
  "Technology",
  "Artificial Intelligence",
  "Software Engineering",
  "Cybersecurity",
  "Cloud Computing",
  "DevOps",
  "Data Science",
  "Finance",
  "Healthcare",
  "HR",
  "Marketing",
  "Product Management",
  "Career Growth",
  "Resume Reviews",
  "Interview Preparation",
  "Freelancing",
  "Remote Jobs",
  "Students",
  "Startups",
  "Networking",
];

export default function CommunitiesPage() {
  const { user } = useAuth();
  const { showNotification } = useNotifications();

  const [communities, setCommunities] = useState<Community[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);

  // Create Community Modal states
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newName, setNewName] = useState("");
  const [newSlug, setNewSlug] = useState("");
  const [newDesc, setNewDesc] = useState("");
  const [newCategory, setNewCategory] = useState("Technology");
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchCommunities();
  }, []);

  const fetchCommunities = async () => {
    setLoading(true);
    try {
      const data = await api.get<Community[]>("/communities");
      setCommunities(data || []);
    } catch (err: any) {
      showNotification(err.message || "Failed to load communities", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCommunity = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim() || !newSlug.trim()) {
      showNotification("Name and slug are required", "error");
      return;
    }

    setSubmitting(true);
    try {
      const newComm = await api.post<Community>("/communities", {
        name: newName,
        slug: newSlug,
        description: newDesc,
        category: newCategory,
      });
      showNotification("Community created successfully!", "success");
      setCommunities((prev) => [...prev, newComm]);
      setShowCreateModal(false);
      // Reset form
      setNewName("");
      setNewSlug("");
      setNewDesc("");
      setNewCategory("Technology");
    } catch (err: any) {
      showNotification(err.message || "Failed to create community", "error");
    } finally {
      setSubmitting(false);
    }
  };

  // Filter list
  const filtered = communities.filter((c) => {
    const matchesSearch =
      c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.description.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCategory = !selectedCategory || c.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  return (
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
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Communities" }]} />

      {/* Hero Banner Section */}
      <section style={{ maxWidth: "920px", margin: "0 auto", padding: "clamp(56px,7vw,100px) 40px clamp(40px,5vw,56px)", textAlign: "center" }}>
        <div style={{ display: "inline-block", fontSize: "13px", fontWeight: 700, letterSpacing: "0.1em", textTransform: "uppercase", color: "#C2683C", background: "rgba(194,104,60,0.12)", padding: "8px 16px", borderRadius: "100px", marginBottom: "26px" }}>
          Quiet Communities
        </div>
        <h1 style={{ fontWeight: 800, fontSize: "clamp(40px,6.5vw,72px)", lineHeight: 1.02, letterSpacing: "-0.025em", margin: "0 auto 22px", maxWidth: "760px" }}>
          Circles where professionals share leads, not selfies.
        </h1>
        <p style={{ fontSize: "clamp(17px,2vw,20px)", lineHeight: 1.6, color: "#5B554C", maxWidth: "600px", margin: "0 auto 34px" }}>
          No algorithms, no performative feeds. Just small, dedicated groups where experienced practitioners support each other with real job referrals and career support.
        </p>
        {user ? (
          <button
            onClick={() => setShowCreateModal(true)}
            style={{ border: "none", background: "#C2683C", color: "#fff", fontSize: "16px", fontWeight: 600, padding: "16px 32px", borderRadius: "100px", cursor: "pointer" }}
          >
            Create a New Circle
          </button>
        ) : (
          <a
            href="/sign-in"
            style={{ background: "#C2683C", color: "#fff", fontSize: "16px", fontWeight: 600, padding: "16px 32px", borderRadius: "100px", display: "inline-block", textDecoration: "none" }}
          >
            Sign In to Join
          </a>
        )}
      </section>

      {/* Search & Filter Toolbar */}
      <section style={{ maxWidth: "1240px", margin: "0 auto", width: "100%", padding: "0 40px clamp(24px,4vw,32px)" }}>
        <div style={{ display: "flex", gap: "16px", flexWrap: "wrap", alignItems: "center", background: "#ffffff", border: "1px solid #EFE7DC", borderRadius: "20px", padding: "18px 24px" }}>
          <input
            type="text"
            placeholder="Search circles by name or topic..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            style={{
              flex: 1,
              minWidth: "260px",
              border: "1px solid #E2D9CC",
              borderRadius: "10px",
              padding: "12px 14px",
              fontSize: "15px",
              background: "#FCFAF7",
              outline: "none",
              fontFamily: "'Public Sans', sans-serif",
            }}
          />
          <button
            onClick={() => {
              setSearchQuery("");
              setSelectedCategory(null);
            }}
            style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "12px 20px", borderRadius: "10px", fontWeight: 600, cursor: "pointer", fontSize: "14px" }}
          >
            Clear Filters
          </button>
        </div>

        {/* Categories Carousel */}
        <div style={{ display: "flex", gap: "8px", overflowX: "auto", padding: "14px 0", scrollbarWidth: "none" }}>
          <button
            onClick={() => setSelectedCategory(null)}
            style={{
              flex: "none",
              border: "none",
              borderRadius: "100px",
              padding: "8px 16px",
              fontSize: "13px",
              fontWeight: 600,
              cursor: "pointer",
              background: selectedCategory === null ? "#C2683C" : "#EFE7DC",
              color: selectedCategory === null ? "#FFFFFF" : "#5B554C",
            }}
          >
            All Circles
          </button>
          {CATEGORIES.map((cat) => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              style={{
                flex: "none",
                border: "none",
                borderRadius: "100px",
                padding: "8px 16px",
                fontSize: "13px",
                fontWeight: 600,
                cursor: "pointer",
                background: selectedCategory === cat ? "#C2683C" : "#EFE7DC",
                color: selectedCategory === cat ? "#FFFFFF" : "#5B554C",
              }}
            >
              {cat}
            </button>
          ))}
        </div>
      </section>

      {/* Dynamic Grid / Loading State */}
      <section style={{ maxWidth: "1240px", margin: "0 auto", width: "100%", padding: "0 40px clamp(48px,6vw,72px)", flex: 1 }}>
        {loading ? (
          <div style={{ display: "flex", justifyContent: "center", alignItems: "center", padding: "64px" }}>
            <CircularProgress style={{ color: "#C2683C" }} />
          </div>
        ) : (
          <div>
            <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(320px, 1fr))", gap: "24px" }}>
              {filtered.map((item) => (
                <a
                  key={item.id}
                  href={`/communities/${item.slug}`}
                  style={{
                    background: "#ffffff",
                    border: "1px solid #EFE7DC",
                    borderRadius: "20px",
                    padding: "28px",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "space-between",
                    gap: "18px",
                    textDecoration: "none",
                    color: "inherit",
                    boxShadow: "0 4px 12px rgba(43, 38, 32, 0.03)",
                    transition: "transform 0.2s, box-shadow 0.2s",
                  }}
                  onMouseOver={(e) => {
                    e.currentTarget.style.transform = "translateY(-4px)";
                    e.currentTarget.style.boxShadow = "0 12px 24px rgba(43, 38, 32, 0.08)";
                  }}
                  onMouseOut={(e) => {
                    e.currentTarget.style.transform = "none";
                    e.currentTarget.style.boxShadow = "0 4px 12px rgba(43, 38, 32, 0.03)";
                  }}
                >
                  <div>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "12px" }}>
                      <span style={{ width: "42px", height: "42px", borderRadius: "10px", background: "rgba(194, 104, 60, 0.12)", color: "#C2683C", display: "flex", alignItems: "center", justifyContent: "center", fontSize: "20px", fontWeight: 700 }}>
                        {item.name.charAt(0)}
                      </span>
                      <span style={{ fontSize: "12px", color: "#4F7C6A", fontWeight: 600, background: "rgba(79,124,106,0.12)", padding: "4px 10px", borderRadius: "100px" }}>
                        {item.category || "General"}
                      </span>
                    </div>
                    <h3 style={{ fontSize: "18px", fontWeight: 700, margin: "0 0 6px 0", color: "#2B2620" }}>{item.name}</h3>
                    <p style={{ fontSize: "14px", lineHeight: "1.5", color: "#8A8175", margin: 0 }}>
                      {item.description || "A professional circle to discuss topics and share referral links."}
                    </p>
                  </div>
                  <div style={{ borderTop: "1px solid #F6EFE6", paddingTop: "14px", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <span style={{ fontSize: "13px", color: "#5B554C" }}>
                      <strong>{item.member_count}</strong> members
                    </span>
                    <span style={{ fontSize: "13px", fontWeight: 600, color: "#C2683C" }}>Enter Circle →</span>
                  </div>
                </a>
              ))}
            </div>
            {filtered.length === 0 && (
              <div style={{ textAlign: "center", padding: "64px 24px", color: "#8A8175" }}>
                <span style={{ fontSize: "40px" }}>◍</span>
                <h3 style={{ margin: "16px 0 6px", fontSize: "18px", fontWeight: 700, color: "#2B2620" }}>No circles found</h3>
                <p style={{ margin: 0, fontSize: "14px" }}>Try clearing your filters or create a new community circle above.</p>
              </div>
            )}
          </div>
        )}
      </section>

      {/* CREATE COMMUNITY MODAL */}
      {showCreateModal && (
        <div style={{ position: "fixed", top: 0, left: 0, right: 0, bottom: 0, background: "rgba(0,0,0,0.5)", display: "flex", justifyContent: "center", alignItems: "center", zIndex: 10000, padding: "20px" }}>
          <form onSubmit={handleCreateCommunity} style={{ background: "#fff", border: "1px solid #EFE7DC", borderRadius: "24px", maxWidth: "520px", width: "100%", padding: "28px", display: "flex", flexDirection: "column", gap: "20px", boxShadow: "0 24px 60px rgba(43,38,32,0.18)" }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <h2 style={{ fontSize: "20px", fontWeight: 700, margin: 0 }}>Create a Professional Circle</h2>
              <button type="button" onClick={() => setShowCreateModal(false)} style={{ border: "none", background: "transparent", fontSize: "24px", cursor: "pointer", color: "#8A8175" }}>×</button>
            </div>

            <div>
              <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Circle Name</label>
              <input
                type="text"
                required
                value={newName}
                onChange={(e) => {
                  setNewName(e.target.value);
                  // Auto-generate slug
                  setNewSlug(e.target.value.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, ""));
                }}
                placeholder="e.g. Senior Product Managers"
                style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
              />
            </div>

            <div>
              <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Unique Slug (URL Handle)</label>
              <input
                type="text"
                required
                value={newSlug}
                onChange={(e) => setNewSlug(e.target.value.toLowerCase().replace(/[^a-z0-9_-]+/g, ""))}
                placeholder="senior-product-managers"
                style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
              />
            </div>

            <div>
              <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Category</label>
              <select
                value={newCategory}
                onChange={(e) => setNewCategory(e.target.value)}
                style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px", fontSize: "15px", outline: "none", background: "#FCFAF7" }}
              >
                {CATEGORIES.map((cat) => (
                  <option key={cat} value={cat}>{cat}</option>
                ))}
              </select>
            </div>

            <div>
              <label style={{ display: "block", fontSize: "13px", fontWeight: 600, color: "#8A8175", marginBottom: "6px" }}>Description</label>
              <textarea
                value={newDesc}
                onChange={(e) => setNewDesc(e.target.value)}
                placeholder="Describe who this circle is for and what rules members should follow..."
                rows={3}
                style={{ width: "100%", border: "1px solid #E2D9CC", borderRadius: "10px", padding: "12px", fontSize: "15px", outline: "none", background: "#FCFAF7", resize: "none", fontFamily: "inherit" }}
              />
            </div>

            <div style={{ display: "flex", justifyContent: "flex-end", gap: "12px", marginTop: "8px" }}>
              <button
                type="button"
                onClick={() => setShowCreateModal(false)}
                style={{ border: "1px solid #E2D9CC", background: "transparent", color: "#5B554C", padding: "12px 24px", borderRadius: "100px", cursor: "pointer", fontSize: "14px", fontWeight: 600 }}
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={submitting}
                style={{ border: "none", background: "#C2683C", color: "#fff", padding: "12px 24px", borderRadius: "100px", cursor: submitting ? "not-allowed" : "pointer", fontSize: "14px", fontWeight: 600 }}
              >
                {submitting ? "Creating..." : "Create Circle"}
              </button>
            </div>
          </form>
        </div>
      )}

      <SiteFooter />
    </div>
  );
}
