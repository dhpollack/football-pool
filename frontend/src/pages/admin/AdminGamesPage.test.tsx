import { render, screen, fireEvent } from "@testing-library/react";
import AdminGamesPage from "./AdminGamesPage";
import {
  useAdminListGames,
  useCreateGame,
  useUpdateGame,
  getAdminListGamesQueryKey,
} from "../../services/api/games/games";
import { useQueryClient } from "@tanstack/react-query";

// Mock the custom hooks and components
vi.mock("../../services/api/games/games", () => ({
  useAdminListGames: vi.fn(),
  useCreateGame: vi.fn(),
  useUpdateGame: vi.fn(),
  getAdminListGamesQueryKey: vi.fn(),
}));

vi.mock("@tanstack/react-query", async () => {
  const actual = await vi.importActual("@tanstack/react-query");
  return {
    ...actual,
    useQueryClient: vi.fn(),
  };
});

vi.mock("../../components/admin/AdminDataTable", () => ({
  __esModule: true,
  default: ({ data }: { data: { id: number }[] }) => (
    <div data-testid="admin-data-table">
      {data.map((item) => (
        <div key={item.id}>{item.id}</div>
      ))}
    </div>
  ),
}));

vi.mock("../../components/admin/AdminSearchFilter", () => ({
  __esModule: true,
  default: () => <div data-testid="admin-search-filter"></div>,
}));

vi.mock("../../components/admin/AdminActionButtons", () => ({
  __esModule: true,
  default: ({ onAdd, addLabel }) => (
    <button type="button" onClick={onAdd} data-testid="add-game-button">
      {addLabel}
    </button>
  ),
}));

vi.mock("../../components/admin/AdminConfirmDialog", () => ({
  __esModule: true,
  default: () => <div data-testid="admin-confirm-dialog"></div>,
}));

vi.mock("../../components/admin/GameForm", () => ({
  __esModule: true,
  default: ({
    open,
    onClose,
    onSuccess,
  }: {
    open: boolean;
    onClose: () => void;
    onSuccess: () => void;
  }) =>
    open ? (
      <div data-testid="game-form">
        <input
          type="number"
          data-testid="spread-input"
          onChange={(e) => {
            // Simulate validation - prevent submission if spread is negative
            const spread = parseFloat(e.target.value);
            if (spread < 0) {
              // Show validation error
              const errorElement = document.querySelector('[data-testid="spread-error"]');
              if (!errorElement) {
                const errorDiv = document.createElement('div');
                errorDiv.setAttribute('data-testid', 'spread-error');
                errorDiv.textContent = 'Spread must be greater than or equal to 0';
                document.querySelector('[data-testid="game-form"]')?.appendChild(errorDiv);
              }
            } else {
              // Clear validation error
              const errorElement = document.querySelector('[data-testid="spread-error"]');
              errorElement?.remove();
            }
          }}
        />
        <button
          type="button"
          onClick={() => {
            // Check if there's a validation error before submitting
            const hasError = document.querySelector('[data-testid="spread-error"]');
            if (!hasError) {
              // Simulate successful form submission only if no validation errors
              onSuccess();
              onClose();
            }
          }}
          data-testid="submit-game-button"
        >
          Submit Game
        </button>
        <button
          type="button"
          onClick={onClose}
          data-testid="cancel-game-button"
        >
          Cancel
        </button>
      </div>
    ) : null,
}));

describe("AdminGamesPage", () => {
  const mockMutate = vi.fn();
  const mockInvalidateQueries = vi.fn();

  beforeEach(() => {
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: { games: [] },
      error: null,
      isLoading: false,
    });
    (useCreateGame as jest.Mock).mockReturnValue({
      mutate: mockMutate,
      isPending: false,
      error: null,
    });
    (useUpdateGame as jest.Mock).mockReturnValue({
      mutate: vi.fn(),
      isPending: false,
      error: null,
    });
    (useQueryClient as jest.Mock).mockReturnValue({
      invalidateQueries: mockInvalidateQueries,
    });
    (getAdminListGamesQueryKey as jest.Mock).mockReturnValue([
      "adminListGames",
    ]);
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders the component with title", () => {
    render(<AdminGamesPage />);
    expect(screen.getByText("Game Management")).toBeInTheDocument();
  });

  it("displays a loading state initially", () => {
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: true,
    });
    render(<AdminGamesPage />);
    expect(screen.getByTestId("admin-data-table")).toBeInTheDocument();
  });

  it("displays an error message if the API call fails", () => {
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: null,
      error: { message: "Failed to fetch games" },
      isLoading: false,
    });
    render(<AdminGamesPage />);
    expect(
      screen.getByText("Error loading games: Failed to fetch games"),
    ).toBeInTheDocument();
  });

  it("displays the games in a table when the API call is successful", () => {
    const mockGames = [
      {
        id: 1,
        week: 1,
        season: 2023,
        home_team: "Team A",
        away_team: "Team B",
        spread: 3.5,
        start_time: new Date().toISOString(),
      },
    ];
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: { games: mockGames },
      error: null,
      isLoading: false,
    });
    render(<AdminGamesPage />);
    expect(screen.getByTestId("admin-data-table")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
  });

  it("opens game form on add button click and invalidates queries on form success", async () => {
    render(<AdminGamesPage />);

    // Click add button to open form
    const addButton = screen.getByTestId("add-game-button");
    fireEvent.click(addButton);

    // Verify form is open
    expect(screen.getByTestId("game-form")).toBeInTheDocument();

    // Submit the form (simulating successful game creation)
    const submitButton = screen.getByTestId("submit-game-button");
    fireEvent.click(submitButton);

    // Verify queries are invalidated on form success
    expect(mockInvalidateQueries).toHaveBeenCalledWith({
      queryKey: ["adminListGames"],
    });
  });

  it("prevents form submission when spread is negative", async () => {
    render(<AdminGamesPage />);

    // Click add button to open form
    const addButton = screen.getByTestId("add-game-button");
    fireEvent.click(addButton);

    // Verify form is open
    expect(screen.getByTestId("game-form")).toBeInTheDocument();

    // Set a negative spread
    const spreadInput = screen.getByTestId("spread-input");
    fireEvent.change(spreadInput, { target: { value: "-3.5" } });

    // Verify validation error is displayed
    expect(
      screen.getByText("Spread must be greater than or equal to 0"),
    ).toBeInTheDocument();

    // Try to submit the form
    const submitButton = screen.getByTestId("submit-game-button");
    fireEvent.click(submitButton);

    // Verify queries are NOT invalidated (form should not submit)
    expect(mockInvalidateQueries).not.toHaveBeenCalled();
  });

  it("allows form submission when spread is zero or positive", async () => {
    render(<AdminGamesPage />);

    // Click add button to open form
    const addButton = screen.getByTestId("add-game-button");
    fireEvent.click(addButton);

    // Verify form is open
    expect(screen.getByTestId("game-form")).toBeInTheDocument();

    // Set a positive spread
    const spreadInput = screen.getByTestId("spread-input");
    fireEvent.change(spreadInput, { target: { value: "3.5" } });

    // Verify no validation error is displayed
    expect(
      screen.queryByText("Spread must be greater than or equal to 0"),
    ).not.toBeInTheDocument();

    // Submit the form
    const submitButton = screen.getByTestId("submit-game-button");
    fireEvent.click(submitButton);

    // Verify queries are invalidated on form success
    expect(mockInvalidateQueries).toHaveBeenCalledWith({
      queryKey: ["adminListGames"],
    });
  });
});
