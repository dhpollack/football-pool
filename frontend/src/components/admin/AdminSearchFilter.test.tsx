import { render, screen, fireEvent } from "@testing-library/react";
import AdminSearchFilter from "./AdminSearchFilter";

const filterFields = [
  {
    name: "status",
    label: "Status",
    type: "select" as const,
    options: [{ value: "active", label: "Active" }],
  },
  { name: "name", label: "Name", type: "text" as const },
];

describe("AdminSearchFilter", () => {
  const props = {
    searchTerm: "",
    onSearchChange: vi.fn(),
    filters: {},
    onFilterChange: vi.fn(),
    filterFields,
    onClearFilters: vi.fn(),
  };

  it("renders the search input and filter fields", () => {
    render(<AdminSearchFilter {...props} />);
    expect(screen.getByPlaceholderText("Search...")).toBeInTheDocument();
    expect(screen.getByLabelText("Status")).toBeInTheDocument();
    expect(screen.getByLabelText("Name")).toBeInTheDocument();
  });

  it("calls onSearchChange when the search input value changes", () => {
    render(<AdminSearchFilter {...props} />);
    const searchInput = screen.getByPlaceholderText("Search...");
    fireEvent.change(searchInput, { target: { value: "test" } });
    expect(props.onSearchChange).toHaveBeenCalledWith("test");
  });

  it("calls onFilterChange when a filter field value changes", () => {
    render(<AdminSearchFilter {...props} />);
    const statusFilter = screen.getByLabelText("Status");
    fireEvent.change(statusFilter, { target: { value: "active" } });
    expect(props.onFilterChange).toHaveBeenCalledWith("status", "active");
  });

  it('calls onClearFilters when the "Clear" button is clicked', () => {
    render(<AdminSearchFilter {...props} searchTerm="test" />);
    const clearButton = screen.getByText("Clear");
    fireEvent.click(clearButton);
    expect(props.onClearFilters).toHaveBeenCalledTimes(1);
  });

  it('disables the "Clear" button when there are no active filters', () => {
    render(<AdminSearchFilter {...props} />);
    const clearButton = screen.getByText("Clear");
    expect(clearButton).toBeDisabled();
  });
});
