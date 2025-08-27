import { test, expect } from '@playwright/test';
import { E2E_CONFIG } from './e2e.config';

// Test user credentials
const ADMIN_USER = {
  email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
  password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
};

test.describe('Authentication', () => {
  test('should allow a user to log in', async ({ page }) => {
    // Navigate to login page
    await page.goto('/login');
    
    // Fill login form
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, ADMIN_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, ADMIN_USER.password);
    
    // Submit form
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    
    // Wait for navigation to home page
    await expect(page).toHaveURL('/');
    
    // Verify we're on home page
    await expect(page.locator('h4')).toBeVisible();
    
    // Verify page title
    await expect(page).toHaveTitle('Football Pool');
  });

  test('should handle invalid login credentials', async ({ page }) => {
    await page.goto('/login');
    
    // Fill with invalid credentials
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, 'nonexistent@example.com');
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, 'wrongpassword');
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    
    // Should show error message and stay on login page
    await expect(page.locator(E2E_CONFIG.SELECTORS.LOGIN.ERROR)).toBeVisible();
    await expect(page).toHaveURL(/\/login/);
    
    const errorMessage = await page.locator(E2E_CONFIG.SELECTORS.LOGIN.ERROR).textContent();
    expect(errorMessage).toBeTruthy();
  });

  test('should redirect authenticated users away from auth pages', async ({ page }) => {
    // First login
    await page.goto('/login');
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, ADMIN_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, ADMIN_USER.password);
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    await expect(page).toHaveURL('/');
    
    // Try to access login page while authenticated
    await page.goto('/login');
    
    // Should be redirected to home page
    await expect(page).toHaveURL('/');
    
    // Try to access register page while authenticated
    await page.goto('/register');
    
    // Should be redirected to home page
    await expect(page).toHaveURL('/');
  });

  test('should maintain session across page navigation', async ({ page }) => {
    // Login
    await page.goto('/login');
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.EMAIL, ADMIN_USER.email);
    await page.fill(E2E_CONFIG.SELECTORS.LOGIN.PASSWORD, ADMIN_USER.password);
    await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);
    await expect(page).toHaveURL('/');
    
    // Navigate to multiple protected routes
    await page.goto('/profile');
    await expect(page.locator('h4')).toBeVisible();
    
    await page.goto('/picks');
    await expect(page.locator('h4')).toBeVisible();
    
    await page.goto('/');
    await expect(page.locator('h4')).toBeVisible();
    
    // Should remain authenticated throughout
    expect(page.url()).not.toContain('/login');
    expect(page.url()).not.toContain('/register');
  });
});