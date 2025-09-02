import { test, expect } from "@playwright/test";

test.describe("Admin Role Tests", () => {
  test.use({ storageState: "playwright/.auth/admin.json" });

  test("should be authenticated as admin user", async ({ page }) => {
    await page.goto("/");
    // Check that we're not redirected to login (authenticated)
    expect(page.url()).not.toContain("/login");
    expect(page.url()).toBe("http://localhost:5173/");
  });

  test("should maintain admin session across navigation", async ({ page }) => {
    await page.goto("/profile");
    // Should be able to access protected routes without redirect
    expect(page.url()).not.toContain("/login");

    await page.goto("/");
    expect(page.url()).not.toContain("/login");

    await page.goto("/picks");
    expect(page.url()).not.toContain("/login");
  });
});

test.describe("User Role Tests", () => {
  test.use({ storageState: "playwright/.auth/user.json" });

  test("should be authenticated as regular user", async ({ page }) => {
    await page.goto("/");
    // Check that we're not redirected to login (authenticated)
    expect(page.url()).not.toContain("/login");
    expect(page.url()).toBe("http://localhost:5173/");
  });

  test("should maintain user session across navigation", async ({ page }) => {
    await page.goto("/profile");
    // Should be able to access protected routes without redirect
    expect(page.url()).not.toContain("/login");

    await page.goto("/");
    expect(page.url()).not.toContain("/login");

    await page.goto("/picks");
    expect(page.url()).not.toContain("/login");
  });
});
