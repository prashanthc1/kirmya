"use client";

import React, { Suspense, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Chip,
  Alert,
  CircularProgress,
} from "@mui/material";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";
import { api, ApiError } from "@/lib/api/client";

interface SearchHit {
  type: string;
  ref_id: string;
  title: string;
  subtitle: string;
  url: string;
  score: number;
}

interface SearchResponse {
  results: SearchHit[];
  engine: string;
}

function SearchResults() {
  const params = useSearchParams();
  const query = params.get("q") ?? "";

  const [results, setResults] = useState<SearchHit[]>([]);
  const [engine, setEngine] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Adjust state when query changes during render to avoid useEffect warnings
  const [prevQuery, setPrevQuery] = useState(query);
  if (query !== prevQuery) {
    setPrevQuery(query);
    setLoading(true);
    setResults([]);
    setError(null);
  }

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const data = await api.get<SearchResponse>(
          `/search?q=${encodeURIComponent(query)}`
        );
        if (active) {
          setResults(data?.results ?? []);
          setEngine(data?.engine ?? "");
        }
      } catch (err) {
        if (active) {
          setError(
            err instanceof ApiError
              ? err.message
              : "Search is unavailable right now. Please try again."
          );
        }
      } finally {
        if (active) setLoading(false);
      }
    })();
    return () => {
      active = false;
    };
  }, [query]);

  return (
    <Container maxWidth="md" sx={{ flex: 1, py: { xs: 6, md: 8 } }}>
      <Box sx={{ mb: 4 }}>
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
          Search
        </Typography>
        <Typography
          variant="h2"
          sx={{
            fontSize: { xs: "2.25rem", md: "3rem" },
            mb: 2,
          }}
        >
          {query ? <>Results for &ldquo;{query}&rdquo;</> : "Search Kirmya"}
        </Typography>

        <Typography variant="body2" color="text.secondary">
          {loading
            ? "Searching…"
            : `${results.length} result${results.length === 1 ? "" : "s"}`}
          {engine && !loading ? ` · Served by ${engine}` : ""}
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 4, borderRadius: 3 }}>
          {error}
        </Alert>
      )}

      {loading && (
        <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
          <CircularProgress color="primary" />
        </Box>
      )}

      {!loading && !error && results.length === 0 && query && (
        <Typography variant="body1" color="text.secondary" sx={{ py: 4 }}>
          Nothing matched &ldquo;{query}&rdquo;. Try a different term.
        </Typography>
      )}

      {!loading && !error && results.length > 0 && (
        <Grid container spacing={2}>
          {results.map((hit) => (
            <Grid item xs={12} key={`${hit.type}-${hit.ref_id}`}>
              <Card
                component={Link}
                href={hit.url || "#"}
                sx={{
                  display: "block",
                  textDecoration: "none",
                  "&:hover": {
                    transform: "translateY(-2px)",
                    boxShadow: "0 10px 30px rgba(43, 38, 32, 0.06)",
                    borderColor: "primary.light",
                  },
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
                    <Chip
                      label={hit.type}
                      sx={{
                        backgroundColor: "rgba(55, 97, 77, 0.08)",
                        color: "secondary.main",
                        fontWeight: 700,
                        fontSize: "0.75rem",
                        textTransform: "uppercase",
                      }}
                    />
                    <Box sx={{ minWidth: 0, flex: 1 }}>
                      <Typography
                        variant="h6"
                        component="div"
                        sx={{
                          fontWeight: 700,
                          color: "text.primary",
                          fontFamily: "var(--font-public-sans)",
                        }}
                      >
                        {hit.title}
                      </Typography>
                      {hit.subtitle && (
                        <Typography variant="body2" color="text.secondary" noWrap>
                          {hit.subtitle}
                        </Typography>
                      )}
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Container>
  );
}

export default function SearchPage() {
  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <SiteNav breadcrumb={[{ label: "Home", href: "/" }, { label: "Search" }]} />
      <Suspense
        fallback={
          <Container maxWidth="md" sx={{ flex: 1, py: 8 }}>
            <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
              <CircularProgress color="primary" />
            </Box>
          </Container>
        }
      >
        <SearchResults />
      </Suspense>
      <SiteFooter />
    </Box>
  );
}
