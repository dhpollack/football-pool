import { render, screen, fireEvent } from "@testing-library/react";
import { vi } from "vitest";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import ProfilePage from "./ProfilePage";

import { getGetProfileResponseMock } from "../services/api/user/user.msw";

// Mock the React Query hooks
const mockProfileData = getGetProfileResponseMock({
  name: "John Doe",
  address: "123 Main St",
});

const mockUpdateProfile = vi.fn();
vi.mock("../services/api/user/user", () => ({
  useGetProfile: () => ({
    data: mockProfileData,
    error: null,
    isLoading: false,
  }),
  useUpdateProfile: () => ({
    mutateAsync: mockUpdateProfile,
    isPending: false,
  }),
}));

// Mock alert
global.alert = vi.fn();

describe("ProfilePage", () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });
    mockUpdateProfile.mockResolvedValue({});
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const renderWithProviders = () => {
    return render(
      <QueryClientProvider client={queryClient}>
        <ProfilePage />
      </QueryClientProvider>,
    );
  };

  it("renders the profile", async () => {
    renderWithProviders();

    expect(await screen.findByDisplayValue(/john doe/i)).toBeInTheDocument();
    expect(await screen.findByDisplayValue(/123 main st/i)).toBeInTheDocument();
  });

  it("updates the profile", async () => {
    renderWithProviders();

    await screen.findByDisplayValue(/john doe/i);

    fireEvent.change(screen.getByLabelText(/name/i), {
      target: { value: "Jane Doe" },
    });
    fireEvent.click(screen.getByText(/save/i));

    // Wait for the async operation to complete
    await vi.waitFor(() => {
      expect(mockUpdateProfile).toHaveBeenCalledWith({
        data: { name: "Jane Doe", address: "123 Main St" },
      });
    });

    expect(global.alert).toHaveBeenCalledWith("Profile updated successfully!");
  });
});
