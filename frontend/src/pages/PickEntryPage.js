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
  Radio,
  RadioGroup,
  FormControlLabel,
  TextField,
  Alert,
} from "@mui/material";

const PickEntryPage = () => {
  const [games, setGames] = useState([]);
  const [picks, setPicks] = useState({});
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [isQuickPick, setIsQuickPick] = useState(false);

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

  const handlePickChange = (gameId, team) => {
    setPicks({ ...picks, [gameId]: { ...picks[gameId], team } });
  };

  const handleRankChange = (gameId, rank) => {
    setPicks({ ...picks, [gameId]: { ...picks[gameId], rank } });
  };

  const handleQuickPick = () => {
    const newPicks = {};
    const ranks = Array.from({ length: games.length }, (_, i) => i + 1);
    games.forEach((game) => {
      const randomTeam = Math.random() > 0.5 ? "favorite" : "underdog";
      const randomRankIndex = Math.floor(Math.random() * ranks.length);
      const randomRank = ranks.splice(randomRankIndex, 1)[0];
      newPicks[game.ID] = { team: randomTeam, rank: randomRank };
    });
    setPicks(newPicks);
    setIsQuickPick(true);
  };

  const handleSubmit = async () => {
    try {
      const ranks = Object.values(picks).map((pick) => parseInt(pick.rank));
      const uniqueRanks = new Set(ranks);
      if (ranks.length !== uniqueRanks.size) {
        setError("Ranks must be unique");
        return;
      }

      const picksToSubmit = Object.keys(picks).map((gameId) => ({
        game_id: parseInt(gameId),
        picked_team: picks[gameId].team,
        rank: parseInt(picks[gameId].rank),
        quick_pick: isQuickPick,
      }));
      await axios.post("/api/picks/submit", picksToSubmit);
      setSuccess("Picks submitted successfully");
    } catch (error) {
      setError("Failed to submit picks");
      console.error("Failed to submit picks", error);
    }
  };

  return (
    <div>
      <h1>Pick Entry</h1>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      <Button variant="contained" onClick={handleQuickPick} style={{ marginBottom: "20px" }}>
        Quick Pick
      </Button>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Game</TableCell>
              <TableCell>Pick</TableCell>
              <TableCell>Rank</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {games.map((game) => (
              <TableRow key={game.ID}>
                <TableCell>
                  {game.FavoriteTeam} vs {game.UnderdogTeam} ({game.Spread})
                </TableCell>
                <TableCell>
                  <RadioGroup
                    row
                    value={picks[game.ID]?.team || ""}
                    onChange={(e) => handlePickChange(game.ID, e.target.value)}
                  >
                    <FormControlLabel
                      value="favorite"
                      control={<Radio />}
                      label={game.FavoriteTeam}
                    />
                    <FormControlLabel
                      value="underdog"
                      control={<Radio />}
                      label={game.UnderdogTeam}
                    />
                  </RadioGroup>
                </TableCell>
                <TableCell>
                  <TextField
                    type="number"
                    inputProps={{ min: 1, max: games.length }}
                    value={picks[game.ID]?.rank || ""}
                    onChange={(e) => handleRankChange(game.ID, e.target.value)}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Button
        variant="contained"
        onClick={handleSubmit}
        style={{ marginTop: "20px" }}
      >
        Submit Picks
      </Button>
    </div>
  );
};

export default PickEntryPage;
