import { useState, useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Typography,
  Alert,
} from "@mui/material";
import { useCreateGame, useUpdateGame } from "../../services/api/games/games";
import type { GameRequest, GameResponse } from "../../services/model";

interface GameFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  game?: GameResponse | null;
}

const GameForm = ({ open, onClose, onSuccess, game }: GameFormProps) => {
  const [formData, setFormData] = useState<GameRequest>({
    week: 1,
    season: new Date().getFullYear(),
    home_team: "",
    away_team: "",
    spread: 0,
    start_time: "",
  });
  const [errors, setErrors] = useState<
    Partial<Record<keyof GameRequest, string>>
  >({});

  const createGameMutation = useCreateGame();
  const updateGameMutation = useUpdateGame();

  useEffect(() => {
    if (game) {
      // Convert ISO datetime to local format for datetime-local input
      const localStartTime = game.start_time
        ? new Date(game.start_time).toISOString().slice(0, 16)
        : "";

      setFormData({
        week: game.week,
        season: game.season,
        home_team: game.home_team,
        away_team: game.away_team,
        spread: game.spread,
        start_time: localStartTime,
      });
    } else {
      setFormData({
        week: 1,
        season: new Date().getFullYear(),
        home_team: "",
        away_team: "",
        spread: 0,
        start_time: "",
      });
    }
    setErrors({});
  }, [game]);

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof GameRequest, string>> = {};

    if (!formData.week || formData.week < 1 || formData.week > 18) {
      newErrors.week = "Week must be between 1 and 18";
    }

    if (!formData.season || formData.season < 2000) {
      newErrors.season = "Season must be a valid year";
    }

    if (!formData.home_team.trim()) {
      newErrors.home_team = "Favorite team is required";
    }

    if (!formData.away_team.trim()) {
      newErrors.away_team = "Underdog team is required";
    } else if (formData.home_team === formData.away_team) {
      newErrors.away_team = "Teams cannot be the same";
    }

    if (formData.spread === undefined || formData.spread === null) {
      newErrors.spread = "Spread is required";
    }

    if (!formData.start_time) {
      newErrors.start_time = "Start time is required";
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
        start_time: new Date(formData.start_time).toISOString(),
      };

      if (game) {
        await updateGameMutation.mutateAsync({
          id: game.id,
          data: submitData,
        });
      } else {
        await createGameMutation.mutateAsync({
          data: [submitData],
        });
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Error saving game:", error);
    }
  };

  const handleInputChange = (
    field: keyof GameRequest,
    value: string | number,
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
    createGameMutation.isPending || updateGameMutation.isPending;
  const error = createGameMutation.error || updateGameMutation.error;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>{game ? "Edit Game" : "Add New Game"}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} mt={1}>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error.message || "An error occurred while saving the game"}
              </Alert>
            )}

            <Box display="flex" gap={2}>
              <TextField
                label="Week"
                type="number"
                value={formData.week}
                onChange={(e) =>
                  handleInputChange("week", parseInt(e.target.value, 10) || 1)
                }
                error={!!errors.week}
                helperText={errors.week}
                fullWidth
                inputProps={{
                  min: 1,
                  max: 18,
                  "data-testid": "week-input",
                }}
              />
              <TextField
                label="Season"
                type="number"
                value={formData.season}
                onChange={(e) =>
                  handleInputChange(
                    "season",
                    parseInt(e.target.value, 10) || new Date().getFullYear(),
                  )
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
                label="Favorite Team"
                value={formData.home_team}
                onChange={(e) =>
                  handleInputChange("home_team", e.target.value)
                }
                error={!!errors.home_team}
                helperText={errors.home_team}
                fullWidth
                inputProps={{
                  "data-testid": "favorite-team-input",
                }}
              />
              <TextField
                label="Underdog Team"
                value={formData.away_team}
                onChange={(e) =>
                  handleInputChange("away_team", e.target.value)
                }
                error={!!errors.away_team}
                helperText={errors.away_team}
                fullWidth
                inputProps={{
                  "data-testid": "underdog-team-input",
                }}
              />
            </Box>

            <Box display="flex" gap={2}>
              <TextField
                label="Spread"
                type="number"
                value={formData.spread}
                onChange={(e) =>
                  handleInputChange("spread", parseFloat(e.target.value) || 0)
                }
                error={!!errors.spread}
                helperText={
                  errors.spread || "Positive number means favorite is favored"
                }
                fullWidth
                inputProps={{
                  step: 0.5,
                  "data-testid": "spread-input",
                }}
              />
              <TextField
                label="Start Time"
                type="datetime-local"
                value={formData.start_time}
                onChange={(e) =>
                  handleInputChange("start_time", e.target.value)
                }
                error={!!errors.start_time}
                helperText={errors.start_time}
                fullWidth
                inputProps={{
                  "data-testid": "start-time-input",
                }}
                InputLabelProps={{
                  shrink: true,
                }}
              />
            </Box>

            {formData.spread !== 0 && (
              <Typography variant="body2" color="text.secondary">
                Spread: {formData.home_team}{" "}
                {formData.spread > 0 ? "-" : "+"} {Math.abs(formData.spread)} vs{" "}
                {formData.away_team}
              </Typography>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={isLoading}>
            {isLoading ? "Saving..." : game ? "Update Game" : "Create Game"}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default GameForm;
