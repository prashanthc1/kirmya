import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import SiteFooter from "@/components/shared/SiteFooter";

vi.mock("next/link", () => ({
  default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...rest}>
      {children}
    </a>
  ),
}));

describe("SiteFooter", () => {
  it("renders the default column headings", () => {
    render(<SiteFooter />);
    expect(screen.getByText("Candidates")).toBeInTheDocument();
    expect(screen.getByText("Recruiters")).toBeInTheDocument();
    expect(screen.getByText("Company")).toBeInTheDocument();
  });

  it("links a default candidate item to the right route", () => {
    render(<SiteFooter />);
    expect(screen.getByRole("link", { name: "Browse jobs" })).toHaveAttribute("href", "/jobs");
  });

  it("renders custom groups when provided", () => {
    render(
      <SiteFooter
        groups={[
          { heading: "One", links: [{ label: "Alpha", href: "/a" }] },
          { heading: "Two", links: [{ label: "Beta", href: "/b" }] },
          { heading: "Three", links: [{ label: "Gamma", href: "/c" }] },
        ]}
      />,
    );
    expect(screen.getByText("One")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Gamma" })).toHaveAttribute("href", "/c");
    expect(screen.queryByText("Candidates")).not.toBeInTheDocument();
  });
});
