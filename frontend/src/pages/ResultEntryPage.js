import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TextField,
  Alert,
} from "@mui/material";

const ResultEntryPage = () => {
  const [games, setGames] = useState([]);
  const [results, setResults] = useState({});
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const fetchGames = async () => {
      try {
        const { data } = await axios.get("/api/games?season=2025&week=1");
        setGames(data);
      } catch (error) {
        setError("Failed to fetch games");
        console.error("Failed to fetch games", error);
      }
    };
    fetchGames();
  }, []);

  const handleScoreChange = (gameId, team, score) => {
    setResults({ ...results, [gameId]: { ...results[gameId], [team]: score } });
  };

  const handleSubmit = async (gameId) => {
    try {
      await axios.post("/api/results", { ...results[gameId], game_id: gameId });
      setSuccess("Result submitted successfully");
    } catch (error) {
      setError("Failed to submit result");
      console.error("Failed to submit result", error);
    }
  };

  return (
    <div>
      <h1>Result Entry</h1>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Game</TableCell>
              <TableCell>Favorite Score</TableCell>
              <TableCell>Underdog Score</TableCell>
              <TableCell>Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {games.map((game) => (
              <TableRow key={game.ID}>
                <TableCell>
                  {game.FavoriteTeam} vs {game.UnderdogTeam}
                </TableCell>
                <TableCell>
                  <TextField
                    type="number"
                    onChange={(e) =>
                      handleScoreChange(
                        game.ID,
                        "favorite_score",
                        parseInt(e.target.value)
                      )
                    }
                  />
                </TableCell>
                <TableCell>
                  <TextField
                    type="number"
                    onChange={(e) =>
                      handleScoreChange(
                        game.ID,
                        "underdog_score",
                        parseInt(e.target.value)
                      )
                    }
                  />
                </TableCell>
                <TableCell>
                  <Button variant="contained" onClick={() => handleSubmit(game.ID)}>
                    Submit
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default ResultEntryPage;
