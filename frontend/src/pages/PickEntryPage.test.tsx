import { render, screen, fireEvent } from "@testing-library/react";
import { vi } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import PickEntryPage from "./PickEntryPage";

// Mock the React Query hooks
const mockGamesData = {
  games: [
    { id: 1, favorite_team: "Team A", underdog_team: "Team B", spread: 3 },
    { id: 2, favorite_team: "Team C", underdog_team: "Team D", spread: 7 },
  ]
};

const mockSubmitPicks = vi.fn();
vi.mock("../services/api/default/default", () => ({
  useGetGames: () => ({
    data: mockGamesData,
    isLoading: false,
    error: null,
  }),
  useSubmitPicks: () => ({
    mutateAsync: mockSubmitPicks,
    isPending: false,
  }),
}));

// Mock alert
global.alert = vi.fn();

describe("PickEntryPage", () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    mockSubmitPicks.mockResolvedValue({});
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const renderWithProviders = () => {
    return render(
      <QueryClientProvider client={queryClient}>
        <PickEntryPage />
      </QueryClientProvider>
    );
  };

  it("renders the games", async () => {
    renderWithProviders();

    expect(await screen.findByText(/team a/i)).toBeInTheDocument();
    expect(await screen.findByText(/team b/i)).toBeInTheDocument();
    expect(await screen.findByText(/team c/i)).toBeInTheDocument();
    expect(await screen.findByText(/team d/i)).toBeInTheDocument();
  });

  it("handles quick pick", async () => {
    renderWithProviders();

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/quick pick/i));

    // After quick pick, the selects should have values
    const selects = await screen.findAllByRole("combobox");
    selects.forEach((select) => {
      expect(select).not.toHaveValue("");
    });
  });

  it("submits picks", async () => {
    renderWithProviders();

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/submit picks/i));

    // Wait for the async operation to complete
    await vi.waitFor(() => {
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