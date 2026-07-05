import { test, expect } from "@playwright/test";
import { login, openAccountMenu, DEMO, DEMO_PASSWORD } from "./helpers";

test("register, log out, and log back in", async ({ page }) => {
  const email = `e2e+${Date.now()}@example.com`;

  // Register a fresh account.
  await page.goto("/sign-up");
  await page.locator("#name").fill("E2E Tester");
  await page.locator("#email").fill(email);
  await page.locator("#password").fill("Password123!");
  await page.getByRole("button", { name: "Create one free" }).click();
  await expect(page).toHaveURL(/\/dashboard/);

  // Log out via the avatar account menu; sign-out returns to the home page.
  await openAccountMenu(page);
  await page.getByRole("menuitem", { name: "Sign out" }).click();
  await expect(page).toHaveURL(/\/$/);

  // Log back in.
  await page.goto("/sign-in");
  await page.locator("#email").fill(email);
  await page.locator("#password").fill("Password123!");
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(page).toHaveURL(/\/dashboard/);
});

test("demo seeker can log in", async ({ page }) => {
  await login(page, DEMO.asha, DEMO_PASSWORD);
  // The logged-in nav renders the account menu button for the signed-in user.
  await expect(page.locator('button[aria-haspopup="menu"]')).toBeVisible();
});
