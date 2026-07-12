import React from "react";
import { render, screen, fireEvent, act } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { CookieProvider, useCookieConsent } from "./CookieContext";

// Mock API Client
vi.mock("@/lib/api/client", () => ({
  api: {
    get: vi.fn(() => Promise.resolve(null)),
    post: vi.fn(() => Promise.resolve({ id: "mocked" })),
  },
}));

// Mock Auth Context
vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({ user: null }),
}));

function ConsumerComponent() {
  const { preferences, acceptAll, rejectNonEssential } = useCookieConsent();
  return (
    <div>
      <span data-testid="essential">{preferences.essential ? "yes" : "no"}</span>
      <span data-testid="functional">{preferences.functional ? "yes" : "no"}</span>
      <span data-testid="analytics">{preferences.analytics ? "yes" : "no"}</span>
      <button onClick={acceptAll} data-testid="accept-btn">
        Accept All
      </button>
      <button onClick={rejectNonEssential} data-testid="reject-btn">
        Reject All
      </button>
    </div>
  );
}

describe("CookieProvider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    document.cookie = "kirmya_cookie_preferences=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  });

  it("provides default preferences with essential enabled and others disabled", () => {
    render(
      <CookieProvider>
        <ConsumerComponent />
      </CookieProvider>
    );

    expect(screen.getByTestId("essential").textContent).toBe("yes");
    expect(screen.getByTestId("functional").textContent).toBe("no");
    expect(screen.getByTestId("analytics").textContent).toBe("no");
  });

  it("updates preferences when acceptAll is clicked", async () => {
    render(
      <CookieProvider>
        <ConsumerComponent />
      </CookieProvider>
    );

    const btn = screen.getByTestId("accept-btn");
    await act(async () => {
      fireEvent.click(btn);
    });

    expect(screen.getByTestId("essential").textContent).toBe("yes");
    expect(screen.getByTestId("functional").textContent).toBe("yes");
    expect(screen.getByTestId("analytics").textContent).toBe("yes");
  });
});
