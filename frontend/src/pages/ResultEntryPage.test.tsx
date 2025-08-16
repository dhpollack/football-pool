import { render, screen, fireEvent } from "@testing-library/react";
import { vi } from "vitest";
import ResultEntryPage from "./ResultEntryPage";

describe("ResultEntryPage", () => {
  const mockGames = [
    { id: 1, favorite_team: "Team A", underdog_team: "Team B", spread: 3 },
    { id: 2, favorite_team: "Team C", underdog_team: "Team D", spread: 7 },
  ];

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockGames),
      }),
    ) as vi.Mock;
  });

  it("renders the games", async () => {
    render(<ResultEntryPage />);

    expect(await screen.findByText(/team a/i)).toBeInTheDocument();
    expect(await screen.findByText(/team b/i)).toBeInTheDocument();
    expect(await screen.findByText(/team c/i)).toBeInTheDocument();
    expect(await screen.findByText(/team d/i)).toBeInTheDocument();
  });

  it("submits the results", async () => {
    render(<ResultEntryPage />);

    await screen.findByText(/team a/i);

    const allSpinbuttons = screen.getAllByRole("spinbutton");
    const favoriteScoreInput = allSpinbuttons.find(
      (input) => input.getAttribute("name") === "favorite_score_1",
    );
    const underdogScoreInput = allSpinbuttons.find(
      (input) => input.getAttribute("name") === "underdog_score_1",
    );

    fireEvent.change(favoriteScoreInput, { target: { value: "10" } });
    fireEvent.change(underdogScoreInput, { target: { value: "7" } });

    fireEvent.click(screen.getByText(/save results/i));

    expect(global.fetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/results",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify([
          { game_id: 1, favorite_score: 10, underdog_score: 7 },
          { game_id: 2, favorite_score: 0, underdog_score: 0 },
        ]),
      }),
    );
  });
});
