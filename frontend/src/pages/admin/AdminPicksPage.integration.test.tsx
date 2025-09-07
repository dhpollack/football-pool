import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import AdminPicksPage from "./AdminPicksPage";
import { useAdminListPicks } from "../../services/api/picks/picks";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { PickResponse } from "../../services/model";

// Mock the custom hooks
vi.mock("../../services/api/picks/picks", () => ({
  useAdminListPicks: vi.fn(),
}));

describe("AdminPicksPage (Integration)", () => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );

  const mockPicks: PickResponse[] = [
    {
      id: 1,
      user: { id: 1, email: "user1@example.com", name: "User One", role: "user", created_at: "", updated_at: "" },
      game: { id: 1, favorite_team: "Team A", underdog_team: "Team B", week: 1, season: 2023, spread: 3.5, start_time: "", created_at: "", updated_at: "" },
      picked: "favorite",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: 2,
      user: { id: 2, email: "user2@example.com", name: "User Two", role: "user", created_at: "", updated_at: "" },
      game: { id: 2, favorite_team: "Team C", underdog_team: "Team D", week: 1, season: 2023, spread: 7, start_time: "", created_at: "", updated_at: "" },
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
      data: { picks: mockPicks, pagination: { total: mockPicks.length, page: 1, limit: 20 } },
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
      data: { picks: mockPicks, pagination: { total: mockPicks.length, page: 1, limit: 20 } },
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
      data: { picks: mockPicks, pagination: { total: mockPicks.length, page: 1, limit: 20 } },
      error: null,
      isLoading: false,
    });
    render(<AdminPicksPage />, { wrapper });

    const searchInput = screen.getByPlaceholderText("Search by user email...");
    fireEvent.change(searchInput, { target: { value: "nonexistent" } });

    await waitFor(() => {
      expect(screen.getByText("No picks match your filter criteria")).toBeInTheDocument();
    });
  });
});