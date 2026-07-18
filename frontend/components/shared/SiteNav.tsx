"use client";

import React, { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";
import { useTheme } from "next-themes";
import { styled, useTheme as useMuiTheme } from "@mui/material/styles";
import useMediaQuery from "@mui/material/useMediaQuery";
import {
  Sun,
  Moon,
  Menu,
  X,
  LogOut,
  Settings,
  User,
  ChevronRight,
  Briefcase,
  Users,
  FileText,
  GraduationCap,
  Compass,
  Bell,
  Shield,
  Search,
  MessageSquare,
  LayoutDashboard,
  BrainCircuit,
  CornerDownLeft,
  CheckCircle,
  Sparkles,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

interface SiteNavProps {
  breadcrumb?: BreadcrumbItem[];
}

const PRIMARY_LINKS = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { label: "Jobs", href: "/jobs", icon: Briefcase },
  { label: "Network", href: "/network", icon: Users },
  { label: "AI Coach", href: "/ai-coach", icon: Sparkles },
];

// Styled Components
const Header = styled("header")(({ theme }) => ({
  position: "sticky",
  top: 0,
  zIndex: 40,
  width: "100%",
  height: "72px",
  backgroundColor: theme.palette.mode === "dark" ? "rgba(9, 9, 11, 0.8)" : "rgba(250, 250, 251, 0.8)",
  backdropFilter: "blur(12px)",
  borderBottom: `1px solid ${theme.palette.divider}`,
  transition: "all 0.3s ease",
}));

const Container = styled("div")({
  marginLeft: "auto",
  marginRight: "auto",
  maxWidth: "1280px",
  paddingLeft: "16px",
  paddingRight: "16px",
  height: "100%",
  "@media (min-width: 640px)": {
    paddingLeft: "24px",
    paddingRight: "24px",
  },
  "@media (min-width: 1024px)": {
    paddingLeft: "32px",
    paddingRight: "32px",
  },
});

const FlexBetween = styled("div")({
  display: "flex",
  height: "100%",
  alignItems: "center",
  justifyContent: "space-between",
  gap: "16px",
});

const LogoText = styled("span")(({ theme }) => ({
  fontSize: "20px",
  fontWeight: 900,
  letterSpacing: "-0.025em",
  background: theme.palette.mode === "dark" 
    ? "linear-gradient(135deg, #60A5FA 0%, #818CF8 50%, #A78BFA 100%)"
    : "linear-gradient(135deg, #2563EB 0%, #4F46E5 50%, #7C3AED 100%)",
  WebkitBackgroundClip: "text",
  WebkitTextFillColor: "transparent",
  transition: "opacity 0.2s",
  "&:hover": {
    opacity: 0.9,
  },
}));

const SearchTrigger = styled("button")(({ theme }) => ({
  width: "100%",
  display: "flex",
  alignItems: "center",
  gap: "10px",
  backgroundColor: theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.04)" : "rgba(15, 23, 42, 0.04)",
  "&:hover": {
    backgroundColor: theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.08)" : "rgba(15, 23, 42, 0.08)",
  },
  border: `1px solid ${theme.palette.divider}`,
  borderRadius: "9999px",
  padding: "6px 16px",
  transition: "all 0.2s",
  color: theme.palette.text.secondary,
  fontSize: "12px",
  textAlign: "left",
  cursor: "pointer",
}));

const SearchTriggerShortcut = styled("span")(({ theme }) => ({
  marginLeft: "auto",
  backgroundColor: theme.palette.mode === "dark" ? "#111318" : "#FFFFFF",
  border: `1px solid ${theme.palette.divider}`,
  fontSize: "10px",
  fontFamily: "monospace",
  padding: "2px 6px",
  borderRadius: "6px",
  color: theme.palette.text.secondary,
  opacity: 0.8,
  transform: "scale(0.9)",
}));

const NavLink = styled(Link, {
  shouldForwardProp: (prop) => prop !== "active",
})<{ active?: boolean }>(({ theme, active }) => ({
  position: "relative",
  padding: "6px 12px",
  borderRadius: "9999px",
  fontSize: "12px",
  fontWeight: 600,
  letterSpacing: "0.025em",
  textDecoration: "none",
  transition: "all 0.2s",
  color: active ? theme.palette.text.primary : theme.palette.text.secondary,
  "&:hover": {
    color: theme.palette.text.primary,
  },
}));

const ActionButton = styled("button")(({ theme }) => ({
  padding: "8px",
  borderRadius: "9999px",
  backgroundColor: "transparent",
  border: "none",
  color: theme.palette.text.secondary,
  cursor: "pointer",
  transition: "all 0.2s",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  "&:hover": {
    color: theme.palette.text.primary,
    backgroundColor: theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.05)" : "rgba(15, 23, 42, 0.05)",
  },
}));

const ProfileMenu = styled(motion.div)(({ theme }) => ({
  position: "absolute",
  right: 0,
  marginTop: "8px",
  width: "224px",
  borderRadius: "12px",
  border: `1px solid ${theme.palette.divider}`,
  backgroundColor: theme.palette.background.paper,
  padding: "6px",
  color: theme.palette.text.primary,
  boxShadow: "0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)",
  outline: "none",
  zIndex: 20,
}));

const MenuItemLink = styled(Link)<{ dest?: boolean }>(({ theme, dest }) => ({
  display: "flex",
  width: "100%",
  alignItems: "center",
  gap: "8px",
  padding: "6px 12px",
  borderRadius: "8px",
  fontSize: "12px",
  fontWeight: dest ? 600 : 500,
  textDecoration: "none",
  color: dest 
    ? theme.palette.error.main 
    : theme.palette.text.secondary,
  transition: "colors 0.2s",
  "&:hover": {
    color: dest ? theme.palette.error.main : theme.palette.text.primary,
    backgroundColor: dest 
      ? (theme.palette.mode === "dark" ? "rgba(239, 68, 68, 0.1)" : "rgba(239, 68, 68, 0.05)")
      : (theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.05)" : "rgba(15, 23, 42, 0.05)"),
  },
}));

const ThemeBadge = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  gap: "2px",
  backgroundColor: theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.05)" : "rgba(15, 23, 42, 0.05)",
  border: `1px solid ${theme.palette.divider}`,
  padding: "4px",
  borderRadius: "9999px",
}));

const BreadcrumbBar = styled("div")(({ theme }) => ({
  borderBottom: `1px solid ${theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.05)" : "rgba(15, 23, 42, 0.05)"}`,
  backgroundColor: theme.palette.mode === "dark" ? "rgba(255, 255, 255, 0.01)" : "rgba(15, 23, 42, 0.01)",
  paddingTop: "6px",
  paddingBottom: "6px",
  paddingLeft: "16px",
  paddingRight: "16px",
  "@media (min-width: 640px)": {
    paddingLeft: "24px",
    paddingRight: "24px",
  },
  "@media (min-width: 1024px)": {
    paddingLeft: "32px",
    paddingRight: "32px",
  },
}));

const MobileNavDrawer = styled(motion.div)(({ theme }) => ({
  position: "fixed",
  left: 0,
  right: 0,
  top: "72px",
  backgroundColor: theme.palette.mode === "dark" ? "rgba(9, 9, 11, 0.95)" : "rgba(250, 250, 251, 0.95)",
  backdropFilter: "blur(16px)",
  borderBottom: `1px solid ${theme.palette.divider}`,
  zIndex: 30,
  boxShadow: "0 10px 15px -3px rgba(0, 0, 0, 0.1)",
  overflowY: "auto",
  maxHeight: "calc(100vh - 72px - 64px)",
  padding: "16px",
}));

const MobileBottomNav = styled("nav")(({ theme }) => ({
  position: "fixed",
  bottom: 0,
  left: 0,
  right: 0,
  height: "64px",
  backgroundColor: theme.palette.mode === "dark" ? "rgba(9, 9, 11, 0.9)" : "rgba(250, 250, 251, 0.9)",
  backdropFilter: "blur(12px)",
  borderTop: `1px solid ${theme.palette.divider}`,
  zIndex: 40,
  display: "grid",
  gridTemplateColumns: "repeat(5, minmax(0, 1fr))",
  alignItems: "center",
  justifyContent: "center",
  paddingBottom: "env(safe-area-inset-bottom, 0px)",
}));

const MobileBottomLink = styled(Link, {
  shouldForwardProp: (prop) => prop !== "active",
})<{ active?: boolean }>(({ theme, active }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  gap: "4px",
  textAlign: "center",
  textDecoration: "none",
  color: active ? theme.palette.primary.main : theme.palette.text.secondary,
  transition: "all 0.2s",
  "& span": {
    fontSize: "10px",
    fontWeight: 500,
  },
}));

const CommandPaletteOverlay = styled("div")({
  position: "fixed",
  inset: 0,
  zIndex: 55,
  overflowY: "auto",
});

const CommandPaletteBackdrop = styled(motion.div)({
  position: "fixed",
  inset: 0,
  backgroundColor: "rgba(0, 0, 0, 0.4)",
  backdropFilter: "blur(4px)",
});

const CommandPaletteBox = styled(motion.div)(({ theme }) => ({
  position: "relative",
  width: "100%",
  maxWidth: "512px",
  borderRadius: "16px",
  border: `1px solid ${theme.palette.divider}`,
  backgroundColor: theme.palette.background.paper,
  padding: "12px",
  boxShadow: "0 25px 50px -12px rgba(0, 0, 0, 0.25)",
  outline: "none",
}));

export default function SiteNav({ breadcrumb }: SiteNavProps) {
  const { user, loading, signOut } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const { theme, setTheme } = useTheme();
  const muiTheme = useMuiTheme();
  const isMobile = useMediaQuery(muiTheme.breakpoints.down("md"));

  const [mounted, setMounted] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  
  // Command Palette State
  const [searchQuery, setSearchQuery] = useState("");
  const [searchTab, setSearchTab] = useState<"all" | "jobs" | "people" | "communities" | "mentors">("all");
  const [selectedIndex, setSelectedIndex] = useState(0);

  const searchInputRef = useRef<HTMLInputElement>(null);
  const commandPaletteRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setMounted(true);
  }, []);

  // Keyboard shortcut listener for Ctrl+K / Cmd+K
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === "k") {
        e.preventDefault();
        setIsSearchOpen((prev) => !prev);
      }
      if (e.key === "Escape") {
        setIsSearchOpen(false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  // Autofocus input when Command Palette opens
  useEffect(() => {
    if (isSearchOpen) {
      setTimeout(() => {
        searchInputRef.current?.focus();
      }, 50);
      setSelectedIndex(0);
      setSearchQuery("");
    }
  }, [isSearchOpen]);

  const handleSignOut = async () => {
    setProfileMenuOpen(false);
    router.push("/");

    if (typeof window !== "undefined") {
      let attempts = 0;
      const interval = setInterval(async () => {
        attempts++;
        if (window.location.pathname === "/" || attempts > 100) {
          clearInterval(interval);
          await signOut();
        }
      }, 50);
    } else {
      await signOut();
    }
  };

  const getInitials = (name: string) => {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return "?";
    if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
    return (
      parts[0].charAt(0) + parts[parts.length - 1].charAt(0)
    ).toUpperCase();
  };

  const formatName = (name: string) => {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return "";
    if (parts.length === 1) return parts[0];
    return `${parts[0]} ${parts[parts.length - 1].charAt(0)}.`;
  };

  // Mock Command Palette items
  const MOCK_ITEMS = [
    { type: "jobs", title: "Senior React Engineer", desc: "Linear • Remote", href: "/jobs" },
    { type: "jobs", title: "Staff Product Designer", desc: "Stripe • SF / Hybrid", href: "/jobs" },
    { type: "people", title: "Asha Rao", desc: "Principal Recruiter at Stripe", href: "/network" },
    { type: "people", title: "Devon Webb", desc: "Senior iOS Lead at Vercel", href: "/network" },
    { type: "communities", title: "Product Designers Circle", desc: "1.2k members", href: "/communities" },
    { type: "communities", title: "Go Backend Guild", desc: "840 members", href: "/communities" },
    { type: "mentors", title: "Sarah Jenkins", desc: "Director of UX at Notion", href: "/mentorship" },
    { type: "mentors", title: "Marcus Chen", desc: "Engineering VP at Linear", href: "/mentorship" },
  ];

  const filteredSearchItems = MOCK_ITEMS.filter((item) => {
    const matchesTab = searchTab === "all" || item.type === searchTab;
    const matchesQuery =
      item.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.desc.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesTab && matchesQuery;
  });

  // Handle Command Palette arrow navigation
  const handleSearchKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setSelectedIndex((prev) => (prev + 1) % Math.max(1, filteredSearchItems.length));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setSelectedIndex((prev) => (prev - 1 + filteredSearchItems.length) % Math.max(1, filteredSearchItems.length));
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (filteredSearchItems[selectedIndex]) {
        router.push(filteredSearchItems[selectedIndex].href);
        setIsSearchOpen(false);
      }
    }
  };

  return (
    <>
      <Header>
        <Container>
          <FlexBetween style={{ display: "flex", justifyContent: "space-between" }}>
            
            {/* LEFT: Logo & Desktop Navigation */}
            <div style={{ display: "flex", alignItems: "center", gap: "24px", flexShrink: 0 }}>
              <Link href={user ? "/dashboard" : "/"} style={{ display: "flex", alignItems: "center", gap: "8px", textDecoration: "none" }}>
                <LogoText>Kirmya</LogoText>
              </Link>

              {user && (
                <>
                  <div style={{ height: "16px", width: "1px", backgroundColor: "rgba(128, 128, 128, 0.25)" }} className="hidden-mobile-divider" />
                  <nav style={{ display: "flex", alignItems: "center", gap: "4px" }} className="desktop-only-nav">
                    {PRIMARY_LINKS.map((link) => {
                      const isActive = pathname === link.href || (link.href !== "/dashboard" && pathname.startsWith(link.href));
                      return (
                        <NavLink
                          key={link.href}
                          href={link.href}
                          active={isActive}
                        >
                          {isActive && (
                            <motion.span
                              layoutId="activeNavBackground"
                              style={{
                                position: "absolute",
                                inset: 0,
                                borderRadius: "9999px",
                                zIndex: -1,
                              }}
                              className="nav-active-pill"
                              transition={{ type: "spring", stiffness: 380, damping: 30 }}
                            />
                          )}
                          {link.label}
                        </NavLink>
                      );
                    })}
                  </nav>
                </>
              )}
            </div>

            {/* CENTER: Command Palette Search Bar */}
            {user && (
              <div style={{ flexGrow: 1, maxWidth: "320px", marginLeft: "16px", marginRight: "16px" }} className="desktop-only-nav">
                <SearchTrigger onClick={() => setIsSearchOpen(true)}>
                  <Search style={{ height: "14px", width: "14px", opacity: 0.6 }} />
                  <span>Search (Ctrl+K)...</span>
                  <SearchTriggerShortcut>⌘K</SearchTriggerShortcut>
                </SearchTrigger>
              </div>
            )}

            {/* RIGHT: Global Actions */}
            <div style={{ display: "flex", alignItems: "center", gap: "12px", flexShrink: 0 }}>
              
              {/* Theme Toggle */}
              {mounted && (
                <ThemeBadge>
                  <ActionButton
                    onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
                    style={{ padding: "6px" }}
                    title={theme === "dark" ? "Switch to Light Mode" : "Switch to Dark Mode"}
                  >
                    {theme === "dark" ? <Sun style={{ height: "14px", width: "14px" }} /> : <Moon style={{ height: "14px", width: "14px" }} />}
                  </ActionButton>
                </ThemeBadge>
              )}

              {loading ? (
                <div style={{ height: "32px", width: "64px", borderRadius: "9999px", backgroundColor: "rgba(128, 128, 128, 0.15)" }} />
              ) : user ? (
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  
                  {/* Messages link */}
                  <ActionButton
                    onClick={() => router.push("/messages")}
                    style={{
                      backgroundColor: pathname.startsWith("/messages") || pathname.startsWith("/inbox") ? "rgba(128, 128, 128, 0.12)" : "transparent",
                      color: pathname.startsWith("/messages") || pathname.startsWith("/inbox") ? "var(--mui-palette-text-primary)" : "inherit"
                    }}
                    title="Messages"
                  >
                    <MessageSquare style={{ height: "16px", width: "16px" }} />
                  </ActionButton>

                  {/* Profile Menu Trigger */}
                  <div style={{ position: "relative" }}>
                    <button
                      onClick={() => setProfileMenuOpen(!profileMenuOpen)}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "8px",
                        background: "none",
                        border: "none",
                        outline: "none",
                        padding: 0,
                        cursor: "pointer"
                      }}
                      aria-label={formatName(user.full_name)}
                      aria-haspopup="menu"
                    >
                      <div style={{
                        height: "32px",
                        width: "32px",
                        borderRadius: "50%",
                        backgroundColor: "rgba(37, 99, 235, 0.1)",
                        border: "1px solid rgba(37, 99, 235, 0.2)",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        color: "#2563EB",
                        fontWeight: "bold",
                        fontSize: "12px",
                        userSelect: "none"
                      }}>
                        {getInitials(user.full_name)}
                      </div>
                      <span className="desktop-only-name" style={{
                        fontSize: "12px",
                        fontWeight: 600,
                        color: "var(--mui-palette-text-secondary)"
                      }}>
                        {formatName(user.full_name)}
                      </span>
                    </button>

                    <AnimatePresence>
                      {profileMenuOpen && (
                        <>
                          <div style={{ position: "fixed", inset: 0, zIndex: 10 }} onClick={() => setProfileMenuOpen(false)} />
                          <ProfileMenu
                            role="menu"
                            initial={{ opacity: 0, y: 10, scale: 0.95 }}
                            animate={{ opacity: 1, y: 0, scale: 1 }}
                            exit={{ opacity: 0, y: 10, scale: 0.95 }}
                            transition={{ duration: 0.12 }}
                          >
                            <div style={{ padding: "8px 12px", fontSize: "12px", borderBottom: "1px solid var(--mui-palette-divider)", marginBottom: "4px" }}>
                              <p style={{ fontWeight: 600, margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{user.full_name}</p>
                              <p style={{ color: "var(--mui-palette-text-secondary)", margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{user.email}</p>
                            </div>

                            <MenuItemLink href="/profile" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <User style={{ height: "14px", width: "14px" }} />
                              My Profile
                            </MenuItemLink>

                            <MenuItemLink href="/resume" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <FileText style={{ height: "14px", width: "14px" }} />
                              My Resume
                            </MenuItemLink>

                            <MenuItemLink href="/jobs/saved" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <Briefcase style={{ height: "14px", width: "14px" }} />
                              Saved Jobs
                            </MenuItemLink>

                            <MenuItemLink href="/applications" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <CheckCircle style={{ height: "14px", width: "14px" }} />
                              Applications
                            </MenuItemLink>

                            <MenuItemLink href="/messages" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <MessageSquare style={{ height: "14px", width: "14px" }} />
                              Messages
                            </MenuItemLink>

                            <MenuItemLink href="/notifications" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <Bell style={{ height: "14px", width: "14px" }} />
                              Notifications
                            </MenuItemLink>

                            <MenuItemLink href="/settings" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <Settings style={{ height: "14px", width: "14px" }} />
                              Settings
                            </MenuItemLink>

                            <MenuItemLink href="/help" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <GraduationCap style={{ height: "14px", width: "14px" }} />
                              Help Center
                            </MenuItemLink>

                            <MenuItemLink href="/privacy" role="menuitem" onClick={() => setProfileMenuOpen(false)}>
                              <Shield style={{ height: "14px", width: "14px" }} />
                              Privacy
                            </MenuItemLink>

                            {user?.roles?.includes("admin") && (
                              <MenuItemLink href="/admin" role="link" onClick={() => setProfileMenuOpen(false)} style={{ color: "var(--mui-palette-primary-main)" }}>
                                <Shield style={{ height: "14px", width: "14px" }} />
                                Admin
                              </MenuItemLink>
                            )}

                            <div style={{ borderTop: "1px solid var(--mui-palette-divider)", margin: "4px 0" }} />

                            <button
                              onClick={handleSignOut}
                              role="menuitem"
                              style={{
                                display: "flex",
                                width: "100%",
                                alignItems: "center",
                                gap: "8px",
                                padding: "6px 12px",
                                borderRadius: "8px",
                                fontSize: "12px",
                                fontWeight: 600,
                                border: "none",
                                outline: "none",
                                background: "none",
                                color: "var(--mui-palette-error-main)",
                                cursor: "pointer",
                                transition: "all 0.2s"
                              }}
                              className="signout-btn"
                            >
                              <LogOut style={{ height: "14px", width: "14px" }} />
                              Sign out
                            </button>
                          </ProfileMenu>
                        </>
                      )}
                    </AnimatePresence>
                  </div>

                </div>
              ) : (
                <div className="desktop-only-nav" style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  <Link
                    href="/sign-in"
                    style={{
                      padding: "6px 16px",
                      borderRadius: "9999px",
                      fontSize: "12px",
                      fontWeight: 600,
                      textDecoration: "none",
                      color: "var(--mui-palette-text-primary)",
                      transition: "background-color 0.2s"
                    }}
                    className="signin-link"
                  >
                    Sign in
                  </Link>
                  <Link
                    href="/sign-up"
                    style={{
                      padding: "6px 16px",
                      borderRadius: "9999px",
                      fontSize: "12px",
                      fontWeight: 600,
                      textDecoration: "none",
                      backgroundColor: "var(--mui-palette-primary-main)",
                      color: "#FFFFFF",
                      boxShadow: "0 1px 2px rgba(0, 0, 0, 0.05)",
                      transition: "opacity 0.2s"
                    }}
                    className="signup-link"
                  >
                    Start comeback
                  </Link>
                </div>
              )}

              {/* Mobile hamburger menu toggle */}
              <button
                onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
                style={{
                  padding: "8px",
                  borderRadius: "50%",
                  border: "none",
                  background: "none",
                  outline: "none",
                  color: "var(--mui-palette-text-secondary)",
                  cursor: "pointer"
                }}
                className="mobile-hamburger-toggle"
                aria-label="Toggle Navigation Menu"
              >
                {mobileMenuOpen ? <X style={{ height: "20px", width: "20px" }} /> : <Menu style={{ height: "20px", width: "20px" }} />}
              </button>
            </div>

          </FlexBetween>
        </Container>
      </Header>



      {/* MOBILE EXPANDED MENU DRAWER */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <MobileNavDrawer
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
          >
            <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
              {user ? (
                <>
                  <div style={{ padding: "8px 12px", borderBottom: "1px solid var(--mui-palette-divider)", marginBottom: "8px" }}>
                    <p style={{ fontSize: "14px", fontWeight: "bold", margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{user.full_name}</p>
                    <p style={{ fontSize: "12px", color: "var(--mui-palette-text-secondary)", margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{user.email}</p>
                  </div>

                  {PRIMARY_LINKS.map((link) => {
                    const Icon = link.icon;
                    return (
                      <Link
                        key={link.href}
                        href={link.href}
                        onClick={() => setMobileMenuOpen(false)}
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: "12px",
                          padding: "8px 12px",
                          borderRadius: "12px",
                          fontSize: "14px",
                          fontWeight: 500,
                          textDecoration: "none",
                          color: "var(--mui-palette-text-secondary)",
                          transition: "all 0.2s"
                        }}
                        className="mobile-drawer-link"
                      >
                        <Icon style={{ height: "18px", width: "18px" }} />
                        {link.label}
                      </Link>
                    );
                  })}

                  <div style={{ borderTop: "1px solid var(--mui-palette-divider)", margin: "8px 0", paddingTop: "8px" }} />

                  <Link
                    href="/profile"
                    onClick={() => setMobileMenuOpen(false)}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "12px",
                      padding: "8px 12px",
                      borderRadius: "12px",
                      fontSize: "14px",
                      fontWeight: 500,
                      textDecoration: "none",
                      color: "var(--mui-palette-text-secondary)",
                      transition: "all 0.2s"
                    }}
                    className="mobile-drawer-link"
                  >
                    <User style={{ height: "18px", width: "18px" }} />
                    My Profile
                  </Link>
                  <Link
                    href="/settings"
                    onClick={() => setMobileMenuOpen(false)}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "12px",
                      padding: "8px 12px",
                      borderRadius: "12px",
                      fontSize: "14px",
                      fontWeight: 500,
                      textDecoration: "none",
                      color: "var(--mui-palette-text-secondary)",
                      transition: "all 0.2s"
                    }}
                    className="mobile-drawer-link"
                  >
                    <Settings style={{ height: "18px", width: "18px" }} />
                    Settings
                  </Link>
                  {user?.roles?.includes("admin") && (
                    <Link
                      href="/admin"
                      onClick={() => setMobileMenuOpen(false)}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "12px",
                        padding: "8px 12px",
                        borderRadius: "12px",
                        fontSize: "14px",
                        fontWeight: 600,
                        textDecoration: "none",
                        color: "var(--mui-palette-primary-main)",
                        transition: "all 0.2s"
                      }}
                      className="mobile-drawer-link-admin"
                    >
                      <Shield style={{ height: "18px", width: "18px" }} />
                      Admin
                    </Link>
                  )}
                  <button
                    onClick={handleSignOut}
                    style={{
                      display: "flex",
                      width: "100%",
                      alignItems: "center",
                      gap: "12px",
                      padding: "8px 12px",
                      borderRadius: "12px",
                      fontSize: "14px",
                      fontWeight: 600,
                      border: "none",
                      outline: "none",
                      background: "none",
                      color: "var(--mui-palette-error-main)",
                      cursor: "pointer",
                      textAlign: "left",
                      transition: "all 0.2s"
                    }}
                    className="mobile-drawer-link-signout"
                  >
                    <LogOut style={{ height: "18px", width: "18px" }} />
                    Sign Out
                  </button>
                </>
              ) : (
                <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "8px", paddingTop: "8px" }}>
                  <Link
                    href="/sign-in"
                    onClick={() => setMobileMenuOpen(false)}
                    style={{
                      display: "flex",
                      justifyContent: "center",
                      alignItems: "center",
                      padding: "10px 0",
                      borderRadius: "12px",
                      fontSize: "14px",
                      fontWeight: 600,
                      textDecoration: "none",
                      color: "var(--mui-palette-text-primary)",
                      border: "1px solid var(--mui-palette-divider)",
                      transition: "background-color 0.2s"
                    }}
                    className="mobile-signin-btn"
                  >
                    Sign In
                  </Link>
                  <Link
                    href="/sign-up"
                    onClick={() => setMobileMenuOpen(false)}
                    style={{
                      display: "flex",
                      justifyContent: "center",
                      alignItems: "center",
                      padding: "10px 0",
                      borderRadius: "12px",
                      fontSize: "14px",
                      fontWeight: 600,
                      textDecoration: "none",
                      backgroundColor: "var(--mui-palette-primary-main)",
                      color: "#FFFFFF",
                      transition: "opacity 0.2s"
                    }}
                    className="mobile-signup-btn"
                  >
                    Join Kirmya
                  </Link>
                </div>
              )}
            </div>
          </MobileNavDrawer>
        )}
      </AnimatePresence>

      {/* MOBILE BOTTOM NAVIGATION BAR (md:hidden) */}
      {user && mounted && isMobile && (
        <MobileBottomNav>
          <MobileBottomLink
            href="/dashboard"
            active={pathname === "/dashboard" || pathname === "/"}
          >
            <LayoutDashboard style={{ height: "20px", width: "20px" }} />
            <span>Home</span>
          </MobileBottomLink>

          <MobileBottomLink
            href="/jobs"
            active={pathname.startsWith("/jobs")}
          >
            <Briefcase style={{ height: "20px", width: "20px" }} />
            <span>Jobs</span>
          </MobileBottomLink>

          <MobileBottomLink
            href="/network"
            active={pathname.startsWith("/network")}
          >
            <Users style={{ height: "20px", width: "20px" }} />
            <span>Network</span>
          </MobileBottomLink>

          <MobileBottomLink
            href="/messages"
            active={pathname.startsWith("/messages") || pathname.startsWith("/inbox")}
          >
            <MessageSquare style={{ height: "20px", width: "20px" }} />
            <span>Messages</span>
          </MobileBottomLink>

          <MobileBottomLink
            href="/profile"
            active={pathname.startsWith("/profile")}
          >
            <User style={{ height: "20px", width: "20px" }} />
            <span>Profile</span>
          </MobileBottomLink>
        </MobileBottomNav>
      )}

      {/* COMMAND PALETTE SEARCH OVERLAY (Ctrl+K / ⌘K) */}
      <AnimatePresence>
        {isSearchOpen && (
          <CommandPaletteOverlay>
            {/* Backdrop */}
            <CommandPaletteBackdrop
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsSearchOpen(false)}
            />

            {/* Dialog Content */}
            <div style={{ display: "flex", minHeight: "100vh", alignItems: "start", justifyContent: "center", padding: "16px", paddingTop: "15vh" }}>
              <CommandPaletteBox
                initial={{ opacity: 0, scale: 0.97, y: -10 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.97, y: -10 }}
                transition={{ duration: 0.15 }}
                ref={commandPaletteRef}
              >
                {/* Search Input */}
                <div style={{ display: "flex", alignItems: "center", gap: "12px", borderBottom: "1px solid var(--mui-palette-divider)", paddingBottom: "10px", paddingLeft: "8px", paddingRight: "8px" }}>
                  <Search style={{ height: "16px", width: "16px", color: "var(--mui-palette-text-secondary)" }} />
                  <input
                    ref={searchInputRef}
                    type="text"
                    value={searchQuery}
                    onChange={(e) => {
                      setSearchQuery(e.target.value);
                      setSelectedIndex(0);
                    }}
                    onKeyDown={handleSearchKeyDown}
                    placeholder="Search jobs, people, communities, mentors..."
                    style={{
                      width: "100%",
                      backgroundColor: "transparent",
                      border: "none",
                      outline: "none",
                      fontSize: "14px",
                      color: "var(--mui-palette-text-primary)"
                    }}
                    className="palette-input"
                  />
                  <span style={{
                    fontSize: "10px",
                    color: "var(--mui-palette-text-secondary)",
                    backgroundColor: "rgba(128, 128, 128, 0.1)",
                    border: "1px solid var(--mui-palette-divider)",
                    padding: "2px 6px",
                    borderRadius: "4px",
                    userSelect: "none"
                  }}>ESC</span>
                </div>

                {/* Filter Tabs */}
                <div style={{ display: "flex", alignItems: "center", gap: "4px", paddingTop: "8px", paddingBottom: "8px", borderBottom: "1px solid var(--mui-palette-divider)", marginBottom: "8px", overflowX: "auto" }}>
                  {(["all", "jobs", "people", "communities", "mentors"] as const).map((tab) => {
                    const isSelected = searchTab === tab;
                    return (
                      <button
                        key={tab}
                        onClick={() => {
                          setSearchTab(tab);
                          setSelectedIndex(0);
                        }}
                        style={{
                          padding: "4px 10px",
                          borderRadius: "6px",
                          fontSize: "11px",
                          fontWeight: 600,
                          textTransform: "capitalize",
                          letterSpacing: "0.025em",
                          border: "1px solid transparent",
                          cursor: "pointer",
                          transition: "all 0.2s",
                          backgroundColor: isSelected ? "rgba(37, 99, 235, 0.1)" : "transparent",
                          color: isSelected ? "#2563EB" : "var(--mui-palette-text-secondary)",
                          borderColor: isSelected ? "rgba(37, 99, 235, 0.2)" : "transparent",
                        }}
                        className={!isSelected ? "palette-tab-inactive" : ""}
                      >
                        {tab}
                      </button>
                    );
                  })}
                </div>

                {/* Search Results list */}
                <div style={{ maxHeight: "256px", overflowY: "auto", paddingBottom: "4px" }}>
                  {filteredSearchItems.length > 0 ? (
                    filteredSearchItems.map((item, idx) => {
                      const isSelected = selectedIndex === idx;
                      return (
                        <div
                          key={idx}
                          onClick={() => {
                            router.push(item.href);
                            setIsSearchOpen(false);
                          }}
                          style={{
                            display: "flex",
                            alignItems: "center",
                            justifyContent: "space-between",
                            padding: "8px 12px",
                            borderRadius: "12px",
                            cursor: "pointer",
                            transition: "all 0.2s",
                            backgroundColor: isSelected ? "rgba(128, 128, 128, 0.1)" : "transparent",
                            color: isSelected ? "var(--mui-palette-text-primary)" : "var(--mui-palette-text-secondary)"
                          }}
                          className={!isSelected ? "palette-item-hover" : ""}
                        >
                          <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                            {item.type === "jobs" && <Briefcase style={{ height: "16px", width: "16px", color: "#2563EB" }} />}
                            {item.type === "people" && <User style={{ height: "16px", width: "16px", color: "#10B981" }} />}
                            {item.type === "communities" && <Users style={{ height: "16px", width: "16px", color: "#8B5CF6" }} />}
                            {item.type === "mentors" && <GraduationCap style={{ height: "16px", width: "16px", color: "#F59E0B" }} />}
                            <div>
                              <p style={{ fontSize: "12px", color: "var(--mui-palette-text-primary)", fontWeight: 600, margin: 0 }}>{item.title}</p>
                              <p style={{ fontSize: "10px", color: "var(--mui-palette-text-secondary)", margin: 0, opacity: 0.8 }}>{item.desc}</p>
                            </div>
                          </div>
                          {isSelected && (
                            <span style={{
                              display: "flex",
                              alignItems: "center",
                              gap: "2px",
                              fontSize: "9px",
                              color: "var(--mui-palette-text-secondary)",
                              backgroundColor: "var(--mui-palette-background-default)",
                              border: "1px solid var(--mui-palette-divider)",
                              padding: "2px 4px",
                              borderRadius: "4px"
                            }}>
                              Enter <CornerDownLeft style={{ height: "8px", width: "8px" }} />
                            </span>
                          )}
                        </div>
                      );
                    })
                  ) : (
                    <div style={{ padding: "32px 0", textAlign: "center", fontSize: "12px", color: "var(--mui-palette-text-secondary)", display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", gap: "6px" }}>
                      <BrainCircuit style={{ height: "24px", width: "24px", opacity: 0.45 }} />
                      <span>No matching results found</span>
                    </div>
                  )}
                </div>
              </CommandPaletteBox>
            </div>
          </CommandPaletteOverlay>
        )}
      </AnimatePresence>
    </>
  );
}
