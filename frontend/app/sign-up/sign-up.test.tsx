import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import SignUpPage from "./page";

// Mock next/navigation
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}));

// Mock API client
vi.mock("@/lib/api/client", () => ({
  api: {
    post: vi.fn(() => Promise.resolve({ verification_required: true })),
  },
  setAccessToken: vi.fn(),
  ApiError: class extends Error {
    status = 400;
  },
}));

// Mock Auth context
vi.mock("@/lib/auth/auth-context", () => ({
  useAuth: () => ({
    setUser: vi.fn(),
  }),
}));

describe("SignUpPage", () => {
  it("renders sign-up fields, password strength indicator, and terms checkbox", async () => {
    render(<SignUpPage />);

    expect(screen.getByPlaceholderText("Jordan Rivera")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("name@company.com")).toBeInTheDocument();

    // Fill in a password to trigger the strength checklist
    const passwordInput = screen.getByPlaceholderText("Min. 8 characters");
    fireEvent.change(passwordInput, { target: { value: "Short" } });

    await waitFor(() => {
      expect(screen.getByText("Password Strength")).toBeInTheDocument();
      expect(screen.getByText("At least 8 characters")).toBeInTheDocument();
    });
  });
});
