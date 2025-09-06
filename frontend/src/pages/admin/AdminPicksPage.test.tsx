import { render, screen } from "@testing-library/react";
import AdminPicksPage from "./AdminPicksPage";
import { useAdminListPicks } from "../../services/api/picks/picks";
import { getGetPicksResponseMock } from "../../services/api/picks/picks.msw";

// Mock the custom hooks
vi.mock("../../services/api/picks/picks", () => ({
  useAdminListPicks: vi.fn(),
}));

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({ data, loading, error }: any) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading...</div>}
      {error && <div>Error: {error}</div>}
      {data && (
        <div>
          {data.map((pick: any) => (
            <div key={pick.id} data-testid="pick-row">
              <span data-testid="pick-user">{pick.user?.email}</span>
              <span data-testid="pick-game">
                {pick.game
                  ? `${pick.game.favorite_team} vs ${pick.game.underdog_team}`
                  : "Unknown"}
              </span>
              <span data-testid="pick-choice">{pick.picked}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  ),
}));

vi.mock("../../components/admin/AdminSearchFilter", () => ({
  default: () => <div data-testid="search-filter">Search Filter</div>,
}));

vi.mock("../../components/admin/AdminConfirmDialog", () => ({
  default: ({ open }: any) =>
    open && <div data-testid="confirm-dialog">Confirm Dialog</div>,
}));

describe("AdminPicksPage", () => {
  beforeEach(() => {
    (useAdminListPicks as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: false,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders the page title", () => {
    render(<AdminPicksPage />);

    expect(screen.getByText("Pick Management")).toBeInTheDocument();
  });

  it("renders loading state", () => {
    (useAdminListPicks as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: true,
    });

    render(<AdminPicksPage />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders picks data when loaded", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: {
            id: 1,
            email: "user1@example.com",
            name: "User One",
            role: "user",
            created_at: "2023-09-10T12:00:00Z",
            updated_at: "2023-09-10T12:00:00Z",
          },
          game: {
            id: 1,
            favorite_team: "Team A",
            underdog_team: "Team B",
            week: 1,
            season: 2023,
            spread: 3.5,
            start_time: "2023-09-10T12:00:00Z",
            created_at: "2023-09-10T12:00:00Z",
            updated_at: "2023-09-10T12:00:00Z",
          },
          picked: "favorite",
          created_at: "2023-09-10T12:00:00Z",
          updated_at: "2023-09-10T12:00:00Z",
        },
      ],
      pagination: {
        page: 1,
        limit: 20,
        total: 1,
      },
    });

    (useAdminListPicks as jest.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />);

    expect(screen.getByText("user1@example.com")).toBeInTheDocument();
    expect(screen.getByText("Team A vs Team B")).toBeInTheDocument();
    expect(screen.getByText("favorite")).toBeInTheDocument();
  });

  it("shows error message when there is an error", () => {
    (useAdminListPicks as jest.Mock).mockReturnValue({
      data: null,
      error: { message: "Failed to load picks" },
      isLoading: false,
    });

    render(<AdminPicksPage />);

    expect(screen.getByText("Error: Failed to load picks")).toBeInTheDocument();
  });

  it("renders search filter component", () => {
    render(<AdminPicksPage />);

    expect(screen.getByTestId("search-filter")).toBeInTheDocument();
  });
});
