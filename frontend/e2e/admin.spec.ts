import { test, expect } from "@playwright/test";
import { login, openAccountMenu, DEMO } from "./helpers";

test("admin sees the console; non-admin does not", async ({ page }) => {
  // Admin: the account menu exposes the Admin link and the console loads.
  await login(page, DEMO.admin);
  await openAccountMenu(page);
  const adminLink = page.getByRole("link", { name: "Admin", exact: true });
  await expect(adminLink).toBeVisible();
  await adminLink.click();
  await expect(page).toHaveURL(/\/admin/);
  await expect(
    page.getByRole("heading", { name: "Admin Console" }),
  ).toBeVisible();
  await expect(page.getByText("Total").first()).toBeVisible(); // a stat card

  // Users tab lists accounts.
  await page.getByRole("tab", { name: "Users" }).click();
  await expect(page.getByText(DEMO.admin).first()).toBeVisible();
});

test("non-admin is redirected away from /admin", async ({ page }) => {
  await login(page, DEMO.asha);
  await page.goto("/admin");
  await expect(page).toHaveURL(/\/dashboard/);
});
