import { test, expect } from "@playwright/test";

test.describe("Main Page", () => {
  test("should load the main page and display the welcome message", async ({
    page,
  }) => {
    // Navigate to the home page
    await page.goto("/");

    // Wait for the React app to load
    await page.waitForLoadState("networkidle");

    // Check if the main heading is visible
    const mainHeading = page.locator("text=Welcome to the Football Pool!");
    await expect(mainHeading).toBeVisible();

    // Check if the heading text is correct
    await expect(mainHeading).toHaveText("Welcome to the Football Pool!");
  });
});
