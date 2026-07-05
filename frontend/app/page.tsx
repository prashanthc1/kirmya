"use client";

import React from "react";
import Link from "next/link";
import {
  Box,
  Container,
  Grid,
  Typography,
  Button,
  Card,
  CardContent,
  Avatar,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ArrowForwardIcon from "@mui/icons-material/ArrowForward";
import SiteNav from "@/components/shared/SiteNav";
import SiteFooter from "@/components/shared/SiteFooter";

const FEATURES = [
  {
    href: "/resume",
    icon: "◈",
    iconBg: "rgba(214, 104, 56, 0.12)",
    iconColor: "#D66838",
    title: "AI Resume Coach",
    desc: "Upload your resume and immediately get an ATS score, missing keyword metrics, and phrasing improvements that capture recruiters' attention.",
  },
  {
    href: "/referrals",
    icon: "↳",
    iconBg: "rgba(55, 97, 77, 0.12)",
    iconColor: "#37614D",
    title: "Real Referrals",
    desc: "Skip the automated hiring system filters. Secure warm introductions directly from employees working inside companies you want to join.",
  },
  {
    href: "/mentorship",
    icon: "✳",
    iconBg: "rgba(43, 38, 32, 0.08)",
    iconColor: "#2B2620",
    title: "Supportive Mentorship",
    desc: "Book free, structured sessions with experienced professionals who have sat in the exact seat you are aiming for.",
  },
  {
    href: "/communities",
    icon: "▦",
    iconBg: "rgba(106, 95, 160, 0.12)",
    iconColor: "#6A5FA0",
    title: "Quiet Communities",
    desc: "Industry-focused micro-circles for tech, operations, logistics, HR, and facilities management where members share real leads, not selfies.",
  },
  {
    href: "/career-paths",
    icon: "↗",
    iconBg: "rgba(55, 97, 77, 0.12)",
    iconColor: "#37614D",
    title: "Interactive Career Paths",
    desc: "Visually trace transition opportunities, compare market demand salary expectations, and isolate specific skill gaps holding you back.",
  },
  {
    href: "/coach",
    icon: "✦",
    iconBg: "rgba(214, 104, 56, 0.12)",
    iconColor: "#D66838",
    title: "AI Coach & Interview Prep",
    desc: "Practice roleplay interview scenarios, refine salary negotiation scripts, and plan your weekly outreach goals with an on-demand virtual partner.",
  },
];

const FAQS = [
  {
    q: "Is Kirmya really free to use?",
    a: "Yes. During career transitions, the last thing you need is another bill. The core tools—resume parsing, AI coaching suggestions, community access, and mentorship sessions—are free. We may introduce premium sponsor features later, but job recovery tools remain accessible to all.",
  },
  {
    q: "How does the referral system prevent spam?",
    a: "We do not allow cold, automated messaging. Referral requests require structured introductions, target specific open roles, and are matched based on mutual interest, protecting referrers from spam while keeping candidate signal high.",
  },
  {
    q: "What makes Kirmya different from LinkedIn?",
    a: "LinkedIn optimizes for screen time, influencers, and vanity metrics. Kirmya is a recovery workspace. There is no public content feed. We focus on one metric: the speed at which you land your next interview.",
  },
];

export default function HomePage() {
  return (
    <Box sx={{ minHeight: "100vh", display: "flex", flexDirection: "column" }}>
      <SiteNav />

      {/* Hero Section */}
      <Box
        component="section"
        sx={{
          py: { xs: 8, md: 14 },
          textAlign: "center",
          background: "radial-gradient(circle at top, rgba(214, 104, 56, 0.04) 0%, rgba(252, 250, 247, 0) 60%)",
        }}
      >
        <Container maxWidth="md" className="animate-fade-in-up">
          <Typography
            component="span"
            sx={{
              display: "inline-block",
              fontSize: "0.85rem",
              fontWeight: 800,
              letterSpacing: "0.15em",
              textTransform: "uppercase",
              color: "primary.main",
              backgroundColor: "rgba(214, 104, 56, 0.08)",
              px: 2.5,
              py: 1,
              borderRadius: 100,
              mb: 4,
            }}
          >
            Built for the moment between jobs
          </Typography>

          <Typography
            variant="h1"
            sx={{
              fontSize: { xs: "2.75rem", sm: "4rem", md: "4.75rem" },
              color: "text.primary",
              mb: 3,
            }}
          >
            You didn&apos;t lose your career.
            <br />
            <Box component="span" sx={{ color: "primary.main" }}>
              You just lost that one job.
            </Box>
          </Typography>

          <Typography
            variant="body1"
            color="text.secondary"
            sx={{
              fontSize: { xs: "1.1rem", md: "1.25rem" },
              maxWidth: 680,
              mx: "auto",
              mb: 5,
            }}
          >
            The gap on your resume isn&apos;t a red flag—it&apos;s a chapter. Kirmya is where
            professionals regroup, refine their materials, coordinate referrals, and come back
            stronger.
          </Typography>

          <Box
            sx={{
              display: "flex",
              gap: 2,
              justifyContent: "center",
              flexWrap: "wrap",
              mb: 3,
            }}
          >
            <Button
              component={Link}
              href="/sign-in"
              variant="contained"
              color="primary"
              size="large"
              sx={{
                py: 2,
                px: 4.5,
                fontSize: "1.05rem",
                fontWeight: 700,
                boxShadow: "0 10px 25px -5px rgba(214, 104, 56, 0.3)",
              }}
            >
              Start your comeback
            </Button>
            <Button
              component={Link}
              href="#how"
              variant="outlined"
              color="primary"
              size="large"
              sx={{
                py: 2,
                px: 4.5,
                fontSize: "1.05rem",
                fontWeight: 700,
                backgroundColor: "background.paper",
              }}
            >
              See how it works
            </Button>
          </Box>

          <Typography variant="body2" color="text.secondary">
            Free to join · No spam · Your data stays yours
          </Typography>
        </Container>
      </Box>

      {/* Features Grid */}
      <Container component="section" sx={{ pb: { xs: 8, md: 12 } }}>
        <Box sx={{ textAlign: "center", maxWidth: 600, mx: "auto", mb: { xs: 6, md: 8 } }}>
          <Typography
            variant="h2"
            sx={{
              fontSize: { xs: "2.25rem", md: "2.75rem" },
              mb: 2,
            }}
          >
            Everything the job hunt actually needs.
          </Typography>
          <Typography variant="body1" color="text.secondary">
            No social feeds to perform on. No metrics to game. Just targeted tools and warm connections
            designed to get you hired.
          </Typography>
        </Box>

        <Grid container spacing={3}>
          {FEATURES.map((feature, i) => (
            <Grid item xs={12} sm={6} md={4} key={i}>
              <Card
                className="glass-card"
                component={Link}
                href={feature.href}
                sx={{
                  display: "block",
                  textDecoration: "none",
                  height: "100%",
                  "&:hover": {
                    transform: "translateY(-4px)",
                    boxShadow: "0 12px 30px rgba(43, 38, 32, 0.08)",
                    borderColor: "primary.light",
                  },
                }}
              >
                <CardContent sx={{ p: 4, height: "100%" }}>
                  <Avatar
                    sx={{
                      bgcolor: feature.iconBg,
                      color: feature.iconColor,
                      fontSize: "1.5rem",
                      fontWeight: 700,
                      borderRadius: 3,
                      width: 48,
                      height: 48,
                      mb: 3,
                    }}
                  >
                    {feature.icon}
                  </Avatar>
                  <Typography
                    variant="h5"
                    component="h3"
                    sx={{
                      fontWeight: 700,
                      fontSize: "1.25rem",
                      mb: 1.5,
                      color: "text.primary",
                    }}
                  >
                    {feature.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {feature.desc}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>

      {/* How it Works Section */}
      <Box
        id="how"
        component="section"
        sx={{
          py: { xs: 8, md: 12 },
          backgroundColor: "#F3ECE2",
          borderTop: "1px solid rgba(43, 38, 32, 0.08)",
          borderBottom: "1px solid rgba(43, 38, 32, 0.08)",
        }}
      >
        <Container maxWidth="lg">
          <Box sx={{ textAlign: "center", maxWidth: 600, mx: "auto", mb: { xs: 6, md: 8 } }}>
            <Typography
              component="span"
              sx={{
                fontSize: "0.8rem",
                fontWeight: 800,
                letterSpacing: "0.12em",
                textTransform: "uppercase",
                color: "primary.main",
                mb: 1.5,
                display: "block",
              }}
            >
              How Kirmya works
            </Typography>
            <Typography variant="h2" sx={{ fontSize: { xs: "2.25rem", md: "2.75rem" } }}>
              The 4 Steps to Recovery
            </Typography>
          </Box>

          <Grid container spacing={4}>
            {[
              {
                num: "01",
                title: "Assess & Score",
                desc: "Import your materials. AI evaluates your ATS index, identifies skill gaps, and recommends high-demand roles matching your skill set.",
              },
              {
                num: "02",
                title: "Build the Path",
                desc: "Follow a time-boxed strategy to address skill deficiencies using free learning resources, guided by your AI Career Coach.",
              },
              {
                num: "03",
                title: "Request Warm Intros",
                desc: "Identify target companies. Skip generic forms and request direct referrals from verified employees inside those organizations.",
              },
              {
                num: "04",
                title: "Nail the Interview",
                desc: "Utilize on-demand AI mock sessions and seek feedback from volunteer mentors to secure your next role.",
              },
            ].map((step, i) => (
              <Grid item xs={12} sm={6} md={3} key={i}>
                <Box
                  sx={{
                    p: 3,
                    height: "100%",
                    display: "flex",
                    flexDirection: "column",
                    gap: 2,
                  }}
                >
                  <Typography
                    variant="h2"
                    sx={{
                      color: "primary.main",
                      opacity: 0.35,
                      fontSize: "3.5rem",
                      fontWeight: 800,
                      lineHeight: 1,
                    }}
                  >
                    {step.num}
                  </Typography>
                  <Typography
                    variant="h5"
                    component="h3"
                    sx={{ fontWeight: 800, color: "text.primary" }}
                  >
                    {step.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {step.desc}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </Container>
      </Box>

      {/* FAQ Section */}
      <Container component="section" sx={{ py: { xs: 8, md: 12 }, maxWidth: "md" }}>
        <Box sx={{ textAlign: "center", mb: { xs: 6, md: 8 } }}>
          <Typography
            variant="h2"
            sx={{
              fontSize: { xs: "2.25rem", md: "2.75rem" },
              mb: 2,
            }}
          >
            Frequently Asked Questions
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Everything you need to know about starting your comeback.
          </Typography>
        </Box>

        <Box sx={{ display: "flex", flexDirection: "column", gap: 1.5 }}>
          {FAQS.map((faq, i) => (
            <Accordion
              key={i}
              elevation={0}
              sx={{
                borderRadius: "16px !important",
                border: "1px solid rgba(43, 38, 32, 0.06)",
                background: "rgba(255, 255, 255, 0.5)",
                "&:before": { display: "none" },
              }}
            >
              <AccordionSummary
                expandIcon={<ExpandMoreIcon sx={{ color: "primary.main" }} />}
                sx={{ px: 3, py: 1 }}
              >
                <Typography variant="subtitle1" sx={{ fontWeight: 700, color: "text.primary" }}>
                  {faq.q}
                </Typography>
              </AccordionSummary>
              <AccordionDetails sx={{ px: 3, pb: 3, pt: 0 }}>
                <Typography variant="body2" color="text.secondary">
                  {faq.a}
                </Typography>
              </AccordionDetails>
            </Accordion>
          ))}
        </Box>
      </Container>

      {/* CTA Footer */}
      <Box
        component="section"
        sx={{
          py: { xs: 8, md: 10 },
          textAlign: "center",
          background: "radial-gradient(circle, rgba(214, 104, 56, 0.05) 0%, rgba(252, 250, 247, 0) 100%)",
          borderTop: "1px solid rgba(43, 38, 32, 0.06)",
        }}
      >
        <Container maxWidth="sm">
          <Typography variant="h2" sx={{ mb: 2, fontSize: { xs: "2.25rem", md: "2.75rem" } }}>
            Ready to rewrite your story?
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
            Join thousands of professionals who refused to let job loss define their trajectory.
          </Typography>
          <Button
            component={Link}
            href="/sign-in"
            variant="contained"
            color="primary"
            size="large"
            endIcon={<ArrowForwardIcon />}
            sx={{
              py: 2,
              px: 4.5,
              fontWeight: 700,
              boxShadow: "0 10px 25px -5px rgba(214, 104, 56, 0.3)",
            }}
          >
            Join Kirmya today
          </Button>
        </Container>
      </Box>

      <SiteFooter />
    </Box>
  );
}
