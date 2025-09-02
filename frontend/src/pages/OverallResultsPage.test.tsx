import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import OverallResultsPage from "./OverallResultsPage";

// Mock the React Query hooks
vi.mock("../services/api/default/default", () => ({
  useDebugGetUsers: vi.fn().mockReturnValue({
    data: null,
    isLoading: true,
    error: null,
  }),
}));

import { useDebugGetUsers } from "../services/api/default/default";

describe("OverallResultsPage", () => {
  it("renders the loading state", () => {
    render(<OverallResultsPage />);

    expect(screen.getByText(/loading overall results.../i)).toBeInTheDocument();
  });

  it("renders the error state", async () => {
    vi.mocked(useDebugGetUsers).mockReturnValue({
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
    const mockUsersData = {
      users: [
        {
          id: 1,
          name: "Player 1",
          total_wins: 100,
          pick_count: 10,
          email: "player1@test.com",
          role: "user",
          created_at: "",
          updated_at: "",
        },
        {
          id: 2,
          name: "Player 2",
          total_wins: 90,
          pick_count: 8,
          email: "player2@test.com",
          role: "user",
          created_at: "",
          updated_at: "",
        },
      ],
      pagination: { total: 2, page: 1, per_page: 10, total_pages: 1 },
    };

    vi.mocked(useDebugGetUsers).mockReturnValue({
      data: mockUsersData,
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
