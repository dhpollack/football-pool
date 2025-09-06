import { test, expect } from "@playwright/test";
import { E2E_CONFIG } from "./e2e.config";

test.describe("Admin Game Management", () => {
  test.use({ storageState: "playwright/.auth/admin.json" });

  test.beforeEach(async ({ page }) => {
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_GAMES);
    await expect(page).toHaveURL(E2E_CONFIG.ROUTES.ADMIN_GAMES);
  });

  test("should navigate to admin games page", async ({ page }) => {
    // Verify we're on the admin games page
    await expect(page.locator("h4")).toContainText(/games/i);
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.SEARCH_INPUT),
    ).toBeVisible();
  });

  test("should display games list with information", async ({ page }) => {
    // Check that games are displayed
    const gameRows = page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_ROW);
    await expect(gameRows.first()).toBeVisible();

    // Check that game information is displayed
    const firstGameMatchup = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_MATCHUP)
      .first();
    await expect(firstGameMatchup).toBeVisible();

    const firstGameWeek = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_WEEK)
      .first();
    await expect(firstGameWeek).toBeVisible();

    const firstGameSeason = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_SEASON)
      .first();
    await expect(firstGameSeason).toBeVisible();

    const firstGameStatus = page
      .locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_STATUS)
      .first();
    await expect(firstGameStatus).toBeVisible();
  });

  test("should search for games", async ({ page }) => {
    const searchInput = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.SEARCH_INPUT,
    );

    // Search for games (assuming some games exist)
    await searchInput.fill("team");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should find games containing "team" in matchup
    const gameMatchups = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_MATCHUP,
    );
    const count = await gameMatchups.count();

    for (let i = 0; i < count; i++) {
      const matchupText = await gameMatchups.nth(i).textContent();
      expect(matchupText?.toLowerCase()).toContain("team");
    }
  });

  test("should open create game form", async ({ page }) => {
    const createButton = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.CREATE_BUTTON,
    );
    await createButton.click();

    // Should show game form with all fields
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.FAVORITE_TEAM),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.UNDERDOG_TEAM),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.SPREAD),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.WEEK),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.SEASON),
    ).toBeVisible();
    await expect(
      page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAME_FORM.START_DATE),
    ).toBeVisible();
  });

  test("should filter games by week", async ({ page }) => {
    // Look for week filter controls
    const weekFilter = page.getByRole("button", { name: /week/i });

    // Try to filter by specific week if the control exists
    if (await weekFilter.isVisible()) {
      await weekFilter.click();

      // Select week 1 if available
      const week1Option = page.getByRole("option", { name: "1" });
      if (await week1Option.isVisible()) {
        await week1Option.click();
        await page.waitForLoadState("networkidle");

        // All visible games should be from week 1
        const gameWeeks = page.locator(
          E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_WEEK,
        );
        const count = await gameWeeks.count();

        for (let i = 0; i < count; i++) {
          const weekText = await gameWeeks.nth(i).textContent();
          expect(weekText).toBe("1");
        }
      }
    }
  });

  test("should filter games by season", async ({ page }) => {
    // Look for season filter controls
    const seasonFilter = page.getByRole("button", { name: /season/i });

    // Try to filter by specific season if the control exists
    if (await seasonFilter.isVisible()) {
      await seasonFilter.click();

      // Select current year if available
      const currentYear = new Date().getFullYear().toString();
      const yearOption = page.getByRole("option", { name: currentYear });
      if (await yearOption.isVisible()) {
        await yearOption.click();
        await page.waitForLoadState("networkidle");

        // All visible games should be from current season
        const gameSeasons = page.locator(
          E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_SEASON,
        );
        const count = await gameSeasons.count();

        for (let i = 0; i < count; i++) {
          const seasonText = await gameSeasons.nth(i).textContent();
          expect(seasonText).toBe(currentYear);
        }
      }
    }
  });

  test("should handle empty search results", async ({ page }) => {
    const searchInput = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.SEARCH_INPUT,
    );

    // Search for non-existent game
    await searchInput.fill("nonexistentgame123456");
    await searchInput.press("Enter");

    // Wait for results to load
    await page.waitForLoadState("networkidle");

    // Should show empty state message
    await expect(page.getByText(/no games found/i)).toBeVisible();
  });

  test("should show game actions buttons", async ({ page }) => {
    const gameRows = page.locator(E2E_CONFIG.SELECTORS.ADMIN.GAMES.GAME_ROW);
    await expect(gameRows.first()).toBeVisible();

    // Check that action buttons are available for each game
    const editButtons = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.EDIT_BUTTON,
    );
    await expect(editButtons.first()).toBeVisible();

    const deleteButtons = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.DELETE_BUTTON,
    );
    await expect(deleteButtons.first()).toBeVisible();

    const addResultButtons = page.locator(
      E2E_CONFIG.SELECTORS.ADMIN.GAMES.ADD_RESULT_BUTTON,
    );
    await expect(addResultButtons.first()).toBeVisible();
  });
});

test.describe("Admin Game Management - Access Control", () => {
  test.use({ storageState: "playwright/.auth/user.json" });

  test("should redirect non-admin users from admin games page", async ({
    page,
  }) => {
    // Try to access admin games page as regular user
    await page.goto(E2E_CONFIG.ROUTES.ADMIN_GAMES);

    // Should be redirected away from admin page
    await expect(page).not.toHaveURL(E2E_CONFIG.ROUTES.ADMIN_GAMES);
  });
});
