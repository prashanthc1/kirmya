import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import SiteNav from "@/components/shared/SiteNav";
import type { AuthUser } from "@/lib/auth/auth-context";

// Mutable auth state shared with the mocked useAuth (hoisted so the vi.mock
// factory can reference it).
const h = vi.hoisted(() => ({
  auth: {
    user: null as AuthUser | null,
    loading: false,
    setUser: () => {},
    refreshUser: async () => {},
    signOut: async () => {},
  },
}));

vi.mock("next/link", () => ({
  default: ({
    href,
    children,
    ...rest
  }: {
    href: string;
    children: React.ReactNode;
  }) => (
    <a href={href} {...rest}>
      {children}
    </a>
  ),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => "/",
}));

vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => h.auth,
  AuthProvider: ({ children }: { children: React.ReactNode }) => children,
}));

const sampleUser: AuthUser = {
  id: "u1",
  email: "ada@example.com",
  full_name: "Ada Lovelace",
  email_verified: true,
  mfa_enabled: false,
  roles: ["job_seeker"],
};

describe("SiteNav", () => {
  beforeEach(() => {
    h.auth.user = null;
    h.auth.loading = false;
  });

  it("renders the brand and the logged-out CTAs when no user", () => {
    render(<SiteNav />);
    expect(screen.getByText("Kirmya")).toBeInTheDocument();
    expect(screen.getByText("Sign in")).toBeInTheDocument();
    expect(screen.getByText("Start comeback")).toBeInTheDocument();
  });

  it("shows a loading skeleton while the session is still loading", () => {
    h.auth.loading = true;
    render(<SiteNav />);
    expect(screen.queryByText("Sign in")).not.toBeInTheDocument();
  });

  it("shows the avatar trigger and opens the account menu when logged in", () => {
    h.auth.user = sampleUser;
    render(<SiteNav />);
    expect(screen.getByText("Ada L.")).toBeInTheDocument(); // name on the trigger
    expect(screen.getByText("AL")).toBeInTheDocument(); // avatar initials
    expect(screen.queryByText("Start your comeback")).not.toBeInTheDocument();

    // Menu is collapsed until the avatar is clicked.
    expect(screen.queryByRole("menuitem", { name: /Sign out/ })).toBeNull();

    fireEvent.click(screen.getByRole("button", { name: /Ada L\./ }));

    expect(
      screen.getByRole("menuitem", { name: /Sign out/ }),
    ).toBeInTheDocument();
    expect(screen.getAllByText("Jobs")[0].closest("a")).toHaveAttribute(
      "href",
      "/jobs",
    );
    expect(screen.getByRole("menuitem", { name: "Settings" })).toHaveAttribute(
      "href",
      "/settings",
    );
  });

  it("does not render a breadcrumb when none is provided", () => {
    render(<SiteNav />);
    expect(screen.queryByLabelText(/breadcrumb/i)).not.toBeInTheDocument();
  });

  it("renders breadcrumb items and marks the last as the current page", () => {
    render(
      <SiteNav
        breadcrumb={[
          { label: "Jobs", href: "/jobs" },
          { label: "Frontend Engineer" },
        ]}
      />,
    );
    expect(screen.getByLabelText(/breadcrumb/i)).toBeInTheDocument();
    expect(screen.getAllByText("Jobs")[0].closest("a")).toHaveAttribute(
      "href",
      "/jobs",
    );
    expect(screen.getByText("Frontend Engineer")).toHaveAttribute(
      "aria-current",
      "page",
    );
  });
});
