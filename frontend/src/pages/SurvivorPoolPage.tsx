import { useState, useEffect, useId } from "react";
import {
  Button,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Box,
} from "@mui/material";
import { useSubmitSurvivorPick } from "../services/api/default/default";
import type { SurvivorPickRequest } from "../services/model";
import axios from "axios";

interface Team {
  id: number;
  name: string;
}

const SurvivorPoolPage = () => {
  const teamSelectLabelId = useId();
  const teamSelectId = useId();
  const [availableTeams, setAvailableTeams] = useState<Team[]>([]);
  const [selectedTeam, setSelectedTeam] = useState<string>("");
  const [error, setError] = useState<string | null>(null);

  // React Query hook for submitting survivor picks
  const { mutateAsync: submitSurvivorPick, isPending: isSubmitting } =
    useSubmitSurvivorPick();

  useEffect(() => {
    const fetchAvailableTeams = async () => {
      try {
        const token = localStorage.getItem("token");
        // Assuming an API endpoint to get available teams for survivor pool
        const response = await fetch(
          `${import.meta.env.VITE_BACKEND_URL}/api/survivor/teams`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        if (!response.ok) {
          throw new Error("Failed to fetch available teams");
        }

        const data = await response.json();
        setAvailableTeams(data);
      } catch (error: unknown) {
        if (error instanceof Error) {
          setError(error.message);
        } else {
          setError("An unknown error occurred");
        }
      }
    };

    fetchAvailableTeams();
  }, []);

  const handleSubmitPick = async () => {
    setError(null);
    try {
      const survivorPickData: SurvivorPickRequest = {
        week: 1, // TODO: Get current week dynamically
        team: selectedTeam,
      };

      await submitSurvivorPick({ data: survivorPickData });
      alert("Survivor pick submitted successfully!");
      setSelectedTeam(""); // Clear selection after successful submission
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        setError(error.response?.data?.message || error.message);
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  return (
    <Box>
      <Typography variant="h4">Survivor Pool</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <FormControl fullWidth sx={{ my: 2 }}>
        <InputLabel id={teamSelectLabelId}>Select a Team</InputLabel>
        <Select
          labelId={teamSelectLabelId}
          id={teamSelectId}
          value={selectedTeam}
          label="Select a Team"
          onChange={(e) => setSelectedTeam(e.target.value as string)}
        >
          {availableTeams.map((team) => (
            <MenuItem key={team.id} value={team.name}>
              {team.name}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
      <Button
        variant="contained"
        onClick={handleSubmitPick}
        disabled={!selectedTeam || isSubmitting}
      >
        {isSubmitting ? "Submitting..." : "Submit Survivor Pick"}
      </Button>
    </Box>
  );
};

export default SurvivorPoolPage;
