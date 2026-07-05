import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

test("search returns seeded results", async ({ page }) => {
  await login(page, DEMO.asha);

  // Go straight to the results page (the Nav box drives the same query).
  await page.goto("/search?q=Operations");

  // A job matching "Operations" should appear.
  await expect(page.getByText("Operations Manager").first()).toBeVisible();
  // The engine indicator shows which backend served it (DB fallback in E2E).
  await expect(page.getByText(/Served by/)).toBeVisible();
});
