import { test, expect } from "vitest";
import puppeteer, { type Browser } from "puppeteer";

test("should allow a user to register", async () => {
  let browser: Browser | undefined;
  try {
    browser = await puppeteer.launch({
      headless: true,
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    });
    const page = await browser.newPage();

    page.on("console", (msg) => {
      console.log(`[Browser Console] ${msg.text()}`);
    });

    // --- Pre-test cleanup: Delete testuser if exists ---
    console.log("Attempting to delete testuser@example.com if it exists...");
    // Login as admin to get a token for the delete request
    await page.goto("http://localhost:13000/login");
    await page.type("input[name='email']", "admin@test.com");
    await page.type("input[name='password']", "adminpassword");
    await page.click("button[type='submit']");
    await page.waitForSelector("h4", { timeout: 1000 }); // Wait for home page

    // Now delete the test user
    await page.evaluate(async () => {
      await fetch(
        "http://localhost:18080/api/debug/users/delete?email=testuser@example.com",
        {
          method: "DELETE",
        },
      );
    });
    console.log("Finished attempting to delete testuser@example.com.");
    // --- End of Pre-test cleanup ---

    // Register a new user
    console.log("Navigating to register page...");
    await page.goto("http://localhost:13000/register");
    await page.type("input[name='name']", "Test User");
    await page.type("input[name='email']", "testuser@example.com");
    await page.type("input[name='password']", "password123");
    await page.click("button[type='submit']");

    // Wait for the login page to appear after registration
    await page.waitForSelector("input[name='email']", { timeout: 1000 });

    // Login as admin to check database
    console.log("Navigating to login page as admin...");
    await page.goto("http://localhost:13000/login");
    await page.type("input[name='email']", "admin@test.com");
    await page.type("input[name='password']", "adminpassword");
    await page.click("button[type='submit']");
    await page.waitForSelector("h4", { timeout: 1000 }); // Wait for home page

    // Fetch users from debug endpoint
    console.log("Fetching users from debug endpoint...");
    const users = await page.evaluate(async () => {
      const response = await fetch("http://localhost:18080/api/debug/users");
      return response.json();
    });
    console.log("Users from debug endpoint:", users);

    // Assert that testuser@example.com is present
    const userExists = users.some(
      (user: any) => user.email === "testuser@example.com",
    );
    expect(userExists).toBe(true);

    console.log("Test completed successfully!");
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}, 30000);
