import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import CommunitiesPage from "./page";

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => "/communities",
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

// Mock API Client
vi.mock("@/lib/api/client", () => ({
  api: {
    get: vi.fn(() =>
      Promise.resolve([
        {
          id: "1",
          slug: "senior-staff-engineers",
          name: "Senior Staff Engineers",
          description: "Staff+ circle for navigators.",
          category: "Technology",
          member_count: 140,
        },
      ])
    ),
  },
}));

// Mock Auth Context with required full_name property
vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({ user: { id: "user-123", full_name: "John Doe" } }),
}));

// Mock Notifications
vi.mock("@/components/shared/Notifications", () => ({
  useNotifications: () => ({
    showNotification: vi.fn(),
  }),
}));

describe("CommunitiesPage", () => {
  it("fetches and renders the list of communities on mount", async () => {
    render(<CommunitiesPage />);

    await waitFor(() => {
      expect(screen.getByText("Senior Staff Engineers")).toBeInTheDocument();
      expect(screen.getByText("Staff+ circle for navigators.")).toBeInTheDocument();
    });
  });
});
