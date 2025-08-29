import { test, expect } from '@playwright/test';
import { E2E_CONFIG } from './e2e.config';

// Test user credentials
const TEST_USER = {
  name: 'Test User',
  email: `testuser-${Date.now()}-${Math.floor(Math.random() * 100000)}@example.com`,
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
    
    // Wait for redirect to login page (successful registration)
    await expect(page).toHaveURL('/login');
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.EMAIL)).toBeVisible();
  });

  test('should handle registration validation errors', async ({ page }) => {
    await page.goto('/register');
    
    // Test with minimal invalid data that passes HTML5 but fails custom validation
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, ' '); // Space only (will be trimmed to empty)
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, 'a'); // Invalid email format
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, 'ab'); // Too short password
    
    // Submit form (this should trigger custom validation)
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Wait for validation to complete
    await page.waitForTimeout(1000);
    
    // Should show custom validation errors
    await expect(page.locator('text=Name is required')).toBeVisible();
    await expect(page.locator('text=Email is invalid')).toBeVisible();
    await expect(page.locator('text=Password must be at least 6 characters')).toBeVisible();
    
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
    
    // Wait for validation to occur
    await page.waitForTimeout(500);
    
    // Should show email validation error
    await expect(page.locator('text=Email is invalid')).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });

  test('should validate password strength', async ({ page }) => {
    await page.goto('/register');
    
    // Fill with weak password
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.NAME, TEST_USER.name);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.EMAIL, TEST_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.REGISTER.PASSWORD, 'weak');
    await page.click(E2E_CONFIG.SELECTORS.REGISTER.SUBMIT);
    
    // Wait for validation to occur
    await page.waitForTimeout(500);
    
    // Should show password validation error
    await expect(page.locator('text=Password must be at least 6 characters')).toBeVisible();
    await expect(page).toHaveURL(/\/register/);
  });
});