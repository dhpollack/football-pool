
import { render, screen } from "@testing-library/react";
import AdminDashboard from "./AdminDashboard";

describe("AdminDashboard", () => {
  it("renders the admin dashboard with title and stats", () => {
    render(<AdminDashboard />);

    // Check for the main title
    expect(
      screen.getByText("Admin Dashboard", { selector: "h4" })
    ).toBeInTheDocument();

    // Check for the welcome message
    expect(
      screen.getByText("Welcome to the Football Pool Administration Panel")
    ).toBeInTheDocument();

    // Check for the stats titles
    expect(screen.getByText("Total Games")).toBeInTheDocument();
    expect(screen.getByText("Total Users")).toBeInTheDocument();
    expect(screen.getByText("Total Picks")).toBeInTheDocument();
    expect(screen.getByText("Active Weeks")).toBeInTheDocument();

    // Check for the stats values (which are all 0 in the mock data)
    expect(screen.getAllByText("0")).toHaveLength(4);

    // Check for the section titles
    expect(screen.getByText("Recent Activity")).toBeInTheDocument();
    expect(screen.getByText("Quick Actions")).toBeInTheDocument();
  });
});
