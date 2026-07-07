import { test, expect } from "@playwright/test";
import { login, DEMO } from "./helpers";

test("profile updates are synced with backend and DB", async ({ page }) => {
  page.on("console", (msg) => {
    console.log(`BROWSER CONSOLE [${msg.type()}]:`, msg.text());
  });
  page.on("pageerror", (err) => {
    console.error("BROWSER PAGE ERROR:", err.message);
  });

  await login(page, DEMO.asha);

  // Navigate to profile edit workspace
  await page.goto("/profile/edit");
  await expect(page).toHaveURL(/\/profile\/edit/);

  // Expand the Identity section if not already expanded
  const identityHeader = page.locator("#section-identity");
  await identityHeader.click();

  // Click the Edit inline button
  const editButton = page.getByRole("button", { name: "Edit inline" }).first();
  await editButton.click();

  // Fill in the Professional Headline input
  const newHeadline = `Senior Operations Specialist - ${Date.now()}`;
  const headlineInput = page.locator(
    'label:has-text("Professional Headline") + input',
  );
  await headlineInput.fill(newHeadline);

  // Click Done to finish editing
  const doneButton = page.getByRole("button", { name: "Done" }).first();
  await doneButton.click();

  // Wait for cloud sync pill to say "Synced to cloud"
  const syncPill = page.locator('text="Synced to cloud"');
  await expect(syncPill).toBeVisible({ timeout: 15000 });

  // Reload the page to test database persistence and fetch on reload
  await page.reload();
  await expect(page).toHaveURL(/\/profile\/edit/);

  // Verify the updated headline is still present in the input field
  await identityHeader.click();
  await editButton.click();
  await expect(headlineInput).toHaveValue(newHeadline);
});
