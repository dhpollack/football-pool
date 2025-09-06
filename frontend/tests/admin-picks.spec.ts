import { test, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

test.describe("Admin Pick Management", () => {
  test.use({ storageState: "playwright/.auth/admin.json" });

  test.beforeEach(async ({ page }) => {
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_PICKS);
    await expect(page).toHaveURL(E2E_CONFIG.ROUTES.ADMIN_PICKS);
  });

  test("should navigate to admin picks page", async ({ page }) => {
    // Verify we're on the admin picks page
    await expect(page.locator("h4")).toContainText(/picks/i);
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.SEARCH_INPUT),
    ).toBeVisible();
  });

  test("should display picks list with information", async ({ page }) => {
    // Check that picks are displayed
    const pickRows = page.locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_ROW);
    await expect(pickRows.first()).toBeVisible({ timeout: 10000 });

    // Check that pick information is displayed
    const firstPickUser = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_USER)
      .first();
    await expect(firstPickUser).toBeVisible();

    const firstPickGame = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_GAME)
      .first();
    await expect(firstPickGame).toBeVisible();

    const firstPickChoice = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_CHOICE)
      .first();
    await expect(firstPickChoice).toBeVisible();

    const firstPickWeek = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_WEEK)
      .first();
    await expect(firstPickWeek).toBeVisible();

    const firstPickSeason = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_SEASON)
      .first();
    await expect(firstPickSeason).toBeVisible();
  });

  test("should search for picks", async ({ page }) => {
    const searchInput = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.SEARCH_INPUT,
    );

    // Search for picks containing specific text
    await searchInput.fill("favorite");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should find picks containing "favorite" in choice
    const pickChoices = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_CHOICE,
    );
    const count = await pickChoices.count();

    if (count > 0) {
      for (let i = 0; i < count; i++) {
        const choiceText = await pickChoices.nth(i).textContent();
        expect(choiceText?.toLowerCase()).toContain("favorite");
      }
    }
  });

  test("should filter picks by user", async ({ page }) => {
    // Look for user filter controls
    const userFilter = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.FILTER_USER,
    );

    // Try to filter by specific user if the control exists
    if (await userFilter.isVisible()) {
      await userFilter.click();

      // Select first user option if available
      const firstUserOption = page.getByRole("option").first();
      if (await firstUserOption.isVisible()) {
        const userName = await firstUserOption.textContent();
        await firstUserOption.click();
        await page.waitForLoadState("networkidle");

        // All visible picks should be from the selected user
        const pickUsers = page.locator(
          E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_USER,
        );
        const count = await pickUsers.count();

        for (let i = 0; i < count; i++) {
          const userText = await pickUsers.nth(i).textContent();
          expect(userText).toBe(userName);
        }
      }
    }
  });

  test("should filter picks by week", async ({ page }) => {
    // Look for week filter controls
    const weekFilter = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.FILTER_WEEK,
    );

    // Try to filter by specific week if the control exists
    if (await weekFilter.isVisible()) {
      await weekFilter.click();

      // Select week 1 if available
      const week1Option = page.getByRole("option", { name: "1" });
      if (await week1Option.isVisible()) {
        await week1Option.click();
        await page.waitForLoadState("networkidle");

        // All visible picks should be from week 1
        const pickWeeks = page.locator(
          E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_WEEK,
        );
        const count = await pickWeeks.count();

        for (let i = 0; i < count; i++) {
          const weekText = await pickWeeks.nth(i).textContent();
          expect(weekText).toBe("1");
        }
      }
    }
  });

  test("should filter picks by season", async ({ page }) => {
    // Look for season filter controls
    const seasonFilter = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.FILTER_SEASON,
    );

    // Try to filter by specific season if the control exists
    if (await seasonFilter.isVisible()) {
      await seasonFilter.click();

      // Select current year if available
      const currentYear = new Date().getFullYear().toString();
      const yearOption = page.getByRole("option", { name: currentYear });
      if (await yearOption.isVisible()) {
        await yearOption.click();
        await page.waitForLoadState("networkidle");

        // All visible picks should be from current season
        const pickSeasons = page.locator(
          E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_SEASON,
        );
        const count = await pickSeasons.count();

        for (let i = 0; i < count; i++) {
          const seasonText = await pickSeasons.nth(i).textContent();
          expect(seasonText).toBe(currentYear);
        }
      }
    }
  });

  test("should handle empty search results", async ({ page }) => {
    const searchInput = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.SEARCH_INPUT,
    );

    // Search for non-existent pick
    await searchInput.fill("nonexistentpick123456");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should show empty state message
    await expect(page.getByText(/no picks found/i)).toBeVisible();
  });

  test("should show delete buttons for picks", async ({ page }) => {
    const pickRows = page.locator(E2E_CONFIG.SELECTORS.ADMIN.PICKS.PICK_ROW);
    await expect(pickRows.first()).toBeVisible({ timeout: 10000 });

    // Check that delete buttons are available for each pick
    const deleteButtons = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.PICKS.DELETE_BUTTON,
    );
    await expect(deleteButtons.first()).toBeVisible();
  });

  test("should navigate to weekly picks view", async ({ page }) => {
    // Look for weekly picks navigation
    const weeklyLink = page.getByRole("link", { name: /weekly/i });

    if (await weeklyLink.isVisible()) {
      await weeklyLink.click();

      // Should navigate to weekly picks view
      await expect(page).toHaveURL(/\/admin\/picks\/week/);
      await expect(page.locator("h4")).toContainText(/week/i);
    }
  });

  test("should navigate to user picks view", async ({ page }) => {
    // Look for user picks navigation
    const userPicksLink = page.getByRole("link", { name: /by user/i });

    if (await userPicksLink.isVisible()) {
      await userPicksLink.click();

      // Should navigate to user picks view
      await expect(page).toHaveURL(/\/admin\/picks\/user/);
      await expect(page.locator("h4")).toContainText(/user/i);
    }
  });
});

test.describe("Admin Pick Management - Access Control", () => {
  test.use({ storageState: "playwright/.auth/user.json" });

  test("should redirect non-admin users from admin picks page", async ({
    page,
  }) => {
    // Try to access admin picks page as regular user
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_PICKS);

    // Should be redirected away from admin page
    await expect(page).not.toHaveURL(E2E_CONFIG.ROUTES.ADMIN_PICKS);
  });
});
