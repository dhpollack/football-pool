
import { render, screen, fireEvent } from "@testing-library/react";
import AdminFormDialog from "./AdminFormDialog";

describe("AdminFormDialog", () => {
  const props = {
    open: true,
    onClose: vi.fn(),
    onSubmit: vi.fn(),
    title: "Test Form",
    children: <div>Form Content</div>,
  };

  it("renders the dialog with title and children when open is true", () => {
    render(<AdminFormDialog {...props} />);
    expect(screen.getByText("Test Form")).toBeInTheDocument();
    expect(screen.getByText("Form Content")).toBeInTheDocument();
  });

  it("does not render the dialog when open is false", () => {
    render(<AdminFormDialog {...props} open={false} />);
    expect(screen.queryByText("Test Form")).not.toBeInTheDocument();
  });

  it("calls the onClose handler when the cancel button is clicked", () => {
    render(<AdminFormDialog {...props} />);
    const cancelButton = screen.getByText("Cancel");
    fireEvent.click(cancelButton);
    expect(props.onClose).toHaveBeenCalledTimes(1);
  });

  it("calls the onSubmit handler when the form is submitted", () => {
    render(<AdminFormDialog {...props} />);
    const submitButton = screen.getByText("Save");
    fireEvent.click(submitButton);
    expect(props.onSubmit).toHaveBeenCalledTimes(1);
  });

  it("displays a loading indicator when loading is true", () => {
    render(<AdminFormDialog {...props} loading={true} />);
    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("displays an error message when error is provided", () => {
    render(<AdminFormDialog {...props} error="An error occurred" />);
    expect(screen.getByText("An error occurred")).toBeInTheDocument();
  });

  it("disables the buttons when loading is true", () => {
    render(<AdminFormDialog {...props} loading={true} />);
    const cancelButton = screen.getByText("Cancel");
    const submitButton = screen.getByText("Save");
    expect(cancelButton).toBeDisabled();
    expect(submitButton).toBeDisabled();
  });
});
