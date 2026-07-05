import { test, expect } from "@playwright/test";
import { login, openAccountMenu, DEMO } from "./helpers";

test("seeded jobs are listed and searchable", async ({ page }) => {
  await login(page, DEMO.asha);
  await openAccountMenu(page);
  await page.getByRole("link", { name: "Jobs", exact: true }).click();
  await expect(page).toHaveURL(/\/jobs/);

  // The recruiter seeded an Operations Manager role.
  await expect(page.getByText("Operations Manager").first()).toBeVisible();
  await expect(page.getByText("Acme Logistics").first()).toBeVisible();
});
