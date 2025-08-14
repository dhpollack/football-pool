import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography } from '@mui/material';

const createData = (id: number, name: string, score: number) => {
  return { id, name, score };
}

const rows = [
  createData(1, 'Player A', 1000),
  createData(2, 'Player B', 900),
  createData(3, 'Player C', 800),
];

const OverallResultsPage = () => {
  return (
    <div>
      <Typography variant="h4">Overall Results</Typography>
      <TableContainer component={Paper} sx={{ my: 2 }}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>Player</TableCell>
              <TableCell>Score</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row) => (
              <TableRow
                key={row.id}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {row.name}
                </TableCell>
                <TableCell>{row.score}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default OverallResultsPage;