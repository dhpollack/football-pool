
import { render, screen, fireEvent } from "@testing-library/react";
import AdminDataTable from "./AdminDataTable";

const columns = [
  { id: "name", label: "Name" },
  { id: "email", label: "Email" },
];

const data = [
  { id: 1, name: "John Doe", email: "john.doe@example.com" },
  { id: 2, name: "Jane Doe", email: "jane.doe@example.com" },
];

describe("AdminDataTable", () => {
  const props = {
    columns,
    data,
    page: 0,
    rowsPerPage: 10,
    totalCount: 2,
    onPageChange: vi.fn(),
    onRowsPerPageChange: vi.fn(),
  };

  it("renders the table with columns and data", () => {
    render(<AdminDataTable {...props} />);
    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Email")).toBeInTheDocument();
    expect(screen.getByText("John Doe")).toBeInTheDocument();
    expect(screen.getByText("jane.doe@example.com")).toBeInTheDocument();
  });

  it("displays a loading indicator when loading is true", () => {
    render(<AdminDataTable {...props} loading={true} />);
    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("displays an error message when error is provided", () => {
    render(<AdminDataTable {...props} error="Failed to fetch data" />);
    expect(screen.getByText("Failed to fetch data")).toBeInTheDocument();
  });

  it("displays an empty message when there is no data", () => {
    render(<AdminDataTable {...props} data={[]} />);
    expect(screen.getByText("No data available")).toBeInTheDocument();
  });

  it("calls onPageChange when the page is changed", () => {
    render(<AdminDataTable {...props} totalCount={11} />);
    const nextPageButton = screen.getByRole("button", { name: /next page/i });
    fireEvent.click(nextPageButton);
    expect(props.onPageChange).toHaveBeenCalledWith(1);
  });

  it("calls onRowsPerPageChange when the rows per page is changed", () => {
    render(<AdminDataTable {...props} />);
    const rowsPerPageSelect = screen.getByLabelText(/rows per page/i);
    fireEvent.mouseDown(rowsPerPageSelect);
    const option = screen.getByRole('option', { name: '25' });
    fireEvent.click(option);
    expect(props.onRowsPerPageChange).toHaveBeenCalledWith(25);
  });
});
