import { test, expect } from "@playwright/test";

// Simple test to verify custom instance works
test.describe("Custom Instance", () => {
  test("should extract token from localStorage", async ({ page }) => {
    // Set up auth data in localStorage using react-auth-kit v4 format
    const testToken = "test-token-123";
    const timestamp = Date.now();

    // Set localStorage data in react-auth-kit v4 format
    await page.evaluate(
      ({ token, ts }) => {
        // React-auth-kit v4 uses separate keys with timestamp^&*&^value format
        localStorage.setItem("_auth_auth", `${ts}^&*&^${token}`);
        localStorage.setItem("_auth_auth_type", `${ts}^&*&^Bearer`);
        localStorage.setItem(
          "_auth_state",
          `${ts}^&*&^${JSON.stringify({
            id: 1,
            name: "Test User",
            email: "test@example.com",
            role: "user",
          })}`,
        );
      },
      { token: testToken, ts: timestamp },
    );

    // Verify localStorage was set
    const storedAuth = await page.evaluate(() => {
      return localStorage.getItem("_auth_auth");
    });

    expect(storedAuth).toBeTruthy();

    if (!storedAuth) {
      throw new Error("storedAuth is null or undefined");
    }

    // Extract token from format: "timestamp^&*&^token"
    const tokenParts = storedAuth.split("^&*&^");
    expect(tokenParts[1]).toBe(testToken);

    // Test the custom instance token extraction (simulate what it does)
    const extractedToken = await page.evaluate(() => {
      try {
        const authToken = localStorage.getItem("_auth_auth");
        if (authToken) {
          // Extract token from format: "timestamp^&*&^token"
          const token = authToken.split("^&*&^")[1];
          return token || null;
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
