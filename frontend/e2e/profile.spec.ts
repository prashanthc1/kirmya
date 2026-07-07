import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

// MVP journey: a logged-in seeker can view their profile and open the edit
// screen. Kept at navigation + landmark level so seeded copy can change freely.
test("seeker can view profile and open the editor", async ({ page }) => {
  await login(page, DEMO.asha);

  await page.goto("/profile");
  await expect(page).toHaveURL(/\/profile/);
  await expect(page.getByRole("heading").first()).toBeVisible();

  await page.goto("/profile/edit");
  await expect(page).toHaveURL(/\/profile\/edit/);
  await expect(page.getByRole("main")).toBeVisible();
});

test("communities page is reachable while authenticated", async ({ page }) => {
  await login(page, DEMO.asha);
  await page.goto("/communities");
  await expect(page).toHaveURL(/\/communities/);
  await expect(page.getByRole("heading").first()).toBeVisible();
});
