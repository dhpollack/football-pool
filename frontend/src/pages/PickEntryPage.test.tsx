import { screen, fireEvent, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import PickEntryPage from "./PickEntryPage";
import { useListWeeks } from "../services/api/admin/admin";

// Mock the React Query hooks
const mockWeeksData = {
  weeks: [{ id: 1, week_number: 1, season: 2023, is_active: true }],
};

const mockGamesData = {
  games: [
    {
      id: 1,
      home_team: "Team A",
      away_team: "Team B",
      spread: 3,
      favorite: "Home",
      underdog: "Away",
    },
    {
      id: 2,
      home_team: "Team C",
      away_team: "Team D",
      spread: 7,
      favorite: "Away",
      underdog: "Home",
    },
  ],
};

const mockSubmitPicks = vi.fn();
vi.mock("../services/api/admin/admin", () => ({
  useListWeeks: vi.fn(),
}));

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
    vi.mocked(useListWeeks).mockReturnValue({
      data: mockWeeksData,
      isLoading: false,
      error: null,
    });
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

  it("shows a message if no active week is found", async () => {
    // Mock the hook to return no active weeks
    vi.mocked(useListWeeks).mockReturnValue({
      data: {
        weeks: [{ id: 1, week_number: 1, season: 2023, is_active: false }],
      },
      isLoading: false,
      error: null,
    });

    render(<PickEntryPage />);

    expect(
      await screen.findByText(/no active week found/i),
    ).toBeInTheDocument();
  });

  it("shows a message if no weeks are returned", async () => {
    // Mock the hook to return no weeks
    vi.mocked(useListWeeks).mockReturnValue({
      data: { weeks: [] },
      isLoading: false,
      error: null,
    });

    render(<PickEntryPage />);

    expect(
      await screen.findByText(/no active week found/i),
    ).toBeInTheDocument();
  });
});
