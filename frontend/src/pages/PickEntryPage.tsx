import React, { useState, useId } from "react";
import { useGetGames } from "../services/api/games/games";
import { useSubmitPicks } from "../services/api/picks/picks";
import { useListWeeks } from "../services/api/admin/admin";
import type {
  GameResponse,
  PickRequest,
  WeekResponse,
} from "../services/model";
import axios from "axios";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Typography,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from "@mui/material";

interface Pick {
  picked: string;
  rank: number;
  quick_pick?: boolean;
}

const PickEntryPage = () => {
  const quickPickId = useId();
  const submitPicksId = useId();
  const [picks, setPicks] = useState<{ [gameId: number]: Pick }>({});
  const [error, setError] = useState<string | null>(null);

  // Use React Query hooks
  const {
    data: weeksData,
    isLoading: weeksLoading,
    error: weeksError,
  } = useListWeeks();

  // Find the active week
  const activeWeek = weeksData?.weeks?.find(
    (week: WeekResponse) => week.is_active,
  );

  const {
    data: gamesData,
    isLoading: gamesLoading,
    error: gamesError,
  } = useGetGames({
    week: activeWeek?.week_number || 1,
    season: activeWeek?.season || new Date().getFullYear(),
  });
  const { mutateAsync: submitPicks, isPending: isSubmitting } =
    useSubmitPicks();

  const games = gamesData?.games || [];
  const sortedGames = React.useMemo(
    () =>
      [...games].sort(
        (a, b) =>
          new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
      ),
    [games],
  );

  // Initialize picks when games load
  React.useEffect(() => {
    if (games.length > 0) {
      const initialPicks: { [gameId: number]: Pick } = {};
      sortedGames.forEach((game: GameResponse) => {
        initialPicks[game.id] = { picked: "", rank: 0 };
      });
      setPicks(initialPicks);
    }
  }, [sortedGames, games.length]);

  // Handle weeks and games fetch errors
  React.useEffect(() => {
    if (weeksError || gamesError) {
      const error = weeksError || gamesError;
      if (axios.isAxiosError(error)) {
        setError(
          error.response?.data?.message ||
            error.message ||
            "Failed to load data",
        );
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  }, [weeksError, gamesError]);

  const handleQuickPick = () => {
    const newPicks: { [gameId: number]: Pick } = {};
    const availableRanks = Array.from(
      { length: games.length },
      (_, i) => i + 1,
    );
    const shuffledRanks = availableRanks.sort(() => Math.random() - 0.5);

    sortedGames.forEach((game, index) => {
      const pickOptions = ["favorite", "underdog"];
      const picked =
        pickOptions[Math.floor(Math.random() * pickOptions.length)];
      newPicks[game.id] = {
        picked: picked,
        rank: shuffledRanks[index],
        quick_pick: true,
      };
    });
    setPicks(newPicks);
  };

  const handlePickChange = (gameId: number, picked: string) => {
    setPicks((prevPicks) => ({
      ...prevPicks,
      [gameId]: { ...prevPicks[gameId], picked: picked },
    }));
  };

  const handleRankChange = (gameId: number, rank: number) => {
    setPicks((prevPicks) => ({
      ...prevPicks,
      [gameId]: { ...prevPicks[gameId], rank: rank },
    }));
  };

  const handleSubmitPicks = async () => {
    setError(null);
    try {
      const picksToSubmit: PickRequest[] = Object.keys(picks).map((gameId) => ({
        game_id: parseInt(gameId, 10),
        picked: picks[parseInt(gameId, 10)].picked,
        rank: picks[parseInt(gameId, 10)].rank,
        quick_pick: picks[parseInt(gameId, 10)].quick_pick || false,
      }));

      await submitPicks({ data: picksToSubmit });
      alert("Picks submitted successfully!");
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

  const isLoading = weeksLoading || gamesLoading;

  if (isLoading) {
    return (
      <div>
        <Typography variant="h4">Pick Entry</Typography>
        <Typography>Loading data...</Typography>
      </div>
    );
  }

  if (!activeWeek) {
    return (
      <div>
        <Typography variant="h4">Pick Entry</Typography>
        <Typography color="error">
          No active week found. Please contact an administrator to set an active
          week.
        </Typography>
      </div>
    );
  }

  return (
    <div>
      <Typography variant="h4">Pick Entry</Typography>
      <Typography variant="h6" color="text.secondary" sx={{ mb: 2 }}>
        Week {activeWeek.week_number} - Season {activeWeek.season}
      </Typography>
      <Button
        variant="contained"
        sx={{ my: 2 }}
        onClick={handleQuickPick}
        id={quickPickId}
      >
        Quick Pick
      </Button>
      <Button
        variant="contained"
        sx={{ my: 2, ml: 2 }}
        onClick={handleSubmitPicks}
        id={submitPicksId}
        disabled={isSubmitting}
      >
        {isSubmitting ? "Submitting..." : "Submit Picks"}
      </Button>
      {error && <Typography color="error">{error}</Typography>}
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>Favorite</TableCell>
              <TableCell>Underdog</TableCell>
              <TableCell>Spread</TableCell>
              <TableCell>Pick</TableCell>
              <TableCell>Rank</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sortedGames.map((game) => (
              <TableRow
                key={game.id}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {game.favorite === "Home" ? game.home_team : game.away_team}
                </TableCell>
                <TableCell>
                  {game.underdog === "Home" ? game.home_team : game.away_team}
                </TableCell>
                <TableCell>{game.spread}</TableCell>
                <TableCell>
                  <FormControl fullWidth>
                    <InputLabel>Pick</InputLabel>
                    <Select
                      value={picks[game.id]?.picked || ""}
                      label="Pick"
                      onChange={(e) =>
                        handlePickChange(game.id, e.target.value as string)
                      }
                      data-testid={`pick-select-${game.id}`}
                    >
                      <MenuItem value="favorite">
                        {game.favorite === "Home"
                          ? game.home_team
                          : game.away_team}
                      </MenuItem>
                      <MenuItem value="underdog">
                        {game.underdog === "Home"
                          ? game.home_team
                          : game.away_team}
                      </MenuItem>
                    </Select>
                  </FormControl>
                </TableCell>
                <TableCell>
                  <FormControl fullWidth>
                    <InputLabel>Rank</InputLabel>
                    <Select
                      value={picks[game.id]?.rank || ""}
                      label="Rank"
                      onChange={(e) =>
                        handleRankChange(game.id, e.target.value as number)
                      }
                      data-testid={`rank-select-${game.id}`}
                    >
                      {Array.from(
                        { length: sortedGames.length },
                        (_, i) => i + 1,
                      ).map((rank) => (
                        <MenuItem
                          key={rank}
                          value={rank}
                          data-testid={`rank-option-${rank}`}
                        >
                          {rank}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default PickEntryPage;
