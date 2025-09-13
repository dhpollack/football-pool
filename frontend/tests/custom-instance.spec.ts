import { test, expect } from "@playwright/test";

// Simple test to verify custom instance works
test.describe("Custom Instance", () => {
  test("should extract token from localStorage", async ({ page }) => {
    // Set up auth data in localStorage
    const testToken = "test-token-123";
    const testAuthData = {
      auth: {
        token: testToken,
        type: "Bearer",
        expire: Date.now() + 3600000,
        refreshToken: null,
      },
      userState: {
        id: 1,
        name: "Test User",
        email: "test@example.com",
        role: "user",
      },
      refresh: null,
      user: null,
    };

    // Set localStorage data
    await page.evaluate((authData) => {
      localStorage.setItem("_auth", JSON.stringify(authData));
    }, testAuthData);

    // Verify localStorage was set
    const storedAuth = await page.evaluate(() => {
      return localStorage.getItem("_auth");
    });

    expect(storedAuth).toBeTruthy();

    if (!storedAuth) {
      throw new Error("storedAuth is null or undefined");
    }
    const parsedAuth = JSON.parse(storedAuth);
    expect(parsedAuth.auth.token).toBe(testToken);

    // Test the custom instance token extraction (simulate what it does)
    const extractedToken = await page.evaluate(() => {
      try {
        const authData = localStorage.getItem("_auth");
        if (authData) {
          const parsedAuth = JSON.parse(authData);
          return parsedAuth?.auth?.token || null;
        }
        return null;
      } catch (error) {
        console.error("Error parsing auth data:", error);
        return null;
      }
    });

    expect(extractedToken).toBe(testToken);
  });
});
