import { render, screen, fireEvent } from "@testing-library/react";
import { BrowserRouter as Router } from "react-router-dom";
import { vi } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import RegisterPage from "./RegisterPage";

// Mock the useNavigate hook
const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Mock the React Query hook
const mockRegisterUser = vi.fn();
vi.mock("../services/api/default/default", () => ({
  useRegisterUser: () => ({
    mutateAsync: mockRegisterUser,
    isPending: false,
  }),
}));

describe("RegisterPage", () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    mockRegisterUser.mockResolvedValue({});
  });

  afterEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockClear();
  });

  const renderWithProviders = () => {
    return render(
      <QueryClientProvider client={queryClient}>
        <Router>
          <RegisterPage />
        </Router>
      </QueryClientProvider>,
    );
  };

  it("renders the registration form", () => {
    renderWithProviders();

    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /register/i }),
    ).toBeInTheDocument();
  });

  it("submits the form", async () => {
    renderWithProviders();

    fireEvent.change(screen.getByLabelText(/name/i), {
      target: { value: "John Doe" },
    });
    fireEvent.change(screen.getByLabelText(/email address/i), {
      target: { value: "john@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: "password" },
    });

    fireEvent.click(screen.getByRole("button", { name: /register/i }));

    // Wait for the async operation to complete
    await vi.waitFor(() => {
      expect(mockRegisterUser).toHaveBeenCalledWith({
        data: {
          name: "John Doe",
          email: "john@example.com",
          password: "password",
        },
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith("/login", {
      state: { message: "Registration successful! Please log in." },
    });
  });
});
