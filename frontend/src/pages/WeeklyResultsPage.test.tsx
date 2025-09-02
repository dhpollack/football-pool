import { screen } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import WeeklyResultsPage from "./WeeklyResultsPage";

// Mock the React Query hooks
vi.mock("../services/api/default/default", () => ({
  useDebugGetUsers: vi.fn().mockReturnValue({
    data: null,
    isLoading: false,
    error: null,
  }),
}));

import { useDebugGetUsers } from "../services/api/default/default";

describe("WeeklyResultsPage", () => {
  const mockResults = [
    { player_name: "Player 1", score: 100 },
    { player_name: "Player 2", score: 90 },
  ];

  it("renders the error state", async () => {
    vi.mocked(useDebugGetUsers).mockReturnValue({
      data: null,
      isLoading: false,
      error: new Error("Failed to fetch weekly results"),
    } as any);

    render(<WeeklyResultsPage />);

    expect(
      await screen.findByText(/failed to fetch weekly results/i),
    ).toBeInTheDocument();
  });

  it("renders the results", async () => {
    const mockUsersData = {
      users: [
        {
          id: 1,
          name: "Player 1",
          total_wins: 5,
          pick_count: 10,
          email: "player1@test.com",
          role: "user",
          created_at: "",
          updated_at: "",
        },
        {
          id: 2,
          name: "Player 2",
          total_wins: 3,
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

    render(<WeeklyResultsPage />);

    expect(await screen.findByText(/player 1/i)).toBeInTheDocument();
    expect(await screen.findByText(/5/i)).toBeInTheDocument(); // total_wins (5) + index (0) = 5
    expect(await screen.findByText(/player 2/i)).toBeInTheDocument();
    expect(await screen.findByText(/4/i)).toBeInTheDocument(); // total_wins (3) + index (1) = 4
  });
});
