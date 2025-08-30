import { render, screen, fireEvent } from "@testing-library/react";
import { BrowserRouter as Router } from "react-router-dom";
import { vi } from "vitest";
import RegisterPage from "./RegisterPage";
import { api } from "../services/api";

// Mock the useNavigate hook
const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock("../services/api");

describe("RegisterPage", () => {
  beforeEach(() => {
    vi.mocked(api.post).mockResolvedValue({});
  });

  afterEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockClear();
  });

  it("renders the registration form", () => {
    render(
      <Router>
        <RegisterPage />
      </Router>,
    );

    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /register/i }),
    ).toBeInTheDocument();
  });

  it("submits the form", async () => {
    render(
      <Router>
        <RegisterPage />
      </Router>,
    );

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
      expect(api.post).toHaveBeenCalledWith("/api/register", {
        name: "John Doe",
        email: "john@example.com",
        password: "password",
      });
    });

    expect(mockNavigate).toHaveBeenCalledWith("/login", {
      state: { message: "Registration successful! Please log in." },
    });
  });
});
