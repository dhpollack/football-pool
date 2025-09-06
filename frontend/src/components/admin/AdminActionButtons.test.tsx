import { render, screen, fireEvent } from "@testing-library/react";
import AdminActionButtons from "./AdminActionButtons";

describe("AdminActionButtons", () => {
  it("renders the Add button when onAdd is provided", () => {
    const onAdd = vi.fn();
    render(<AdminActionButtons onAdd={onAdd} />);
    const addButton = screen.getByText("Add New");
    expect(addButton).toBeInTheDocument();
    fireEvent.click(addButton);
    expect(onAdd).toHaveBeenCalledTimes(1);
  });

  it("renders the Edit button when onEdit is provided", () => {
    const onEdit = vi.fn();
    render(<AdminActionButtons onEdit={onEdit} />);
    const editButton = screen.getByText("Edit");
    expect(editButton).toBeInTheDocument();
    fireEvent.click(editButton);
    expect(onEdit).toHaveBeenCalledTimes(1);
  });

  it("renders the Delete button when onDelete is provided", () => {
    const onDelete = vi.fn();
    render(<AdminActionButtons onDelete={onDelete} />);
    const deleteButton = screen.getByText("Delete");
    expect(deleteButton).toBeInTheDocument();
    fireEvent.click(deleteButton);
    expect(onDelete).toHaveBeenCalledTimes(1);
  });

  it("renders the View button when onView is provided", () => {
    const onView = vi.fn();
    render(<AdminActionButtons onView={onView} />);
    const viewButton = screen.getByText("View");
    expect(viewButton).toBeInTheDocument();
    fireEvent.click(viewButton);
    expect(onView).toHaveBeenCalledTimes(1);
  });

  it("disables the Add button when disabled.add is true", () => {
    const onAdd = vi.fn();
    render(<AdminActionButtons onAdd={onAdd} disabled={{ add: true }} />);
    const addButton = screen.getByText("Add New");
    expect(addButton).toBeDisabled();
  });

  it("renders custom labels for the buttons", () => {
    render(
      <AdminActionButtons
        onAdd={() => {}}
        addLabel="Create New"
        onEdit={() => {}}
        editLabel="Modify"
      />,
    );
    expect(screen.getByText("Create New")).toBeInTheDocument();
    expect(screen.getByText("Modify")).toBeInTheDocument();
  });
});
