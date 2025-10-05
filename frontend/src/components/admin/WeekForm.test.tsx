import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import WeekForm from "./WeekForm";
import { useCreateWeek, useUpdateWeek } from "../../services/api/admin/admin";
import type { WeekResponse } from "../../services/model";

// Mock the custom hooks
vi.mock("../../services/api/admin/admin", () => ({
  useCreateWeek: vi.fn(),
  useUpdateWeek: vi.fn(),
}));

describe("WeekForm", () => {
  const mockOnClose = vi.fn();
  const mockOnSuccess = vi.fn();
  const mockMutateAsync = vi.fn();

  beforeEach(() => {
    (useCreateWeek as jest.Mock).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
      error: null,
    });
    (useUpdateWeek as jest.Mock).mockReturnValue({
      mutateAsync: mockMutateAsync,
      isPending: false,
      error: null,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe("Add Week Form", () => {
    it("renders the add week form with default values", () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      expect(screen.getByText("Add New Week")).toBeInTheDocument();
      expect(screen.getByTestId("week-number-input")).toHaveValue(1);
      expect(screen.getByTestId("season-input")).toHaveValue(
        new Date().getFullYear(),
      );
      expect(screen.getByTestId("week-start-time-input")).toHaveValue("");
      expect(screen.getByTestId("week-end-time-input")).toHaveValue("");
      expect(screen.getByRole("checkbox")).not.toBeChecked();
      expect(screen.getByText("Create Week")).toBeInTheDocument();
    });

    it("calls onClose when cancel button is clicked", () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      fireEvent.click(screen.getByText("Cancel"));
      expect(mockOnClose).toHaveBeenCalledTimes(1);
    });

    it("validates required fields and shows errors", async () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      fireEvent.click(screen.getByText("Create Week"));

      await waitFor(() => {
        expect(
          screen.getByText("Week start time is required"),
        ).toBeInTheDocument();
        expect(
          screen.getByText("Week end time is required"),
        ).toBeInTheDocument();
      });
    });

    it("validates week number range", async () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      const weekNumberInput = screen.getByTestId("week-number-input");
      const form = screen.getByTestId("week-form");

      // Test invalid week number 0
      fireEvent.change(weekNumberInput, { target: { value: "0" } });

      // Check that the input value changed
      expect(weekNumberInput).toHaveValue(0);

      // Submit the form
      fireEvent.submit(form);

      // Check for the validation error message
      await waitFor(() => {
        expect(
          screen.getByText("Week number must be between 1 and 18"),
        ).toBeInTheDocument();
      });

      // Test invalid week number 19
      fireEvent.change(weekNumberInput, { target: { value: "19" } });
      expect(weekNumberInput).toHaveValue(19);

      fireEvent.submit(form);

      await waitFor(() => {
        expect(
          screen.getByText("Week number must be between 1 and 18"),
        ).toBeInTheDocument();
      });
    });

    it("validates season year", async () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      const seasonInput = screen.getByTestId("season-input");
      const form = screen.getByTestId("week-form");

      fireEvent.change(seasonInput, { target: { value: "1999" } });

      // Check that the input value changed
      expect(seasonInput).toHaveValue(1999);

      fireEvent.submit(form);

      await waitFor(() => {
        expect(
          screen.getByText("Season must be a valid year"),
        ).toBeInTheDocument();
      });
    });

    it("validates end time is after start time", async () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      const startTimeInput = screen.getByTestId("week-start-time-input");
      const endTimeInput = screen.getByTestId("week-end-time-input");

      // Set start time to a future time
      const startTime = "2025-10-05T12:00";
      const endTime = "2025-10-05T11:00"; // Before start time

      fireEvent.change(startTimeInput, { target: { value: startTime } });
      fireEvent.change(endTimeInput, { target: { value: endTime } });
      fireEvent.click(screen.getByText("Create Week"));

      await waitFor(() => {
        // Look for helper text in Material-UI TextField
        const helperText = screen.getByText(
          "Week end time must be after start time",
        );
        expect(helperText).toBeInTheDocument();
      });
    });

    it("submits form with valid data and calls onSuccess", async () => {
      mockMutateAsync.mockResolvedValueOnce({});

      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      const weekNumberInput = screen.getByTestId("week-number-input");
      const seasonInput = screen.getByTestId("season-input");
      const startTimeInput = screen.getByTestId("week-start-time-input");
      const endTimeInput = screen.getByTestId("week-end-time-input");

      fireEvent.change(weekNumberInput, { target: { value: "5" } });
      fireEvent.change(seasonInput, { target: { value: "2025" } });
      fireEvent.change(startTimeInput, {
        target: { value: "2025-10-05T12:00" },
      });
      fireEvent.change(endTimeInput, { target: { value: "2025-10-05T18:00" } });

      fireEvent.click(screen.getByText("Create Week"));

      await waitFor(() => {
        // Use a more flexible assertion that checks the structure and relative timing
        // instead of exact UTC timestamps which are timezone-dependent
        expect(mockMutateAsync).toHaveBeenCalledWith({
          data: {
            week_number: 5,
            season: 2025,
            week_start_time: expect.stringMatching(/^2025-10-05T\d{2}:00:00\.000Z$/),
            week_end_time: expect.stringMatching(/^2025-10-05T\d{2}:00:00\.000Z$/),
            is_active: false,
          },
        });

        // Verify the timestamps are 6 hours apart (12:00 to 18:00 local time)
        const callData = mockMutateAsync.mock.calls[0][0];
        const startTime = new Date(callData.data.week_start_time);
        const endTime = new Date(callData.data.week_end_time);
        const timeDiff = endTime.getTime() - startTime.getTime();
        expect(timeDiff).toBe(6 * 60 * 60 * 1000); // 6 hours in milliseconds

        expect(mockOnSuccess).toHaveBeenCalledTimes(1);
        expect(mockOnClose).toHaveBeenCalledTimes(1);
      });
    });

    it("shows loading state during submission", async () => {
      (useCreateWeek as jest.Mock).mockReturnValue({
        mutateAsync: mockMutateAsync,
        isPending: true,
        error: null,
      });

      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      expect(screen.getByText("Saving...")).toBeInTheDocument();
      expect(screen.getByText("Cancel")).toBeDisabled();
      expect(screen.getByText("Saving...")).toBeDisabled();
    });

    it("shows error message when submission fails", async () => {
      const mockError = { message: "Failed to create week" };
      (useCreateWeek as jest.Mock).mockReturnValue({
        mutateAsync: mockMutateAsync,
        isPending: false,
        error: mockError,
      });

      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      expect(screen.getByText("Failed to create week")).toBeInTheDocument();
    });
  });

  describe("Edit Week Form", () => {
    const mockWeek: WeekResponse = {
      id: 1,
      week_number: 3,
      season: 2024,
      week_start_time: "2024-09-05T17:00:00Z",
      week_end_time: "2024-09-09T23:59:59Z",
      is_active: true,
    };

    it("renders the edit week form with week data", () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
          week={mockWeek}
        />,
      );

      expect(screen.getByText("Edit Week")).toBeInTheDocument();
      expect(screen.getByTestId("week-number-input")).toHaveValue(3);
      expect(screen.getByTestId("season-input")).toHaveValue(2024);
      expect(screen.getByRole("checkbox")).toBeChecked();
      expect(screen.getByText("Update Week")).toBeInTheDocument();
    });

    it("submits updated week data", async () => {
      mockMutateAsync.mockResolvedValueOnce({});

      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
          week={mockWeek}
        />,
      );

      const weekNumberInput = screen.getByTestId("week-number-input");
      fireEvent.change(weekNumberInput, { target: { value: "4" } });

      fireEvent.click(screen.getByText("Update Week"));

      await waitFor(() => {
        // Use flexible assertions for timezone-dependent timestamps
        expect(mockMutateAsync).toHaveBeenCalledWith({
          id: 1,
          data: {
            week_number: 4,
            season: 2024,
            week_start_time: expect.stringMatching(/^2024-09-05T\d{2}:00:00\.000Z$/),
            week_end_time: expect.stringMatching(/^2024-09-09T\d{2}:59:00\.000Z$/),
            is_active: true,
          },
        });
        expect(mockOnSuccess).toHaveBeenCalledTimes(1);
        expect(mockOnClose).toHaveBeenCalledTimes(1);
      });
    });

    it("clears errors when input values change", () => {
      render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
          week={mockWeek}
        />,
      );

      // First trigger validation errors
      const startTimeInput = screen.getByTestId("week-start-time-input");
      fireEvent.change(startTimeInput, { target: { value: "" } });
      fireEvent.click(screen.getByText("Update Week"));

      // Then fix the error
      fireEvent.change(startTimeInput, {
        target: { value: "2024-09-05T17:00" },
      });

      // Error should be cleared
      expect(
        screen.queryByText("Week start time is required"),
      ).not.toBeInTheDocument();
    });
  });

  describe("Form Behavior", () => {
    it("does not render when open is false", () => {
      render(
        <WeekForm
          open={false}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      expect(screen.queryByText("Add New Week")).not.toBeInTheDocument();
      expect(screen.queryByText("Edit Week")).not.toBeInTheDocument();
    });

    it("resets form when switching between add and edit modes", () => {
      const { rerender } = render(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />,
      );

      expect(screen.getByText("Add New Week")).toBeInTheDocument();

      const mockWeek: WeekResponse = {
        id: 1,
        week_number: 3,
        season: 2024,
        week_start_time: "2024-09-05T17:00:00Z",
        week_end_time: "2024-09-09T23:59:59Z",
        is_active: true,
      };

      rerender(
        <WeekForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
          week={mockWeek}
        />,
      );

      expect(screen.getByText("Edit Week")).toBeInTheDocument();
      expect(screen.getByTestId("week-number-input")).toHaveValue(3);
    });
  });
});
