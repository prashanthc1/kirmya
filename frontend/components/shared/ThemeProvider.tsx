"use client";

import React, { useEffect, useState } from "react";
import { ThemeProvider as NextThemesProvider, useTheme } from "next-themes";
import {
  ThemeProvider as MuiThemeProvider,
  createTheme,
} from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { NotificationsProvider } from "@/components/shared/Notifications";

interface ThemeProviderProps {
  children: React.ReactNode;
}

function applyGlobalAccessibilityStyles(acc: any) {
  if (typeof document === "undefined" || !acc) return;
  let styleTag = document.getElementById("accessibility-styles");
  if (!styleTag) {
    styleTag = document.createElement("style");
    styleTag.id = "accessibility-styles";
    document.head.appendChild(styleTag);
  }

  let css = "";
  if (acc.font_size === "small") {
    css += "body, html { font-size: 14px !important; }";
  } else if (acc.font_size === "large") {
    css += "body, html { font-size: 18px !important; }";
  } else if (acc.font_size === "extra-large") {
    css += "body, html { font-size: 20px !important; }";
  } else {
    css += "body, html { font-size: 16px !important; }";
  }

  if (acc.high_contrast) {
    css += `
      body, html {
        filter: contrast(1.25) !important;
        background: #FFFFFF !important;
        color: #000000 !important;
      }
      button, input, select, textarea {
        border: 2px solid #000000 !important;
        color: #000000 !important;
        background: #FFFFFF !important;
      }
    `;
  }

  if (acc.reduced_motion) {
    css += `
      *, *::before, *::after {
        animation-delay: -1ms !important;
        animation-duration: 1ms !important;
        animation-iteration-count: 1 !important;
        background-attachment: initial !important;
        scroll-behavior: auto !important;
        transition-duration: 0s !important;
        transition-delay: 0s !important;
      }
    `;
  }
  styleTag.innerHTML = css;
}

function ThemeProviderInner({ children }: { children: React.ReactNode }) {
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  // Sync accessibility styles from localStorage and setup event listener
  useEffect(() => {
    if (typeof window !== "undefined") {
      const stored = window.localStorage.getItem("kirmya-accessibility");
      if (stored) {
        try {
          const parsed = JSON.parse(stored);
          applyGlobalAccessibilityStyles(parsed);
        } catch (e) {
          /* ignore */
        }
      }

      const handler = (e: any) => {
        if (e.detail) {
          applyGlobalAccessibilityStyles(e.detail);
        }
      };

      window.addEventListener("kirmya-accessibility-changed", handler);
      return () => {
        window.removeEventListener("kirmya-accessibility-changed", handler);
      };
    }
  }, []);

  // Sync MUI palette mode with next-themes resolved mode
  const dynamicMuiTheme = React.useMemo(() => {
    const isDark = resolvedTheme === "dark";
    return createTheme({
      palette: {
        mode: isDark ? "dark" : "light",
        primary: {
          main: isDark ? "#60A5FA" : "#2563EB",
        },
        background: {
          default: isDark ? "#09090B" : "#FAFAFB",
          paper: isDark ? "#111318" : "#FFFFFF",
        },
        text: {
          primary: isDark ? "#F8FAFC" : "#0F172A",
          secondary: isDark ? "#CBD5E1" : "#64748B",
        },
        divider: isDark
          ? "rgba(255, 255, 255, 0.08)"
          : "rgba(15, 23, 42, 0.08)",
      },
      typography: {
        fontFamily: "var(--font-public-sans), sans-serif",
        button: {
          textTransform: "none",
          fontWeight: 600,
        },
      },
      shape: {
        borderRadius: 12,
      },
    });
  }, [resolvedTheme]);

  useEffect(() => {
    setMounted(true);
  }, []);

  return (
    <MuiThemeProvider theme={dynamicMuiTheme}>
      <CssBaseline />
      <NotificationsProvider>{children}</NotificationsProvider>
    </MuiThemeProvider>
  );
}

export default function ThemeProvider({ children }: ThemeProviderProps) {
  return (
    <NextThemesProvider attribute="class" defaultTheme="system" enableSystem>
      <ThemeProviderInner>{children}</ThemeProviderInner>
    </NextThemesProvider>
  );
}
