import { useState } from "react";
import { Typography, Box, IconButton, Alert } from "@mui/material";
import { Edit, Delete } from "@mui/icons-material";
import {
  useAdminListGames,
  getAdminListGamesQueryKey,
} from "../../services/api/games/games";
import AdminDataTable from "../../components/admin/AdminDataTable";
import AdminSearchFilter from "../../components/admin/AdminSearchFilter";
import AdminActionButtons from "../../components/admin/AdminActionButtons";
import AdminConfirmDialog from "../../components/admin/AdminConfirmDialog";
import GameForm from "../../components/admin/GameForm";
import type { GameResponse } from "../../services/model";
import { useQueryClient } from "@tanstack/react-query";

const AdminGamesPage = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<
    Record<string, string | number | undefined>
  >({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [gameFormOpen, setGameFormOpen] = useState(false);
  const [selectedGame, setSelectedGame] = useState<GameResponse | null>(null);

  const queryClient = useQueryClient();
  const { data, error, isLoading } = useAdminListGames();

  const games = data?.games || [];

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const handleRowsPerPageChange = (newRowsPerPage: number) => {
    setRowsPerPage(newRowsPerPage);
    setPage(0);
  };

  const handleSearchChange = (value: string) => {
    setSearchTerm(value);
    setPage(0);
  };

  const handleFilterChange = (name: string, value: string | number) => {
    setFilters((prev) => ({ ...prev, [name]: value }));
    setPage(0);
  };

  const handleClearFilters = () => {
    setFilters({});
    setSearchTerm("");
    setPage(0);
  };

  const handleDeleteClick = (game: GameResponse) => {
    setSelectedGame(game);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    // TODO: Implement delete mutation
    console.log("Delete game:", selectedGame?.id);
    setDeleteDialogOpen(false);
    setSelectedGame(null);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setSelectedGame(null);
  };

  const handleAddGame = () => {
    setGameFormOpen(true);
  };

  const handleGameFormClose = () => {
    setGameFormOpen(false);
  };

  const handleGameFormSuccess = () => {
    queryClient.invalidateQueries({
      queryKey: getAdminListGamesQueryKey(),
    });
    setGameFormOpen(false);
  };

  const columns = [
    {
      id: "id",
      label: "ID",
      align: "center" as const,
    },
    {
      id: "week",
      label: "Week",
      align: "center" as const,
    },
    {
      id: "season",
      label: "Season",
      align: "center" as const,
    },
    {
      id: "matchup",
      label: "Matchup",
      format: (game: GameResponse) =>
        `${game.home_team} vs ${game.away_team}`,
    },
    {
      id: "spread",
      label: "Spread",
      align: "center" as const,
      format: (game: GameResponse) =>
        game.spread > 0 ? `+${game.spread}` : game.spread.toString(),
    },
    {
      id: "start_time",
      label: "Start Time",
      format: (game: GameResponse) =>
        new Date(game.start_time).toLocaleDateString(),
    },
    {
      id: "actions",
      label: "Actions",
      align: "center" as const,
      format: (game: GameResponse) => (
        <Box display="flex" gap={1} justifyContent="center">
          <IconButton size="small" color="primary">
            <Edit />
          </IconButton>
          <IconButton
            size="small"
            color="error"
            onClick={() => handleDeleteClick(game)}
          >
            <Delete />
          </IconButton>
        </Box>
      ),
    },
  ];

  const filterFields = [
    {
      name: "week",
      label: "Week",
      type: "number" as const,
    },
    {
      name: "season",
      label: "Season",
      type: "number" as const,
    },
  ];

  // Apply simple client-side filtering for now
  const filteredGames = games.filter((game) => {
    const matchesSearch = searchTerm
      ? game.home_team.toLowerCase().includes(searchTerm.toLowerCase()) ||
        game.away_team.toLowerCase().includes(searchTerm.toLowerCase())
      : true;

    const matchesWeek = filters.week
      ? game.week === parseInt(filters.week as string, 10)
      : true;
    const matchesSeason = filters.season
      ? game.season === parseInt(filters.season as string, 10)
      : true;

    return matchesSearch && matchesWeek && matchesSeason;
  });

  const paginatedGames = filteredGames.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage,
  );

  return (
    <Box>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={3}
      >
        <Typography variant="h4">Game Management</Typography>
        <AdminActionButtons onAdd={handleAddGame} addLabel="Add Game" />
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading games: {error.message}
        </Alert>
      )}

      <AdminSearchFilter
        searchTerm={searchTerm}
        onSearchChange={handleSearchChange}
        filters={filters}
        onFilterChange={handleFilterChange}
        filterFields={filterFields}
        onClearFilters={handleClearFilters}
        placeholder="Search by team name..."
      />

      <AdminDataTable
        columns={columns}
        data={paginatedGames}
        loading={isLoading}
        error={error?.message || null}
        page={page}
        rowsPerPage={rowsPerPage}
        totalCount={filteredGames.length}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        emptyMessage={
          searchTerm || Object.keys(filters).length > 0
            ? "No games match your search criteria"
            : "No games available"
        }
      />

      <AdminConfirmDialog
        open={deleteDialogOpen}
        onClose={handleDeleteCancel}
        onConfirm={handleDeleteConfirm}
        title="Delete Game"
        message={`Are you sure you want to delete the game between ${selectedGame?.home_team} and ${selectedGame?.away_team}? This action cannot be undone.`}
        confirmLabel="Delete"
        severity="error"
      />

      <GameForm
        open={gameFormOpen}
        onClose={handleGameFormClose}
        onSuccess={handleGameFormSuccess}
      />
    </Box>
  );
};

export default AdminGamesPage;
