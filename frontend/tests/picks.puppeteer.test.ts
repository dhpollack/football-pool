import { test, expect } from "vitest";
import puppeteer, { type Browser } from "puppeteer";

test("should allow a user to submit picks", async () => {
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
    await page.type("input[name='email']", "admin@test.com");
    await page.type("input[name='password']", "adminpassword");
    await page.click("button[type='submit']");

    await page.waitForNavigation({ waitUntil: "networkidle2" });

    console.log("Navigating to picks page...");
    await page.goto("http://localhost:13000/picks");
    console.log("On picks page.");

    await page.waitForSelector("h4");

    await page.waitForSelector("#quick-pick-button");

    console.log("Clicking Quick Pick button...");
    await page.click("#quick-pick-button");
    console.log("Clicked Quick Pick button.");

    console.log("Clicking Submit Picks button...");
    await page.click("#submit-picks-button");
    console.log("Clicked Submit Picks button.");

    await new Promise((resolve) => setTimeout(resolve, 1000));

    const dialogMessage = await page.evaluate(() => {
      return new Promise((resolve) => {
        const originalAlert = window.alert;
        window.alert = (message) => {
          originalAlert(message);
          resolve(message);
        };
      });
    });

    expect(dialogMessage).toBe("Picks submitted successfully!");
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}, 30000);
