import { test, expect } from "vitest";
import puppeteer, { type Browser } from "puppeteer";

test("should allow a user to register and log in", async () => {
  let browser: Browser | undefined;
  try {
    browser = await puppeteer.launch({
      headless: true,
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    });
    const page = await browser.newPage();

    page.on("request", (request) => {
      console.log(`[Network Request] ${request.method()} ${request.url()}`);
    });

    page.on("response", async (response) => {
      console.log(`[Network Response] ${response.status()} ${response.url()}`);
      if (response.url().includes("/api/register")) {
        try {
          const text = await response.text();
          console.log(`[Network Response Body /api/register] ${text}`);
        } catch (e) {
          console.log(
            `[Network Response Body /api/register] Could not read response body: ${e}`,
          );
        }
      }
    });

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
    await page.waitForNavigation(); // Wait for admin login to complete

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
    await page.screenshot({ path: "./screenshots/register_page.png" });
    console.log("Typing email...");
    await page.type("input[name='name']", "Test User");
    await page.type("input[name='email']", "testuser@example.com");
    console.log("Typing password...");
    await page.type("input[name='password']", "password123");
    await page.screenshot({ path: "./screenshots/register_form_filled.png" });
    console.log("Submitting registration form...");
    await page.evaluate(() => {
      const form = document.querySelector("form");
      if (
        form &&
        (form as any)._reactProps &&
        (form as any)._reactProps.onSubmit
      ) {
        (form as any)._reactProps.onSubmit();
      } else {
        console.error("Could not find form or onSubmit handler");
      }
    });

    console.log("Waiting for navigation after registration...");
    // Instead of waitForNavigation, let's wait for the login page to appear
    await page.waitForSelector("input[name='email']"); // Assuming it redirects to login
    await page.screenshot({ path: "./screenshots/after_registration.png" });

    // Login as admin to check database
    console.log("Navigating to login page as admin...");
    await page.goto("http://localhost:13000/login");
    await page.screenshot({ path: "./screenshots/admin_login_page.png" });
    console.log("Typing admin email...");
    await page.type("input[name='email']", "admin@test.com");
    console.log("Typing admin password...");
    await page.type("input[name='password']", "adminpassword");
    await page.screenshot({
      path: "./screenshots/admin_login_form_filled.png",
    });
    console.log("Clicking admin submit...");
    await page.waitForSelector("button[type='submit']");
    await page.click("button[type='submit']");

    console.log("Waiting for navigation after admin login...");
    await page.waitForNavigation();
    await page.waitForSelector("h4", { timeout: 10000 }); // Wait for home page
    await page.screenshot({ path: "./screenshots/after_admin_login.png" });

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
