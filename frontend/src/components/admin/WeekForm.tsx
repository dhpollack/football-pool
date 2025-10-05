import { useState, useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Alert,
  FormControlLabel,
  Checkbox,
} from "@mui/material";
import { useCreateWeek, useUpdateWeek } from "../../services/api/admin/admin";
import type { WeekRequest, WeekResponse } from "../../services/model";

interface WeekFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  week?: WeekResponse | null;
}

const WeekForm = ({ open, onClose, onSuccess, week }: WeekFormProps) => {
  const [formData, setFormData] = useState<WeekRequest>({
    week_number: 1,
    season: new Date().getFullYear(),
    week_start_time: "",
    week_end_time: "",
    is_active: false,
  });
  const [errors, setErrors] = useState<
    Partial<Record<keyof WeekRequest, string>>
  >({});

  const createWeekMutation = useCreateWeek();
  const updateWeekMutation = useUpdateWeek();

  useEffect(() => {
    if (week) {
      // Convert ISO datetime to local format for datetime-local input
      const localStartTime = week.week_start_time
        ? new Date(week.week_start_time).toISOString().slice(0, 16)
        : "";
      const localEndTime = week.week_end_time
        ? new Date(week.week_end_time).toISOString().slice(0, 16)
        : "";

      setFormData({
        week_number: week.week_number,
        season: week.season,
        week_start_time: localStartTime,
        week_end_time: localEndTime,
        is_active: week.is_active,
      });
    } else {
      setFormData({
        week_number: 1,
        season: new Date().getFullYear(),
        week_start_time: "",
        week_end_time: "",
        is_active: false,
      });
    }
    setErrors({});
  }, [week]);

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof WeekRequest, string>> = {};

    if (
      !formData.week_number ||
      formData.week_number < 1 ||
      formData.week_number > 18
    ) {
      newErrors.week_number = "Week number must be between 1 and 18";
    }

    if (!formData.season || formData.season < 2000) {
      newErrors.season = "Season must be a valid year";
    }

    if (!formData.week_start_time) {
      newErrors.week_start_time = "Week start time is required";
    }

    if (!formData.week_end_time) {
      newErrors.week_end_time = "Week end time is required";
    }

    if (formData.week_start_time && formData.week_end_time) {
      const startTime = new Date(formData.week_start_time);
      const endTime = new Date(formData.week_end_time);
      if (endTime <= startTime) {
        newErrors.week_end_time = "Week end time must be after start time";
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      // Convert datetime-local to ISO 8601 format with timezone
      const submitData = {
        ...formData,
        week_start_time: new Date(formData.week_start_time).toISOString(),
        week_end_time: new Date(formData.week_end_time).toISOString(),
      };

      if (week) {
        await updateWeekMutation.mutateAsync({
          id: week.id,
          data: submitData,
        });
      } else {
        await createWeekMutation.mutateAsync({
          data: submitData,
        });
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Error saving week:", error);
    }
  };

  const handleInputChange = (
    field: keyof WeekRequest,
    value: string | number | boolean,
  ) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));

    if (errors[field]) {
      setErrors((prev) => ({
        ...prev,
        [field]: undefined,
      }));
    }
  };

  const isLoading =
    createWeekMutation.isPending || updateWeekMutation.isPending;
  const error = createWeekMutation.error || updateWeekMutation.error;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>{week ? "Edit Week" : "Add New Week"}</DialogTitle>
      <form onSubmit={handleSubmit} data-testid="week-form">
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} mt={1}>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error.message || "An error occurred while saving the week"}
              </Alert>
            )}

            <Box display="flex" gap={2}>
              <TextField
                label="Week Number"
                type="number"
                value={formData.week_number}
                onChange={(e) =>
                  handleInputChange(
                    "week_number",
                    parseInt(e.target.value, 10) || 0,
                  )
                }
                error={!!errors.week_number}
                helperText={errors.week_number}
                fullWidth
                inputProps={{
                  min: 1,
                  max: 18,
                  "data-testid": "week-number-input",
                }}
              />
              <TextField
                label="Season"
                type="number"
                value={formData.season}
                onChange={(e) =>
                  handleInputChange("season", parseInt(e.target.value, 10) || 0)
                }
                error={!!errors.season}
                helperText={errors.season}
                fullWidth
                inputProps={{
                  min: 2000,
                  max: 2100,
                  "data-testid": "season-input",
                }}
              />
            </Box>

            <Box display="flex" gap={2}>
              <TextField
                label="Week Start Time"
                type="datetime-local"
                value={formData.week_start_time}
                onChange={(e) =>
                  handleInputChange("week_start_time", e.target.value)
                }
                error={!!errors.week_start_time}
                helperText={errors.week_start_time}
                fullWidth
                inputProps={{
                  "data-testid": "week-start-time-input",
                }}
                InputLabelProps={{
                  shrink: true,
                }}
              />
              <TextField
                label="Week End Time"
                type="datetime-local"
                value={formData.week_end_time}
                onChange={(e) =>
                  handleInputChange("week_end_time", e.target.value)
                }
                error={!!errors.week_end_time}
                helperText={errors.week_end_time}
                fullWidth
                inputProps={{
                  "data-testid": "week-end-time-input",
                }}
                InputLabelProps={{
                  shrink: true,
                }}
              />
            </Box>

            <FormControlLabel
              control={
                <Checkbox
                  checked={formData.is_active}
                  onChange={(e) =>
                    handleInputChange("is_active", e.target.checked)
                  }
                />
              }
              label="Active Week"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={isLoading}>
            {isLoading ? "Saving..." : week ? "Update Week" : "Create Week"}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default WeekForm;
