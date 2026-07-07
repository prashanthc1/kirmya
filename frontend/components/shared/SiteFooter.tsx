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

export default function SiteFooter({
  groups = DEFAULT_GROUPS,
}: SiteFooterProps) {
  return (
    <footer className="w-full bg-card border-t border-border/40 py-12 md:py-16 mt-auto">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-8 md:gap-12">
          {/* Brand Info */}
          <div className="lg:col-span-2 space-y-4">
            <span className="text-xl font-bold tracking-tight bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 bg-clip-text text-transparent">
              Kirmya
            </span>
            <p className="text-sm text-muted-foreground leading-relaxed max-w-sm">
              Kirmya. Sanskrit for <em>the instrument of purposeful action</em>.
              The AI career operating system designed to guide you through
              transition and come back stronger.
            </p>
          </div>

          {/* Links columns */}
          {groups.map((group) => (
            <div key={group.heading} className="space-y-4">
              <h4 className="text-xs font-semibold text-foreground uppercase tracking-wider">
                {group.heading}
              </h4>
              <ul className="space-y-2">
                {group.links.map((link, idx) => (
                  <li key={idx}>
                    <Link
                      href={link.href}
                      className="text-sm text-muted-foreground hover:text-foreground transition-colors duration-200"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Footer Bottom */}
        <div className="border-t border-border/40 mt-12 pt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-xs text-muted-foreground">
            &copy; {new Date().getFullYear()} Kirmya. Built for your comeback.
          </p>
          <div className="flex gap-6">
            <Link
              href="/"
              className="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Privacy Policy
            </Link>
            <Link
              href="/"
              className="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              Terms of Service
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
