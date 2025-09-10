import type { GameRequest } from "../src/services/model";
import { E2E_CONFIG } from "./e2e.config";

/**
 * Utility to create test games for E2E tests using React Query patterns
 */
export const createTestGame = async (page: any, gameData: GameRequest): Promise<number> => {
  try {
    // Use React Query's createGame function directly in browser context
    const result = await page.evaluate(async (gameData) => {
      try {
        // Dynamically import the createGame function from React Query client
        // Use relative path that works in browser context
        const { createGame } = await import('../../src/services/api/games/games');
        
        // Call the createGame function with the game data
        const response = await createGame([gameData]);
        
        if (response && response.length > 0 && response[0].id) {
          console.log("Game created successfully with ID:", response[0].id);
          return response[0].id;
        }
        
        throw new Error("Failed to create test game: No game returned in response");
      } catch (error) {
        console.error("Error creating test game in browser context:", error);
        console.error("Error details:", error.response?.data, error.response?.status, error.config);
        throw error;
      }
    }, gameData);
    
    return result;
  } catch (error) {
    console.error("Error creating test game with React Query:", error);
    throw error;
  }
};

/**
 * Utility to delete test games for cleanup using React Query patterns
 */
export const deleteTestGame = async (page: any, gameId: number): Promise<void> => {
  try {
    // Use React Query's deleteGame function directly in browser context
    await page.evaluate(async (gameId) => {
      try {
        // Dynamically import the deleteGame function from React Query client
        // Use relative path that works in browser context
        const { deleteGame } = await import('../../src/services/api/games/games');
        
        // Call the deleteGame function
        await deleteGame(gameId);
        console.log("Game deleted successfully with ID:", gameId);
      } catch (error) {
        // Handle 404 gracefully - the game might not exist (e.g., mock ID)
        if (error.message && error.message.includes("404")) {
          console.log("Game not found for deletion (may be mock ID):", gameId);
          return;
        }
        console.error("Error deleting test game in browser context:", error);
        throw error;
      }
    }, gameId);
  } catch (error) {
    console.error("Error deleting test game with React Query:", error);
    // Don't throw error for cleanup operations to avoid breaking test cleanup
    console.log("Continuing with test cleanup despite deletion error");
  }
};

/**
 * Creates a standard test game for E2E testing
 */
export const createStandardTestGame = async (page: any): Promise<number> => {
  const testGame: GameRequest = {
    week: 1,
    season: new Date().getFullYear(),
    favorite_team: "Test Favorite Team",
    underdog_team: "Test Underdog Team",
    spread: 3.5,
    start_time: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(), // Tomorrow
  };
  
  return createTestGame(page, testGame);
};

/**
 * Creates multiple test games for scenarios that need multiple games
 */
export const createMultipleTestGames = async (page: any, count: number = 2): Promise<number[]> => {
  const gameIds: number[] = [];
  
  for (let i = 0; i < count; i++) {
    const testGame: GameRequest = {
      week: i + 1,
      season: new Date().getFullYear(),
      favorite_team: `Test Favorite ${i + 1}`,
      underdog_team: `Test Underdog ${i + 1}`,
      spread: 2.5 + i,
      start_time: new Date(Date.now() + (i + 1) * 24 * 60 * 60 * 1000).toISOString(),
    };
    
    const gameId = await createTestGame(page, testGame);
    gameIds.push(gameId);
  }
  
  return gameIds;
};

/**
 * Cleans up all test games
 */
export const cleanupTestGames = async (page: any, gameIds: number[]): Promise<void> => {
  const deletePromises = gameIds.map(id => deleteTestGame(page, id));
  await Promise.allSettled(deletePromises);
};