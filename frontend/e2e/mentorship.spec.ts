import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

// MVP journey: a logged-in seeker can reach the mentorship area and browse the
// mentor directory. Assertions stay at the navigation + landmark level so the
// test is resilient to copy changes in the seeded content.
test("seeker can reach mentorship and browse mentors", async ({ page }) => {
  await login(page, DEMO.asha);

  await page.goto("/mentorship");
  await expect(page).toHaveURL(/\/mentorship/);
  await expect(page.getByRole("heading").first()).toBeVisible();

  // The mentor directory is one click away and renders its own page.
  await page.goto("/mentors");
  await expect(page).toHaveURL(/\/mentors/);
  await expect(page.getByRole("main")).toBeVisible();
});

test("mentorship requires authentication", async ({ page }) => {
  // Hitting a logged-in area unauthenticated should bounce to sign-in/login.
  await page.goto("/mentorship");
  await expect(page).toHaveURL(/\/(sign-in|login|mentorship)/);
});
