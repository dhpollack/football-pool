import { render, screen, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { BrowserRouter } from "react-router-dom";
import Layout from "./Layout";

// Mock the useAuth hook
vi.mock("../hooks/useAuth", () => ({
  useAuth: vi.fn(),
}));

import { useAuth } from "../hooks/useAuth";

// Mock user data
const mockUser = {
  id: 1,
  name: "Test User",
  email: "test@example.com",
  role: "player",
};

const mockAdminUser = {
  id: 2,
  name: "Admin User",
  email: "admin@example.com",
  role: "admin",
};

describe("Layout", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows login button when user is not authenticated", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: false,
      isAdmin: false,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("login-button")).toBeInTheDocument();
    });
  });

  it("shows user menu when user is authenticated", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: true,
      isAdmin: false,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByTestId("user-menu-button")).toBeInTheDocument();
    });
  });

  it("shows loading state properly", () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: false,
      isAdmin: false,
      loading: true,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    // Should not show login button or user menu while loading
    expect(screen.queryByTestId("login-button")).not.toBeInTheDocument();
    expect(screen.queryByTestId("user-menu-button")).not.toBeInTheDocument();
  });

  it("shows register link in sidebar when not authenticated", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: false,
      isAdmin: false,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Register")).toBeInTheDocument();
    });
  });

  it("shows protected navigation items when authenticated as admin", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: mockAdminUser,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: true,
      isAdmin: true,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Picks")).toBeInTheDocument();
      expect(screen.getByText("Results")).toBeInTheDocument();
      expect(screen.getByText("Weekly Results")).toBeInTheDocument();
      expect(screen.getByText("Overall Results")).toBeInTheDocument();
      expect(screen.getByText("Survivor Pool")).toBeInTheDocument();
    });
  });

  it("hides results link when authenticated as regular user", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: true,
      isAdmin: false,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.getByText("Picks")).toBeInTheDocument();
      expect(screen.queryByText("Results")).not.toBeInTheDocument();
      expect(screen.getByText("Weekly Results")).toBeInTheDocument();
      expect(screen.getByText("Overall Results")).toBeInTheDocument();
      expect(screen.getByText("Survivor Pool")).toBeInTheDocument();
    });
  });

  it("hides register link when authenticated", async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      login: vi.fn(),
      logout: vi.fn(),
      isAuthenticated: true,
      isAdmin: false,
      loading: false,
    });

    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>,
    );

    await waitFor(() => {
      expect(screen.queryByText("Register")).not.toBeInTheDocument();
    });
  });
});
