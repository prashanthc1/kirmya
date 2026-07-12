import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import AdminPage from "./page";

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  usePathname: () => "/admin",
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
    get: vi.fn((url) => {
      if (url.includes("/admin/stats")) {
        return Promise.resolve({
          users: { total: 100, active: 85, new_7d: 10 },
          jobs: { total: 20, applications: 40 },
          referrals: { total: 15, accepted: 5, hired: 2 },
          communities: { total: 5, posts: 12 },
          reports: { open: 1 },
        });
      }
      if (url.includes("/admin/users")) {
        return Promise.resolve({
          users: [
            {
              id: "user-999",
              email: "admin@kirmya.com",
              full_name: "Super Admin User",
              status: "active",
              roles: ["admin"],
            },
          ],
        });
      }
      if (url.includes("/admin/reports")) {
        return Promise.resolve([]);
      }
      return Promise.resolve([]);
    }),
    post: vi.fn(() => Promise.resolve({})),
    patch: vi.fn(() => Promise.resolve({})),
    delete: vi.fn(() => Promise.resolve({})),
  },
  ApiError: class extends Error {
    status = 500;
  },
}));

// Mock Auth Context with required roles list containing admin
vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({
    user: { id: "user-999", full_name: "Super Admin User", roles: ["admin"] },
    loading: false,
  }),
}));

// Mock Notifications
vi.mock("@/components/shared/Notifications", () => ({
  useNotifications: () => ({
    showNotification: vi.fn(),
  }),
}));

describe("AdminPage", () => {
  it("renders the admin panel console sidebar and content tabs", async () => {
    render(<AdminPage />);

    await waitFor(() => {
      expect(screen.getByText("System Admin")).toBeInTheDocument();
      expect(screen.getByText("Total members")).toBeInTheDocument();
      expect(screen.getByText("85")).toBeInTheDocument();
    });
  });
});
