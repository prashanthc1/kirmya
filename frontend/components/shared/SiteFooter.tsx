import Link from "next/link";

interface FooterLink {
  label: string;
  href: string;
}

interface FooterGroup {
  heading: string;
  links: FooterLink[];
}

interface SiteFooterProps {
  groups?: [FooterGroup, FooterGroup, FooterGroup];
}

const DEFAULT_GROUPS: [FooterGroup, FooterGroup, FooterGroup] = [
  {
    heading: "Candidates",
    links: [
      { label: "Browse jobs", href: "/jobs" },
      { label: "Mentorship", href: "/mentorship" },
      { label: "Resume tools", href: "/resume" },
      { label: "FAQ", href: "/faq" },
    ],
  },
  {
    heading: "Recruiters",
    links: [
      { label: "Post a role", href: "/recruiter" },
      { label: "How it works", href: "/about" },
      { label: "Pricing", href: "/pricing" },
    ],
  },
  {
    heading: "Company",
    links: [
      { label: "About", href: "/about" },
      { label: "Contact", href: "/about" },
      { label: "Help center", href: "/faq" },
    ],
  },
];

export default function SiteFooter({ groups = DEFAULT_GROUPS }: SiteFooterProps) {
  return (
    <footer style={{ background: "#2B2620", color: "#C9C2B8" }}>
      <div
        style={{
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "clamp(48px,5vw,68px) 40px 36px",
          display: "grid",
          gridTemplateColumns: "1.4fr repeat(3,1fr)",
          gap: "40px",
        }}
      >
        <div style={{ minWidth: "220px" }}>
          <div
            style={{
              fontFamily: "'Public Sans',sans-serif",
              fontSize: "24px",
              fontWeight: 800,
              color: "#fff",
              letterSpacing: "-0.02em",
              marginBottom: "14px",
            }}
          >
            Kirmya
          </div>
          <p
            style={{
              fontSize: "15px",
              lineHeight: 1.6,
              color: "#9C958A",
              margin: 0,
              maxWidth: "320px",
            }}
          >
            Kirmya. Sanskrit for{" "}
            <em>the instrument of purposeful action</em>. Built for your
            comeback.
          </p>
        </div>
        {groups.map((group) => (
          <div key={group.heading}>
            <div
              style={{
                fontSize: "13px",
                fontWeight: 700,
                letterSpacing: "0.08em",
                textTransform: "uppercase",
                color: "#fff",
                marginBottom: "16px",
              }}
            >
              {group.heading}
            </div>
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                gap: "11px",
                fontSize: "15px",
              }}
            >
              {group.links.map((link, j) => (
                <Link key={j} href={link.href} style={{ color: "#C9C2B8" }}>
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
      <div
        style={{
          borderTop: "1px solid #3D362F",
          maxWidth: "1240px",
          margin: "0 auto",
          padding: "22px 40px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          gap: "16px",
          flexWrap: "wrap",
          fontSize: "14px",
          color: "#8A8175",
        }}
      >
        <span>© 2026 Kirmya. Built for your comeback.</span>
        <div style={{ display: "flex", gap: "24px" }}>
          <Link href="/" style={{ color: "#8A8175" }}>
            Privacy
          </Link>
          <Link href="/" style={{ color: "#8A8175" }}>
            Terms
          </Link>
        </div>
      </div>
    </footer>
  );
}
