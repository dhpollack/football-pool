import { render, screen, fireEvent } from "@testing-library/react";
import AdminConfirmDialog from "./AdminConfirmDialog";

describe("AdminConfirmDialog", () => {
  const props = {
    open: true,
    onClose: vi.fn(),
    onConfirm: vi.fn(),
    title: "Confirm Action",
    message: "Are you sure you want to perform this action?",
  };

  it("renders the dialog with title and message when open is true", () => {
    render(<AdminConfirmDialog {...props} />);
    expect(screen.getByText("Confirm Action")).toBeInTheDocument();
    expect(
      screen.getByText("Are you sure you want to perform this action?"),
    ).toBeInTheDocument();
  });

  it("does not render the dialog when open is false", () => {
    render(<AdminConfirmDialog {...props} open={false} />);
    expect(screen.queryByText("Confirm Action")).not.toBeInTheDocument();
  });

  it("calls the onClose handler when the cancel button is clicked", () => {
    render(<AdminConfirmDialog {...props} />);
    const cancelButton = screen.getByText("Cancel");
    fireEvent.click(cancelButton);
    expect(props.onClose).toHaveBeenCalledTimes(1);
  });

  it("calls the onConfirm handler when the confirm button is clicked", () => {
    render(<AdminConfirmDialog {...props} />);
    const confirmButton = screen.getByText("Confirm");
    fireEvent.click(confirmButton);
    expect(props.onConfirm).toHaveBeenCalledTimes(1);
  });

  it("displays a loading indicator when loading is true", () => {
    render(<AdminConfirmDialog {...props} loading={true} />);
    expect(screen.getByRole("progressbar")).toBeInTheDocument();
  });

  it("displays an error message when error is provided", () => {
    render(<AdminConfirmDialog {...props} error="An error occurred" />);
    expect(screen.getByText("An error occurred")).toBeInTheDocument();
  });

  it("disables the buttons when loading is true", () => {
    render(<AdminConfirmDialog {...props} loading={true} />);
    const cancelButton = screen.getByText("Cancel");
    const confirmButton = screen.getByText("Confirm");
    expect(cancelButton).toBeDisabled();
    expect(confirmButton).toBeDisabled();
  });
});
