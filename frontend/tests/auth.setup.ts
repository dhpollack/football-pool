import { test as setup, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

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
setup("authenticate as user", async ({ page, request }) => {
  // Create a non-admin user for testing
  const userEmail = `testuser_${Date.now()}@test.com`;

  // Register a new user
  const registerResponse = await request.post(
    `${E2E_CONFIG.BACKEND_URL}/api/register`,
    {
      data: {
        email: userEmail,
        password: E2E_CONFIG.TEST_USER_TEMPLATE.password,
      },
    },
  );

  // Check if registration was successful
  if (registerResponse.status() !== 201) {
    throw new Error(
      `User registration failed with status: ${registerResponse.status()}`,
    );
  }

  // Login with the new user
  await page.goto("/login");
  await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, userEmail);
  await page.fill(
    E2E_CONFIG.SELECTORS.LOGIN.PASSWORD,
    E2E_CONFIG.TEST_USER_TEMPLATE.password,
  );
  await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

  // Wait for successful login
  await page.waitForURL("/");

  // Save authentication state for user
  await page.context().storageState({ path: userFile });
});
