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

    // --- Direct backend API calls for setup/teardown ---
    // Obtain admin token directly from backend
    const adminLoginResponse = await fetch("http://localhost:18080/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: "admin@test.com",
        password: "adminpassword",
      }),
    });
    const adminLoginData = await adminLoginResponse.json();
    const adminToken = adminLoginData.token;

    // Delete testuser if exists
    await fetch(
      "http://localhost:18080/api/debug/users/delete?email=testuser@example.com",
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${adminToken}`,
        },
      },
    );
    // --- End of Direct backend API calls ---

    // Register a new user
    await page.goto("http://localhost:13000/register");
    await page.type("input[name='name']", "Test User");
    await page.type("input[name='email']", "testuser@example.com");
    await page.type("input[name='password']", "password123");
    await page.click("button[type='submit']");

    // Wait for the login page to appear after registration
    await page.waitForSelector("input[name='email']", { timeout: 1000 });

    // Fetch users from debug endpoint directly from Node.js environment
    const usersResponse = await fetch(
      `http://localhost:18080/api/debug/users?timestamp=${new Date().getTime()}`,
      {
        headers: {
          Authorization: `Bearer ${adminToken}`,
        },
      },
    );
    const users = await usersResponse.json();

    // Assert that testuser@example.com is present
    const targetEmail = "testuser@example.com";
    const userExists = users.some((user: any) => {
      return user.Email === targetEmail;
    });
    expect(userExists).toBe(true);
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}, 30000);
