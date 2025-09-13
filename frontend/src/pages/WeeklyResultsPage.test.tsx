import { screen } from "@testing-library/react";
import { vi } from "vitest";
import type { DefinedUseQueryResult } from "@tanstack/react-query";
import { render } from "../test-utils";
import WeeklyResultsPage from "./WeeklyResultsPage";
import type { WeeklyResult } from "../services/api/results/results";

// Mock the React Query hooks
vi.mock("../services/api/results/results", () => ({
  useGetWeeklyResults: vi.fn().mockReturnValue({
    data: null,
    isLoading: false,
    error: null,
  }),
}));

import { useGetWeeklyResults } from "../services/api/results/results";

type MockQueryResult = Partial<DefinedUseQueryResult<WeeklyResult[], Error>>;

describe("WeeklyResultsPage", () => {
  const _mockResults = [
    { player_name: "Player 1", score: 100 },
    { player_name: "Player 2", score: 90 },
  ];

  it("renders the error state", async () => {
    vi.mocked(useGetWeeklyResults).mockReturnValue({
      data: null,
      isLoading: false,
      error: new Error("Failed to fetch weekly results"),
    } as MockQueryResult);

    render(<WeeklyResultsPage />);

    expect(
      await screen.findByText(/failed to fetch weekly results/i),
    ).toBeInTheDocument();
  });

  it("renders the results", async () => {
    const mockWeeklyResults = [
      { player_name: "Player 1", score: 100 },
      { player_name: "Player 2", score: 90 },
    ];

    vi.mocked(useGetWeeklyResults).mockReturnValue({
      data: mockWeeklyResults,
      isLoading: false,
      error: null,
    } as MockQueryResult);

    render(<WeeklyResultsPage />);

    expect(await screen.findByText(/player 1/i)).toBeInTheDocument();
    expect(await screen.findByText(/100/i)).toBeInTheDocument();
    expect(await screen.findByText(/player 2/i)).toBeInTheDocument();
    expect(await screen.findByText(/90/i)).toBeInTheDocument();
  });
});
