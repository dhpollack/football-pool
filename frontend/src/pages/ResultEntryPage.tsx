import { useState, useEffect } from "react";
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

interface Game {
  id: number;
  favorite_team: string;
  underdog_team: string;
  spread: number;
}

interface ScoreEntry {
  game_id: number;
  favorite_score: number;
  underdog_score: number;
}

const ResultEntryPage = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [scores, setScores] = useState<{
    [key: number]: { favorite_score: number; underdog_score: number };
  }>({});
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGames = async () => {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("http://localhost:8080/api/games", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch games");
        }

        const data = await response.json();
        setGames(data);
        const initialScores: {
          [key: number]: { favorite_score: number; underdog_score: number };
        } = {};
        data.forEach((game: Game) => {
          initialScores[game.id] = { favorite_score: 0, underdog_score: 0 };
        });
        setScores(initialScores);
      } catch (error: any) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchGames();
  }, []);

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

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);

    try {
      const token = localStorage.getItem("token");
      const resultsToSubmit: ScoreEntry[] = games.map((game) => ({
        game_id: game.id,
        favorite_score: scores[game.id]?.favorite_score || 0,
        underdog_score: scores[game.id]?.underdog_score || 0,
      }));

      const response = await fetch("http://localhost:8080/api/results", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(resultsToSubmit),
      });

      if (!response.ok) {
        throw new Error("Failed to save results");
      }

      alert("Results saved successfully!");
    } catch (error: any) {
      setError(error.message);
    }
  };

  if (loading) {
    return <Typography>Loading games...</Typography>;
  }

  return (
    <div>
      <Typography variant="h4">Result Entry</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <Box component="form" onSubmit={handleSubmit}>
        <Button type="submit" variant="contained" sx={{ my: 2 }}>
          Save Results
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
