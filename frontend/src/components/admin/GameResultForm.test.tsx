
import { render, screen, fireEvent } from "@testing-library/react";
import GameResultForm from "./GameResultForm";
import { useSubmitResult } from "../../services/api/default/default";

// Mock the custom hooks
vi.mock("../../services/api/default/default", () => ({
  useSubmitResult: vi.fn(),
}));

describe("GameResultForm", () => {
  const mockSubmitResult = vi.fn();

  beforeEach(() => {
    (useSubmitResult as jest.Mock).mockReturnValue({ mutateAsync: mockSubmitResult, isPending: false, error: null });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const game = {
    id: 1,
    week: 1,
    season: 2023,
    favorite_team: "Team A",
    underdog_team: "Team B",
    spread: -3.5,
    start_time: "2023-09-10T12:00:00Z",
  };

  const props = {
    open: true,
    onClose: vi.fn(),
    onSuccess: vi.fn(),
    game,
  };

  it("renders the form with the correct title and game information", () => {
    render(<GameResultForm {...props} />);
    expect(screen.getByText("Add Game Result - Team A vs Team B")).toBeInTheDocument();
    expect(screen.getByText("Week 1, Season 2023")).toBeInTheDocument();
  });

  it("calculates and displays the outcome of the game based on the scores", () => {
    render(<GameResultForm {...props} />);
    const favoriteScoreInput = screen.getByLabelText("Team A Score");
    const underdogScoreInput = screen.getByLabelText("Team B Score");

    fireEvent.change(favoriteScoreInput, { target: { value: "24" } });
    fireEvent.change(underdogScoreInput, { target: { value: "20" } });

    expect(screen.getByText("Calculated Outcome: FAVORITE (Team A covers)")).toBeInTheDocument();
  });

  it("shows validation errors for invalid input", async () => {
    render(<GameResultForm {...props} />);
    const submitButton = screen.getByText("Submit Result");
    fireEvent.click(submitButton);

    expect(await screen.findByText("Outcome is required")).toBeInTheDocument();
  });

  it("calls the useSubmitResult mutation when the form is submitted", async () => {
    render(<GameResultForm {...props} />);
    const favoriteScoreInput = screen.getByLabelText("Team A Score");
    const underdogScoreInput = screen.getByLabelText("Team B Score");
    
    fireEvent.change(favoriteScoreInput, { target: { value: "24" } });
    fireEvent.change(underdogScoreInput, { target: { value: "20" } });

    const outcomeSelect = screen.getByRole("combobox");
    fireEvent.mouseDown(outcomeSelect);
    const favoriteWinsOption = await screen.findByRole('option', { name: 'Favorite Wins' });
    fireEvent.click(favoriteWinsOption);

    const submitButton = screen.getByText("Submit Result");
    fireEvent.click(submitButton);

    await vi.waitFor(() => {
        expect(mockSubmitResult).toHaveBeenCalledTimes(1);
    });
  });
});
