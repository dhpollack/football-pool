import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Button, Typography, TextField } from '@mui/material';

const createData = (id: number, favorite: string, underdog: string, spread: number) => {
  return { id, favorite, underdog, spread };
}

const rows = [
  createData(1, 'Team A', 'Team B', 7),
  createData(2, 'Team C', 'Team D', 3),
  createData(3, 'Team E', 'Team F', 10),
];

const ResultEntryPage = () => {
  return (
    <div>
      <Typography variant="h4">Result Entry</Typography>
      <Button variant="contained" sx={{ my: 2 }}>Save Results</Button>
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
            {rows.map((row) => (
              <TableRow
                key={row.id}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {row.favorite}
                </TableCell>
                <TableCell>{row.underdog}</TableCell>
                <TableCell>{row.spread}</TableCell>
                <TableCell>
                  <TextField type="number" />
                </TableCell>
                <TableCell>
                  <TextField type="number" />
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