import { test, expect } from "vitest";
import puppeteer, { type Browser, type Page } from "puppeteer";

// Test user credentials
const TEST_USER = {
  name: "E2E Test User",
  email: `e2e-test-${Date.now()}@example.com`,
  password: "SecurePassword123!",
};

// Helper functions for authentication flow
export class AuthFlowHelpers {
  static async setupBrowser(): Promise<{ browser: Browser; page: Page }> {
    const browser = await puppeteer.launch({
      headless: true,
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    });
    const page = await browser.newPage();
    
    // Configure page settings
    await page.setViewport({ width: 1280, height: 800 });
    
    // Capture browser console logs
    page.on("console", (msg) => {
      console.log(`[Browser Console] ${msg.text()}`);
    });
    
    return { browser, page };
  }
  
  static async cleanupTestUser(email: string): Promise<void> {
    try {
      // Get admin token for cleanup
      const adminLoginResponse = await fetch("http://localhost:18080/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: "admin@test.com",
          password: "adminpassword",
        }),
      });
      
      if (adminLoginResponse.ok) {
        const adminLoginData = await adminLoginResponse.json();
        const adminToken = adminLoginData.token;
        
        // Delete test user
        await fetch(
          `http://localhost:18080/api/debug/users/delete?email=${encodeURIComponent(email)}`,
          {
            method: "DELETE",
            headers: { Authorization: `Bearer ${adminToken}` },
          }
        );
      }
    } catch (error) {
      console.warn("Cleanup failed:", error);
    }
  }
  
  static async registerUser(page: Page, user: typeof TEST_USER): Promise<void> {
    await page.goto("http://localhost:13000/register");
    
    // Fill registration form
    await page.type("input[name='name']", user.name);
    await page.type("input[name='email']", user.email);
    await page.type("input[name='password']", user.password);
    
    // Submit form
    await page.click("button[type='submit']");
    
    // Wait for redirect to login page with success message
    await page.waitForSelector("input[name='email']", { timeout: 5000 });
    
    // Verify we're on login page
    const pageTitle = await page.title();
    expect(pageTitle).toBe("Football Pool");
    
    // Check for success message in URL state
    const currentUrl = page.url();
    expect(currentUrl).toContain("/login");
  }
  
  static async loginUser(page: Page, email: string, password: string): Promise<void> {
    await page.goto("http://localhost:13000/login");
    
    // Fill login form
    await page.type("input[name='email']", email);
    await page.type("input[name='password']", password);
    
    // Submit form
    await page.click("button[type='submit']");
    
    // Wait for redirect to home page
    await page.waitForSelector("h4", { timeout: 5000 });
    
    // Verify we're on home page
    const pageTitle = await page.title();
    expect(pageTitle).toBe("Football Pool");
    
    // Verify URL is home page
    const currentUrl = page.url();
    expect(currentUrl).toBe("http://localhost:13000/");
  }
  
  static async logoutUser(page: Page): Promise<void> {
    // Assuming there's a logout button in the layout
    await page.click("button[data-testid='logout-button']");
    
    // Wait for redirect to login page
    await page.waitForSelector("input[name='email']", { timeout: 5000 });
    
    // Verify we're on login page
    const currentUrl = page.url();
    expect(currentUrl).toContain("/login");
  }
  
  static async verifyProtectedRouteAccess(page: Page): Promise<void> {
    // Try to access a protected route
    await page.goto("http://localhost:13000/profile");
    
    // Should either be redirected to login or see the profile page
    const currentUrl = page.url();
    
    if (currentUrl.includes("/login")) {
      // Not authenticated - should see login page
      await page.waitForSelector("input[name='email']", { timeout: 3000 });
    } else {
      // Authenticated - should see profile page
      await page.waitForSelector("h4", { timeout: 3000 });
      const heading = await page.$eval("h4", el => el.textContent);
      expect(heading).toContain("Profile");
    }
  }
}

// Main test suite for authentication flow
describe("Authentication Flow E2E Tests", () => {
  let browser: Browser;
  let page: Page;
  
  beforeAll(async () => {
    // Cleanup any existing test user
    await AuthFlowHelpers.cleanupTestUser(TEST_USER.email);
  });
  
  beforeEach(async () => {
    const setup = await AuthFlowHelpers.setupBrowser();
    browser = setup.browser;
    page = setup.page;
  });
  
  afterEach(async () => {
    if (browser) {
      await browser.close();
    }
  });
  
  afterAll(async () => {
    // Final cleanup
    await AuthFlowHelpers.cleanupTestUser(TEST_USER.email);
  });
  
  test("should complete full registration -> login -> logout flow", async () => {
    // 1. Registration
    await AuthFlowHelpers.registerUser(page, TEST_USER);
    
    // 2. Login with newly created account
    await AuthFlowHelpers.loginUser(page, TEST_USER.email, TEST_USER.password);
    
    // 3. Verify access to protected routes
    await page.goto("http://localhost:13000/profile");
    await page.waitForSelector("h4", { timeout: 3000 });
    const profileHeading = await page.$eval("h4", el => el.textContent);
    expect(profileHeading).toContain("Profile");
    
    // 4. Logout
    // Note: This assumes a logout button exists in the UI
    // For now, we'll test logout by clearing cookies and verifying redirect
    await page.deleteCookie({ name: 'token', domain: 'localhost' });
    await page.reload();
    
    // Should be redirected to login
    await page.waitForSelector("input[name='email']", { timeout: 3000 });
    const currentUrl = page.url();
    expect(currentUrl).toContain("/login");
  }, 30000);
  
  test("should handle invalid login credentials", async () => {
    await page.goto("http://localhost:13000/login");
    
    // Fill with invalid credentials
    await page.type("input[name='email']", "nonexistent@example.com");
    await page.type("input[name='password']", "wrongpassword");
    await page.click("button[type='submit']");
    
    // Should show error message and stay on login page
    await page.waitForSelector("[data-testid='error-message']", { timeout: 3000 });
    
    const errorMessage = await page.$eval(
      "[data-testid='error-message']",
      el => el.textContent
    );
    
    expect(errorMessage).toBeTruthy();
    expect(page.url()).toContain("/login");
  }, 15000);
  
  test("should handle registration validation errors", async () => {
    await page.goto("http://localhost:13000/register");
    
    // Try to submit empty form
    await page.click("button[type='submit']");
    
    // Should show validation errors
    await page.waitForSelector(".MuiFormHelperText-root", { timeout: 3000 });
    
    const errorElements = await page.$$(".MuiFormHelperText-root");
    expect(errorElements.length).toBeGreaterThan(0);
    
    // Should stay on registration page
    expect(page.url()).toContain("/register");
  }, 15000);
  
  test("should redirect authenticated users away from auth pages", async () => {
    // First login
    await AuthFlowHelpers.loginUser(page, "admin@test.com", "adminpassword");
    
    // Try to access login page while authenticated
    await page.goto("http://localhost:13000/login");
    
    // Should be redirected to home page
    await page.waitForSelector("h4", { timeout: 3000 });
    expect(page.url()).toBe("http://localhost:13000/");
    
    // Try to access register page while authenticated
    await page.goto("http://localhost:13000/register");
    
    // Should be redirected to home page
    await page.waitForSelector("h4", { timeout: 3000 });
    expect(page.url()).toBe("http://localhost:13000/");
  }, 20000);
  
  test("should maintain session across page navigation", async () => {
    // Login
    await AuthFlowHelpers.loginUser(page, "admin@test.com", "adminpassword");
    
    // Navigate to multiple protected routes
    await page.goto("http://localhost:13000/profile");
    await page.waitForSelector("h4", { timeout: 3000 });
    
    await page.goto("http://localhost:13000/picks");
    await page.waitForSelector("h4", { timeout: 3000 });
    
    await page.goto("http://localhost:13000/");
    await page.waitForSelector("h4", { timeout: 3000 });
    
    // Should remain authenticated throughout
    const currentUrl = page.url();
    expect(currentUrl).not.toContain("/login");
    expect(currentUrl).not.toContain("/register");
  }, 20000);
});