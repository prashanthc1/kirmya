"use client";

import React, { useState } from "react";
import Link from "next/link";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Box,
  Container,
  Breadcrumbs,
  Link as MuiLink,
  Divider,
  ListItemIcon,
  ListItemText,
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import SettingsIcon from "@mui/icons-material/Settings";
import LogoutIcon from "@mui/icons-material/Logout";
import KeyboardArrowRightIcon from "@mui/icons-material/KeyboardArrowRight";
import PersonIcon from "@mui/icons-material/Person";
import ArrowForwardIcon from "@mui/icons-material/ArrowForward";

const MENU_LINKS = [
  { label: "Jobs", href: "/jobs", icon: "◎" },
  { label: "Referrals", href: "/referrals", icon: "↳" },
  { label: "Mentorship", href: "/mentorship", icon: "✳" },
  { label: "Communities", href: "/communities", icon: "▦" },
  { label: "Career Paths", href: "/career-paths", icon: "↗" },
  { label: "Coach", href: "/coach", icon: "✦" },
  { label: "Resume", href: "/resume", icon: "▤" },
];

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

interface SiteNavProps {
  breadcrumb?: BreadcrumbItem[];
}

export default function SiteNav({ breadcrumb }: SiteNavProps) {
  const { user, loading, signOut } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  // Mobile menu control
  const [mobileAnchor, setMobileAnchor] = useState<null | HTMLElement>(null);
  // Profile menu control
  const [profileAnchor, setProfileAnchor] = useState<null | HTMLElement>(null);

  const handleOpenMobile = (event: React.MouseEvent<HTMLElement>) => {
    setMobileAnchor(event.currentTarget);
  };
  const handleCloseMobile = () => {
    setMobileAnchor(null);
  };

  const handleOpenProfile = (event: React.MouseEvent<HTMLElement>) => {
    setProfileAnchor(event.currentTarget);
  };
  const handleCloseProfile = () => {
    setProfileAnchor(null);
  };

  const handleSignOut = async () => {
    handleCloseProfile();
    await signOut();
    router.push("/");
  };

  // Get initials for Avatar fallback
  const getInitials = (name: string) => {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return "?";
    if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
    return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
  };

  const firstName = user?.full_name?.trim().split(/\s+/)[0] || "User";
  const lastInitial =
    user?.full_name?.trim().split(/\s+/).slice(1).join(" ").charAt(0) || "";

  return (
    <AppBar
      position="sticky"
      elevation={0}
      className="glass-nav"
      sx={{
        background: "rgba(252, 250, 247, 0.8)",
        backdropFilter: "blur(20px) saturate(180%)",
        borderBottom: "1px solid rgba(43, 38, 32, 0.06)",
        top: 0,
        zIndex: 1100,
      }}
    >
      <Container maxWidth="lg">
        <Toolbar
          disableGutters
          sx={{
            height: 72,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          {/* Brand Logo */}
          <Typography
            variant="h5"
            component={Link}
            href="/"
            sx={{
              fontFamily: "var(--font-bricolage), sans-serif",
              fontWeight: 800,
              color: "text.primary",
              textDecoration: "none",
              letterSpacing: "-0.02em",
              display: "flex",
              alignItems: "center",
              mr: 2,
            }}
          >
            Kirmya
          </Typography>

          {/* Navigation Links for Logged-In Users on Desktop */}
          {user && (
            <Box sx={{ display: { xs: "none", md: "flex" }, gap: 1 }}>
              {MENU_LINKS.map((link) => {
                const isActive = pathname === link.href;
                return (
                  <Button
                    key={link.href}
                    component={Link}
                    href={link.href}
                    sx={{
                      color: isActive ? "primary.main" : "text.secondary",
                      fontWeight: isActive ? 700 : 500,
                      fontSize: "0.95rem",
                      px: 2,
                      py: 1,
                      borderRadius: 100,
                      backgroundColor: isActive ? "rgba(214, 104, 56, 0.06)" : "transparent",
                      "&:hover": {
                        backgroundColor: isActive
                          ? "rgba(214, 104, 56, 0.1)"
                          : "rgba(43, 38, 32, 0.04)",
                        color: isActive ? "primary.main" : "text.primary",
                      },
                    }}
                  >
                    {link.label}
                  </Button>
                );
              })}
            </Box>
          )}

          {/* Right Action Area */}
          <Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
            {loading ? (
              // Loading state placeholder
              <Box sx={{ width: 100, height: 36, borderRadius: 100, bgcolor: "rgba(43, 38, 32, 0.05)" }} />
            ) : !user ? (
              // Logged out area
              <>
                <Button
                  component={Link}
                  href="/sign-in"
                  variant="text"
                  sx={{
                    color: "text.primary",
                    fontWeight: 600,
                    px: 3,
                    py: 1,
                  }}
                >
                  Sign in
                </Button>
                <Button
                  component={Link}
                  href="/sign-up"
                  variant="contained"
                  color="primary"
                  endIcon={<ArrowForwardIcon />}
                  sx={{
                    fontWeight: 600,
                    boxShadow: "0 4px 14px rgba(214, 104, 56, 0.2)",
                    px: 3.5,
                    py: 1.2,
                  }}
                >
                  Start comeback
                </Button>
              </>
            ) : (
              // Logged in user profile menu
              <>
                {/* Mobile Burger Menu button */}
                <IconButton
                  color="inherit"
                  aria-label="open mobile navigation"
                  edge="start"
                  onClick={handleOpenMobile}
                  sx={{ display: { xs: "flex", md: "none" }, color: "text.primary" }}
                >
                  <MenuIcon />
                </IconButton>

                {/* Profile Trigger */}
                <Button
                  onClick={handleOpenProfile}
                  aria-controls={Boolean(profileAnchor) ? "account-menu" : undefined}
                  aria-haspopup="menu"
                  aria-expanded={Boolean(profileAnchor) ? "true" : undefined}
                  sx={{
                    px: 1.5,
                    py: 0.75,
                    borderRadius: 100,
                    border: "1px solid rgba(43, 38, 32, 0.08)",
                    color: "text.primary",
                    textTransform: "none",
                    gap: 1.5,
                    backgroundColor: Boolean(profileAnchor) ? "rgba(43, 38, 32, 0.04)" : "transparent",
                    "&:hover": {
                      backgroundColor: "rgba(43, 38, 32, 0.04)",
                      borderColor: "rgba(43, 38, 32, 0.15)",
                    },
                  }}
                >
                  <Typography variant="subtitle2" sx={{ fontWeight: 700, display: { xs: "none", sm: "block" } }}>
                    {firstName}
                    {lastInitial ? ` ${lastInitial}.` : ""}
                  </Typography>
                  <Avatar
                    sx={{
                      width: 32,
                      height: 32,
                      bgcolor: "secondary.main",
                      color: "primary.contrastText",
                      fontFamily: "var(--font-bricolage), sans-serif",
                      fontWeight: 700,
                      fontSize: "0.85rem",
                    }}
                  >
                    {getInitials(user.full_name)}
                  </Avatar>
                </Button>
              </>
            )}
          </Box>

          {/* Desktop/Mobile Profile Dropdown Menu */}
          <Menu
            anchorEl={profileAnchor}
            id="account-menu"
            open={Boolean(profileAnchor)}
            onClose={handleCloseProfile}
            onClick={handleCloseProfile}
            transformOrigin={{ horizontal: "right", vertical: "top" }}
            anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
            PaperProps={{
              elevation: 0,
              sx: {
                overflow: "visible",
                filter: "drop-shadow(0px 8px 30px rgba(43, 38, 32, 0.1))",
                mt: 1.5,
                borderRadius: 4,
                width: 280,
                border: "1px solid rgba(43, 38, 32, 0.06)",
                padding: "8px",
              },
            }}
          >
            {/* Header info */}
            {user && (
              <Box sx={{ p: 2, display: "flex", flexDirection: "column", gap: 0.5 }}>
                <Typography variant="body1" sx={{ fontWeight: 800, fontFamily: "var(--font-bricolage)" }}>
                  {user.full_name}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ textTransform: "capitalize" }}>
                  {user.roles?.[0]?.replace(/[_-]+/g, " ") || "Member"}
                </Typography>
              </Box>
            )}
            <Divider sx={{ my: 1 }} />

            {/* Menu Links with role="link" for E2E backward compatibility */}
            {MENU_LINKS.map((link) => (
              <MenuItem
                key={link.href}
                component={Link}
                href={link.href}
                role="link"
                sx={{
                  borderRadius: 2,
                  py: 1.0,
                  display: { xs: "flex", md: "none" }, // Hide on desktop because they are already in the top nav
                }}
              >
                <ListItemIcon aria-hidden="true">
                  <Typography variant="body1">{link.icon}</Typography>
                </ListItemIcon>
                <ListItemText primary={link.label} primaryTypographyProps={{ fontWeight: 500 }} />
              </MenuItem>
            ))}

            <Divider sx={{ my: 1, display: { xs: "block", md: "none" } }} />

            <MenuItem component={Link} href="/profile" sx={{ borderRadius: 2, py: 1.2 }}>
              <ListItemIcon>
                <PersonIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="My Profile" primaryTypographyProps={{ fontWeight: 600 }} />
            </MenuItem>

            <MenuItem component={Link} href="/settings" sx={{ borderRadius: 2, py: 1.2 }}>
              <ListItemIcon>
                <SettingsIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Settings" primaryTypographyProps={{ fontWeight: 600 }} />
            </MenuItem>

            {user?.roles?.includes("admin") && (
              <MenuItem
                component={Link}
                href="/admin"
                role="link"
                sx={{ borderRadius: 2, py: 1.2, color: "secondary.main" }}
              >
                <ListItemIcon aria-hidden="true" sx={{ color: "secondary.main" }}>
                  <Typography variant="subtitle2" sx={{ fontWeight: 800, ml: 0.5 }}>◆</Typography>
                </ListItemIcon>
                <ListItemText primary="Admin" primaryTypographyProps={{ fontWeight: 700 }} />
              </MenuItem>
            )}

            <Divider sx={{ my: 1 }} />

            <MenuItem onClick={handleSignOut} sx={{ borderRadius: 2, py: 1.2, color: "error.main" }}>
              <ListItemIcon sx={{ color: "error.main" }}>
                <LogoutIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText primary="Sign out" primaryTypographyProps={{ fontWeight: 600 }} />
            </MenuItem>
          </Menu>

          {/* Mobile Navigation Dropdown Menu */}
          <Menu
            anchorEl={mobileAnchor}
            id="mobile-nav-menu"
            open={Boolean(mobileAnchor)}
            onClose={handleCloseMobile}
            onClick={handleCloseMobile}
            transformOrigin={{ horizontal: "right", vertical: "top" }}
            anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
            PaperProps={{
              elevation: 0,
              sx: {
                overflow: "visible",
                filter: "drop-shadow(0px 8px 30px rgba(43, 38, 32, 0.1))",
                mt: 1.5,
                borderRadius: 4,
                width: 240,
                border: "1px solid rgba(43, 38, 32, 0.06)",
                padding: "8px",
              },
            }}
          >
            {MENU_LINKS.map((link) => (
              <MenuItem
                key={link.href}
                component={Link}
                href={link.href}
                onClick={handleCloseMobile}
                sx={{ borderRadius: 2, py: 1.2 }}
              >
                <ListItemIcon>
                  <Typography variant="body1">{link.icon}</Typography>
                </ListItemIcon>
                <ListItemText primary={link.label} primaryTypographyProps={{ fontWeight: 600 }} />
              </MenuItem>
            ))}
          </Menu>
        </Toolbar>
      </Container>

      {/* Breadcrumb Navigation Bar */}
      {breadcrumb && breadcrumb.length > 0 && (
        <Box
          sx={{
            borderTop: "1px solid rgba(43, 38, 32, 0.06)",
            background: "rgba(252, 250, 247, 0.5)",
            py: 1.5,
          }}
        >
          <Container maxWidth="lg">
            <Breadcrumbs
              separator={<KeyboardArrowRightIcon sx={{ fontSize: 16, color: "text.disabled" }} />}
              aria-label="breadcrumb"
            >
              {breadcrumb.map((item, index) => {
                const isLast = index === breadcrumb.length - 1;
                return isLast ? (
                  <Typography
                    key={index}
                    variant="body2"
                    aria-current="page"
                    sx={{ color: "text.primary", fontWeight: 700 }}
                  >
                    {item.label}
                  </Typography>
                ) : (
                  <MuiLink
                    key={index}
                    component={Link}
                    href={item.href || "#"}
                    underline="hover"
                    variant="body2"
                    sx={{ color: "text.secondary", fontWeight: 500 }}
                  >
                    {item.label}
                  </MuiLink>
                );
              })}
            </Breadcrumbs>
          </Container>
        </Box>
      )}
    </AppBar>
  );
}
