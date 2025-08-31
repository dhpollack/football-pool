import type React from "react";
import { useState, useEffect, useMemo } from "react";
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
  TextField,
  Box,
} from "@mui/material";
import { useGetGames, useSubmitResult } from "../services/api/default/default";
import type { GameResponse, ResultRequest } from "../services/model";
import axios from "axios";

interface ScoreEntry {
  favorite_score: number;
  underdog_score: number;
}

const ResultEntryPage = () => {
  const [scores, setScores] = useState<{
    [key: number]: ScoreEntry;
  }>({});
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // React Query hooks
  const {
    data: gamesData,
    isLoading,
    error: gamesError,
  } = useGetGames({
    week: 1,
    season: 2025,
  });
  const { mutateAsync: submitResult } = useSubmitResult();

  const games = useMemo(() => gamesData?.games || [], [gamesData]);

  // Initialize scores when games load
  useEffect(() => {
    if (games.length > 0) {
      const initialScores: { [key: number]: ScoreEntry } = {};
      games.forEach((game: GameResponse) => {
        initialScores[game.id] = { favorite_score: 0, underdog_score: 0 };
      });
      setScores(initialScores);
    }
  }, [games]);

  // Handle games fetch error
  useEffect(() => {
    if (gamesError) {
      if (axios.isAxiosError(gamesError)) {
        setError(
          gamesError.response?.data?.message ||
            gamesError.message ||
            "Failed to load games",
        );
      } else if (gamesError instanceof Error) {
        setError(gamesError.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  }, [gamesError]);

  const handleScoreChange = (
    gameId: number,
    type: "favorite_score" | "underdog_score",
    value: string,
  ) => {
    setScores((prevScores) => ({
      ...prevScores,
      [gameId]: {
        ...prevScores[gameId],
        [type]: parseInt(value, 10) || 0,
      },
    }));
  };

  const determineOutcome = (
    favoriteScore: number,
    underdogScore: number,
    spread: number,
  ): string => {
    const adjustedUnderdogScore = underdogScore + spread;
    if (favoriteScore > adjustedUnderdogScore) {
      return "favorite";
    } else if (favoriteScore < adjustedUnderdogScore) {
      return "underdog";
    } else {
      return "push";
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      // Submit results for each game individually
      const promises = games.map(async (game) => {
        const gameScores = scores[game.id];
        if (!gameScores) return;

        const favoriteScore = gameScores.favorite_score;
        const underdogScore = gameScores.underdog_score;
        const outcome = determineOutcome(
          favoriteScore,
          underdogScore,
          game.spread,
        );

        const resultData: ResultRequest = {
          game_id: game.id,
          favorite_score: favoriteScore,
          underdog_score: underdogScore,
          outcome: outcome,
        };

        return submitResult({ data: resultData });
      });

      await Promise.all(promises.filter(Boolean));
      alert("Results saved successfully!");
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        setError(error.response?.data?.message || error.message);
      } else if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return <Typography>Loading games...</Typography>;
  }

  return (
    <div>
      <Typography variant="h4">Result Entry</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <Box component="form" onSubmit={handleSubmit}>
        <Button
          type="submit"
          variant="contained"
          sx={{ my: 2 }}
          disabled={isSubmitting}
        >
          {isSubmitting ? "Saving Results..." : "Save Results"}
        </Button>
        <TableContainer component={Paper}>
          <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
              <TableRow>
                <TableCell>Favorite</TableCell>
                <TableCell>Underdog</TableCell>
                <TableCell>Spread</TableCell>
                <TableCell>Favorite Score</TableCell>
                <TableCell>Underdog Score</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {games.map((game) => (
                <TableRow
                  key={game.id}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <TableCell component="th" scope="row">
                    {game.favorite_team}
                  </TableCell>
                  <TableCell>{game.underdog_team}</TableCell>
                  <TableCell>{game.spread}</TableCell>
                  <TableCell>
                    <TextField
                      type="number"
                      name={`favorite_score_${game.id}`}
                      value={scores[game.id]?.favorite_score || ""}
                      onChange={(e) =>
                        handleScoreChange(
                          game.id,
                          "favorite_score",
                          e.target.value,
                        )
                      }
                    />
                  </TableCell>
                  <TableCell>
                    <TextField
                      type="number"
                      name={`underdog_score_${game.id}`}
                      value={scores[game.id]?.underdog_score || ""}
                      onChange={(e) =>
                        handleScoreChange(
                          game.id,
                          "underdog_score",
                          e.target.value,
                        )
                      }
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    </div>
  );
};

export default ResultEntryPage;
