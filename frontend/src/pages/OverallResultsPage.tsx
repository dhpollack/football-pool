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
import { useGetSeasonResults } from "../services/api/results/results";
import type { SeasonResult } from "../services/model";

const OverallResultsPage = () => {
  const {
    data: seasonResultsData,
    isLoading,
    error,
  } = useGetSeasonResults({
    season: 2025,
  });

  const seasonResults: SeasonResult[] = seasonResultsData || [];

  if (isLoading) {
    return <Typography>Loading overall results...</Typography>;
  }

  return (
    <div>
      <Typography variant="h4">Overall Results</Typography>
      {error && (
        <Typography color="error">
          {error instanceof Error
            ? error.message
            : "Failed to load overall results"}
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
            {seasonResults?.map((row, index) => (
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

export default OverallResultsPage;
