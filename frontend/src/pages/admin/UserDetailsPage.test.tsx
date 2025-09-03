import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import UserDetailsPage from "./UserDetailsPage";
import { useAdminGetUser, useAdminGetPicksByUser } from "../../services/api/default/default";
import { getGetProfileResponseMock, getGetPicksResponseMock } from "../../services/api/default/default.msw";

// Mock the custom hooks
vi.mock("../../services/api/default/default", () => ({
  useAdminGetUser: vi.fn(),
  useAdminGetPicksByUser: vi.fn(),
}));

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({ data, loading, error }: any) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading picks...</div>}
      {error && <div>Error: {error}</div>}
      {data && (
        <div>
          {data.map((pick: any) => (
            <div key={pick.id} data-testid="pick-row">
              <span data-testid="pick-game">
                {pick.game ? `${pick.game.favorite_team} vs ${pick.game.underdog_team}` : "Unknown"}
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
    useParams: () => ({ userId: "1" }),
  };
});

describe("UserDetailsPage", () => {
  beforeEach(() => {
    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: null, 
      error: null, 
      isLoading: false 
    });
    (useAdminGetPicksByUser as jest.Mock).mockReturnValue({ 
      data: null, 
      error: null, 
      isLoading: false 
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const renderWithRouter = () => {
    return render(
      <MemoryRouter initialEntries={["/admin/users/1"]}>
        <Routes>
          <Route path="/admin/users/:userId" element={<UserDetailsPage />} />
        </Routes>
      </MemoryRouter>
    );
  };

  it("renders breadcrumbs and loading state", () => {
    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: null, 
      error: null, 
      isLoading: true 
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
        address: "123 Test St"
      },
      created_at: "2023-09-10T12:00:00Z"
    });

    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: mockUser, 
      error: null, 
      isLoading: false 
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
    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: null, 
      error: { message: "Failed to load user" }, 
      isLoading: false 
    });

    renderWithRouter();
    
    expect(screen.getByText("Error loading user: Failed to load user")).toBeInTheDocument();
  });

  it("renders user picks when available", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user"
    });

    const mockPicks = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: { 
            id: 1, 
            email: "test@example.com", 
            name: "Test User", 
            role: "user"
          },
          game: {
            id: 1,
            favorite_team: "Team A",
            underdog_team: "Team B", 
            week: 1,
            season: 2023
          },
          picked: "favorite",
          created_at: "2023-09-10T12:00:00Z"
        }
      ],
      pagination: { page: 1, limit: 10, total: 1 }
    });

    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: mockUser, 
      error: null, 
      isLoading: false 
    });
    (useAdminGetPicksByUser as jest.Mock).mockReturnValue({ 
      data: mockPicks, 
      error: null, 
      isLoading: false 
    });

    renderWithRouter();
    
    expect(screen.getByText("Team A vs Team B")).toBeInTheDocument();
    expect(screen.getByText("favorite")).toBeInTheDocument();
  });

  it("shows error message when picks loading fails", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user"
    });

    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: mockUser, 
      error: null, 
      isLoading: false 
    });
    (useAdminGetPicksByUser as jest.Mock).mockReturnValue({ 
      data: null, 
      error: { message: "Failed to load picks" }, 
      isLoading: false 
    });

    renderWithRouter();
    
    expect(screen.getByText("Error loading picks: Failed to load picks")).toBeInTheDocument();
  });

  it("shows statistics tab content", () => {
    const mockUser = getGetProfileResponseMock({
      id: 1,
      email: "test@example.com",
      name: "Test User",
      role: "user"
    });

    (useAdminGetUser as jest.Mock).mockReturnValue({ 
      data: mockUser, 
      error: null, 
      isLoading: false 
    });

    renderWithRouter();
    
    // The statistics tab content is hidden by default, but we can check if the tab exists
    expect(screen.getByText("Statistics")).toBeInTheDocument();
  });
});