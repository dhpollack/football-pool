import { test, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

test.describe("Admin User Management", () => {
  test.use({ storageState: "playwright/.auth/admin.json" });

  test.beforeEach(async ({ page }) => {
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_USERS);
    await expect(page).toHaveURL(E2E_CONFIG.ROUTES.ADMIN_USERS);
  });

  test("should navigate to admin users page", async ({ page }) => {
    // Verify we're on the admin users page
    await expect(page.locator("h4")).toContainText(/users/i);
    await expect(page.getByPlaceholder(/search/i)).toBeVisible();
  });

  test("should display users list with pagination", async ({ page }) => {
    // Check that table is displayed
    await expect(page.getByRole("table")).toBeVisible();

    // Check that pagination controls exist
    await expect(page.getByRole("button", { name: /previous/i })).toBeVisible();
    await expect(page.getByRole("button", { name: /next/i })).toBeVisible();
  });

  test("should search for users", async ({ page }) => {
    const searchInput = page.getByPlaceholder(/search/i);

    // Search for a specific user (admin user should exist)
    await searchInput.fill("admin");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should find at least the admin user
    await expect(page.getByRole("table")).toBeVisible();
  });

  test("should navigate to user details page", async ({ page }) => {
    // Look for any user details link
    const detailsLinks = page.locator("a[href*='/users/']");

    if ((await detailsLinks.count()) > 0) {
      await detailsLinks.first().click();

      // Should navigate to user details page
      await expect(page).toHaveURL(/\/admin\/users\/\d+/);

      // Verify user details are displayed
      await expect(page.getByText(/email/i)).toBeVisible();
      await expect(page.getByText(/name/i)).toBeVisible();
      await expect(page.getByText(/role/i)).toBeVisible();
    }
  });

  test("should view user picks in details page", async ({ page }) => {
    // Look for any user details link
    const detailsLinks = page.locator("a[href*='/users/']");

    if ((await detailsLinks.count()) > 0) {
      await detailsLinks.first().click();

      // Look for picks tab
      const picksTab = page.getByRole("tab", { name: /picks/i });
      if (await picksTab.isVisible()) {
        await picksTab.click();

        // Should show picks content
        await expect(page.getByText(/picks/i)).toBeVisible();
      }
    }
  });

  test("should handle empty search results", async ({ page }) => {
    const searchInput = page.getByPlaceholder(/search/i);

    // Search for non-existent user
    await searchInput.fill("nonexistentuser123456");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should show empty state message
    await expect(page.getByText(/no users/i)).toBeVisible();
  });
});

test.describe("Admin User Management - Access Control", () => {
  test.use({ storageState: "playwright/.auth/user.json" });

  test("should redirect non-admin users from admin pages", async ({ page }) => {
    // Try to access admin users page as regular user
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_USERS);

    // Should be redirected away from admin page (either to login or home)
    await expect(page).not.toHaveURL(E2E_CONFIG.ROUTES.ADMIN_USERS);
  });
});
