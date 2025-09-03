
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router-dom";
import AdminLayout from "./AdminLayout";
import { useAuth } from "../../hooks/useAuth";

// Mock the useAuth hook
vi.mock("../../hooks/useAuth", () => ({
  useAuth: vi.fn(),
}));

// Mock the AdminNavigation component
vi.mock("./AdminNavigation", () => ({
  __esModule: true,
  default: () => <div data-testid="admin-navigation"></div>,
}));

describe("AdminLayout", () => {
  afterEach(() => {
    vi.clearAllMocks();
  });

  it('renders "Access Denied" if the user is not an admin', () => {
    (useAuth as jest.Mock).mockReturnValue({ isAdmin: false, loading: false });
    render(
      <MemoryRouter>
        <AdminLayout />
      </MemoryRouter>
    );
    expect(screen.getByText("Access Denied")).toBeInTheDocument();
    expect(
      screen.getByText("You do not have permission to access the admin area.")
    ).toBeInTheDocument();
  });

  it("renders the admin layout if the user is an admin", () => {
    (useAuth as jest.Mock).mockReturnValue({ isAdmin: true, loading: false });
    render(
      <MemoryRouter>
        <Routes>
          <Route path="/" element={<AdminLayout />}>
            <Route index element={<div>Admin Content</div>} />
          </Route>
        </Routes>
      </MemoryRouter>
    );

    expect(screen.getByText("Football Pool Admin")).toBeInTheDocument();
    expect(screen.getByTestId("admin-navigation")).toBeInTheDocument();
    expect(screen.getByText("Admin Content")).toBeInTheDocument();
  });

  it("renders nothing while loading and not admin", () => {
    (useAuth as jest.Mock).mockReturnValue({ isAdmin: false, loading: true });
    const { container } = render(
      <MemoryRouter>
        <AdminLayout />
      </MemoryRouter>
    );
    expect(container.firstChild).toBeNull();
  });
});
