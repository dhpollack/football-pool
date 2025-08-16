import { render, screen, fireEvent } from "@testing-library/react";
import { vi } from "vitest";
import SurvivorPoolPage from "./SurvivorPoolPage";

describe("SurvivorPoolPage", () => {
  const mockTeams = [
    { id: 1, name: "Team A" },
    { id: 2, name: "Team B" },
  ];

  beforeEach(() => {
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockTeams),
      }),
    ) as vi.Mock;
  });

  it("renders available teams", async () => {
    render(<SurvivorPoolPage />);

    fireEvent.mouseDown(screen.getByLabelText(/select a team/i));
    expect(
      await screen.findByRole("option", { name: /team a/i }),
    ).toBeInTheDocument();
    expect(
      await screen.findByRole("option", { name: /team b/i }),
    ).toBeInTheDocument();
  });

  it("submits a survivor pick", async () => {
    render(<SurvivorPoolPage />);

    fireEvent.mouseDown(screen.getByLabelText(/select a team/i));
    fireEvent.click(await screen.findByRole("option", { name: /team a/i }));

    fireEvent.click(screen.getByText(/submit survivor pick/i));

    expect(global.fetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/survivor/picks",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ team: "Team A" }),
      }),
    );
  });
});
