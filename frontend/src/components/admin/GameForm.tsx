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
    favorite_team: "",
    underdog_team: "",
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
      setFormData({
        week: game.week,
        season: game.season,
        favorite_team: game.favorite_team,
        underdog_team: game.underdog_team,
        spread: game.spread,
        start_time: game.start_time,
      });
    } else {
      setFormData({
        week: 1,
        season: new Date().getFullYear(),
        favorite_team: "",
        underdog_team: "",
        spread: 0,
        start_time: "",
      });
    }
    setErrors({});
  }, [game, open]);

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof GameRequest, string>> = {};

    if (!formData.week || formData.week < 1 || formData.week > 18) {
      newErrors.week = "Week must be between 1 and 18";
    }

    if (!formData.season || formData.season < 2000) {
      newErrors.season = "Season must be a valid year";
    }

    if (!formData.favorite_team.trim()) {
      newErrors.favorite_team = "Favorite team is required";
    }

    if (!formData.underdog_team.trim()) {
      newErrors.underdog_team = "Underdog team is required";
    } else if (formData.favorite_team === formData.underdog_team) {
      newErrors.underdog_team = "Teams cannot be the same";
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
      if (game) {
        await updateGameMutation.mutateAsync({
          id: game.id,
          data: formData,
        });
      } else {
        await createGameMutation.mutateAsync({
          data: [formData],
        });
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Error saving game:", error);
    }
  };

  const handleInputChange = (field: keyof GameRequest, value: any) => {
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
                }}
              />
            </Box>

            <Box display="flex" gap={2}>
              <TextField
                label="Favorite Team"
                value={formData.favorite_team}
                onChange={(e) =>
                  handleInputChange("favorite_team", e.target.value)
                }
                error={!!errors.favorite_team}
                helperText={errors.favorite_team}
                fullWidth
              />
              <TextField
                label="Underdog Team"
                value={formData.underdog_team}
                onChange={(e) =>
                  handleInputChange("underdog_team", e.target.value)
                }
                error={!!errors.underdog_team}
                helperText={errors.underdog_team}
                fullWidth
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
                InputLabelProps={{
                  shrink: true,
                }}
              />
            </Box>

            {formData.spread !== 0 && (
              <Typography variant="body2" color="text.secondary">
                Spread: {formData.favorite_team}{" "}
                {formData.spread > 0 ? "-" : "+"} {Math.abs(formData.spread)} vs{" "}
                {formData.underdog_team}
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
