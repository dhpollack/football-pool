import { test, expect } from '@playwright/test';
import { E2E_CONFIG } from './e2e.config';

// Test user credentials
const TEST_USER = {
  name: 'Test User',
  email: `testuser-${Date.now()}@example.com`,
  password: 'password123',
};

// Helper function to cleanup test user
async function cleanupTestUser(email: string) {
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
      await fetch(
        `${E2E_CONFIG.BACKEND_URL}/api/debug/users/delete?email=${encodeURIComponent(email)}`,
        {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${adminToken}` },
        }
      );
    }
  } catch (error) {
    console.warn('Cleanup failed:', error);
  }
}

test.describe('Registration', () => {
  test.beforeEach(async () => {
    // Cleanup test user before each test
    await cleanupTestUser(TEST_USER.email);
  });

  test.afterEach(async () => {
    // Cleanup test user after each test
    await cleanupTestUser(TEST_USER.email);
  });

  test('should allow a user to register', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/register');
    
    // Fill registration form
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, TEST_USER.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, TEST_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, TEST_USER.password);
    
    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Wait for redirect to login page
    await expect(page).toHaveURL('/login');
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL)).toBeVisible();
    
    // Verify user was created in backend
    const adminLoginResponse = await fetch(`${E2E_CONFIG.BACKEND_URL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
        password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
      }),
    });
    
    expect(adminLoginResponse.ok).toBeTruthy();
    const adminLoginData = await adminLoginResponse.json();
    const adminToken = adminLoginData.token;
    
    // Fetch users from debug endpoint
    const usersResponse = await fetch(
      `${E2E_CONFIG.BACKEND_URL}/api/debug/users`,
      {
        headers: {
          Authorization: `Bearer ${adminToken}`,
        },
      }
    );
    
    expect(usersResponse.ok).toBeTruthy();
    const users = await usersResponse.json();
    
    // Assert that test user is present
    const userExists = users.some((user: { Email: string }) => 
      user.Email === TEST_USER.email
    );
    expect(userExists).toBe(true);
  });

  test('should handle registration validation errors', async ({ page }) => {
    await page.goto('/register');
    
    // Try to submit empty form
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Should show validation errors
    await expect(page.locator(E2E_CONFIG.SELECTORS.REGISTER.ERROR)).toBeVisible();
    
    const errorElements = await page.locator(E2E_CONFIG.SELECTORS.REGISTER.ERROR).all();
    expect(errorElements.length).toBeGreaterThan(0);
    
    // Should stay on registration page
    await expect(page).toHaveURL(/\/register/);
  });

  test('should validate email format', async ({ page }) => {
    await page.goto('/register');
    
    // Fill with invalid email
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, TEST_USER.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, 'invalid-email');
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, TEST_USER.password);
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Should show validation error
    await expect(page.locator(E2E_CONFIG.SELECTORS.REGISTER.ERROR)).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });

  test('should validate password strength', async ({ page }) => {
    await page.goto('/register');
    
    // Fill with weak password
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, TEST_USER.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, TEST_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, 'weak');
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Should show validation error
    await expect(page.locator(E2E_CONFIG.SELECTORS.REGISTER.ERROR)).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });
});