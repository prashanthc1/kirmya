import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

test("admin sees the console; non-admin does not", async ({ page }) => {
  // Admin: the nav link is visible and the console loads with analytics.
  await login(page, DEMO.admin);
  const adminLink = page.getByRole("link", { name: "Admin" });
  await expect(adminLink).toBeVisible();
  await adminLink.click();
  await expect(page).toHaveURL(/\/admin/);
  await expect(page.getByRole("heading", { name: "Admin Console" })).toBeVisible();
  await expect(page.getByText("Total").first()).toBeVisible(); // a stat card

  // Users tab lists accounts.
  await page.getByRole("button", { name: "Users" }).click();
  await expect(page.getByText(DEMO.admin).first()).toBeVisible();
});

test("non-admin is redirected away from /admin", async ({ page }) => {
  await login(page, DEMO.asha);
  await page.goto("/admin");
  await expect(page).toHaveURL(/\/dashboard/);
});
