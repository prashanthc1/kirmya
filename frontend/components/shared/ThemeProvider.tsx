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

function ThemeProviderInner({ children }: { children: React.ReactNode }) {
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

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
