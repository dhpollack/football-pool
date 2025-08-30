import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { BrowserRouter } from "react-router-dom";
import Layout from "./Layout";
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

// Helper component to provide auth context
const TestWrapper = ({
  children,
  user = null,
  loading = false,
  isAuthenticated = false,
}: {
  children: React.ReactNode;
  user?: any;
  loading?: boolean;
  isAuthenticated?: boolean;
}) => {
  const mockAuthContext = {
    user,
    login: vi.fn(),
    logout: vi.fn(),
    isAuthenticated,
    loading,
  };

  return (
    <BrowserRouter>
      <AuthContext.Provider value={mockAuthContext}>
        {children}
      </AuthContext.Provider>
    </BrowserRouter>
  );
};

describe("Layout", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows login button when user is not authenticated", async () => {
    render(
      <TestWrapper user={null} loading={false} isAuthenticated={false}>
        <Layout />
      </TestWrapper>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("login-button")).toBeInTheDocument();
    });
  });

  it("shows user menu when user is authenticated", async () => {
    render(
      <TestWrapper user={mockUser} loading={false} isAuthenticated={true}>
        <Layout />
      </TestWrapper>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("user-menu-button")).toBeInTheDocument();
    });
  });

  it("shows loading state properly", () => {
    render(
      <TestWrapper user={null} loading={true} isAuthenticated={false}>
        <Layout />
      </TestWrapper>,
    );

    // Should not show login button or user menu while loading
    expect(screen.queryByTestId("login-button")).not.toBeInTheDocument();
    expect(screen.queryByTestId("user-menu-button")).not.toBeInTheDocument();
  });

  it("shows register link in sidebar when not authenticated", async () => {
    render(
      <TestWrapper user={null} loading={false} isAuthenticated={false}>
        <Layout />
      </TestWrapper>,
    );

    await waitFor(() => {
      expect(screen.getByText("Register")).toBeInTheDocument();
    });
  });

  it("shows protected navigation items when authenticated", async () => {
    render(
      <TestWrapper user={mockUser} loading={false} isAuthenticated={true}>
        <Layout />
      </TestWrapper>,
    );

    await waitFor(() => {
      expect(screen.getByText("Picks")).toBeInTheDocument();
      expect(screen.getByText("Results")).toBeInTheDocument();
      expect(screen.getByText("Weekly Results")).toBeInTheDocument();
      expect(screen.getByText("Overall Results")).toBeInTheDocument();
      expect(screen.getByText("Survivor Pool")).toBeInTheDocument();
    });
  });

  it("hides register link when authenticated", async () => {
    render(
      <TestWrapper user={mockUser} loading={false} isAuthenticated={true}>
        <Layout />
      </TestWrapper>,
    );

    await waitFor(() => {
      expect(screen.queryByText("Register")).not.toBeInTheDocument();
    });
  });
});
