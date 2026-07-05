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
  Chip,
  Button,
  TextField,
  CircularProgress,
  Alert,
  InputAdornment,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import ArrowForwardIcon from "@mui/icons-material/ArrowForward";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface Job {
  id: string;
  title: string;
  company: string;
  location: string;
  salary: string;
  job_type: string;
  description: string;
}

const JOB_TYPES = ["All Types", "Full-time", "Contract", "Remote", "Part-time"];

export default function JobsPage() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filter & Search states
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState("All Types");

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await api.get<{ jobs: Job[] }>("/jobs");
        if (active) {
          const list = data?.jobs ?? [];
          setJobs(list);
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError ? err.message : "Could not load jobs. Please try again."
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, []);

  // Filter Handler (Computed value during render)
  const filteredJobs = React.useMemo(() => {
    let result = jobs;

    // Apply search filter
    if (searchQuery.trim() !== "") {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (job) =>
          job.title.toLowerCase().includes(q) ||
          job.company.toLowerCase().includes(q) ||
          job.description?.toLowerCase().includes(q) ||
          job.location?.toLowerCase().includes(q)
      );
    }

    // Apply job type filter
    if (selectedType !== "All Types") {
      result = result.filter(
        (job) => job.job_type?.toLowerCase() === selectedType.toLowerCase()
      );
    }

    return result;
  }, [searchQuery, selectedType, jobs]);

  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Jobs" }]} />

      <Container
        component="main"
        maxWidth="lg"
        sx={{
          flex: 1,
          py: { xs: 6, md: 8 },
        }}
      >
        {/* Page Header */}
        <Box sx={{ mb: 6 }}>
          <Typography
            component="span"
            sx={{
              fontSize: "0.8rem",
              fontWeight: 800,
              letterSpacing: "0.12em",
              textTransform: "uppercase",
              color: "primary.main",
              display: "block",
              mb: 1.5,
            }}
          >
            {loading ? "Discovering Opportunities..." : `${filteredJobs.length} active openings`}
          </Typography>
          <Typography
            variant="h2"
            sx={{
              fontSize: { xs: "2.25rem", md: "3rem" },
              mb: 2,
            }}
          >
            Roles that value your experience
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ maxWidth: 600 }}>
            Hand-vetted career opportunities from recruiters actively looking for proven skills,
            resilience, and transition readiness.
          </Typography>
        </Box>

        {/* Search & Filter bar */}
        <Box
          className="glass-card"
          sx={{
            p: 3,
            borderRadius: 4,
            mb: 5,
            display: "flex",
            flexDirection: "column",
            gap: 2.5,
            boxShadow: "0 10px 30px -5px rgba(43, 38, 32, 0.04)",
          }}
        >
          {/* Search bar */}
          <TextField
            fullWidth
            placeholder="Search by title, company, skills, or location..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon sx={{ color: "text.disabled" }} />
                </InputAdornment>
              ),
            }}
          />

          {/* Tag filter row */}
          <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap", alignItems: "center" }}>
            <Typography variant="body2" sx={{ fontWeight: 700, mr: 1, color: "text.secondary" }}>
              Filter Type:
            </Typography>
            {JOB_TYPES.map((type) => (
              <Chip
                key={type}
                label={type}
                clickable
                onClick={() => setSelectedType(type)}
                color={selectedType === type ? "primary" : "default"}
                variant={selectedType === type ? "filled" : "outlined"}
                sx={{
                  fontWeight: 600,
                  fontSize: "0.85rem",
                  px: 1,
                  py: 0.5,
                  borderRadius: "100px",
                  borderColor: "rgba(43, 38, 32, 0.12)",
                  transition: "all 0.2s ease",
                  "&:hover": {
                    borderColor: "primary.main",
                  },
                }}
              />
            ))}
          </Box>
        </Box>

        {/* Loading & Error States */}
        {loading && (
          <Box sx={{ display: "flex", justifyContent: "center", py: 10 }}>
            <CircularProgress color="primary" />
          </Box>
        )}

        {error && (
          <Alert severity="error" sx={{ mb: 4, borderRadius: 3 }}>
            {error}
          </Alert>
        )}

        {/* Job Listings Grid */}
        {!loading && !error && (
          <>
            {filteredJobs.length === 0 ? (
              <Box sx={{ textAlign: "center", py: 8 }}>
                <Typography variant="h6" color="text.secondary">
                  No matching roles found. Try adjusting your search query or filters.
                </Typography>
              </Box>
            ) : (
              <Grid container spacing={3}>
                {filteredJobs.map((job) => (
                  <Grid item xs={12} key={job.id}>
                    <Card
                      sx={{
                        "&:hover": {
                          transform: "translateY(-2px)",
                          boxShadow: "0 10px 30px rgba(43, 38, 32, 0.06)",
                          borderColor: "primary.light",
                        },
                      }}
                    >
                      <CardContent sx={{ p: 4 }}>
                        <Box
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "flex-start",
                            gap: 2,
                            flexWrap: "wrap",
                            mb: 2,
                          }}
                        >
                          <Box>
                            <Typography
                              variant="h5"
                              component="h3"
                              sx={{
                                fontWeight: 800,
                                fontFamily: "var(--font-bricolage)",
                                mb: 0.5,
                                color: "text.primary",
                              }}
                            >
                              {job.title}
                            </Typography>
                            <Typography
                              variant="subtitle1"
                              color="text.secondary"
                              sx={{ fontWeight: 500 }}
                            >
                              {job.company} {job.location ? `· ${job.location}` : ""}
                            </Typography>
                          </Box>

                          {job.job_type && (
                            <Chip
                              label={job.job_type}
                              sx={{
                                backgroundColor: "rgba(55, 97, 77, 0.08)",
                                color: "secondary.main",
                                fontWeight: 700,
                                fontSize: "0.75rem",
                                textTransform: "capitalize",
                              }}
                            />
                          )}
                        </Box>

                        {job.description && (
                          <Typography
                            variant="body2"
                            color="text.secondary"
                            sx={{
                              mb: 3,
                              display: "-webkit-box",
                              WebkitLineClamp: 3,
                              WebkitBoxOrient: "vertical",
                              overflow: "hidden",
                              lineHeight: 1.6,
                            }}
                          >
                            {job.description}
                          </Typography>
                        )}

                        <Box
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "center",
                            flexWrap: "wrap",
                            gap: 2,
                            borderTop: "1px solid rgba(43, 38, 32, 0.06)",
                            pt: 2.5,
                          }}
                        >
                          {job.salary ? (
                            <Typography
                              variant="subtitle1"
                              sx={{ fontWeight: 800, color: "secondary.main" }}
                            >
                              {job.salary}
                            </Typography>
                          ) : (
                            <Box />
                          )}

                          <Button
                            component={Link}
                            href="/jobs/detail"
                            variant="text"
                            color="primary"
                            endIcon={<ArrowForwardIcon sx={{ fontSize: 16 }} />}
                            sx={{ fontWeight: 700, p: 0, minWidth: "auto", "&:hover": { background: "none" } }}
                          >
                            View details
                          </Button>
                        </Box>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            )}
          </>
        )}
      </Container>

      <SiteFooter />
    </Box>
  );
}
