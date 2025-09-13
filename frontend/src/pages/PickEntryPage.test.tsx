import { screen, fireEvent, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import PickEntryPage from "./PickEntryPage";

// Mock the React Query hooks
const mockGamesData = {
  games: [
    { id: 1, favorite_team: "Team A", underdog_team: "Team B", spread: 3 },
    { id: 2, favorite_team: "Team C", underdog_team: "Team D", spread: 7 },
  ],
};

const mockSubmitPicks = vi.fn();
vi.mock("../services/api/games/games", () => ({
  useGetGames: () => ({
    data: mockGamesData,
    isLoading: false,
    error: null,
  }),
}));

vi.mock("../services/api/picks/picks", () => ({
  useSubmitPicks: () => ({
    mutateAsync: mockSubmitPicks,
    isPending: false,
  }),
}));

// Mock alert
global.alert = vi.fn();

describe("PickEntryPage", () => {
  beforeEach(() => {
    mockSubmitPicks.mockResolvedValue({});
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders the games", async () => {
    render(<PickEntryPage />);

    expect(await screen.findByText(/team a/i)).toBeInTheDocument();
    expect(await screen.findByText(/team b/i)).toBeInTheDocument();
    expect(await screen.findByText(/team c/i)).toBeInTheDocument();
    expect(await screen.findByText(/team d/i)).toBeInTheDocument();
  });

  it("handles quick pick", async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/quick pick/i));

    // After quick pick, the selects should have values
    const selects = await screen.findAllByRole("combobox");
    selects.forEach((select) => {
      expect(select).not.toHaveValue("");
    });
  });

  it("submits picks", async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/submit picks/i));

    // Wait for the async operation to complete
    await waitFor(() => {
      expect(mockSubmitPicks).toHaveBeenCalledWith({
        data: expect.arrayContaining([
          expect.objectContaining({
            game_id: expect.any(Number),
            picked: expect.any(String),
            rank: expect.any(Number),
            quick_pick: expect.any(Boolean),
          }),
        ]),
      });
    });

    expect(global.alert).toHaveBeenCalledWith("Picks submitted successfully!");
  });
});
