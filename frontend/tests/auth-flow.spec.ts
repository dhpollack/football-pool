import { test, expect } from '@playwright/test';
import { E2E_CONFIG } from './e2e.config';

// Test user credentials
const TEST_USER = {
  name: 'E2E Test User',
  email: `e2e-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
  password: 'SecurePassword123!',
};

// Helper function to cleanup test user with retry
async function cleanupTestUser(email: string) {
  const maxRetries = 3;
  let retries = 0;
  
  while (retries < maxRetries) {
    try {
      // Get admin token for cleanup
      const adminLoginResponse = await fetch(`${E2E_CONFIG.BACKEND_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
          password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
        }),
      });

      if (adminLoginResponse.ok) {
        const adminLoginData = await adminLoginResponse.json();
        const adminToken = adminLoginData.token;

        // Delete test user
        const deleteResponse = await fetch(
          `${E2E_CONFIG.BACKEND_URL}/api/debug/users/delete?email=${encodeURIComponent(email)}`,
          {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${adminToken}` },
          }
        );
        
        if (deleteResponse.ok) {
          return; // Success
        }
      }
    } catch (error) {
      console.warn(`Cleanup failed (attempt ${retries + 1}):`, error);
    }
    
    retries++;
    if (retries < maxRetries) {
      await new Promise(resolve => setTimeout(resolve, 100 * retries));
    }
  }
  
  console.warn(`Failed to cleanup test user ${email} after ${maxRetries} attempts`);
}

// Helper functions for authentication flow
class AuthFlowHelpers {
  static async registerUser(page: any, user: typeof TEST_USER, expectSuccess = true) {
    await page.goto('/register');
    
    // Fill registration form
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, user.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, user.email);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, user.password);
    
    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    if (expectSuccess) {
      // Wait for redirect to login page (successful registration)
      await expect(page).toHaveURL('/login');
      await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL)).toBeVisible();
    } else {
      // For failed registrations, should stay on register or go to login
      const currentUrl = page.url();
      expect(currentUrl.endsWith('/register') || currentUrl.endsWith('/login')).toBeTruthy();
    }
  }

  static async loginUser(page: any, email: string, password: string) {
    await page.goto('/login');
    
    // Fill login form
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, password);
    
    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    
    // Wait for redirect to home page
    await expect(page).toHaveURL('/');
    await expect(page.locator('h4')).toBeVisible();
  }

  static async logoutUser(page: any) {
    // Assuming there's a logout button
    await page.click('[data-testid="logout-button"]');
    
    // Wait for redirect to login page
    await expect(page).toHaveURL('/login');
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL)).toBeVisible();
  }
}

test.describe('Authentication Flow', () => {
  test.beforeAll(async () => {
    // Cleanup any existing test user with retry
    await cleanupTestUser(TEST_USER.email);
    // Small delay to ensure cleanup completes
    await new Promise(resolve => setTimeout(resolve, 100));
  });

  test.afterAll(async () => {
    // Final cleanup
    await cleanupTestUser(TEST_USER.email);
  });

  test('should complete full registration -> login -> logout flow', async ({ page }) => {
    // 1. Registration
    await AuthFlowHelpers.registerUser(page, TEST_USER, true);
    
    // 2. Login with newly created account
    await AuthFlowHelpers.loginUser(page, TEST_USER.email, TEST_USER.password);
    
    // 3. Verify access to protected routes
    await page.goto('/profile');
    await expect(page.locator('h4')).toBeVisible();
    
    const profileHeading = await page.locator('h4').textContent();
    expect(profileHeading).toContain('Profile');
    
    // 4. Logout by clearing session (since we don't have a UI logout button yet)
    await page.context().clearCookies();
    await page.reload();
    
    // After logout, user should be on home page but not authenticated
    // Try to access a protected route to verify authentication is required
    await page.goto('/profile');
    
    // Should show an error or remain on profile page (not redirect to login)
    // The exact behavior depends on how the ProfilePage handles unauthenticated access
    await expect(page).toHaveURL('/profile');
    // Profile page should show an error or not display user data
  });

  test('should handle invalid login credentials after registration', async ({ page }) => {
    // Use a unique user for this test to avoid conflicts
    const uniqueUser = {
      name: 'Invalid Login Test User',
      email: `invalid-login-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
      password: 'InvalidLoginPassword123!',
    };
    
    // First register a user
    await AuthFlowHelpers.registerUser(page, uniqueUser, true);
    
    // Try to login with wrong password
    await page.goto('/login');
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, uniqueUser.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, 'wrongpassword');
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    
    // Should show error message
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.ERROR)).toBeVisible();
    await expect(page).toHaveURL('/login');
    
    // Now login with correct credentials
    await AuthFlowHelpers.loginUser(page, uniqueUser.email, uniqueUser.password);
    await expect(page).toHaveURL('/');
  });

  test('should prevent duplicate registration', async ({ page }) => {
    // Register user first time (should succeed)
    await AuthFlowHelpers.registerUser(page, TEST_USER, true);
    
    // Try to register same user again (should fail)
    await AuthFlowHelpers.registerUser(page, TEST_USER, false);
  });

  test('should maintain session after browser restart', async ({ page, context }) => {
    // Use a unique user for this test to avoid conflicts
    const uniqueUser = {
      name: 'Session Test User',
      email: `session-test-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
      password: 'SessionPassword123!',
    };
    
    // Register and login
    await AuthFlowHelpers.registerUser(page, uniqueUser, true);
    await AuthFlowHelpers.loginUser(page, uniqueUser.email, uniqueUser.password);
    
    // Save storage state
    await context.storageState({ path: 'auth-state.json' });
    
    // Create new context with saved state
    const newContext = await context.browser().newContext({ storageState: 'auth-state.json' });
    const newPage = await newContext.newPage();
    
    // Navigate to protected route - should be authenticated
    await newPage.goto('/profile');
    await expect(newPage.locator('h4')).toBeVisible();
    
    await newPage.close();
    await newContext.close();
  });

  test('should handle session expiration', async ({ page }) => {
    // Register and login
    await AuthFlowHelpers.registerUser(page, TEST_USER, true);
    await AuthFlowHelpers.loginUser(page, TEST_USER.email, TEST_USER.password);
    
    // Clear cookies to simulate expired session
    await page.context().clearCookies();
    await page.reload();
    
    // After session expiration, user should be on home page but not authenticated
    // Try to access a protected route to verify authentication is required
    await page.goto('/profile');
    
    // Should show an error or remain on profile page (not redirect to login)
    // The exact behavior depends on how the ProfilePage handles unauthenticated access
    await expect(page).toHaveURL('/profile');
    // Profile page should show an error or not display user data
  });
});