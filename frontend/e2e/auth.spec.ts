import { test, expect } from "@playwright/test";
import { login, DEMO, DEMO_PASSWORD } from "./helpers";

test("register, log out, and log back in", async ({ page }) => {
  const email = `e2e+${Date.now()}@example.com`;

  // Register a fresh account.
  await page.goto("/register");
  await page.locator("#name").fill("E2E Tester");
  await page.locator("#email").fill(email);
  await page.locator("#password").fill("Password123!");
  await page.getByRole("button", { name: "Sign up" }).click();
  await expect(page).toHaveURL(/\/dashboard/);

  // Log out.
  await page.getByRole("button", { name: "Log out" }).click();
  await expect(page).toHaveURL(/\/login/);

  // Log back in.
  await page.locator("#email").fill(email);
  await page.locator("#password").fill("Password123!");
  await page.getByRole("button", { name: "Log in" }).click();
  await expect(page).toHaveURL(/\/dashboard/);
});

test("demo seeker can log in", async ({ page }) => {
  await login(page, DEMO.asha, DEMO_PASSWORD);
  await expect(page.getByRole("navigation")).toContainText("Dashboard");
});
