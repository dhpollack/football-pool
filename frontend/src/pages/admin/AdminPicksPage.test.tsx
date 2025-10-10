import { screen, fireEvent } from "@testing-library/react";
import { render } from "../../test-utils";
import type { PickResponse } from "../../services/model";
import AdminPicksPage from "./AdminPicksPage";
import { useAdminListPicks } from "../../services/api/picks/picks";
import { getGetPicksResponseMock } from "../../services/api/picks/picks.msw";

// Mock the custom hooks
vi.mock("../../services/api/picks/picks", async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    useAdminListPicks: vi.fn(),
    useAdminDeletePick: vi.fn(),
  };
});

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({
    data,
    loading,
    error,
    onPageChange,
    onRowsPerPageChange,
    page,
    rowsPerPage,
    totalCount,
    columns,
  }: {
    data: PickResponse[];
    loading: boolean;
    error: string | null;
    onPageChange: (page: number) => void;
    onRowsPerPageChange: (rowsPerPage: number) => void;
    page: number;
    rowsPerPage: number;
    totalCount: number;
    columns: Array<{
      id: string;
      format?: (value: PickResponse) => React.ReactNode;
    }>;
  }) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading...</div>}
      {error && <div>Error: {error}</div>}
      {data && (
        <div>
          {data.map((pick: PickResponse) => (
            <div key={pick.id} data-testid="pick-row">
              <span data-testid="pick-user">{pick.user?.email}</span>
              <span data-testid="pick-game">
                {pick.game
                  ? `${pick.game.home_team} vs ${pick.game.away_team}`
                  : "Unknown"}
              </span>
              <span data-testid="pick-choice">{pick.picked}</span>
              {columns.find((c) => c.id === "actions")?.format?.(pick)}
            </div>
          ))}
        </div>
      )}
      <div>
        <span>Rows per page:</span>
        <select
          data-testid="rows-per-page-select"
          value={rowsPerPage}
          onChange={(e) => onRowsPerPageChange(Number(e.target.value))}
        >
          <option value={10}>10</option>
          <option value={25}>25</option>
          <option value={50}>50</option>
        </select>
        <span>
          {page * rowsPerPage + 1}-
          {Math.min((page + 1) * rowsPerPage, totalCount)} of {totalCount}
        </span>
        <button
          type="button"
          onClick={() => onPageChange(page - 1)}
          disabled={page === 0}
        >
          Previous Page
        </button>
        <button
          type="button"
          onClick={() => onPageChange(page + 1)}
          disabled={(page + 1) * rowsPerPage >= totalCount}
        >
          Next Page
        </button>
      </div>
    </div>
  ),
}));

vi.mock("../../components/admin/AdminSearchFilter", () => ({
  default: () => <div data-testid="search-filter">Search Filter</div>,
}));

vi.mock("../../components/admin/AdminConfirmDialog", () => ({
  default: ({
    open,
    onConfirm,
    onClose,
  }: {
    open: boolean;
    onConfirm: () => void;
    onClose: () => void;
  }) =>
    open ? (
      <div data-testid="confirm-dialog">
        <button type="button" onClick={onConfirm}>
          Confirm
        </button>
        <button type="button" onClick={onClose}>
          Cancel
        </button>
      </div>
    ) : null,
}));

describe("AdminPicksPage", () => {
  beforeEach(() => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
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
    (useAdminListPicks as vi.Mock).mockReturnValue({
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

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    render(<AdminPicksPage />);

    expect(screen.getByText("user1@example.com")).toBeInTheDocument();
    expect(screen.getByText("Team A vs Team B")).toBeInTheDocument();
  });

  it("shows error message when there is an error", () => {
    (useAdminListPicks as vi.Mock).mockReturnValue({
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

  it("handles page changes", () => {
    const mockData = getGetPicksResponseMock({
      picks: new Array(30).fill(null).map(
        (_, i) =>
          ({
            id: i,
            user: { email: `user${i}@example.com` },
          }) as PickResponse,
      ),
      pagination: { page: 1, limit: 25, total: 30 },
    });

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    const { rerender } = render(<AdminPicksPage />);

    const nextPageButton = screen.getByRole("button", { name: /Next Page/i });
    fireEvent.click(nextPageButton);

    rerender(<AdminPicksPage />);

    expect(useAdminListPicks).toHaveBeenCalledWith({
      request: {
        params: {
          user_id: undefined,
          game_id: undefined,
          week: undefined,
          season: undefined,
          page: 2,
          limit: 25,
        },
      },
    });
  });

  it("handles rows per page changes", () => {
    const mockData = getGetPicksResponseMock({
      picks: new Array(30).fill(null).map(
        (_, i) =>
          ({
            id: i,
            user: { email: `user${i}@example.com` },
          }) as PickResponse,
      ),
      pagination: { page: 1, limit: 25, total: 30 },
    });

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    const { rerender } = render(<AdminPicksPage />);

    const rowsPerPageSelect = screen.getByTestId("rows-per-page-select");
    fireEvent.change(rowsPerPageSelect, { target: { value: "10" } });

    rerender(<AdminPicksPage />);

    expect(useAdminListPicks).toHaveBeenCalledWith({
      request: {
        params: {
          user_id: undefined,
          game_id: undefined,
          week: undefined,
          season: undefined,
          page: 1,
          limit: 10,
        },
      },
    });
  });

  it("opens delete confirmation dialog on delete button click", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: { email: "user1@example.com" },
          game: { home_team: "Team A", away_team: "Team B" },
          picked: "favorite",
        } as PickResponse,
      ],
      pagination: { page: 1, limit: 25, total: 1 },
    });

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    const { rerender } = render(<AdminPicksPage />);

    const deleteButton = screen.getByTestId("DeleteIcon");
    fireEvent.click(deleteButton);

    rerender(<AdminPicksPage />);

    expect(screen.getByTestId("confirm-dialog")).toBeInTheDocument();
  });

  it("calls handleDeleteConfirm on confirm button click", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: { email: "user1@example.com" },
          game: { home_team: "Team A", away_team: "Team B" },
          picked: "favorite",
        } as PickResponse,
      ],
      pagination: { page: 1, limit: 25, total: 1 },
    });

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    const { rerender } = render(<AdminPicksPage />);

    const deleteButton = screen.getByTestId("DeleteIcon");
    fireEvent.click(deleteButton);

    rerender(<AdminPicksPage />);

    const confirmButton = screen.getByRole("button", { name: /confirm/i });
    fireEvent.click(confirmButton);
  });

  it("closes delete confirmation dialog on cancel button click", () => {
    const mockData = getGetPicksResponseMock({
      picks: [
        {
          id: 1,
          user: { email: "user1@example.com" },
          game: { home_team: "Team A", away_team: "Team B" },
          picked: "favorite",
        } as PickResponse,
      ],
      pagination: { page: 1, limit: 25, total: 1 },
    });

    (useAdminListPicks as vi.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    const { rerender } = render(<AdminPicksPage />);

    const deleteButton = screen.getByTestId("DeleteIcon");
    fireEvent.click(deleteButton);

    rerender(<AdminPicksPage />);

    const cancelButton = screen.getByRole("button", { name: /cancel/i });
    fireEvent.click(cancelButton);

    rerender(<AdminPicksPage />);

    expect(screen.queryByTestId("confirm-dialog")).not.toBeInTheDocument();
  });
});
