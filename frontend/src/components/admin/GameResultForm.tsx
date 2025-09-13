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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import { useSubmitResult } from "../../services/api/results/results";
import type { ResultRequest, GameResponse } from "../../services/model";

interface GameResultFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  game: GameResponse;
}

const GameResultForm = ({
  open,
  onClose,
  onSuccess,
  game,
}: GameResultFormProps) => {
  const [formData, setFormData] = useState<ResultRequest>({
    game_id: game.id,
    favorite_score: 0,
    underdog_score: 0,
    outcome: "",
  });
  const [errors, setErrors] = useState<
    Partial<Record<keyof ResultRequest, string>>
  >({});
  const [calculatedOutcome, setCalculatedOutcome] = useState<string>("");

  const submitResultMutation = useSubmitResult();

  useEffect(() => {
    if (open) {
      setFormData({
        game_id: game.id,
        favorite_score: 0,
        underdog_score: 0,
        outcome: "",
      });
      setErrors({});
      setCalculatedOutcome("");
    }
  }, [open, game]);

  useEffect(() => {
    if (
      formData.favorite_score !== undefined &&
      formData.underdog_score !== undefined
    ) {
      const favScore = formData.favorite_score;
      const udScore = formData.underdog_score;
      const spread = game.spread;

      if (favScore === udScore) {
        setCalculatedOutcome("TIE");
      } else {
        const favAdjusted = favScore + spread;
        if (favAdjusted > udScore) {
          setCalculatedOutcome("FAVORITE");
        } else if (favAdjusted < udScore) {
          setCalculatedOutcome("UNDERDOG");
        } else {
          setCalculatedOutcome("PUSH");
        }
      }
    }
  }, [formData.favorite_score, formData.underdog_score, game.spread]);

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof ResultRequest, string>> = {};

    if (formData.favorite_score === undefined || formData.favorite_score < 0) {
      newErrors.favorite_score = "Favorite score must be a non-negative number";
    }

    if (formData.underdog_score === undefined || formData.underdog_score < 0) {
      newErrors.underdog_score = "Underdog score must be a non-negative number";
    }

    if (!formData.outcome) {
      newErrors.outcome = "Outcome is required";
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
      await submitResultMutation.mutateAsync({
        data: formData,
      });
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Error submitting result:", error);
    }
  };

  const handleInputChange = (
    field: keyof ResultRequest,
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

  const isLoading = submitResultMutation.isPending;
  const error = submitResultMutation.error;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        Add Game Result - {game.favorite_team} vs {game.underdog_team}
      </DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} mt={1}>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error.message ||
                  "An error occurred while submitting the result"}
              </Alert>
            )}

            <Typography variant="body2" color="text.secondary">
              Week {game.week}, Season {game.season}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Spread: {game.favorite_team} {game.spread > 0 ? "-" : "+"}{" "}
              {Math.abs(game.spread)}
            </Typography>

            <Box display="flex" gap={2}>
              <TextField
                label={`${game.favorite_team} Score`}
                type="number"
                value={formData.favorite_score}
                onChange={(e) =>
                  handleInputChange(
                    "favorite_score",
                    parseInt(e.target.value, 10) || 0,
                  )
                }
                error={!!errors.favorite_score}
                helperText={errors.favorite_score}
                fullWidth
                inputProps={{
                  min: 0,
                }}
              />
              <TextField
                label={`${game.underdog_team} Score`}
                type="number"
                value={formData.underdog_score}
                onChange={(e) =>
                  handleInputChange(
                    "underdog_score",
                    parseInt(e.target.value, 10) || 0,
                  )
                }
                error={!!errors.underdog_score}
                helperText={errors.underdog_score}
                fullWidth
                inputProps={{
                  min: 0,
                }}
              />
            </Box>

            {calculatedOutcome && (
              <Alert severity="info">
                Calculated Outcome: {calculatedOutcome}
                {calculatedOutcome === "FAVORITE" &&
                  ` (${game.favorite_team} covers)`}
                {calculatedOutcome === "UNDERDOG" &&
                  ` (${game.underdog_team} covers)`}
                {calculatedOutcome === "PUSH" && " (Spread exactly matches)"}
                {calculatedOutcome === "TIE" && " (Scores are tied)"}
              </Alert>
            )}

            <FormControl fullWidth error={!!errors.outcome}>
              <InputLabel>Outcome</InputLabel>
              <Select
                value={formData.outcome}
                label="Outcome"
                onChange={(e) => handleInputChange("outcome", e.target.value)}
              >
                <MenuItem value="FAVORITE">Favorite Wins</MenuItem>
                <MenuItem value="UNDERDOG">Underdog Wins</MenuItem>
                <MenuItem value="PUSH">Push (Spread exactly matches)</MenuItem>
                <MenuItem value="TIE">Tie</MenuItem>
              </Select>
              {errors.outcome && (
                <Typography variant="caption" color="error">
                  {errors.outcome}
                </Typography>
              )}
            </FormControl>

            <Typography variant="body2" color="text.secondary">
              Final Score: {game.favorite_team} {formData.favorite_score} -{" "}
              {game.underdog_team} {formData.underdog_score}
            </Typography>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={isLoading}>
            {isLoading ? "Submitting..." : "Submit Result"}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default GameResultForm;
