import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import OverallResultsPage from "./OverallResultsPage";

// Mock the React Query hooks
vi.mock("../services/api/results/results", () => ({
  useGetSeasonResults: vi.fn().mockReturnValue({
    data: null,
    isLoading: true,
    error: null,
  }),
}));

import { useGetSeasonResults } from "../services/api/results/results";

describe("OverallResultsPage", () => {
  it("renders the loading state", () => {
    render(<OverallResultsPage />);

    expect(screen.getByText(/loading overall results.../i)).toBeInTheDocument();
  });

  it("renders the error state", async () => {
    vi.mocked(useGetSeasonResults).mockReturnValue({
      data: null,
      isLoading: false,
      error: new Error("Failed to fetch overall results"),
    } as any);

    render(<OverallResultsPage />);

    expect(
      await screen.findByText(/failed to fetch overall results/i),
    ).toBeInTheDocument();
  });

  it("renders the results", async () => {
    const mockSeasonResults = [
      { player_name: "Player 1", score: 100 },
      { player_name: "Player 2", score: 90 },
    ];

    vi.mocked(useGetSeasonResults).mockReturnValue({
      data: mockSeasonResults,
      isLoading: false,
      error: null,
    } as any);

    render(<OverallResultsPage />);

    expect(await screen.findByText(/player 1/i)).toBeInTheDocument();
    expect(await screen.findByText(/100/i)).toBeInTheDocument();
    expect(await screen.findByText(/player 2/i)).toBeInTheDocument();
    expect(await screen.findByText(/90/i)).toBeInTheDocument();
  });
});
