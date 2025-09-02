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
import { useDebugGetUsers } from "../services/api/default/default";
import type { UserWithStats } from "../services/model";

const OverallResultsPage = () => {
  const { data: usersData, isLoading, error } = useDebugGetUsers();

  // Transform the user data into the format expected by the UI
  const seasonResults =
    usersData?.users
      ?.filter((user: UserWithStats) => user.total_wins > 0)
      ?.map((user: UserWithStats) => ({
        player_name: user.name,
        score: user.total_wins,
      }))
      ?.sort((a, b) => b.score - a.score) || [];

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
