import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import {
  Typography,
  Box,
  Alert,
  Chip,
  Breadcrumbs,
  TextField,
  Button,
  Grid,
  Card,
} from "@mui/material";
import { Home, EmojiEvents, Search } from "@mui/icons-material";
import { useAdminGetPicksByWeek } from "../../services/api/picks/picks";
import AdminDataTable from "../../components/admin/AdminDataTable";
import type { PickResponse } from "../../services/model";

const WeeklyPicksPage = () => {
  const { week: weekParam } = useParams<{ week: string }>();
  const [week, setWeek] = useState(weekParam ? parseInt(weekParam, 10) : 1);
  const [searchTerm, setSearchTerm] = useState("");
  const [season, setSeason] = useState(2024); // Default season

  const {
    data: picksData,
    error: picksError,
    isLoading: picksLoading,
    refetch,
  } = useAdminGetPicksByWeek(week, {
    request: {
      params: {
        season: season,
      },
    },
  });

  const picks = picksData?.picks || [];

  const handleWeekChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newWeek = parseInt(event.target.value, 10);
    if (!Number.isNaN(newWeek) && newWeek > 0) {
      setWeek(newWeek);
    }
  };

  const handleSeasonChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newSeason = parseInt(event.target.value, 10);
    if (!Number.isNaN(newSeason) && newSeason > 0) {
      setSeason(newSeason);
    }
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value);
  };

  const handleLoadPicks = () => {
    refetch();
  };

  // Filter picks by search term
  const filteredPicks = picks.filter((pick) => {
    if (!searchTerm) return true;

    const searchLower = searchTerm.toLowerCase();
    return (
      pick.user?.email?.toLowerCase().includes(searchLower) ||
      pick.user?.name?.toLowerCase().includes(searchLower) ||
      pick.game?.home_team?.toLowerCase().includes(searchLower) ||
      pick.game?.away_team?.toLowerCase().includes(searchLower)
    );
  });

  const columns = [
    {
      id: "user",
      label: "User",
      format: (pick: PickResponse) => pick.user?.email || "Unknown",
    },
    {
      id: "game",
      label: "Game",
      format: (pick: PickResponse) =>
        pick.game
          ? `${pick.game.home_team} vs ${pick.game.away_team}`
          : "Unknown",
    },
    {
      id: "picked",
      label: "Pick",
      align: "center" as const,
      format: (pick: PickResponse) => (
        <Chip
          label={pick.picked}
          color={pick.picked === "favorite" ? "primary" : "secondary"}
          size="small"
        />
      ),
    },
    {
      id: "spread",
      label: "Spread",
      align: "center" as const,
      format: (pick: PickResponse) => pick.game?.spread || "-",
    },
    {
      id: "start_time",
      label: "Game Time",
      format: (pick: PickResponse) =>
        pick.game?.start_time
          ? new Date(pick.game.start_time).toLocaleString()
          : "-",
    },
    {
      id: "created_at",
      label: "Submitted",
      format: (pick: PickResponse) =>
        new Date(pick.created_at).toLocaleDateString(),
    },
  ];

  return (
    <Box>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 3 }}>
        <Link
          to="/admin"
          style={{
            display: "flex",
            alignItems: "center",
            textDecoration: "none",
            color: "inherit",
          }}
        >
          <Home sx={{ mr: 0.5 }} fontSize="inherit" />
          Admin
        </Link>
        <Link
          to="/admin/picks"
          style={{
            display: "flex",
            alignItems: "center",
            textDecoration: "none",
            color: "inherit",
          }}
        >
          <EmojiEvents sx={{ mr: 0.5 }} fontSize="inherit" />
          Picks
        </Link>
        <Typography color="text.primary">Week {week} Picks</Typography>
      </Breadcrumbs>

      <Typography variant="h4" gutterBottom>
        Week {week} Picks
      </Typography>

      {picksError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading picks: {picksError.message}
        </Alert>
      )}

      {/* Filters */}
      <Card sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid size={{ xs: 12, sm: 6, md: 2 }}>
            <TextField
              fullWidth
              label="Week"
              type="number"
              value={week}
              onChange={handleWeekChange}
              inputProps={{ min: 1 }}
            />
          </Grid>
          <Grid size={{ xs: 12, sm: 6, md: 2 }}>
            <TextField
              fullWidth
              label="Season"
              type="number"
              value={season}
              onChange={handleSeasonChange}
              inputProps={{ min: 2020, max: 2030 }}
            />
          </Grid>
          <Grid size={{ xs: 12, sm: 8, md: 4 }}>
            <TextField
              fullWidth
              label="Search"
              placeholder="Search by user, team..."
              value={searchTerm}
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: (
                  <Search sx={{ mr: 1, color: "text.secondary" }} />
                ),
              }}
            />
          </Grid>
          <Grid size={{ xs: 12, sm: 4, md: 2 }}>
            <Button
              fullWidth
              variant="contained"
              onClick={handleLoadPicks}
              disabled={picksLoading}
              sx={{ height: "56px" }}
              data-testid="load-picks-button"
            >
              {picksLoading ? "Loading..." : "Load Picks"}
            </Button>
          </Grid>
        </Grid>
      </Card>

      {/* Stats Summary */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          Summary
        </Typography>
        <Grid container spacing={2}>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Box
              sx={{
                textAlign: "center",
                p: 2,
                bgcolor: "grey.100",
                borderRadius: 1,
              }}
              data-testid="total-picks-box"
            >
              <Typography variant="h4" color="primary">
                {picks.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Total Picks
              </Typography>
            </Box>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Box
              sx={{
                textAlign: "center",
                p: 2,
                bgcolor: "grey.100",
                borderRadius: 1,
              }}
              data-testid="unique-users-box"
            >
              <Typography variant="h4" color="primary">
                {new Set(picks.map((pick) => pick.user_id)).size}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Unique Users
              </Typography>
            </Box>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Box
              sx={{
                textAlign: "center",
                p: 2,
                bgcolor: "grey.100",
                borderRadius: 1,
              }}
              data-testid="favorite-picks-box"
            >
              <Typography variant="h4" color="primary">
                {picks.filter((pick) => pick.picked === "favorite").length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Favorite Picks
              </Typography>
            </Box>
          </Grid>
          <Grid size={{ xs: 6, sm: 3 }}>
            <Box
              sx={{
                textAlign: "center",
                p: 2,
                bgcolor: "grey.100",
                borderRadius: 1,
              }}
              data-testid="underdog-picks-box"
            >
              <Typography variant="h4" color="primary">
                {picks.filter((pick) => pick.picked === "underdog").length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Underdog Picks
              </Typography>
            </Box>
          </Grid>
        </Grid>
      </Box>

      {/* Picks Table */}
      <AdminDataTable
        columns={columns}
        data={filteredPicks}
        loading={picksLoading}
        error={picksError?.message || null}
        page={0}
        rowsPerPage={20}
        totalCount={filteredPicks.length}
        onPageChange={() => {}}
        onRowsPerPageChange={() => {}}
        emptyMessage={
          picksLoading
            ? "Loading picks..."
            : week
              ? `No picks found for week ${week}`
              : "Please select a week to view picks"
        }
      />
    </Box>
  );
};

export default WeeklyPicksPage;
