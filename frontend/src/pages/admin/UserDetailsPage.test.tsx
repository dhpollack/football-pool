import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import UserDetailsPage from "./UserDetailsPage";
import { getGetProfileResponseMock } from "../../services/api/user/user.msw";
import { getGetPicksResponseMock } from "../../services/api/picks/picks.msw";

// Mock the custom hooks with proper return values to avoid hoisting issues
vi.mock("../../services/api/user/user", () => ({
  useAdminGetUser: vi.fn(() => ({
    data: null,
    error: null,
    isLoading: false,
  })),
}));

vi.mock("../../services/api/picks/picks", () => ({
  useAdminGetPicksByUser: vi.fn(() => ({
    data: null,
    error: null,
    isLoading: false,
  })),
}));

// Import the hooks after mocking to get the mocked versions
import { useAdminGetUser } from "../../services/api/user/user";
import { useAdminGetPicksByUser } from "../../services/api/picks/picks";

// Create typed mock references using vi.mocked()
const mockUseAdminGetUser = vi.mocked(useAdminGetUser);
const mockUseAdminGetPicksByUser = vi.mocked(useAdminGetPicksByUser);

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({ data, loading, error }: any) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading picks...</div>}
      {error && <div>Error: {error}</div>}
      {!loading && data && data.length > 0 && (
        <div>
          {data.map((pick: any) => (
            <div key={pick.id} data-testid="pick-row">
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
      {!loading && (!data || data.length === 0) && (
        <div>No picks found for this user</div>
      )}
    </div>
  ),
}));

// Mock react-router-dom useParams
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useParams: () => ({ userId: "1" }),
  };
});

describe("UserDetailsPage", () => {
  beforeEach(() => {
    // Reset mocks before each test
    mockUseAdminGetUser.mockClear();
    mockUseAdminGetPicksByUser.mockClear();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
  });

  const renderWithRouter = () => {
    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/admin/users/1"]}>
          <Routes>
            <Route path="/admin/users/:userId" element={<UserDetailsPage />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>,
    );
  };

  it("renders breadcrumbs and loading state", () => {
    mockUseAdminGetUser.mockReturnValue({
      data: null,
      error: null,
      isLoading: true,
    });

    renderWithRouter();

    expect(screen.getByText("Admin")).toBeInTheDocument();
    expect(screen.getByText("Users")).toBeInTheDocument();
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders user information when loaded", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user",
      player: {
        id: 1,
        user_id: 1,
        name: "Test Player",
        address: "123 Test St",
      },
      created_at: "2023-09-10T12:00:00Z",
    });

    mockUseAdminGetUser.mockReturnValue({
      data: mockUser,
      error: null,
      isLoading: false,
    });

    renderWithRouter();

    // Check for email in user info section (not breadcrumbs)
    const userInfoSections = screen.getAllByText("test@example.com");
    expect(userInfoSections.length).toBeGreaterThan(0);
    expect(screen.getByText("Test User")).toBeInTheDocument();
    expect(screen.getByText("user")).toBeInTheDocument();
    expect(screen.getByText("Test Player")).toBeInTheDocument();
    expect(screen.getByText("123 Test St")).toBeInTheDocument();
  });

  it("shows error message when user loading fails", () => {
    mockUseAdminGetUser.mockReturnValue({
      data: null,
      error: { message: "Failed to load user" },
      isLoading: false,
    });

    renderWithRouter();

    expect(
      screen.getByText("Error loading user: Failed to load user"),
    ).toBeInTheDocument();
  });

  it("renders user picks when available", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user",
    });

    const mockPicks = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: {
            id: 1,
            email: "test@example.com",
            name: "Test User",
            role: "user",
          },
          game: {
            id: 1,
            favorite_team: "Team A",
            underdog_team: "Team B",
            week: 1,
            season: 2023,
          },
          picked: "favorite",
          created_at: "2023-09-10T12:00:00Z",
        },
      ],
      pagination: { page: 1, limit: 10, total: 1 },
    });

    // Set up mock return values before rendering
    mockUseAdminGetUser.mockReturnValue({
      data: mockUser,
      error: null,
      isLoading: false,
    });
    mockUseAdminGetPicksByUser.mockReturnValue({
      data: mockPicks,
      error: null,
      isLoading: false,
    });

    renderWithRouter();

    expect(screen.getByTestId("pick-game")).toHaveTextContent("Team A vs Team B");
    expect(screen.getByTestId("pick-choice")).toHaveTextContent("favorite");
  });

  it("shows error message when picks loading fails", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user",
    });

    mockUseAdminGetUser.mockReturnValue({
      data: mockUser,
      error: null,
      isLoading: false,
    });
    mockUseAdminGetPicksByUser.mockReturnValue({
      data: null,
      error: { message: "Failed to load picks" },
      isLoading: false,
    });

    renderWithRouter();

    expect(
      screen.getByText("Error loading picks: Failed to load picks"),
    ).toBeInTheDocument();
  });

  it("shows statistics tab content", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user",
    });

    mockUseAdminGetUser.mockReturnValue({
      data: mockUser,
      error: null,
      isLoading: false,
    });

    renderWithRouter();

    // The statistics tab content is hidden by default, but we can check if the tab exists
    expect(screen.getByText("Statistics")).toBeInTheDocument();
  });
});
