import { render, screen, fireEvent } from "@testing-library/react";
import { vi } from "vitest";
import PickEntryPage from "./PickEntryPage";

describe("PickEntryPage", () => {
  const mockGames = [
    { id: 1, favorite_team: "Team A", underdog_team: "Team B", spread: 3 },
    { id: 2, favorite_team: "Team C", underdog_team: "Team D", spread: 7 },
  ];

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockGames),
      }),
    ) as vi.Mock;

    Object.defineProperty(window, "localStorage", {
      value: {
        getItem: vi.fn(() => "test-token"),
      },
      writable: true,
    });
  });

  it("renders the games", async () => {
    render(<PickEntryPage />);

    expect(await screen.findByText(/team a/i)).toBeInTheDocument();
    expect(await screen.findByText(/team b/i)).toBeInTheDocument();
    expect(await screen.findByText(/team c/i)).toBeInTheDocument();
    expect(await screen.findByText(/team d/i)).toBeInTheDocument();
  });

  it("handles quick pick", async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/quick pick/i));

    // After quick pick, the selects should have values
    const selects = await screen.findAllByRole("combobox");
    selects.forEach((select) => {
      expect(select).not.toHaveValue("");
    });
  });

  it("submits correct payload after quick pick", async () => {
    render(<PickEntryPage />);

    await screen.findByText(/team a/i);

    fireEvent.click(screen.getByText(/quick pick/i));
    fireEvent.click(screen.getByText(/submit picks/i));

    const fetchCalls = (global.fetch as vi.Mock).mock.calls;
    const picksCall = fetchCalls.find((call) => call[0].includes("/api/picks"));
    const actualPayload = JSON.parse(picksCall[1].body);

    actualPayload.sort((a: any, b: any) => a.game_id - b.game_id);

    const expectedPayload = mockGames.map((game) => ({
      game_id: game.id,
      picked: expect.stringMatching(/favorite|underdog/),
      rank: expect.any(Number),
      quick_pick: true,
    }));

    expect(actualPayload).toEqual(expectedPayload);

    expect(global.fetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/picks",
      expect.objectContaining({
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer test-token",
        },
      }),
    );
  });
});
