import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { vi } from "vitest";
import { AuthProvider } from "../contexts/AuthContext";
import LoginPage from "./LoginPage";

// Mock the React Query hooks used by AuthContext
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

describe("LoginPage", () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
  });

  it("renders the login form", () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <AuthProvider>
            <LoginPage />
          </AuthProvider>
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /sign in/i }),
    ).toBeInTheDocument();
  });
});