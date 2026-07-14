import { Page, expect } from "@playwright/test";

export const DEMO_PASSWORD = "Password123!";
export const DEMO = {
  admin: "admin@demo.kirmya.io",
  asha: "asha.rao@demo.kirmya.io",
};

/** Log in through the UI and wait for the dashboard. */
export async function login(
  page: Page,
  email: string,
  password = DEMO_PASSWORD,
) {
  await page.goto("/sign-in");
  await page.locator("#email").fill(email);
  await page.locator("#password").fill(password);
  await page.getByRole("button", { name: "Sign in", exact: true }).click();
  await expect(page).toHaveURL(/\/dashboard/, { timeout: 30000 });
}

/**
 * Open the avatar account menu in the top nav. The primary navigation links
 * (Jobs, Referrals, Admin, Sign out, …) live inside this dropdown.
 */
export async function openAccountMenu(page: Page) {
  await page.locator('button[aria-haspopup="menu"]').click();
  await expect(page.getByRole("menu")).toBeVisible();
}
