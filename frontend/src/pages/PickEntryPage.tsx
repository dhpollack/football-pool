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
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from "@mui/material";

interface Game {
  id: number;
  favorite_team: string;
  underdog_team: string;
  spread: number;
}

interface Pick {
  picked_team: string;
  rank: number;
}

const PickEntryPage = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [picks, setPicks] = useState<{ [gameId: number]: Pick }>({});
  const [error, setError] = useState(null);

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
        const initialPicks: { [gameId: number]: Pick } = {};
        data.forEach((game: Game) => {
          initialPicks[game.id] = { picked_team: "", rank: 0 };
        });
        setPicks(initialPicks);
      } catch (error: any) {
        setError(error.message);
      }
    };

    fetchGames();
  }, []);

  const handleQuickPick = () => {
    const newPicks: { [gameId: number]: Pick } = {};
    const availableRanks = Array.from(
      { length: games.length },
      (_, i) => i + 1,
    );
    const shuffledRanks = availableRanks.sort(() => Math.random() - 0.5);

    games.forEach((game, index) => {
      const pickOptions = [game.favorite_team, game.underdog_team];
      const pickedTeam =
        pickOptions[Math.floor(Math.random() * pickOptions.length)];
      newPicks[game.id] = {
        picked_team: pickedTeam,
        rank: shuffledRanks[index],
      };
    });
    setPicks(newPicks);
  };

  const handlePickChange = (gameId: number, pickedTeam: string) => {
    setPicks((prevPicks) => ({
      ...prevPicks,
      [gameId]: { ...prevPicks[gameId], picked_team: pickedTeam },
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
      const token = localStorage.getItem("token");
      const picksToSubmit = Object.keys(picks).map((gameId) => ({
        game_id: parseInt(gameId, 10),
        picked_team: picks[parseInt(gameId, 10)].picked_team,
        rank: picks[parseInt(gameId, 10)].rank,
      }));

      const response = await fetch("http://localhost:8080/api/picks", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(picksToSubmit),
      });

      if (!response.ok) {
        throw new Error("Failed to submit picks");
      }

      alert("Picks submitted successfully!");
    } catch (error: any) {
      setError(error.message);
    }
  };

  return (
    <div>
      <Typography variant="h4">Pick Entry</Typography>
      <Button variant="contained" sx={{ my: 2 }} onClick={handleQuickPick}>
        Quick Pick
      </Button>
      <Button
        variant="contained"
        sx={{ my: 2, ml: 2 }}
        onClick={handleSubmitPicks}
      >
        Submit Picks
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
                  <FormControl fullWidth>
                    <InputLabel>Pick</InputLabel>
                    <Select
                      value={picks[game.id]?.picked_team || ""}
                      label="Pick"
                      onChange={(e) =>
                        handlePickChange(game.id, e.target.value as string)
                      }
                    >
                      <MenuItem value={game.favorite_team}>
                        {game.favorite_team}
                      </MenuItem>
                      <MenuItem value={game.underdog_team}>
                        {game.underdog_team}
                      </MenuItem>
                    </Select>
                  </FormControl>
                </TableCell>
                <TableCell>
                  <FormControl fullWidth>
                    <InputLabel>Rank</InputLabel>
                    <Select
                      value={picks[game.id]?.rank || 0}
                      label="Rank"
                      onChange={(e) =>
                        handleRankChange(game.id, e.target.value as number)
                      }
                    >
                      {Array.from(
                        { length: games.length },
                        (_, i) => i + 1,
                      ).map((rank) => (
                        <MenuItem key={rank} value={rank}>
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
