"use client";

import React, { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
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
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import PersonAddIcon from "@mui/icons-material/PersonAdd";
import HourglassEmptyIcon from "@mui/icons-material/HourglassEmpty";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { profileClient, Profile } from "@/lib/api/profile";
import { networkClient, ConnectionStatusResponse, Connection } from "@/lib/api/network";
import { ApiError } from "@/lib/api/client";

export default function OtherProfilePage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [profile, setProfile] = useState<Profile | null>(null);
  const [currentUserID, setCurrentUserID] = useState<string | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatusResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    (async () => {
      try {
        // Fetch current user ID
        const me = await profileClient.getMe();
        setCurrentUserID(me.user_id);

        if (me.user_id === id) {
          router.replace("/profile");
          return;
        }

        // Fetch dynamic profile
        const data = await profileClient.getByID(id);
        setProfile(data);

        // Fetch relationship status
        const conn = await networkClient.getConnectionStatus(id);
        setConnectionStatus(conn);
      } catch (err) {
        setError(
          err instanceof ApiError
            ? err.message
            : "Could not load profile. It might be private or not exist."
        );
      } finally {
        setLoading(false);
      }
    })();
  }, [id, router]);

  const handleConnect = async () => {
    if (!id || actionLoading) return;
    setActionLoading(true);
    try {
      await networkClient.sendRequest(id);
      const conn = await networkClient.getConnectionStatus(id);
      setConnectionStatus(conn);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleAcceptConnection = async (reqID: string) => {
    if (!reqID || actionLoading) return;
    setActionLoading(true);
    try {
      await networkClient.acceptRequest(reqID);
      const conn = await networkClient.getConnectionStatus(id);
      setConnectionStatus(conn);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleRejectConnection = async (reqID: string) => {
    if (!reqID || actionLoading) return;
    setActionLoading(true);
    try {
      await networkClient.rejectRequest(reqID);
      const conn = await networkClient.getConnectionStatus(id);
      setConnectionStatus(conn);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const renderConnectionButton = () => {
    if (!connectionStatus) return null;

    const { status, requester_id } = connectionStatus;

    if (status === "accepted") {
      return (
        <Button
          variant="outlined"
          color="secondary"
          startIcon={<CheckCircleOutlineIcon />}
          disabled
          sx={{
            borderRadius: "100px",
            textTransform: "none",
            fontWeight: 600,
            fontFamily: "var(--font-public-sans)",
            borderColor: "secondary.main",
            "&.Mui-disabled": {
              color: "secondary.main",
              borderColor: "secondary.main",
            },
          }}
        >
          Connected
        </Button>
      );
    }

    if (status === "pending") {
      if (requester_id === currentUserID) {
        return (
          <Button
            variant="outlined"
            color="warning"
            startIcon={<HourglassEmptyIcon />}
            disabled
            sx={{
              borderRadius: "100px",
              textTransform: "none",
              fontWeight: 600,
              fontFamily: "var(--font-public-sans)",
              borderColor: "warning.main",
              "&.Mui-disabled": {
                color: "warning.main",
                borderColor: "warning.main",
              },
            }}
          >
            Pending Request
          </Button>
        );
      } else {
        // Find the incoming request ID. Since we don't have it directly in Status, we fetch incoming requests
        return (
          <Box sx={{ display: "flex", gap: 1 }}>
            <Button
              variant="contained"
              color="primary"
              disabled={actionLoading}
              onClick={async () => {
                try {
                  const reqs = await networkClient.getIncomingRequests();
                  const found = reqs.find((r) => r.requester_id === id);
                  if (found) {
                    await handleAcceptConnection(found.id);
                  }
                } catch (e) {
                  console.error(e);
                }
              }}
              sx={{
                borderRadius: "100px",
                textTransform: "none",
                fontWeight: 600,
                fontFamily: "var(--font-public-sans)",
                px: 3,
              }}
            >
              Accept
            </Button>
            <Button
              variant="outlined"
              color="inherit"
              disabled={actionLoading}
              onClick={async () => {
                try {
                  const reqs = await networkClient.getIncomingRequests();
                  const found = reqs.find((r) => r.requester_id === id);
                  if (found) {
                    await handleRejectConnection(found.id);
                  }
                } catch (e) {
                  console.error(e);
                }
              }}
              sx={{
                borderRadius: "100px",
                textTransform: "none",
                fontWeight: 600,
                fontFamily: "var(--font-public-sans)",
              }}
            >
              Ignore
            </Button>
          </Box>
        );
      }
    }

    return (
      <Button
        variant="contained"
        color="primary"
        startIcon={<PersonAddIcon />}
        disabled={actionLoading}
        onClick={handleConnect}
        sx={{
          borderRadius: "100px",
          textTransform: "none",
          fontWeight: 600,
          fontFamily: "var(--font-public-sans)",
          px: 4,
          py: 1.5,
          backgroundColor: "#C2683C",
          "&:hover": {
            backgroundColor: "#a8562f",
          },
        }}
      >
        Connect
      </Button>
    );
  };

  if (loading) {
    return (
      <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column", background: "#FBF7F2" }}>
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]} />
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
        <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Profile" }]} />
        <Container maxWidth="md" sx={{ flex: 1, py: 8 }}>
          <Alert severity="error" sx={{ borderRadius: 3 }}>
            {error || "Profile could not be loaded."}
          </Alert>
        </Container>
        <SiteFooter />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        background: "#FBF7F2",
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        fontFamily: "var(--font-public-sans)",
      }}
    >
      <SiteNav
        breadcrumb={[
          { label: "Home", href: "/" },
          { label: "Profiles", href: "/search?type=user" },
          { label: profile.headline || "Profile" },
        ]}
      />

      <Container maxWidth="md" sx={{ flex: 1, py: { xs: 6, md: 8 } }}>
        {/* Core Identity Panel */}
        <Card
          sx={{
            borderRadius: 6,
            border: "1px solid #EFE7DC",
            boxShadow: "none",
            p: { xs: 3, md: 4 },
            mb: 4,
            background: "#fff",
          }}
        >
          <Grid container spacing={3} alignItems="flex-start">
            <Grid item xs={12} sm="auto">
              <Avatar
                src={profile.photo_url || "/assets/avatar-marcus.svg"}
                alt={profile.headline}
                sx={{
                  width: 96,
                  height: 96,
                  borderRadius: 5,
                  backgroundColor: "#F3E7DC",
                }}
              />
            </Grid>
            <Grid item xs={12} sm={8} sx={{ flex: 1 }}>
              <Box sx={{ display: "flex", alignItems: "center", gap: 1.5, flexWrap: "wrap", mb: 1 }}>
                <Typography
                  variant="h4"
                  sx={{
                    fontWeight: 800,
                    fontFamily: "var(--font-bricolage)",
                    color: "#2B2620",
                  }}
                >
                  {profile.headline || "Professional"}
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
                      textTransform: "capitalize",
                    }}
                  />
                )}
              </Box>

              <Typography variant="body1" color="text.secondary" sx={{ mb: 2, fontSize: "1.1rem" }}>
                {profile.bio || "Career recovery & transition professional"}
              </Typography>

              <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap", mb: 3 }}>
                {profile.location && (
                  <Chip
                    label={`📍 ${profile.location}`}
                    size="small"
                    sx={{ background: "#F3ECE2", color: "#5B554C", fontWeight: 500 }}
                  />
                )}
                {profile.open_to_remote && (
                  <Chip
                    label="💻 Open to Remote"
                    size="small"
                    sx={{ background: "#F3ECE2", color: "#5B554C", fontWeight: 500 }}
                  />
                )}
                {profile.willing_to_mentor && (
                  <Chip
                    label="🤝 Willing to Mentor"
                    size="small"
                    sx={{ background: "#F3ECE2", color: "#5B554C", fontWeight: 500 }}
                  />
                )}
              </Box>

              {/* Action Button */}
              <Box sx={{ display: "flex", gap: 2, alignItems: "center" }}>
                {renderConnectionButton()}
              </Box>
            </Grid>
          </Grid>
        </Card>

        {/* Dynamic Detail Sections */}
        <Grid container spacing={4}>
          <Grid item xs={12} md={8}>
            {/* About / Summary */}
            {profile.about && (
              <Box sx={{ mb: 4 }}>
                <Typography
                  variant="h6"
                  sx={{
                    fontFamily: "var(--font-bricolage)",
                    fontWeight: 700,
                    mb: 1.5,
                    color: "#2B2620",
                  }}
                >
                  About
                </Typography>
                <Typography variant="body1" sx={{ color: "#4A443B", lineHeight: 1.7 }}>
                  {profile.about}
                </Typography>
              </Box>
            )}

            {/* Experience */}
            {profile.experiences && profile.experiences.length > 0 && (
              <Box sx={{ mb: 4 }}>
                <Typography
                  variant="h6"
                  sx={{
                    fontFamily: "var(--font-bricolage)",
                    fontWeight: 700,
                    mb: 2,
                    color: "#2B2620",
                  }}
                >
                  Experience
                </Typography>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
                  {profile.experiences.map((exp, index) => (
                    <Box key={exp.id || index} sx={{ display: "flex", gap: 2 }}>
                      <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
                        <Box
                          sx={{
                            width: 12,
                            height: 12,
                            borderRadius: "50%",
                            background: index === 0 ? "#C2683C" : "#4F7C6A",
                            mt: 0.5,
                          }}
                        />
                        {index < profile.experiences.length - 1 && (
                          <Box sx={{ width: 2, flex: 1, background: "#EFE7DC", my: 0.5 }} />
                        )}
                      </Box>
                      <Box>
                        <Typography variant="subtitle1" sx={{ fontWeight: 700, color: "#2B2620" }}>
                          {exp.title}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                          {exp.company} &bull; {exp.start_date} &ndash; {exp.is_current ? "Present" : exp.end_date}
                        </Typography>
                        {exp.description && (
                          <Typography variant="body2" sx={{ color: "#5B554C", lineHeight: 1.6 }}>
                            {exp.description}
                          </Typography>
                        )}
                      </Box>
                    </Box>
                  ))}
                </Box>
              </Box>
            )}

            {/* Education */}
            {profile.educations && profile.educations.length > 0 && (
              <Box sx={{ mb: 4 }}>
                <Typography
                  variant="h6"
                  sx={{
                    fontFamily: "var(--font-bricolage)",
                    fontWeight: 700,
                    mb: 2,
                    color: "#2B2620",
                  }}
                >
                  Education
                </Typography>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
                  {profile.educations.map((edu, index) => (
                    <Box key={edu.id || index}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 700, color: "#2B2620" }}>
                        {edu.school}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {edu.degree} &bull; {edu.field_of_study}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {edu.start_date} &ndash; {edu.end_date}
                      </Typography>
                    </Box>
                  ))}
                </Box>
              </Box>
            )}
          </Grid>

          {/* Sidebar */}
          <Grid item xs={12} md={4}>
            {/* Skills */}
            {profile.skills && profile.skills.length > 0 && (
              <Card
                sx={{
                  borderRadius: 4,
                  border: "1px solid #EFE7DC",
                  boxShadow: "none",
                  p: 3,
                  mb: 3,
                  background: "#fff",
                }}
              >
                <Typography
                  variant="subtitle1"
                  sx={{
                    fontFamily: "var(--font-bricolage)",
                    fontWeight: 700,
                    mb: 2,
                    color: "#2B2620",
                  }}
                >
                  Skills
                </Typography>
                <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                  {profile.skills.map((sk) => (
                    <Chip
                      key={sk.name}
                      label={sk.name}
                      variant="outlined"
                      size="small"
                      sx={{
                        borderRadius: "100px",
                        borderColor: "#EFE7DC",
                        color: "#2B2620",
                      }}
                    />
                  ))}
                </Box>
              </Card>
            )}

            {/* Languages */}
            {profile.languages && profile.languages.length > 0 && (
              <Card
                sx={{
                  borderRadius: 4,
                  border: "1px solid #EFE7DC",
                  boxShadow: "none",
                  p: 3,
                  background: "#fff",
                }}
              >
                <Typography
                  variant="subtitle1"
                  sx={{
                    fontFamily: "var(--font-bricolage)",
                    fontWeight: 700,
                    mb: 2,
                    color: "#2B2620",
                  }}
                >
                  Languages
                </Typography>
                <Box sx={{ display: "flex", flexDirection: "column", gap: 1 }}>
                  {profile.languages.map((l) => (
                    <Box key={l.name} sx={{ display: "flex", justifyContent: "space-between" }}>
                      <Typography variant="body2" sx={{ fontWeight: 600, color: "#2B2620" }}>
                        {l.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {l.proficiency}
                      </Typography>
                    </Box>
                  ))}
                </Box>
              </Card>
            )}
          </Grid>
        </Grid>
      </Container>

      <SiteFooter />
    </Box>
  );
}
