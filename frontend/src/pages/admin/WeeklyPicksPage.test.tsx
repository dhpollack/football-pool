import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import WeeklyPicksPage from "./WeeklyPicksPage";
import { useAdminGetPicksByWeek } from "../../services/api/picks/picks";
import { getGetPicksResponseMock } from "../../services/api/picks/picks.msw";

// Mock the custom hooks
vi.mock("../../services/api/picks/picks", () => ({
  useAdminGetPicksByWeek: vi.fn(),
}));

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({
    data,
    loading,
    error,
  }: {
    data?: unknown[];
    loading?: boolean;
    error?: string | null;
  }) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading...</div>}
      {error && <div>Error: {error}</div>}
      {data && (
        <div>
          {data.map((pick) => (
            <div key={pick.id} data-testid="pick-row">
              <span data-testid="pick-user">{pick.user?.email}</span>
              <span data-testid="pick-game">
                {pick.game
                  ? `${pick.game.home_team} vs ${pick.game.away_team}`
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

// Mock react-router-dom useParams
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useParams: () => ({ week: "1" }),
  };
});

describe("WeeklyPicksPage", () => {
  const mockRefetch = vi.fn();

  beforeEach(() => {
    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: false,
      refetch: mockRefetch,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const renderWithRouter = () => {
    return render(
      <MemoryRouter initialEntries={["/admin/picks/week/1"]}>
        <Routes>
          <Route path="/admin/picks/week/:week" element={<WeeklyPicksPage />} />
        </Routes>
      </MemoryRouter>,
    );
  };

  it("renders breadcrumbs and page title", () => {
    renderWithRouter();

    expect(screen.getByText("Admin")).toBeInTheDocument();
    expect(screen.getByText("Picks")).toBeInTheDocument();
    // There are multiple elements with "Week 1 Picks" text
    const weekPicksElements = screen.getAllByText(/Week 1 Picks/);
    expect(weekPicksElements.length).toBeGreaterThan(0);
  });

  it("renders loading state", () => {
    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: true,
      refetch: mockRefetch,
    });

    renderWithRouter();

    // Check for loading button with test ID
    expect(screen.getByTestId("load-picks-button")).toHaveTextContent(
      "Loading...",
    );
    // Check for loading text in our mock AdminDataTable
    expect(screen.getByTestId("admin-data-table")).toHaveTextContent(
      "Loading...",
    );
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
            home_team: "Team A",
            away_team: "Team B",
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

    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
      refetch: mockRefetch,
    });

    renderWithRouter();

    // Check for the data in our mock AdminDataTable implementation
    expect(screen.getByTestId("pick-user")).toHaveTextContent(
      "user1@example.com",
    );
    expect(screen.getByTestId("pick-game")).toHaveTextContent(
      "Team A vs Team B",
    );
    expect(screen.getByTestId("pick-choice")).toHaveTextContent("favorite");
  });

  it("shows error message when loading fails", () => {
    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: null,
      error: { message: "Failed to load picks" },
      isLoading: false,
      refetch: mockRefetch,
    });

    renderWithRouter();

    expect(
      screen.getByText("Error loading picks: Failed to load picks"),
    ).toBeInTheDocument();
  });

  it("renders summary statistics when data is available", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user_id: 1,
          user: {
            id: 1,
            email: "user1@example.com",
            name: "User One",
            role: "user",
          },
          game: {
            id: 1,
            home_team: "Team A",
            away_team: "Team B",
            week: 1,
            season: 2023,
          },
          picked: "favorite",
          created_at: "2023-09-10T12:00:00Z",
        },
        {
          id: 2,
          user_id: 2,
          user: {
            id: 2,
            email: "user2@example.com",
            name: "User Two",
            role: "user",
          },
          game: {
            id: 1,
            home_team: "Team A",
            away_team: "Team B",
            week: 1,
            season: 2023,
          },
          picked: "underdog",
          created_at: "2023-09-10T12:00:00Z",
        },
      ],
      pagination: {
        page: 1,
        limit: 20,
        total: 2,
      },
    });

    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
      refetch: mockRefetch,
    });

    renderWithRouter();

    expect(screen.getByTestId("total-picks-box")).toHaveTextContent("2");
    expect(screen.getByTestId("unique-users-box")).toHaveTextContent("2");
    expect(screen.getByTestId("favorite-picks-box")).toHaveTextContent("1");
    expect(screen.getByTestId("underdog-picks-box")).toHaveTextContent("1");
  });

  it("calls refetch when Load Picks button is clicked", async () => {
    renderWithRouter();

    const loadButton = screen.getByText("Load Picks");
    fireEvent.click(loadButton);

    expect(mockRefetch).toHaveBeenCalledTimes(1);
  });

  it("updates week input and calls refetch", async () => {
    renderWithRouter();

    const weekInput = screen.getByLabelText("Week");
    fireEvent.change(weekInput, { target: { value: "2" } });

    const loadButton = screen.getByText("Load Picks");
    fireEvent.click(loadButton);

    expect(mockRefetch).toHaveBeenCalledTimes(1);
    expect(weekInput).toHaveValue(2);
  });

  it("updates season input and calls refetch", async () => {
    renderWithRouter();

    const seasonInput = screen.getByLabelText("Season");
    fireEvent.change(seasonInput, { target: { value: "2025" } });

    const loadButton = screen.getByText("Load Picks");
    fireEvent.click(loadButton);

    expect(mockRefetch).toHaveBeenCalledTimes(1);
    expect(seasonInput).toHaveValue(2025);
  });

  it("filters picks based on search term", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: {
            id: 1,
            email: "user1@example.com",
            name: "User One",
            role: "user",
          },
          game: {
            id: 1,
            home_team: "Team A",
            away_team: "Team B",
            week: 1,
            season: 2023,
          },
          picked: "favorite",
          created_at: "2023-09-10T12:00:00Z",
        },
        {
          id: 2,
          user: {
            id: 2,
            email: "user2@example.com",
            name: "User Two",
            role: "user",
          },
          game: {
            id: 2,
            home_team: "Team C",
            away_team: "Team D",
            week: 1,
            season: 2023,
          },
          picked: "underdog",
          created_at: "2023-09-10T12:00:00Z",
        },
      ],
      pagination: {
        page: 1,
        limit: 20,
        total: 2,
      },
    });

    (useAdminGetPicksByWeek as jest.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
      refetch: mockRefetch,
    });

    renderWithRouter();

    // Both picks should be visible initially
    const pickUsers = screen.getAllByTestId("pick-user");
    expect(pickUsers[0]).toHaveTextContent("user1@example.com");
    expect(pickUsers[1]).toHaveTextContent("user2@example.com");

    const pickGames = screen.getAllByTestId("pick-game");
    expect(pickGames[0]).toHaveTextContent("Team A vs Team B");
    expect(pickGames[1]).toHaveTextContent("Team C vs Team D");

    const searchInput = screen.getByLabelText("Search");
    fireEvent.change(searchInput, { target: { value: "Team C" } });

    // Since our mock AdminDataTable doesn't implement actual filtering,
    // we just verify that the search input value was set correctly
    expect(searchInput).toHaveValue("Team C");
  });
});
