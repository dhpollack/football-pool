import { test, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

test.describe("Logout UI", () => {
  test("should show user menu and allow logout with admin user", async ({
    page,
  }) => {
    // Login with admin user (who has a player profile)
    await page.goto("/login");
    await page.fill(
      E2E_CONFIG.SELECTORS.LOGIN.EMAIL,
      E2E_CONFIG.ADMIN_CREDENTIALS.email,
    );
    await page.fill(
      E2E_CONFIG.SELECTORS.LOGIN.PASSWORD,
      E2E_CONFIG.ADMIN_CREDENTIALS.password,
    );
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

    // Wait for navigation to home page
    await expect(page).toHaveURL("/");

    // Wait a moment for the login to complete, then reload to force auth state detection
    await page.waitForTimeout(2000);
    await page.reload();
    await expect(page).toHaveURL("/");

    // Wait for authentication state to update after reload
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.LOGIN.BUTTON),
    ).not.toBeVisible({ timeout: 10000 });

    // Verify user menu appears
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.BUTTON),
    ).toBeVisible({ timeout: 10000 });

    // Verify login button is hidden
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.LOGIN.BUTTON),
    ).not.toBeVisible();

    // Click user menu button
    await page.click(E2E_CONFIG.SELECTORS.USER_MENU.BUTTON);

    // Wait for menu to appear
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.MENU),
    ).toBeVisible();

    // Verify menu contains profile link and logout button
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.PROFILE_LINK),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.LOGOUT_BUTTON),
    ).toBeVisible();

    // Click logout button
    await page.click(E2E_CONFIG.SELECTORS.USER_MENU.LOGOUT_BUTTON);

    // Wait for logout to complete and check if we're redirected
    // Note: The app might stay on home page but show login button
    await page.waitForTimeout(2000);

    // After logout, login button should be visible again
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.BUTTON)).toBeVisible({ timeout: 10000 });

    // Verify the user menu is gone
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.BUTTON),
    ).not.toBeVisible();
  });
});
