"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Chip,
  Avatar,
  CircularProgress,
  Alert,
  Divider,
} from "@mui/material";
import CheckIcon from "@mui/icons-material/Check";
import CloseIcon from "@mui/icons-material/Close";
import MessageIcon from "@mui/icons-material/Message";
import GroupIcon from "@mui/icons-material/Group";
import EditIcon from "@mui/icons-material/Edit";
import VisibilityIcon from "@mui/icons-material/Visibility";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api } from "@/lib/api/client";
import { profileClient, Profile } from "@/lib/api/profile";
import { networkClient, Connection } from "@/lib/api/network";

interface DashboardSummary {
  unread_notifications: number;
  job_seeker: {
    applications: number;
    saved_jobs: number;
    outgoing_referrals: number;
  };
  recruiter: {
    posted_jobs: number;
    total_applicants: number;
    incoming_referrals: number;
  };
  mentor: {
    upcoming_sessions: number;
    pending_requests: number;
    completed_sessions: number;
  };
}

export default function DashboardPage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [incomingRequests, setIncomingRequests] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadData = async () => {
    try {
      const p = await profileClient.getMe();
      setProfile(p);

      const s = await api.get<DashboardSummary>("/me/dashboard");
      setSummary(s);

      const conns = await networkClient.getConnections();
      setConnections(conns);

      const reqs = await networkClient.getIncomingRequests();
      setIncomingRequests(reqs);
    } catch (err) {
      setError("Failed to load dashboard data.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAccept = async (id: string) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      await networkClient.acceptRequest(id);
      await loadData();
    } catch (e) {
      console.error(e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async (id: string) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      await networkClient.rejectRequest(id);
      await loadData();
    } catch (e) {
      console.error(e);
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column", background: "#FBF7F2" }}>
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]} />
        <Box sx={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center", py: 8 }}>
          <CircularProgress color="primary" />
        </Box>
        <SiteFooter />
      </Box>
    );
  }

  if (error || !profile) {
    return (
      <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column", background: "#FBF7F2" }}>
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]} />
        <Container maxWidth="lg" sx={{ flex: 1, py: 8 }}>
          <Alert severity="error" sx={{ borderRadius: 3 }}>
            {error || "Could not load dashboard."}
          </Alert>
        </Container>
        <SiteFooter />
      </Box>
    );
  }

  // Determine the display name
  const displayName = profile.headline || "Professional";
  const firstName = displayName.split(" ")[0];

  return (
    <Box
      sx={{
        background: "#FBF7F2",
        fontFamily: "var(--font-public-sans), sans-serif",
        color: "#2B2620",
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Dashboard" }]} />

      {/* Header section */}
      <Box sx={{ maxWidth: 1240, width: "100%", mx: "auto", px: 5, pt: 6, pb: 4 }}>
        <Typography
          sx={{
            fontSize: "0.8rem",
            fontWeight: 700,
            letterSpacing: "0.12em",
            textTransform: "uppercase",
            color: "#C2683C",
            mb: 1.5,
          }}
        >
          Welcome back
        </Typography>
        <Typography
          variant="h2"
          sx={{
            fontFamily: "var(--font-bricolage)",
            fontWeight: 800,
            fontSize: { xs: "2rem", md: "2.8rem" },
            lineHeight: 1.05,
            letterSpacing: "-0.025em",
            mb: 1,
          }}
        >
          Good morning, {firstName}.
        </Typography>
        <Typography variant="body1" sx={{ color: "#5B554C", fontSize: "1.1rem" }}>
          Here’s what’s moving today in your professional recovery network.
        </Typography>
      </Box>

      {/* Stats Summary cards */}
      <Container maxWidth="lg" sx={{ px: { xs: 3, md: 5 }, mb: 4 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6} md={3}>
            <Card sx={{ borderRadius: 5, border: "1px solid #EFE7DC", boxShadow: "none", p: 3, background: "#fff" }}>
              <Typography variant="body2" sx={{ color: "#8A8175", mb: 1, fontWeight: 500 }}>
                Active applications
              </Typography>
              <Typography variant="h3" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 800 }}>
                {summary?.job_seeker.applications ?? 0}
              </Typography>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card sx={{ borderRadius: 5, border: "1px solid #EFE7DC", boxShadow: "none", p: 3, background: "#fff" }}>
              <Typography variant="body2" sx={{ color: "#8A8175", mb: 1, fontWeight: 500 }}>
                Saved Jobs
              </Typography>
              <Typography variant="h3" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 800 }}>
                {summary?.job_seeker.saved_jobs ?? 0}
              </Typography>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card sx={{ borderRadius: 5, border: "1px solid #EFE7DC", boxShadow: "none", p: 3, background: "#fff" }}>
              <Typography variant="body2" sx={{ color: "#8A8175", mb: 1, fontWeight: 500 }}>
                Connections
              </Typography>
              <Typography variant="h3" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 800, color: "#C2683C" }}>
                {connections.length}
              </Typography>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card sx={{ borderRadius: 5, border: "1px solid #EFE7DC", boxShadow: "none", p: 3, background: "#fff" }}>
              <Typography variant="body2" sx={{ color: "#8A8175", mb: 1, fontWeight: 500 }}>
                Profile strength
              </Typography>
              <Typography variant="h3" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 800, color: "#4F7C6A" }}>
                {profile.profile_completeness_score}%
              </Typography>
            </Card>
          </Grid>
        </Grid>
      </Container>

      {/* Main dashboard content */}
      <Container maxWidth="lg" sx={{ px: { xs: 3, md: 5 }, pb: 8 }}>
        <Grid container spacing={4} alignItems="start">
          <Grid item xs={12} md={8}>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 4 }}>
              {/* Profile Snapshot card */}
              <Card sx={{ borderRadius: 6, border: "1px solid #EFE7DC", boxShadow: "none", p: 4, background: "#fff" }}>
                <Grid container spacing={3} alignItems="center">
                  <Grid item xs={12} sm="auto">
                    <Avatar
                      src={profile.photo_url || "/assets/avatar-marcus.svg"}
                      alt={profile.headline}
                      sx={{ width: 76, height: 76, borderRadius: 4, backgroundColor: "#F3E7DC" }}
                    />
                  </Grid>
                  <Grid item xs={12} sm sx={{ flex: 1 }}>
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1.5, flexWrap: "wrap", mb: 0.5 }}>
                      <Typography variant="h5" sx={{ fontWeight: 800, fontFamily: "var(--font-bricolage)" }}>
                        {profile.headline || "Add a headline"}
                      </Typography>
                      {profile.career_status && (
                        <Chip
                          label={profile.career_status.replace("_", " ")}
                          size="small"
                          sx={{
                            backgroundColor: "rgba(79, 124, 106, 0.12)",
                            color: "#4F7C6A",
                            fontWeight: 700,
                            fontSize: "0.75rem",
                          }}
                        />
                      )}
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                      📍 {profile.location || "Add your location"} &bull; {profile.open_to_remote ? "Open to Remote" : ""}
                    </Typography>
                    <Box sx={{ display: "flex", gap: 1.5, flexWrap: "wrap" }}>
                      <Button
                        component={Link}
                        href="/profile/edit"
                        variant="contained"
                        size="small"
                        startIcon={<EditIcon />}
                        sx={{
                          borderRadius: "100px",
                          textTransform: "none",
                          fontWeight: 600,
                          backgroundColor: "#C2683C",
                          "&:hover": { backgroundColor: "#a8562f" },
                        }}
                      >
                        Edit profile
                      </Button>
                      <Button
                        component={Link}
                        href="/profile"
                        variant="outlined"
                        size="small"
                        startIcon={<VisibilityIcon />}
                        sx={{
                          borderRadius: "100px",
                          textTransform: "none",
                          fontWeight: 600,
                          borderColor: "#D8CFC2",
                          color: "#2B2620",
                          "&:hover": { borderColor: "#2B2620" },
                        }}
                      >
                        Public view
                      </Button>
                    </Box>
                  </Grid>
                </Grid>
              </Card>

              {/* Incoming Connection Requests Panel */}
              {incomingRequests.length > 0 && (
                <Card sx={{ borderRadius: 6, border: "1px solid #EFE7DC", boxShadow: "none", p: 4, background: "#fff" }}>
                  <Typography
                    variant="h5"
                    sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 700, mb: 3, color: "#2B2620" }}
                  >
                    Pending Connection Requests ({incomingRequests.length})
                  </Typography>
                  <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
                    {incomingRequests.map((req) => {
                      const requesterName = req.requester_name || "Professional";
                      const requesterHeadline = req.requester_headline || "Kirmya Member";
                      return (
                        <Box
                          key={req.id}
                          sx={{
                            display: "flex",
                            alignItems: "center",
                            justifyContent: "space-between",
                            p: 2.5,
                            border: "1px solid #EFE7DC",
                            borderRadius: 4,
                            flexWrap: "wrap",
                            gap: 2,
                          }}
                        >
                          <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                            <Avatar
                              src={req.requester_photo_url || "/assets/avatar-marcus.svg"}
                              sx={{ borderRadius: 3, width: 48, height: 48 }}
                            />
                            <Box>
                              <Typography
                                component={Link}
                                href={`/profile/${req.requester_id}`}
                                variant="subtitle1"
                                sx={{
                                  fontWeight: 700,
                                  color: "#2B2620",
                                  textDecoration: "none",
                                  "&:hover": { color: "#C2683C" },
                                }}
                              >
                                {requesterName}
                              </Typography>
                              <Typography variant="body2" color="text.secondary">
                                {requesterHeadline}
                              </Typography>
                            </Box>
                          </Box>
                          <Box sx={{ display: "flex", gap: 1 }}>
                            <Button
                              variant="contained"
                              color="primary"
                              disabled={actionLoading}
                              onClick={() => handleAccept(req.id)}
                              startIcon={<CheckIcon />}
                              sx={{
                                borderRadius: "100px",
                                textTransform: "none",
                                fontWeight: 600,
                                px: 3,
                              }}
                            >
                              Accept
                            </Button>
                            <Button
                              variant="outlined"
                              color="inherit"
                              disabled={actionLoading}
                              onClick={() => handleReject(req.id)}
                              startIcon={<CloseIcon />}
                              sx={{
                                borderRadius: "100px",
                                textTransform: "none",
                                fontWeight: 600,
                              }}
                            >
                              Ignore
                            </Button>
                          </Box>
                        </Box>
                      );
                    })}
                  </Box>
                </Card>
              )}

              {/* Connections List Card */}
              <Card sx={{ borderRadius: 6, border: "1px solid #EFE7DC", boxShadow: "none", p: 4, background: "#fff" }}>
                <Box sx={{ display: "flex", alignItems: "center", gap: 1.5, mb: 3 }}>
                  <GroupIcon color="primary" />
                  <Typography
                    variant="h5"
                    sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 700, color: "#2B2620", m: 0 }}
                  >
                    Your Connections ({connections.length})
                  </Typography>
                </Box>

                {connections.length === 0 ? (
                  <Box sx={{ py: 4, textAlign: "center" }}>
                    <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
                      You haven’t connected with anyone yet.
                    </Typography>
                    <Button
                      component={Link}
                      href="/search?type=user"
                      variant="outlined"
                      sx={{ borderRadius: "100px", textTransform: "none", color: "#C2683C", borderColor: "#C2683C" }}
                    >
                      Search professionals to connect
                    </Button>
                  </Box>
                ) : (
                  <Grid container spacing={2.5}>
                    {connections.map((c) => {
                      const isRequester = c.requester_id === profile.user_id;
                      const connName = isRequester ? c.receiver_name : c.requester_name;
                      const connHeadline = isRequester ? c.receiver_headline : c.requester_headline;
                      const connPhoto = isRequester ? c.receiver_photo_url : c.requester_photo_url;
                      const connID = isRequester ? c.receiver_id : c.requester_id;

                      return (
                        <Grid item xs={12} sm={6} key={c.id}>
                          <Box
                            sx={{
                              p: 2.5,
                              border: "1px solid #EFE7DC",
                              borderRadius: 4,
                              height: "100%",
                              display: "flex",
                              flexDirection: "column",
                              justifyContent: "space-between",
                            }}
                          >
                            <Box sx={{ display: "flex", alignItems: "flex-start", gap: 2, mb: 2 }}>
                              <Avatar src={connPhoto || "/assets/avatar-marcus.svg"} sx={{ borderRadius: 3, width: 44, height: 44 }} />
                              <Box sx={{ minWidth: 0 }}>
                                <Typography
                                  component={Link}
                                  href={`/profile/${connID}`}
                                  variant="subtitle1"
                                  sx={{
                                    fontWeight: 700,
                                    color: "#2B2620",
                                    textDecoration: "none",
                                    "&:hover": { color: "#C2683C" },
                                    display: "block",
                                    overflow: "hidden",
                                    textOverflow: "ellipsis",
                                    whiteSpace: "nowrap",
                                  }}
                                >
                                  {connName || "Kirmya Member"}
                                </Typography>
                                <Typography
                                  variant="body2"
                                  color="text.secondary"
                                  sx={{
                                    overflow: "hidden",
                                    textOverflow: "ellipsis",
                                    whiteSpace: "nowrap",
                                  }}
                                >
                                  {connHeadline || "Professional"}
                                </Typography>
                              </Box>
                            </Box>

                            <Button
                              component={Link}
                              href="/inbox"
                              variant="outlined"
                              color="inherit"
                              fullWidth
                              startIcon={<MessageIcon sx={{ fontSize: 16 }} />}
                              sx={{
                                borderRadius: "100px",
                                textTransform: "none",
                                fontWeight: 600,
                                fontSize: "0.85rem",
                                borderColor: "#D8CFC2",
                                color: "#5B554C",
                                "&:hover": { borderColor: "#2B2620", color: "#2B2620" },
                              }}
                            >
                              Message
                            </Button>
                          </Box>
                        </Grid>
                      );
                    })}
                  </Grid>
                )}
              </Card>
            </Box>
          </Grid>

          {/* Right sidebar */}
          <Grid item xs={12} md={4}>
            <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
              {/* Profile Strength Card */}
              <Card sx={{ borderRadius: 6, border: "1px solid #EFE7DC", boxShadow: "none", p: 4, background: "#fff" }}>
                <Typography variant="subtitle1" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 700, mb: 1.5 }}>
                  Profile completeness
                </Typography>
                <Box sx={{ height: 10, background: "#F3ECE2", borderRadius: 100, overflow: "hidden", mb: 1 }}>
                  <Box sx={{ width: `${profile.profile_completeness_score}%`, height: "100%", background: "#4F7C6A" }} />
                </Box>
                <Typography variant="body2" color="text.secondary">
                  {profile.profile_completeness_score}% complete. Keep adding experiences to showcase your career comeback story.
                </Typography>
              </Card>

              {/* Recommended Jobs */}
              <Card sx={{ borderRadius: 6, border: "1px solid #EFE7DC", boxShadow: "none", p: 4, background: "#fff" }}>
                <Typography variant="subtitle1" sx={{ fontFamily: "var(--font-bricolage)", fontWeight: 700, mb: 2 }}>
                  Recommended roles
                </Typography>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
                  <Box component={Link} href="/jobs" sx={{ textDecoration: "none", color: "inherit", display: "block" }}>
                    <Typography variant="subtitle2" sx={{ fontWeight: 700, "&:hover": { color: "#C2683C" } }}>
                      VP, Supply Chain Operations
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Atlas Co &bull; Bangalore / Remote
                    </Typography>
                  </Box>
                  <Divider />
                  <Box component={Link} href="/jobs" sx={{ textDecoration: "none", color: "inherit", display: "block" }}>
                    <Typography variant="subtitle2" sx={{ fontWeight: 700, "&:hover": { color: "#C2683C" } }}>
                      Director of Logistics & Operations
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Vertex &bull; Mumbai / Hybrid
                    </Typography>
                  </Box>
                </Box>
              </Card>
            </Box>
          </Grid>
        </Grid>
      </Container>

      <SiteFooter />
    </Box>
  );
}
