import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

test("seeker sees their seeded referral request", async ({ page }) => {
  await login(page, DEMO.asha);
  await page.getByRole("link", { name: "Referrals", exact: true }).click();
  await expect(page).toHaveURL(/\/referrals/);

  // Asha has a seeded outgoing referral directed at BuildCo.
  await expect(page.getByText("Your requests")).toBeVisible();
  await expect(page.getByText("BuildCo").first()).toBeVisible();
});
