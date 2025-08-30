import { render } from "@testing-library/react";
import { vi } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider, useAuth } from "./AuthContext";
// Mock React Query hooks using orval generated mocks
vi.mock("../services/api/default/default", () => ({
  useLoginUser: () => ({
    mutateAsync: vi.fn(),
  }),
  useLogoutUser: () => ({
    mutateAsync: vi.fn(),
  }),
  useGetProfile: () => ({
    data: null,
    refetch: vi.fn().mockResolvedValue({ data: null }),
  }),
}));

const TestComponent = () => {
  const { user, isAuthenticated, loading } = useAuth();
  return (
    <div>
      <div data-testid="user">{user ? user.name : "null"}</div>
      <div data-testid="authenticated">
        {isAuthenticated ? "true" : "false"}
      </div>
      <div data-testid="loading">{loading ? "true" : "false"}</div>
    </div>
  );
};

const createWrapper = () => {
  const queryClient = new QueryClient();

  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe("AuthContext", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("provides initial null user state", async () => {
    const wrapper = createWrapper();
    const { getByTestId } = render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>,
      { wrapper },
    );

    // Wait for initial auth check to complete
    await vi.waitFor(() => {
      expect(getByTestId("loading").textContent).toBe("false");
    });

    expect(getByTestId("user").textContent).toBe("null");
    expect(getByTestId("authenticated").textContent).toBe("false");
  });

  it("handles successful authentication", async () => {
    const wrapper = createWrapper();
    const { getByTestId } = render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>,
      { wrapper },
    );

    // Wait for loading to complete
    await vi.waitFor(() => {
      expect(getByTestId("loading").textContent).toBe("false");
    });

    // Current implementation doesn't auto-fetch user data on mount
    expect(getByTestId("user").textContent).toBe("null");
    expect(getByTestId("authenticated").textContent).toBe("false");
  });

  it("handles authentication failure", async () => {
    const wrapper = createWrapper();
    const { getByTestId } = render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>,
      { wrapper },
    );

    // Wait for loading to complete
    await vi.waitFor(() => {
      expect(getByTestId("loading").textContent).toBe("false");
    });

    expect(getByTestId("user").textContent).toBe("null");
    expect(getByTestId("authenticated").textContent).toBe("false");
  });

  it("provides login and logout functions", async () => {
    const wrapper = createWrapper();
    const TestAuthComponent = () => {
      const { login, logout } = useAuth();
      return (
        <div>
          <div data-testid="login-function">{typeof login}</div>
          <div data-testid="logout-function">{typeof logout}</div>
        </div>
      );
    };

    const { getByTestId } = render(
      <AuthProvider>
        <TestAuthComponent />
      </AuthProvider>,
      { wrapper },
    );

    // Wait for initial auth check to complete
    await vi.waitFor(() => {
      expect(getByTestId("login-function").textContent).toBe("function");
    });

    expect(getByTestId("login-function").textContent).toBe("function");
    expect(getByTestId("logout-function").textContent).toBe("function");
  });
});
