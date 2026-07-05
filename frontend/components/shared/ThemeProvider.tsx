"use client";

import React from "react";
import { ThemeProvider as MuiThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import GlobalStyles from "@mui/material/GlobalStyles";
import theme from "@/lib/theme";
import { NotificationsProvider } from "@/components/shared/Notifications";

interface ThemeProviderProps {
  children: React.ReactNode;
}

export default function ThemeProvider({ children }: ThemeProviderProps) {
  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <GlobalStyles
        styles={{
          ":root": {
            "--color-sunset-coral": "#D66838",
            "--color-eucalyptus-green": "#37614D",
            "--color-alabaster-warm": "#FCFAF7",
          },
          body: {
            backgroundColor: "#FCFAF7",
            color: "#2B2620",
            minHeight: "100vh",
            fontFamily: "var(--font-public-sans), sans-serif",
            WebkitFontSmoothing: "antialiased",
            MozOsxFontSmoothing: "grayscale",
          },
          "@keyframes fadeInUp": {
            from: {
              opacity: 0,
              transform: "translateY(16px)",
            },
            to: {
              opacity: 1,
              transform: "translateY(0)",
            },
          },
          "@keyframes pulseGlow": {
            "0%": {
              boxShadow: "0 0 0 0 rgba(214, 104, 56, 0.4)",
            },
            "70%": {
              boxShadow: "0 0 0 10px rgba(214, 104, 56, 0)",
            },
            "100%": {
              boxShadow: "0 0 0 0 rgba(214, 104, 56, 0)",
            },
          },
          ".animate-fade-in-up": {
            animation: "fadeInUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards",
          },
          ".pulse-indicator": {
            animation: "pulseGlow 2s infinite",
            borderRadius: "50%",
          },
          // Custom glassmorphic utilities for CSS overrides if needed
          ".glass-nav": {
            background: "rgba(252, 250, 247, 0.8) !important",
            backdropFilter: "blur(20px) saturate(180%)",
            borderBottom: "1px solid rgba(43, 38, 32, 0.06) !important",
          },
          ".glass-card": {
            background: "rgba(255, 255, 255, 0.7) !important",
            backdropFilter: "blur(12px)",
            border: "1px solid rgba(43, 38, 32, 0.05) !important",
          },
        }}
      />
      <NotificationsProvider>{children}</NotificationsProvider>
    </MuiThemeProvider>
  );
}
