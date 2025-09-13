import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import AdminPicksPage from "./AdminPicksPage";
import { useAdminListPicks } from "../../services/api/picks/picks";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { PickResponse } from "../../services/model";

// Mock the custom hooks
vi.mock("../../services/api/picks/picks", () => ({
  useAdminListPicks: vi.fn(),
  useAdminDeletePick: vi.fn(),
}));

describe("AdminPicksPage (Integration)", () => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );

  const mockPicks: PickResponse[] = [
    {
      id: 1,
      user: {
        id: 1,
        email: "user1@example.com",
        name: "User One",
        role: "user",
        created_at: "",
        updated_at: "",
      },
      game: {
        id: 1,
        favorite_team: "Team A",
        underdog_team: "Team B",
        week: 1,
        season: 2023,
        spread: 3.5,
        start_time: "",
        created_at: "",
        updated_at: "",
      },
      picked: "favorite",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: 2,
      user: {
        id: 2,
        email: "user2@example.com",
        name: "User Two",
        role: "user",
        created_at: "",
        updated_at: "",
      },
      game: {
        id: 2,
        favorite_team: "Team C",
        underdog_team: "Team D",
        week: 1,
        season: 2023,
        spread: 7,
        start_time: "",
        created_at: "",
        updated_at: "",
      },
      picked: "underdog",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
  ];

  afterEach(() => {
    vi.clearAllMocks();
    queryClient.clear();
  });

  it("should show the real loading spinner and then hide it", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: null,
      isLoading: true,
    });

    const { rerender } = render(<AdminPicksPage />, { wrapper });

    expect(screen.getByRole("progressbar")).toBeInTheDocument();

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: {
        picks: mockPicks,
        pagination: { total: mockPicks.length, page: 1, limit: 20 },
      },
      error: null,
      isLoading: false,
    });

    rerender(<AdminPicksPage />);

    await waitFor(() => {
      expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();
    });

    expect(screen.getByText("user1@example.com")).toBeInTheDocument();
  });

  it("should filter picks based on search term", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: {
        picks: mockPicks,
        pagination: { total: mockPicks.length, page: 1, limit: 20 },
      },
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    expect(screen.getByText("user1@example.com")).toBeInTheDocument();
    expect(screen.getByText("user2@example.com")).toBeInTheDocument();

    const searchInput = screen.getByPlaceholderText("Search by user email...");
    fireEvent.change(searchInput, { target: { value: "user1" } });

    await waitFor(() => {
      expect(screen.getByText("user1@example.com")).toBeInTheDocument();
      expect(screen.queryByText("user2@example.com")).not.toBeInTheDocument();
    });
  });

  it("should show an empty message when client-side search yields no results", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: {
        picks: mockPicks,
        pagination: { total: mockPicks.length, page: 1, limit: 20 },
      },
      error: null,
      isLoading: false,
    });
    render(<AdminPicksPage />, { wrapper });

    const searchInput = screen.getByPlaceholderText("Search by user email...");
    fireEvent.change(searchInput, { target: { value: "nonexistent" } });

    await waitFor(() => {
      expect(
        screen.getByText("No picks match your filter criteria"),
      ).toBeInTheDocument();
    });
  });

  it("should show empty state when API returns no picks", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: { picks: [], pagination: { total: 0, page: 1, limit: 20 } },
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    await waitFor(() => {
      expect(screen.getByText("No picks available")).toBeInTheDocument();
    });
  });

  it("should render table with correct column structure when picks exist", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: {
        picks: mockPicks,
        pagination: { total: mockPicks.length, page: 1, limit: 20 },
      },
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    // Verify table headers - use getAllByText since there are multiple elements with same text
    expect(screen.getAllByText("ID").length).toBeGreaterThan(0);
    expect(screen.getAllByText("User").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Game").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Week").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Season").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Pick").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Submitted").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Actions").length).toBeGreaterThan(0);

    // Verify pick data is rendered
    await waitFor(() => {
      expect(screen.getByText("user1@example.com")).toBeInTheDocument();
      expect(screen.getByText("Team A vs Team B")).toBeInTheDocument();
      expect(screen.getAllByText("1").length).toBeGreaterThan(0);
      expect(screen.getAllByText("2023").length).toBeGreaterThan(0);
    });

    // Verify pick chips are rendered
    expect(screen.getByText("favorite")).toBeInTheDocument();
    expect(screen.getByText("underdog")).toBeInTheDocument();
  });

  it("should handle API errors gracefully", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: { message: "Failed to load picks" },
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    await waitFor(() => {
      expect(
        screen.getByText("Error loading picks: Failed to load picks"),
      ).toBeInTheDocument();
    });
  });

  it("should show loading state when data is being fetched", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: null,
      isLoading: true,
    });

    render(<AdminPicksPage />, { wrapper });

    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("should handle API returning empty picks array", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: { picks: [], pagination: { total: 0, page: 1, limit: 20 } },
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    // Should show empty state message
    await waitFor(() => {
      expect(screen.getByText("No picks available")).toBeInTheDocument();
    });

    // Should not show loading spinner
    expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();

    // Should not show error message
    expect(screen.queryByText(/Error loading picks/i)).not.toBeInTheDocument();
  });

  it("should handle authentication errors from API", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: {
        message: "Unauthorized",
        response: {
          status: 401,
          data: { message: "Authentication required" },
        },
      },
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    // Should show authentication error message
    await waitFor(() => {
      expect(
        screen.getByText("Error loading picks: Unauthorized"),
      ).toBeInTheDocument();
    });

    // Should not show loading spinner
    expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();

    // Should not show empty state
    expect(screen.queryByText("No picks available")).not.toBeInTheDocument();
  });

  it("should handle permission errors from API", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: {
        message: "Forbidden",
        response: {
          status: 403,
          data: { message: "Insufficient permissions" },
        },
      },
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    // Should show permission error message
    await waitFor(() => {
      expect(
        screen.getByText("Error loading picks: Forbidden"),
      ).toBeInTheDocument();
    });

    // Should not show loading spinner
    expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();

    // Should not show empty state
    expect(screen.queryByText("No picks available")).not.toBeInTheDocument();
  });

  it("should handle network errors from API", async () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: undefined,
      error: {
        message: "Network Error",
        code: "NETWORK_ERROR",
      },
      isLoading: false,
    });

    render(<AdminPicksPage />, { wrapper });

    // Should show network error message
    await waitFor(() => {
      expect(
        screen.getByText("Error loading picks: Network Error"),
      ).toBeInTheDocument();
    });

    // Should not show loading spinner
    expect(screen.queryByRole("progressbar")).not.toBeInTheDocument();

    // Should not show empty state
    expect(screen.queryByText("No picks available")).not.toBeInTheDocument();
  });
});
