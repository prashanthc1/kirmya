import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import NetworkPage from "./page";

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn() }),
  usePathname: () => "/network",
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
      if (url.includes("/network/connections")) {
        return Promise.resolve([
          {
            id: "conn-1",
            requester_id: "user-123",
            receiver_id: "user-456",
            status: "accepted",
            origin: "manual_request",
            receiver_name: "Jane Smith",
            receiver_headline: "VP of Product",
          },
        ]);
      }
      if (url.includes("/network/requests/incoming")) {
        return Promise.resolve([]);
      }
      if (url.includes("/search")) {
        return Promise.resolve({
          results: [
            {
              id: "user-456",
              title: "Jane Smith",
              description: "jane@example.com",
              snippet: "VP of Product",
            },
          ],
        });
      }
      return Promise.resolve([]);
    }),
    post: vi.fn(() => Promise.resolve({})),
    put: vi.fn(() => Promise.resolve({})),
    delete: vi.fn(() => Promise.resolve({})),
  },
}));

// Mock Auth Context
vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({ user: { id: "user-123", full_name: "John Doe" } }),
}));

// Mock Notifications
vi.mock("@/components/shared/Notifications", () => ({
  useNotifications: () => ({
    showNotification: vi.fn(),
  }),
}));

describe("NetworkPage", () => {
  it("renders the networking dashboard and list of recommended professionals", async () => {
    render(<NetworkPage />);

    await waitFor(() => {
      expect(screen.getByText("Connect with verified practitioners.")).toBeInTheDocument();
      expect(screen.getByText("Jane Smith")).toBeInTheDocument();
      expect(screen.getByText("VP of Product")).toBeInTheDocument();
    });
  });
});
