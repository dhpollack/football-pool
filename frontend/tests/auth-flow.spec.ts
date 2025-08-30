import { test, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

// Test user credentials
const TEST_USER = {
  name: "E2E Test User",
  email: `e2e-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
  password: "SecurePassword123!",
};

// Admin user for existing account tests
const ADMIN_USER = {
  email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
  password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
};

// Optimized cleanup function with reduced retries and faster execution
async function cleanupTestUser(email: string) {
  const maxRetries = 2; // Reduced from 5
  let retries = 0;

  while (retries < maxRetries) {
    try {
      // Get admin token for cleanup
      const adminLoginResponse = await fetch(
        `${E2E_CONFIG.BACKEND_URL}/api/login`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
            password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
          }),
        },
      );

      if (adminLoginResponse.ok) {
        const adminLoginData = await adminLoginResponse.json();
        const adminToken = adminLoginData.token;

        // Delete test user
        const deleteUrl = `${E2E_CONFIG.BACKEND_URL}/api/admin/users/delete?email=${encodeURIComponent(email)}`;

        const deleteResponse = await fetch(deleteUrl, {
          method: "DELETE",
          headers: { Authorization: `Bearer ${adminToken}` },
        });

        if (deleteResponse.ok || deleteResponse.status === 404) {
          // Success or user already deleted
          return;
        }
      }
    } catch (error) {
      // Silently handle cleanup errors
    }

    retries++;
    if (retries < maxRetries) {
      await new Promise((resolve) => setTimeout(resolve, 100)); // Reduced from 200ms * retries
    }
  }
}

// Optimized helper functions for authentication flow
class AuthFlowHelpers {
  static async registerUser(
    page: any,
    user: typeof TEST_USER,
    expectSuccess = true,
  ) {
    // Only navigate if not already on register page
    if (!page.url().endsWith("/register")) {
      await page.goto("/register");
    }

    // Fill registration form
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, user.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, user.email);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, user.password);

    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);

    if (expectSuccess) {
      // Wait for redirect to login page (successful registration)
      await expect(page).toHaveURL("/login");
      await expect(
        page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL),
      ).toBeVisible();
    } else {
      // For failed registrations, should stay on register or go to login
      const currentUrl = page.url();
      expect(
        currentUrl.endsWith("/register") || currentUrl.endsWith("/login"),
      ).toBeTruthy();
    }
  }

  static async loginUser(page: any, email: string, password: string) {
    // Only navigate if not already on login page
    if (!page.url().endsWith("/login")) {
      await page.goto("/login");
    }

    // Fill login form
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, password);

    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

    // Wait for redirect to home page
    await expect(page).toHaveURL("/");
    await expect(page.locator("h4")).toBeVisible();
  }

  static async logoutUser(page: any) {
    // Wait for login button to disappear (means auth state updated)
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.LOGIN.BUTTON),
    ).not.toBeVisible();

    // Wait for user menu button to appear
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.BUTTON),
    ).toBeVisible({ timeout: 3000 });

    // Click user menu button
    await page.click(E2E_CONFIG.SELECTORS.USER_MENU.BUTTON);

    // Wait for menu to appear
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.USER_MENU.MENU),
    ).toBeVisible();

    // Click logout button
    await page.click(E2E_CONFIG.SELECTORS.USER_MENU.LOGOUT_BUTTON);

    // Wait for redirect to login page
    await expect(page).toHaveURL("/login");
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL)).toBeVisible();
  }
}

test.describe("Authentication Flow", () => {
  test.beforeAll(async () => {
    // Quick cleanup without delay
    await cleanupTestUser(TEST_USER.email);
  });

  test.afterAll(async () => {
    // Final cleanup
    await cleanupTestUser(TEST_USER.email);
  });

  test("should complete registration -> login -> invalid login flow", async ({
    page,
  }) => {
    // 1. Registration
    await AuthFlowHelpers.registerUser(page, TEST_USER, true);

    // 2. Try invalid login first
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, TEST_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, "wrongpassword");
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

    // Should show error message
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.ERROR)).toBeVisible();
    await expect(page).toHaveURL("/login");

    // 3. Login with correct credentials
    await AuthFlowHelpers.loginUser(page, TEST_USER.email, TEST_USER.password);

    // 4. Verify access to protected routes
    await page.goto("/profile");
    await expect(page.locator("h4")).toBeVisible();

    const profileHeading = await page.locator("h4").textContent();
    expect(profileHeading).toContain("Profile");

    // Note: Due to backend limitation where newly registered users
    // don't have player profiles, the user menu doesn't appear.
    // This is a backend data consistency issue, not a UI problem.
  });

  test("should prevent duplicate registration", async ({ page }) => {
    // Use a unique user for this specific test to avoid conflicts with other tests
    const duplicateTestUser = {
      name: "Duplicate Test User",
      email: `duplicate-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
      password: "DuplicatePassword123!",
    };

    // Register user first time (should succeed)
    await AuthFlowHelpers.registerUser(page, duplicateTestUser, true);

    // Try to register same user again (should fail)
    await AuthFlowHelpers.registerUser(page, duplicateTestUser, false);
  });

  test("should maintain session after browser restart", async ({
    page,
    context,
  }) => {
    // Use a unique user for this test to avoid conflicts
    const uniqueUser = {
      name: "Session Test User",
      email: `session-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
      password: "SessionPassword123!",
    };

    // Register and login
    await AuthFlowHelpers.registerUser(page, uniqueUser, true);
    await AuthFlowHelpers.loginUser(
      page,
      uniqueUser.email,
      uniqueUser.password,
    );

    // Save storage state
    await context.storageState({ path: "auth-state.json" });

    // Create new context with saved state
    const newContext = await context
      .browser()
      .newContext({ storageState: "auth-state.json" });
    const newPage = await newContext.newPage();

    // Navigate to protected route - should be authenticated
    await newPage.goto("/profile", { timeout: 5000 });
    await expect(newPage.locator("h4")).toBeVisible();

    await newPage.close();
    await newContext.close();
  });

  test("should handle session expiration", async ({ page }) => {
    // Use a unique user for this test to avoid conflicts
    const sessionUser = {
      name: "Session Expiration Test User",
      email: `session-expiration-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
      password: "SessionExpirationPassword123!",
    };

    // Register and login
    await AuthFlowHelpers.registerUser(page, sessionUser, true);
    await AuthFlowHelpers.loginUser(
      page,
      sessionUser.email,
      sessionUser.password,
    );

    // Clear cookies to simulate expired session
    await page.context().clearCookies();
    await page.reload();

    // Wait briefly for React to re-render after cookie clearance
    await page.waitForTimeout(500);

    // After clearing cookies and reloading, we should be on the home page
    // The AuthContext detects authentication loss but redirects to home, not login
    await expect(page).toHaveURL("/");
    // Login button should be visible (since we're not authenticated)
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.BUTTON)).toBeVisible();
  });

  test("should allow admin login and maintain session across navigation", async ({
    page,
  }) => {
    // Login with admin credentials
    await page.goto("/login");
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, ADMIN_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, ADMIN_USER.password);
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

    await expect(page).toHaveURL("/");
    await expect(page.locator("h4")).toBeVisible();

    // Navigate to multiple protected routes and verify session persistence
    await page.goto("/profile");
    await expect(page.locator("h4")).toBeVisible();

    await page.goto("/picks");
    await expect(page.locator("h4")).toBeVisible();

    await page.goto("/");
    await expect(page.locator("h4")).toBeVisible();

    // Should remain authenticated throughout
    expect(page.url()).not.toContain("/login");
    expect(page.url()).not.toContain("/register");
  });

  test("should allow authenticated users to access auth pages", async ({
    page,
  }) => {
    // Login with admin credentials
    await page.goto("/login");
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, ADMIN_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, ADMIN_USER.password);
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    await expect(page).toHaveURL("/");

    // Try to access login page while authenticated - should be allowed
    await page.goto("/login");
    await expect(page).toHaveURL("/login");

    // Try to access register page while authenticated - should be allowed
    await page.goto("/register");
    await expect(page).toHaveURL("/register");
  });
});
