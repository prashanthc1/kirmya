import { Page, expect } from "@playwright/test";

export const DEMO_PASSWORD = "Password123!";
export const DEMO = {
  admin: "admin@demo.kirmya.io",
  asha: "asha.rao@demo.kirmya.io",
};

/** Log in through the UI and wait for the dashboard. */
export async function login(page: Page, email: string, password = DEMO_PASSWORD) {
  await page.goto("/login");
  await page.locator("#email").fill(email);
  await page.locator("#password").fill(password);
  await page.getByRole("button", { name: "Log in" }).click();
  await expect(page).toHaveURL(/\/dashboard/);
}
