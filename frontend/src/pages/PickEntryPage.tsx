import { useState, useEffect } from 'react';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Button, Typography } from '@mui/material';

interface Game {
  id: number;
  favorite_team: string;
  underdog_team: string;
  spread: number;
}

const PickEntryPage = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchGames = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/games');
        if (!response.ok) {
          throw new Error('Failed to fetch games');
        }
        const data = await response.json();
        setGames(data);
      } catch (error: any) {
        setError(error.message);
      }
    };

    fetchGames();
  }, []);

  return (
    <div>
      <Typography variant="h4">Pick Entry</Typography>
      <Button variant="contained" sx={{ my: 2 }}>Quick Pick</Button>
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
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {game.favorite_team}
                </TableCell>
                <TableCell>{game.underdog_team}</TableCell>
                <TableCell>{game.spread}</TableCell>
                <TableCell></TableCell>
                <TableCell></TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default PickEntryPage;