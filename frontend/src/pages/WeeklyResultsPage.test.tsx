import { render, screen } from "@testing-library/react";
import { vi } from "vitest";
import WeeklyResultsPage from "./WeeklyResultsPage";

describe("WeeklyResultsPage", () => {
  const mockResults = [
    { player_name: "Player 1", score: 100 },
    { player_name: "Player 2", score: 90 },
  ];

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockResults),
      }),
    ) as vi.Mock;
  });

  it("renders the error state", async () => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: false,
        json: () =>
          Promise.resolve({ message: "Failed to fetch weekly results" }),
      }),
    ) as vi.Mock;

    render(<WeeklyResultsPage />);

    expect(
      await screen.findByText(/failed to fetch weekly results/i),
    ).toBeInTheDocument();
  });

  it("renders the results", async () => {
    render(<WeeklyResultsPage />);

    expect(await screen.findByText(/player 1/i)).toBeInTheDocument();
    expect(await screen.findByText(/100/i)).toBeInTheDocument();
    expect(await screen.findByText(/player 2/i)).toBeInTheDocument();
    expect(await screen.findByText(/90/i)).toBeInTheDocument();
  });
});
