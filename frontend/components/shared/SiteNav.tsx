"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";
import { useTheme } from "next-themes";
import { 
  Sun, 
  Moon, 
  Laptop, 
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
  Sparkles, 
  Compass,
  Bell
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const MENU_LINKS = [
  { label: "Jobs", href: "/jobs", icon: Briefcase },
  { label: "Referrals", href: "/referrals", icon: Compass },
  { label: "Mentorship", href: "/mentorship", icon: GraduationCap },
  { label: "Communities", href: "/communities", icon: Users },
  { label: "AI Coach", href: "/coach", icon: Sparkles },
  { label: "Resume", href: "/resume", icon: FileText },
];

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

interface SiteNavProps {
  breadcrumb?: BreadcrumbItem[];
}

export default function SiteNav({ breadcrumb }: SiteNavProps) {
  const { user, signOut } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const { theme, setTheme, resolvedTheme } = useTheme();
  
  const [mounted, setMounted] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const handleSignOut = async () => {
    setProfileMenuOpen(false);
    await signOut();
    router.push("/");
  };

  const getInitials = (name: string) => {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return "?";
    if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
    return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
  };

  const currentTheme = mounted ? theme : "system";

  return (
    <header className="sticky top-0 z-50 w-full glass-nav transition-all duration-300">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Left Side: Brand Logo and Breadcrumbs */}
          <div className="flex items-center gap-6">
            <Link href={user ? "/dashboard" : "/"} className="flex items-center gap-2">
              <span className="text-xl font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 bg-clip-text text-transparent">
                Kirmya
              </span>
            </Link>

            {/* Breadcrumb section */}
            {breadcrumb && breadcrumb.length > 0 && (
              <nav className="hidden sm:flex items-center space-x-1 text-sm text-muted-foreground">
                <ChevronRight className="h-4 w-4" />
                {breadcrumb.map((item, idx) => (
                  <React.Fragment key={idx}>
                    {idx > 0 && <ChevronRight className="h-4 w-4" />}
                    {item.href ? (
                      <Link href={item.href} className="hover:text-foreground transition-colors font-medium">
                        {item.label}
                      </Link>
                    ) : (
                      <span className="text-foreground font-semibold">{item.label}</span>
                    )}
                  </React.Fragment>
                ))}
              </nav>
            )}
          </div>

          {/* Desktop Navigation Center */}
          <nav className="hidden md:flex items-center space-x-1">
            {user && MENU_LINKS.map((link) => {
              const Icon = link.icon;
              const isActive = pathname.startsWith(link.href);
              return (
                <Link
                  key={link.href}
                  href={link.href}
                  className={`flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-sm font-medium transition-all duration-200 ${
                    isActive 
                      ? "bg-primary text-primary-foreground shadow-sm shadow-blue-500/10" 
                      : "text-muted-foreground hover:bg-secondary hover:text-foreground"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {link.label}
                </Link>
              );
            })}
          </nav>

          {/* Right Side: Theme toggler, Inbox, User profile / Auth buttons */}
          <div className="flex items-center gap-3">
            {/* Theme Selector widget */}
            {mounted && (
              <div className="flex items-center gap-1 bg-secondary border border-border/40 p-1 rounded-full">
                <button
                  onClick={() => setTheme("light")}
                  className={`p-1.5 rounded-full transition-all duration-200 ${
                    theme === "light" 
                      ? "bg-background text-foreground shadow-sm" 
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                  title="Light mode"
                >
                  <Sun className="h-3.5 w-3.5" />
                </button>
                <button
                  onClick={() => setTheme("dark")}
                  className={`p-1.5 rounded-full transition-all duration-200 ${
                    theme === "dark" 
                      ? "bg-background text-foreground shadow-sm" 
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                  title="Dark mode"
                >
                  <Moon className="h-3.5 w-3.5" />
                </button>
                <button
                  onClick={() => setTheme("system")}
                  className={`p-1.5 rounded-full transition-all duration-200 ${
                    theme === "system" 
                      ? "bg-background text-foreground shadow-sm" 
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                  title="System preference"
                >
                  <Laptop className="h-3.5 w-3.5" />
                </button>
              </div>
            )}

            {user ? (
              <div className="flex items-center gap-2">
                {/* Inbox Quick Link */}
                <Link 
                  href="/inbox" 
                  className={`p-2 rounded-full text-muted-foreground hover:text-foreground hover:bg-secondary transition-all ${
                    pathname.startsWith("/inbox") ? "text-foreground bg-secondary" : ""
                  }`}
                  title="Messages"
                >
                  <Bell className="h-4.5 w-4.5" />
                </Link>

                {/* Profile menu toggle dropdown */}
                <div className="relative">
                  <button
                    onClick={() => setProfileMenuOpen(!profileMenuOpen)}
                    className="flex items-center gap-1.5 focus:outline-none"
                  >
                    <div className="h-8 w-8 rounded-full bg-primary/10 border border-primary/20 flex items-center justify-center text-primary font-bold text-sm select-none">
                      {getInitials(user.full_name)}
                    </div>
                  </button>

                  <AnimatePresence>
                    {profileMenuOpen && (
                      <>
                        <div 
                          className="fixed inset-0 z-10" 
                          onClick={() => setProfileMenuOpen(false)}
                        />
                        <motion.div
                          initial={{ opacity: 0, y: 10, scale: 0.95 }}
                          animate={{ opacity: 1, y: 0, scale: 1 }}
                          exit={{ opacity: 0, y: 10, scale: 0.95 }}
                          transition={{ duration: 0.15 }}
                          className="absolute right-0 mt-2 w-56 rounded-2xl border border-border/80 bg-card p-2 text-card-foreground shadow-lg shadow-black/5 ring-1 ring-black/5 focus:outline-none z-20"
                        >
                          <div className="px-3 py-2 text-xs border-b border-border/40 mb-1">
                            <p className="font-semibold text-foreground truncate">{user.full_name}</p>
                            <p className="text-muted-foreground truncate">{user.email}</p>
                          </div>
                          
                          <Link
                            href="/profile"
                            onClick={() => setProfileMenuOpen(false)}
                            className="flex w-full items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-secondary/60 transition-colors"
                          >
                            <User className="h-4 w-4" />
                            My Profile
                          </Link>
                          
                          <Link
                            href="/settings"
                            onClick={() => setProfileMenuOpen(false)}
                            className="flex w-full items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-secondary/60 transition-colors"
                          >
                            <Settings className="h-4 w-4" />
                            Settings
                          </Link>

                          <button
                            onClick={handleSignOut}
                            className="flex w-full items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-destructive hover:bg-destructive/10 transition-colors"
                          >
                            <LogOut className="h-4 w-4" />
                            Log Out
                          </button>
                        </motion.div>
                      </>
                    )}
                  </AnimatePresence>
                </div>
              </div>
            ) : (
              <div className="hidden sm:flex items-center gap-2">
                <Link
                  href="/sign-in"
                  className="px-4 py-1.5 rounded-full text-sm font-semibold text-foreground hover:bg-secondary transition-colors"
                >
                  Sign In
                </Link>
                <Link
                  href="/sign-up"
                  className="px-4 py-1.5 rounded-full text-sm font-semibold bg-primary text-primary-foreground hover:bg-primary/95 transition-all shadow-sm hover:shadow"
                >
                  Join Kirmya
                </Link>
              </div>
            )}

            {/* Mobile Hamburger menu toggle */}
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="p-2 rounded-full text-muted-foreground hover:text-foreground hover:bg-secondary md:hidden"
            >
              {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Menu Panel */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="md:hidden border-t border-border/40 bg-background/95 backdrop-blur-lg overflow-hidden"
          >
            <div className="space-y-1 px-4 py-3 pb-4">
              {user ? (
                <>
                  <div className="px-3 py-2 border-b border-border/40 mb-2">
                    <p className="text-sm font-bold text-foreground truncate">{user.full_name}</p>
                    <p className="text-xs text-muted-foreground truncate">{user.email}</p>
                  </div>
                  
                  {MENU_LINKS.map((link) => {
                    const Icon = link.icon;
                    return (
                      <Link
                        key={link.href}
                        href={link.href}
                        onClick={() => setMobileMenuOpen(false)}
                        className="flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-secondary"
                      >
                        <Icon className="h-4.5 w-4.5" />
                        {link.label}
                      </Link>
                    );
                  })}
                  
                  <div className="border-t border-border/40 my-2 pt-2" />
                  
                  <Link
                    href="/profile"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-secondary"
                  >
                    <User className="h-4.5 w-4.5" />
                    My Profile
                  </Link>
                  <Link
                    href="/settings"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-secondary"
                  >
                    <Settings className="h-4.5 w-4.5" />
                    Settings
                  </Link>
                  <button
                    onClick={handleSignOut}
                    className="flex w-full items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium text-destructive hover:bg-destructive/10"
                  >
                    <LogOut className="h-4.5 w-4.5" />
                    Log Out
                  </button>
                </>
              ) : (
                <div className="grid grid-cols-2 gap-2 pt-2">
                  <Link
                    href="/sign-in"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex justify-center items-center py-2.5 rounded-xl text-sm font-semibold border border-border hover:bg-secondary"
                  >
                    Sign In
                  </Link>
                  <Link
                    href="/sign-up"
                    onClick={() => setMobileMenuOpen(false)}
                    className="flex justify-center items-center py-2.5 rounded-xl text-sm font-semibold bg-primary text-primary-foreground hover:bg-primary/95"
                  >
                    Join Kirmya
                  </Link>
                </div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
}
