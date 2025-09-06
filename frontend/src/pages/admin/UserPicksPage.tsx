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
import { Home, People, Search } from "@mui/icons-material";
import { useAdminGetPicksByUser } from "../../services/api/picks/picks";
import { useAdminGetUser } from "../../services/api/user/user";
import AdminDataTable from "../../components/admin/AdminDataTable";
import { PickResponse, UserResponse } from "../../services/model";

const UserPicksPage = () => {
  const { userId: userIdParam } = useParams<{ userId: string }>();
  const [userId, setUserId] = useState(userIdParam ? parseInt(userIdParam) : 0);
  const [searchTerm, setSearchTerm] = useState("");
  const [season, setSeason] = useState(2024); // Default season

  const { data: userData, error: userError } = useAdminGetUser(userId);

  const {
    data: picksData,
    error: picksError,
    isLoading: picksLoading,
    refetch,
  } = useAdminGetPicksByUser(userId, {
    request: {
      params: {
        season: season,
      },
    },
  });

  const user = userData || ({} as UserResponse);
  const picks = picksData?.picks || [];

  const handleUserIdChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newUserId = parseInt(event.target.value);
    if (!isNaN(newUserId) && newUserId > 0) {
      setUserId(newUserId);
    }
  };

  const handleSeasonChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newSeason = parseInt(event.target.value);
    if (!isNaN(newSeason) && newSeason > 0) {
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
      pick.game?.favorite_team?.toLowerCase().includes(searchLower) ||
      pick.game?.underdog_team?.toLowerCase().includes(searchLower) ||
      pick.game?.week?.toString().includes(searchTerm) ||
      pick.picked?.toLowerCase().includes(searchLower)
    );
  });

  const columns = [
    {
      id: "game",
      label: "Game",
      format: (pick: PickResponse) =>
        pick.game
          ? `${pick.game.favorite_team} vs ${pick.game.underdog_team}`
          : "Unknown",
    },
    {
      id: "week",
      label: "Week",
      align: "center" as const,
      format: (pick: PickResponse) => pick.game?.week || "-",
    },
    {
      id: "season",
      label: "Season",
      align: "center" as const,
      format: (pick: PickResponse) => pick.game?.season || "-",
    },
    {
      id: "picked",
      label: "Pick",
      align: "center" as const,
      format: (picked: string) => (
        <Chip
          label={picked}
          color={picked === "favorite" ? "primary" : "secondary"}
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
      format: (dateString: string) => new Date(dateString).toLocaleDateString(),
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
          <People sx={{ mr: 0.5 }} fontSize="inherit" />
          Picks
        </Link>
        <Typography color="text.primary">User {userId} Picks</Typography>
      </Breadcrumbs>

      <Typography variant="h4" gutterBottom>
        User Picks: {user.email || `User ${userId}`}
      </Typography>

      {userError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading user: {userError.message}
        </Alert>
      )}

      {picksError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading picks: {picksError.message}
        </Alert>
      )}

      {/* User Info */}
      {user.id && (
        <Card sx={{ p: 3, mb: 3 }}>
          <Grid container spacing={2}>
            <Grid size={{ xs: 12, sm: 6, md: 3 }}>
              <Typography variant="subtitle2" color="text.secondary">
                User ID
              </Typography>
              <Typography variant="body1">{user.id}</Typography>
            </Grid>
            <Grid size={{ xs: 12, sm: 6, md: 3 }}>
              <Typography variant="subtitle2" color="text.secondary">
                Email
              </Typography>
              <Typography variant="body1">{user.email}</Typography>
            </Grid>
            <Grid size={{ xs: 12, sm: 6, md: 3 }}>
              <Typography variant="subtitle2" color="text.secondary">
                Name
              </Typography>
              <Typography variant="body1">{user.name || "-"}</Typography>
            </Grid>
            <Grid size={{ xs: 12, sm: 6, md: 3 }}>
              <Typography variant="subtitle2" color="text.secondary">
                Role
              </Typography>
              <Chip
                label={user.role || "user"}
                color={user.role === "admin" ? "primary" : "default"}
                size="small"
              />
            </Grid>
          </Grid>
        </Card>
      )}

      {/* Filters */}
      <Card sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid size={{ xs: 12, sm: 6, md: 2 }}>
            <TextField
              fullWidth
              label="User ID"
              type="number"
              value={userId}
              onChange={handleUserIdChange}
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
              placeholder="Search by team, week, pick..."
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
              disabled={picksLoading || userId === 0}
              sx={{ height: "56px" }}
            >
              {picksLoading ? "Loading..." : "Load Picks"}
            </Button>
          </Grid>
        </Grid>
      </Card>

      {/* Stats Summary */}
      {picks.length > 0 && (
        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>
            Pick Summary
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
              >
                <Typography variant="h4" color="primary">
                  {picks.filter((pick) => pick.picked === "underdog").length}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Underdog Picks
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
              >
                <Typography variant="h4" color="primary">
                  {new Set(picks.map((pick) => pick.game?.week)).size}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Weeks Played
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Box>
      )}

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
            : userId
              ? `No picks found for user ${userId}`
              : "Please enter a user ID to view picks"
        }
      />
    </Box>
  );
};

export default UserPicksPage;
