"use client";

import React, { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth/auth-context";

/** Links shown in the avatar dropdown, matching the Kirmya design. */
const MENU_LINKS: { label: string; href: string; icon: string }[] = [
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
  return (
    <header
      style={{
        position: "sticky",
        top: 0,
        zIndex: 50,
        background: "rgba(251,247,242,0.86)",
        backdropFilter: "blur(10px)",
        borderBottom: "1px solid #EFE7DC",
      }}
    >
      <nav
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: "24px",
          padding: "18px 40px",
          maxWidth: "1240px",
          margin: "0 auto",
        }}
      >
        <Link
          href="/"
          style={{
            fontFamily: "'Bricolage Grotesque',sans-serif",
            fontSize: "24px",
            fontWeight: 800,
            letterSpacing: "-0.02em",
            color: "#2B2620",
          }}
        >
          Kirmya
        </Link>
        <NavAuthArea />
      </nav>
      {breadcrumb && breadcrumb.length > 0 && (
        <div
          style={{
            borderTop: "1px solid #EFE7DC",
            background: "rgba(251,247,242,0.6)",
          }}
        >
          <nav
            aria-label="Breadcrumb"
            style={{
              maxWidth: "1240px",
              margin: "0 auto",
              padding: "12px 40px",
            }}
          >
            <ol
              style={{
                display: "flex",
                alignItems: "center",
                gap: "10px",
                flexWrap: "wrap",
                listStyle: "none",
                margin: 0,
                padding: 0,
              }}
            >
              {breadcrumb.map((item, i) => {
                const isLast = i === breadcrumb.length - 1;
                return (
                  <React.Fragment key={i}>
                    {i > 0 && (
                      <li
                        aria-hidden="true"
                        style={{ color: "#C9BEAD", fontSize: "15px" }}
                      >
                        {"›"}
                      </li>
                    )}
                    <li
                      aria-current={isLast ? "page" : undefined}
                      style={{
                        fontSize: "14px",
                        color: isLast ? "#2B2620" : "#8A8175",
                        fontWeight: isLast ? 600 : 500,
                      }}
                    >
                      {item.href ? (
                        <Link
                          href={item.href}
                          style={{ color: "#8A8175", fontWeight: 500 }}
                        >
                          {item.label}
                        </Link>
                      ) : (
                        item.label
                      )}
                    </li>
                  </React.Fragment>
                );
              })}
            </ol>
          </nav>
        </div>
      )}
    </header>
  );
}

/** Two-letter initials from a full name, for the avatar fallback. */
function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
  return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
}

/**
 * Right-hand nav slot: marketing CTAs when logged out (or while the session is
 * still being restored), and the user's avatar dropdown menu when logged in.
 */
function NavAuthArea() {
  const { user, loading, signOut } = useAuth();
  const router = useRouter();

  if (loading || !user) {
    return (
      <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
        <Link
          href="/sign-in"
          style={{ fontSize: "15px", color: "#2B2620", fontWeight: 500 }}
        >
          Sign in
        </Link>
        <Link
          href="/sign-up"
          style={{
            background: "#C2683C",
            color: "#fff",
            fontSize: "14px",
            fontWeight: 600,
            padding: "11px 22px",
            borderRadius: "100px",
          }}
        >
          Start your comeback
        </Link>
      </div>
    );
  }

  return <ProfileMenu user={user} signOut={signOut} router={router} />;
}

/** A short label for the dropdown subtitle, derived from the user's role. */
function roleLabel(roles: string[]): string {
  const r = roles?.[0];
  if (!r) return "Member";
  return r.replace(/[_-]+/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

/**
 * Avatar button that opens the account dropdown menu (matches the Kirmya
 * design): profile header, primary navigation links, then Settings + Sign out.
 * Closes on outside-click and Escape.
 */
function ProfileMenu({
  user,
  signOut,
  router,
}: {
  user: { full_name: string; roles?: string[] };
  signOut: () => Promise<void>;
  router: ReturnType<typeof useRouter>;
}) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const firstName = user.full_name.trim().split(/\s+/)[0] || "there";
  const lastInitial =
    user.full_name.trim().split(/\s+/).slice(1).join(" ").charAt(0) || "";

  useEffect(() => {
    if (!open) return;
    function onPointer(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  async function handleSignOut() {
    setOpen(false);
    await signOut();
    router.push("/");
  }

  const avatar = (size: number, font: number) => (
    <span
      aria-hidden="true"
      style={{
        flex: "none",
        width: size + "px",
        height: size + "px",
        borderRadius: "50%",
        background: "#4F7C6A",
        color: "#fff",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        fontFamily: "'Bricolage Grotesque',sans-serif",
        fontWeight: 700,
        fontSize: font + "px",
      }}
    >
      {initials(user.full_name)}
    </span>
  );

  return (
    <div ref={ref} style={{ position: "relative" }}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-haspopup="menu"
        aria-expanded={open}
        title={user.full_name}
        style={{
          display: "flex",
          alignItems: "center",
          gap: "10px",
          border: "none",
          background: "transparent",
          cursor: "pointer",
          padding: "4px",
          borderRadius: "100px",
        }}
      >
        <span style={{ fontSize: "15px", color: "#2B2620", fontWeight: 600 }}>
          {firstName}
          {lastInitial ? " " + lastInitial + "." : ""}
        </span>
        {avatar(36, 14)}
      </button>

      {open && (
        <div
          role="menu"
          style={{
            position: "absolute",
            top: "calc(100% + 12px)",
            right: 0,
            width: "288px",
            background: "#fff",
            border: "1px solid #EFE7DC",
            borderRadius: "18px",
            boxShadow: "0 18px 50px rgba(43,38,32,0.16)",
            padding: "10px",
            zIndex: 60,
          }}
        >
          <Link
            href="/profile"
            role="menuitem"
            onClick={() => setOpen(false)}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "12px",
              padding: "12px 12px 14px",
              textDecoration: "none",
            }}
          >
            {avatar(44, 16)}
            <span style={{ minWidth: 0 }}>
              <span
                style={{
                  display: "block",
                  fontFamily: "'Bricolage Grotesque',sans-serif",
                  fontWeight: 800,
                  fontSize: "16px",
                  color: "#2B2620",
                  letterSpacing: "-0.01em",
                }}
              >
                {user.full_name}
              </span>
              <span style={{ display: "block", fontSize: "13px", color: "#8A8175" }}>
                {roleLabel(user.roles ?? [])}
              </span>
            </span>
          </Link>

          <div style={{ height: "1px", background: "#EFE7DC", margin: "4px 6px 8px" }} />

          {MENU_LINKS.map((item) => (
            <MenuRow
              key={item.href}
              href={item.href}
              icon={item.icon}
              label={item.label}
              onSelect={() => setOpen(false)}
            />
          ))}

          <div style={{ height: "1px", background: "#EFE7DC", margin: "8px 6px" }} />

          <MenuRow
            href="/settings"
            icon={"⚙"}
            label="Settings"
            onSelect={() => setOpen(false)}
          />
          <button
            type="button"
            role="menuitem"
            onClick={handleSignOut}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "12px",
              width: "100%",
              border: "none",
              background: "transparent",
              cursor: "pointer",
              padding: "10px 12px",
              borderRadius: "12px",
              fontSize: "15px",
              fontWeight: 600,
              color: "#C2683C",
              textAlign: "left",
            }}
          >
            <span aria-hidden="true" style={{ width: "20px", fontSize: "15px" }}>
              {"⏻"}
            </span>
            Sign out
          </button>
        </div>
      )}
    </div>
  );
}

/** A single hoverable row in the account dropdown. */
function MenuRow({
  href,
  icon,
  label,
  onSelect,
}: {
  href: string;
  icon: string;
  label: string;
  onSelect: () => void;
}) {
  const [hover, setHover] = useState(false);
  return (
    <Link
      href={href}
      role="menuitem"
      onClick={onSelect}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        display: "flex",
        alignItems: "center",
        gap: "12px",
        padding: "10px 12px",
        borderRadius: "12px",
        fontSize: "15px",
        fontWeight: 500,
        color: "#2B2620",
        textDecoration: "none",
        background: hover ? "#F3ECE2" : "transparent",
      }}
    >
      <span aria-hidden="true" style={{ width: "20px", color: "#8A8175", fontSize: "15px" }}>
        {icon}
      </span>
      {label}
    </Link>
  );
}
