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

// Mock connections hooks
vi.mock("@/hooks/useConnections", () => ({
  useConnectionsStore: (cb: any) => cb({
    statusOverrides: {},
    setStatusOverride: vi.fn(),
  }),
  useConnections: () => ({
    data: [
      {
        id: "conn-1",
        user_a_id: "user-123",
        user_b_id: "user-456",
        status: "accepted",
        requested_by: "user-123",
        created_at: "2026-07-06T00:00:00Z",
        updated_at: "2026-07-06T00:00:00Z",
        user: {
          id: "user-456",
          name: "Jane Smith",
          headline: "VP of Product",
          avatar_url: "",
        },
      },
    ],
    isLoading: false,
  }),
  usePendingRequests: () => ({
    data: [],
    isLoading: false,
  }),
  useSuggestions: () => ({
    data: [
      {
        user: {
          id: "user-789",
          name: "Alice Johnson",
          headline: "Software Engineer",
          avatar_url: "",
        },
        mutual_connection_count: 2,
        reason: "Similar industry",
      },
    ],
    isLoading: false,
  }),
  useSendConnectionRequest: () => ({ mutate: vi.fn(), isPending: false }),
  useAcceptConnection: () => ({ mutate: vi.fn(), isPending: false }),
  useDeclineConnection: () => ({ mutate: vi.fn(), isPending: false }),
  useRemoveConnection: () => ({ mutate: vi.fn(), isPending: false }),
  useBlockUser: () => ({ mutate: vi.fn(), isPending: false }),
  useUnblockUser: () => ({ mutate: vi.fn(), isPending: false }),
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
