# Refactoring Playwright Authentication to Use Worker Fixtures

## 1. Goal

This document outlines a plan to refactor the Playwright E2E test suite's authentication mechanism. The current approach uses a global setup file (`tests/auth.setup.ts`) to generate authentication state files on disk. This can lead to race conditions and failures when running tests with multiple parallel workers.

The goal is to replace this file-based approach with a more robust, in-memory **Worker Fixture**. This will improve reliability, simplify the configuration, and align with modern Playwright best practices.

## 2. Context: The Current Setup

Currently, the authentication process is handled by a combination of three files:

### `playwright.config.ts`

The config file defines a special `setup` project that all other test projects depend on. This forces the setup project to run once before any tests begin.

```typescript
// frontend/playwright.config.ts (snippet)

  // Configure projects for major browsers
  projects: [
    // Setup project - runs authentication setup first
    {
      name: "setup",
      testMatch: /.*\.setup\.ts/,
    },

    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
      dependencies: ["setup"],
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
      dependencies: ["setup"],
    },
    // ... other projects also depend on "setup"
  ],
```

### `tests/auth.setup.ts`

This file is matched by the `setup` project. It contains logic to log in as different users and save their browser state (cookies, `localStorage`) to JSON files. This file has a race condition because the multiple `setup()` calls within it can run in parallel and interfere with each other.

```typescript
// frontend/tests/auth.setup.ts (snippet)

import { test as setup, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

const adminFile = "playwright/.auth/admin.json";
setup("authenticate as admin", async ({ page }) => {
  await page.goto("/login");
  // ... login logic ...
  await page.context().storageState({ path: adminFile });
});

const userFile = "playwright/.auth/user.json";
setup("authenticate as user", async ({ page }) => {
  await page.goto("/login");
  // ... login logic ...
  await page.context().storageState({ path: userFile });
});
```

### `tests/auth-roles.spec.ts`

The actual test files use `test.use({ storageState: '...' })` to load the previously saved authentication files, allowing the tests to start in an already-logged-in state.

```typescript
// frontend/tests/auth-roles.spec.ts (snippet)

import { test, expect } from "@playwright/test";

test.describe("Admin Role Tests", () => {
  test.use({ storageState: "playwright/.auth/admin.json" });

  test("should be authenticated as admin user", async ({ page }) => {
    await page.goto("/");
    expect(page.url()).not.toContain("/login");
  });
});
```

## 3. Detailed Refactoring Plan

This plan will be executed in four steps.

### Step 1: Create the Authentication Fixture

We will create a new file, `frontend/tests/auth.fixture.ts`, to define our worker-scoped fixture. This fixture will be responsible for creating an authenticated page for our tests. It will be parameterized to handle different user roles.

**Action:** Create the file `frontend/tests/auth.fixture.ts` with the following content:

```typescript
import { test as baseTest } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

// Define the types for our fixture options
type AuthFixtureOptions = {
  role: 'admin' | 'user';
};

// Extend the base Playwright test object with our custom fixture
export const test = baseTest.extend<AuthFixtureOptions & { page: void }>(
  {
    // Define the role option with a default value
    role: ['user', { option: true }],

    // Override the built-in 'page' fixture.
    // This new fixture will be worker-scoped, meaning it runs once per worker.
    page: [async ({ browser }, use, workerInfo) => {
      // Get the role from the test options
      const { role } = workerInfo.project.use as AuthFixtureOptions;

      // Create a new browser context and page
      const context = await browser.newContext();
      const page = await context.newPage();

      // Perform login based on the specified role
      await page.goto("/login");

      // For this plan, both 'admin' and 'user' use the same credentials
      // as in the original auth.setup.ts. This can be adjusted later.
      await page.fill(
        E2E_CONFIG.SELECTORS.LOGIN.EMAIL,
        E2E_CONFIG.ADMIN_CREDENTIALS.email,
      );
      await page.fill(
        E2E_CONFIG.SELECTORS.LOGIN.PASSWORD,
        E2E_CONFIG.ADMIN_CREDENTIALS.password,
      );
      await page.click(E2E_CONFIG.SELECTORS.LOGIN.SUBMIT);

      // Wait for login to complete
      await page.waitForURL("/");

      // Provide the authenticated page to the test
      await use(page);

      // Clean up after the test
      await context.close();
    }, { scope: 'worker' }],
  }
);

export { expect } from "@playwright/test";
```

### Step 2: Update Role-Based Tests

Next, we will update `frontend/tests/auth-roles.spec.ts` to use our new fixture. We will remove the `test.use({ storageState: ... })` directive and import our custom `test` object.

**Action:** Modify `frontend/tests/auth-roles.spec.ts` to look like this:

```typescript
// Import the custom test object and expect from our new fixture file
import { test, expect } from "./auth.fixture";

test.describe("Admin Role Tests", () => {
  // Set the role for this group of tests
  test.use({ role: 'admin' });

  test("should be authenticated as admin user", async ({ page }) => {
    await page.goto("/");
    // The fixture ensures we are already logged in, so we shouldn't be on the login page
    expect(page.url()).not.toContain("/login");
    expect(page.url()).toBe("http://localhost:5173/");
  });

  test("should maintain admin session across navigation", async ({ page }) => {
    await page.goto("/profile");
    expect(page.url()).not.toContain("/login");

    await page.goto("/picks");
    expect(page.url()).not.toContain("/login");
  });
});

test.describe("User Role Tests", () => {
  // Set the role for this group of tests
  test.use({ role: 'user' });

  test("should be authenticated as regular user", async ({ page }) => {
    await page.goto("/");
    expect(page.url()).not.toContain("/login");
    expect(page.url()).toBe("http://localhost:5173/");
  });

  test("should maintain user session across navigation", async ({ page }) => {
    await page.goto("/profile");
    expect(page.url()).not.toContain("/login");

    await page.goto("/picks");
    expect(page.url()).not.toContain("/login");
  });
});
```

### Step 3: Update Playwright Configuration

Now we must simplify the `playwright.config.ts` file by removing the `setup` project and its dependencies, as it is no longer needed.

**Action:** Modify the `projects` array in `frontend/playwright.config.ts`:

**_Before:_**
```typescript
  projects: [
    { name: 'setup', testMatch: /.*\.setup\.ts/ },
    { name: 'chromium', use: { ...devices["Desktop Chrome"] }, dependencies: ['setup'] },
    { name: 'firefox', use: { ...devices["Desktop Firefox"] }, dependencies: ['setup'] },
    { name: 'Mobile Chrome', use: { ...devices["Pixel 5"] }, dependencies: ['setup'] },
    // ...
  ],
```

**_After:_**
```typescript
  projects: [
    { name: 'chromium', use: { ...devices["Desktop Chrome"] } },
    { name: 'firefox', use: { ...devices["Desktop Firefox"] } },
    { name: 'Mobile Chrome', use: { ...devices["Pixel 5"] } },
    // ...
  ],
```

### Step 4: Clean Up Old Files

Finally, we can remove the now-redundant files.

**Action:**
1.  Delete the `frontend/tests/auth.setup.ts` file.
2.  Delete the `frontend/playwright/.auth` directory and its contents (`admin.json`, `user.json`).

## 4. Summary of Benefits

- **No Race Conditions**: Each worker manages its own authenticated state in memory, eliminating file-based race conditions.
- **Better Encapsulation**: Authentication logic is encapsulated in a reusable fixture, making tests cleaner and easier to understand.
- **No File I/O**: The setup no longer depends on writing to or reading from the disk.
- **Simplified Config**: The `playwright.config.ts` becomes simpler by removing the setup project dependency.
