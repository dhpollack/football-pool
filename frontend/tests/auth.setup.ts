import { test as setup, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

// Test user credentials for registration
const TEST_USER = {
  name: "Test User",
  email: `test-user-${Date.now()}@example.com`,
  password: "SecurePassword123!",
};

const adminFile = "playwright/.auth/admin.json";
setup("authenticate as admin", async ({ page }) => {
  // Login with admin credentials
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

  // Wait for successful login and redirect to home
  await page.waitForURL("/");
  await expect(page.locator("h4")).toBeVisible();

  // Save authentication state for admin
  await page.context().storageState({ path: adminFile });
});

const userFile = "playwright/.auth/user.json";
setup("authenticate as user", async ({ page }) => {
  // Use admin credentials for now to test the setup flow
  // We'll implement proper user registration later

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

  // Wait for successful login
  await page.waitForURL("/");

  // Save authentication state for user
  await page.context().storageState({ path: userFile });
});
