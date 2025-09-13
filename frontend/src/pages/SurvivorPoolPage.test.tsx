import { screen, fireEvent, waitFor } from "@testing-library/react";
import { vi } from "vitest";
import { render } from "../test-utils";
import SurvivorPoolPage from "./SurvivorPoolPage";

// Mock fetch
global.fetch = vi.fn();

// Mock the React Query hook
const mockSubmitSurvivorPick = vi.fn();
vi.mock("../services/api/survivor/survivor", () => ({
  useSubmitSurvivorPick: () => ({
    mutateAsync: mockSubmitSurvivorPick,
    isPending: false,
  }),
}));

describe("SurvivorPoolPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Mock localStorage
    Object.defineProperty(window, "localStorage", {
      value: {
        getItem: vi.fn(() => "fake-token"),
      },
      writable: true,
    });
  });

  it("renders available teams", async () => {
    // Mock fetch response for teams
    (global.fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => [
        { id: 1, name: "Team A" },
        { id: 2, name: "Team B" },
      ],
    });

    render(<SurvivorPoolPage />);

    // Wait for teams to load
    await waitFor(() => {
      fireEvent.mouseDown(screen.getByLabelText(/select a team/i));
    });

    expect(
      await screen.findByRole("option", { name: /team a/i }),
    ).toBeInTheDocument();
    expect(
      await screen.findByRole("option", { name: /team b/i }),
    ).toBeInTheDocument();
  });

  it("submits a survivor pick", async () => {
    // Mock fetch response for teams
    (global.fetch as any).mockResolvedValueOnce({
      ok: true,
      json: async () => [
        { id: 1, name: "Team A" },
        { id: 2, name: "Team B" },
      ],
    });

    mockSubmitSurvivorPick.mockResolvedValueOnce({});

    render(<SurvivorPoolPage />);

    // Wait for teams to load and select a team
    await waitFor(() => {
      fireEvent.mouseDown(screen.getByLabelText(/select a team/i));
    });

    fireEvent.click(await screen.findByRole("option", { name: /team a/i }));

    await waitFor(() => {
      fireEvent.click(screen.getByText(/submit survivor pick/i));
    });

    expect(mockSubmitSurvivorPick).toHaveBeenCalledWith({
      data: {
        week: 1,
        team: "Team A",
      },
    });
  });
});
