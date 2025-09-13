import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
} from "@mui/material";
import { useGetWeeklyResults } from "../services/api/results/results";
import type { WeeklyResult } from "../services/model";

const WeeklyResultsPage = () => {
  const {
    data: weeklyResultsData,
    isLoading,
    error,
  } = useGetWeeklyResults({
    week: 1,
    season: 2025,
  });

  const weeklyResults: WeeklyResult[] = weeklyResultsData || [];

  if (isLoading) {
    return <Typography>Loading weekly results...</Typography>;
  }

  return (
    <div>
      <Typography variant="h4">Weekly Results</Typography>
      {error && (
        <Typography color="error">
          {error instanceof Error
            ? error.message
            : "Failed to load weekly results"}
        </Typography>
      )}
      <TableContainer component={Paper} sx={{ my: 2 }}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell>Player</TableCell>
              <TableCell>Score</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {weeklyResults?.map((row, index) => (
              <TableRow
                key={`${row.player_name}-${index}`}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <TableCell component="th" scope="row">
                  {row.player_name}
                </TableCell>
                <TableCell>{row.score}</TableCell>
              </TableRow>
            )) || []}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default WeeklyResultsPage;
