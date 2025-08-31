import { defineConfig, devices } from "@playwright/test";
import { E2E_CONFIG } from "./tests/e2e.config";

export default defineConfig({
  // Test directory
  testDir: "./tests",

  // Timeout for each test - balanced for reliability and speed
  timeout: 25000,

  // Expect timeout - reduced for faster assertions
  expect: {
    timeout: 3000,
  },

  // Run tests in parallel
  fullyParallel: true,

  // Fail the build on CI if you accidentally left test.only in the source code
  forbidOnly: !!process.env.CI,

  // Retry on CI only
  retries: process.env.CI ? 2 : 0,

  // Opt out of parallel tests on CI
  workers: process.env.CI ? 1 : undefined,

  // Reporter to use
  reporter: process.env.CLAUDECODE ? "line" : "html",

  // Configure web servers to start automatically
  webServer: [
    {
      command: "just -f ../backend/justfile run",
      url: `${E2E_CONFIG.BACKEND_URL}/api/health`,
      timeout: 60000,
      reuseExistingServer: !process.env.CI,
      env: {
        FOOTBALL_POOL_DSN: "file::memory:",
        FOOTBALL_POOL_PORT: "18080",
      },
    },
    {
      command: "npm run dev -- --strictPort",
      url: E2E_CONFIG.FRONTEND_URL,
      timeout: 60000,
      reuseExistingServer: !process.env.CI,
      env: {
        VITE_BACKEND_URL: E2E_CONFIG.BACKEND_URL,
      },
    },
  ],

  // Shared settings for all projects
  use: {
    // Base URL to use in actions like `await page.goto('/')`
    baseURL: E2E_CONFIG.FRONTEND_URL,

    // Collect trace when retrying the failed test
    trace: "on-first-retry",

    // Capture screenshot on failure
    screenshot: "only-on-failure",

    // Video recording
    video: "retain-on-failure",
  },

  // Configure projects for major browsers
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
    },

    // Test against mobile viewports
    {
      name: "Mobile Chrome",
      use: { ...devices["Pixel 5"] },
    },

    // Test in headed mode for debugging
    {
      name: "chromium-headed",
      use: {
        ...devices["Desktop Chrome"],
        headless: false,
      },
    },
  ],

  // No web server configuration - we're starting servers manually
});
