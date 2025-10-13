import { render, screen, fireEvent } from "@testing-library/react";
import GameForm from "./GameForm";
import { useCreateGame, useUpdateGame } from "../../services/api/games/games";

// Mock the custom hooks
vi.mock("../../services/api/games/games", () => ({
  useCreateGame: vi.fn(),
  useUpdateGame: vi.fn(),
}));

describe("GameForm", () => {
  const mockCreateGame = vi.fn();
  const mockUpdateGame = vi.fn();

  beforeEach(() => {
    (useCreateGame as jest.Mock).mockReturnValue({
      mutateAsync: mockCreateGame,
      isPending: false,
      error: null,
    });
    (useUpdateGame as jest.Mock).mockReturnValue({
      mutateAsync: mockUpdateGame,
      isPending: false,
      error: null,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const props = {
    open: true,
    onClose: vi.fn(),
    onSuccess: vi.fn(),
  };

  it("renders the form with the correct title for creating a new game", () => {
    render(<GameForm {...props} />);
    expect(screen.getByText("Add New Game")).toBeInTheDocument();
  });

  it("renders the form with the correct title and pre-filled data for editing a game", () => {
    const game = {
      id: 1,
      week: 1,
      season: 2023,
      home_team: "Team A",
      away_team: "Team B",
      spread: 3.5,
      start_time: "2023-09-10T12:00:00Z",
    };
    render(<GameForm {...props} game={game} />);
    expect(screen.getByText("Edit Game")).toBeInTheDocument();
    expect(screen.getByLabelText("Week")).toHaveValue(1);
    expect(screen.getByLabelText("Home Team")).toHaveValue("Team A");
  });

  it("shows validation errors for invalid input", async () => {
    render(<GameForm {...props} />);
    const createButton = screen.getByText("Create Game");
    fireEvent.click(createButton);

    const errorMessages = await screen.findAllByText(/required/i);
    expect(errorMessages).toHaveLength(3);
  });

  it("calls the useCreateGame mutation when creating a new game", async () => {
    render(<GameForm {...props} />);
    const weekInput = screen.getByLabelText("Week");
    const seasonInput = screen.getByLabelText("Season");
    const homeTeamInput = screen.getByLabelText("Home Team");
    const awayTeamInput = screen.getByLabelText("Away Team");
    const spreadInput = screen.getByLabelText("Spread");
    const startTimeInput = screen.getByLabelText("Start Time");

    fireEvent.change(weekInput, { target: { value: "1" } });
    fireEvent.change(seasonInput, { target: { value: "2023" } });
    fireEvent.change(homeTeamInput, { target: { value: "Team A" } });
    fireEvent.change(awayTeamInput, { target: { value: "Team B" } });
    fireEvent.change(spreadInput, { target: { value: "3.5" } });
    fireEvent.change(startTimeInput, { target: { value: "2023-09-10T12:00" } });

    const createButton = screen.getByText("Create Game");
    fireEvent.click(createButton);

    await vi.waitFor(() => {
      expect(mockCreateGame).toHaveBeenCalledTimes(1);
    });
  });
});
