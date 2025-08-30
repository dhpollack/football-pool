import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { BrowserRouter } from "react-router-dom";
import UserMenu from "./UserMenu";
import { AuthContext } from "../contexts/AuthContext";

// Mock the API module
vi.mock("../services/api", () => ({
  api: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

// Mock user data
const mockUser = {
  id: 1,
  name: "Test User",
  email: "test@example.com",
};

// Helper component to provide context
const TestWrapper = ({
  children,
  user = mockUser,
}: {
  children: React.ReactNode;
  user?: any;
}) => {
  // Mock AuthContext value
  const mockAuthContext = {
    user,
    login: vi.fn(),
    logout: vi.fn(),
    isAuthenticated: !!user,
    loading: false,
  };

  return (
    <BrowserRouter>
      <AuthContext.Provider value={mockAuthContext}>
        {children}
      </AuthContext.Provider>
    </BrowserRouter>
  );
};

describe("UserMenu", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders user menu button when user is authenticated", () => {
    render(
      <TestWrapper>
        <UserMenu />
      </TestWrapper>,
    );

    const userMenuButton = screen.getByTestId("user-menu-button");
    expect(userMenuButton).toBeInTheDocument();
  });

  it("opens menu when user menu button is clicked", async () => {
    render(
      <TestWrapper>
        <UserMenu />
      </TestWrapper>,
    );

    const userMenuButton = screen.getByTestId("user-menu-button");
    fireEvent.click(userMenuButton);

    await waitFor(() => {
      expect(screen.getByTestId("user-menu")).toBeVisible();
    });
  });

  it("displays user information in the menu", async () => {
    render(
      <TestWrapper>
        <UserMenu />
      </TestWrapper>,
    );

    const userMenuButton = screen.getByTestId("user-menu-button");
    fireEvent.click(userMenuButton);

    await waitFor(() => {
      expect(screen.getByText(mockUser.name)).toBeInTheDocument();
      expect(screen.getByText(mockUser.email)).toBeInTheDocument();
    });
  });

  it("shows profile link in the menu", async () => {
    render(
      <TestWrapper>
        <UserMenu />
      </TestWrapper>,
    );

    const userMenuButton = screen.getByTestId("user-menu-button");
    fireEvent.click(userMenuButton);

    await waitFor(() => {
      expect(screen.getByTestId("profile-link")).toBeInTheDocument();
    });
  });

  it("shows logout button in the menu", async () => {
    render(
      <TestWrapper>
        <UserMenu />
      </TestWrapper>,
    );

    const userMenuButton = screen.getByTestId("user-menu-button");
    fireEvent.click(userMenuButton);

    await waitFor(() => {
      expect(screen.getByTestId("logout-button")).toBeInTheDocument();
    });
  });
});
