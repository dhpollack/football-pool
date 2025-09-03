
import { render, screen, fireEvent } from "@testing-library/react";
import AdminGamesPage from "./AdminGamesPage";
import { useAdminListGames } from "../../services/api/default/default";

// Mock the custom hooks and components
vi.mock("../../services/api/default/default", () => ({
  useAdminListGames: vi.fn(),
}));

vi.mock("../../components/admin/AdminDataTable", () => ({
  __esModule: true,
  default: ({ data }: { data: any[] }) => (
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
  default: () => <div data-testid="admin-action-buttons"></div>,
}));

vi.mock("../../components/admin/AdminConfirmDialog", () => ({
  __esModule: true,
  default: () => <div data-testid="admin-confirm-dialog"></div>,
}));

describe("AdminGamesPage", () => {
  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders the component with title", () => {
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: { games: [] },
      error: null,
      isLoading: false,
    });
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
    // In the actual component, the loading state is passed to AdminDataTable.
    // We can check if the mocked data table is rendered.
    expect(screen.getByTestId("admin-data-table")).toBeInTheDocument();
  });

  it("displays an error message if the API call fails", () => {
    (useAdminListGames as jest.Mock).mockReturnValue({
      data: null,
      error: { message: "Failed to fetch games" },
      isLoading: false,
    });
    render(<AdminGamesPage />);
    expect(screen.getByText("Error loading games: Failed to fetch games")).toBeInTheDocument();
  });

  it("displays the games in a table when the API call is successful", () => {
    const mockGames = [
      {
        id: 1,
        week: 1,
        season: 2023,
        favorite_team: "Team A",
        underdog_team: "Team B",
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
});
