import { screen, fireEvent, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import ResultEntryPage from "./ResultEntryPage";

// Mock the React Query hooks
const mockSubmitResult = vi.fn();
const mockGamesData = {
  games: [
    { id: 1, favorite_team: "Team A", underdog_team: "Team B", spread: 3 },
    { id: 2, favorite_team: "Team C", underdog_team: "Team D", spread: 7 },
  ],
};

vi.mock("../services/api/default/default", () => ({
  useGetGames: () => ({
    data: mockGamesData,
    isLoading: false,
    error: null,
  }),
  useSubmitResult: () => ({
    mutateAsync: mockSubmitResult,
  }),
}));

describe("ResultEntryPage", () => {
  it("renders the games", async () => {
    const { unmount } = render(<ResultEntryPage />);

    try {
      expect(await screen.findByText(/team a/i)).toBeInTheDocument();
      expect(await screen.findByText(/team b/i)).toBeInTheDocument();
      expect(await screen.findByText(/team c/i)).toBeInTheDocument();
      expect(await screen.findByText(/team d/i)).toBeInTheDocument();
    } finally {
      unmount();
    }
  });

  it("submits the results", async () => {
    mockSubmitResult.mockResolvedValueOnce({});

    render(<ResultEntryPage />);

    // Wait for games to load
    await screen.findByText(/team a/i);
    
    // Get all spinbuttons and find the specific ones
    const allSpinbuttons = screen.getAllByRole("spinbutton");
    const favoriteScoreInput = allSpinbuttons.find(
      (input) => input.getAttribute("name") === "favorite_score_1",
    );
    const underdogScoreInput = allSpinbuttons.find(
      (input) => input.getAttribute("name") === "underdog_score_1",
    );
    
    if (!favoriteScoreInput || !underdogScoreInput) {
      throw new Error("Could not find score inputs");
    }

    fireEvent.change(favoriteScoreInput, { target: { value: "10" } });
    fireEvent.change(underdogScoreInput, { target: { value: "7" } });

    fireEvent.click(screen.getByText(/save results/i));

    await waitFor(() => {
      expect(mockSubmitResult).toHaveBeenCalledWith({
        data: {
          game_id: 1,
          favorite_score: 10,
          underdog_score: 7,
          outcome: "push",
        },
      });
    });
  });
});
