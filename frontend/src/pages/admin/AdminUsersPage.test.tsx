import { render, screen } from "@testing-library/react";
import AdminUsersPage from "./AdminUsersPage";
import { useAdminListUsers } from "../../services/api/user/user";
import { getAdminListUsersResponseMock } from "../../services/api/user/user.msw";

// Mock the custom hooks
vi.mock("../../services/api/user/user", async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    useAdminListUsers: vi.fn(),
    useCreateUsers: vi.fn(),
  };
});

// Mock the admin components with simple implementations
vi.mock("../../components/admin/AdminDataTable", () => ({
  default: ({ data, loading, error }: any) => (
    <div data-testid="admin-data-table">
      {loading && <div>Loading...</div>}
      {error && <div>Error: {error}</div>}
      {data && (
        <div>
          {data.map((user: any) => (
            <div key={user.id} data-testid="user-row">
              <span data-testid="user-email">{user.email}</span>
              <span data-testid="user-name">{user.name}</span>
              <span data-testid="user-role">{user.role}</span>
              <span data-testid="user-picks">{user.pick_count}</span>
              <span data-testid="user-wins">{user.total_wins}</span>
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

vi.mock("../../components/admin/AdminActionButtons", () => ({
  default: () => <div data-testid="action-buttons">Action Buttons</div>,
}));

vi.mock("../../components/admin/AdminConfirmDialog", () => ({
  default: ({ open }: any) =>
    open && <div data-testid="confirm-dialog">Confirm Dialog</div>,
}));

describe("AdminUsersPage", () => {
  beforeEach(() => {
    (useAdminListUsers as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: false,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders the page title", () => {
    render(<AdminUsersPage />);

    expect(screen.getByText("User Management")).toBeInTheDocument();
  });

  it("renders loading state", () => {
    (useAdminListUsers as jest.Mock).mockReturnValue({
      data: null,
      error: null,
      isLoading: true,
    });

    render(<AdminUsersPage />);

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders users data when loaded", () => {
    const mockData = getAdminListUsersResponseMock({
      users: [
        {
          id: 1,
          email: "admin@example.com",
          name: "Admin User",
          role: "admin",
          pick_count: 10,
          total_wins: 5,
          created_at: "2023-09-10T12:00:00Z",
          updated_at: "2023-09-10T12:00:00Z",
        },
        {
          id: 2,
          email: "user@example.com",
          name: "Regular User",
          role: "user",
          pick_count: 8,
          total_wins: 3,
          created_at: "2023-09-10T13:00:00Z",
          updated_at: "2023-09-10T13:00:00Z",
        },
      ],
      pagination: {
        page: 1,
        limit: 10,
        total: 2,
      },
    });

    (useAdminListUsers as jest.Mock).mockReturnValue({
      data: mockData,
      error: null,
      isLoading: false,
    });

    render(<AdminUsersPage />);

    expect(screen.getByText("admin@example.com")).toBeInTheDocument();
    expect(screen.getByText("user@example.com")).toBeInTheDocument();
    expect(screen.getByText("Admin User")).toBeInTheDocument();
    expect(screen.getByText("Regular User")).toBeInTheDocument();
  });

  it("shows error message when there is an error", () => {
    (useAdminListUsers as jest.Mock).mockReturnValue({
      data: null,
      error: { message: "Failed to load users" },
      isLoading: false,
    });

    render(<AdminUsersPage />);

    expect(screen.getByText("Error: Failed to load users")).toBeInTheDocument();
  });

  it("renders search filter component", () => {
    render(<AdminUsersPage />);

    expect(screen.getByTestId("search-filter")).toBeInTheDocument();
  });

  it("renders action buttons component", () => {
    render(<AdminUsersPage />);

    expect(screen.getByTestId("action-buttons")).toBeInTheDocument();
  });
});
