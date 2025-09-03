
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import AdminNavigation from "./AdminNavigation";

describe("AdminNavigation", () => {
  const renderWithRouter = (path: string) => {
    return render(
      <MemoryRouter initialEntries={[path]}>
        <AdminNavigation />
      </MemoryRouter>
    );
  };

  it("renders the navigation links", () => {
    renderWithRouter("/admin");
    expect(screen.getByText("Dashboard")).toBeInTheDocument();
    expect(screen.getByText("Game Management")).toBeInTheDocument();
    expect(screen.getByText("User Management")).toBeInTheDocument();
    expect(screen.getByText("Pick Management")).toBeInTheDocument();
    expect(screen.getByText("Statistics")).toBeInTheDocument();
    expect(screen.getByText("Back to Main Site")).toBeInTheDocument();
  });

  it("highlights the active link", () => {
    renderWithRouter("/admin/games");
    const activeLink = screen.getByText("Game Management").closest("a");
    expect(activeLink).toHaveClass("Mui-selected");
  });

  it("links to the correct paths", () => {
    renderWithRouter("/admin");
    expect(screen.getByText("Dashboard").closest("a")).toHaveAttribute(
      "href",
      "/admin"
    );
    expect(screen.getByText("Game Management").closest("a")).toHaveAttribute(
      "href",
      "/admin/games"
    );
  });
});
