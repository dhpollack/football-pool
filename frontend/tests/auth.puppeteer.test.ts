import { test, expect } from "vitest";
import puppeteer, { type Browser } from "puppeteer";

test("should allow a user to log in", async () => {
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

    await page.goto("http://localhost:13000/login");
    console.log("Navigated to login page");
    await page.type("input[name='email']", "admin@test.com");
    await page.type("input[name='password']", "adminpassword");
    await page.click("button[type='submit']");
    console.log("Clicked submit button");

    await page.waitForSelector("h4", { timeout: 1000 }); // Wait for home page
    console.log("h4 selector found");

    const pageTitle = await page.title();
    expect(pageTitle).toBe("Football Pool");

    console.log("Test completed successfully!");
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}, 30000);
