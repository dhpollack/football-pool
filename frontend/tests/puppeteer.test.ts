
import { test, expect } from 'vitest';
import puppeteer from 'puppeteer';

test('puppeteer should be working', async () => {
  let browser;
  try {
    browser = await puppeteer.launch({
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });
    const page = await browser.newPage();
    await page.goto('about:blank');
    const title = await page.title();
    expect(title).toBe('');
  } finally {
    if (browser) {
      await browser.close();
    }
  }
});
