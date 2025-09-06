import { useState } from "react";
import { Typography, Box, IconButton, Alert } from "@mui/material";
import { Edit, Delete } from "@mui/icons-material";
import { useAdminListGames } from "../../services/api/games/games";
import AdminDataTable from "../../components/admin/AdminDataTable";
import AdminSearchFilter from "../../components/admin/AdminSearchFilter";
import AdminActionButtons from "../../components/admin/AdminActionButtons";
import AdminConfirmDialog from "../../components/admin/AdminConfirmDialog";
import { GameResponse } from "../../services/model";

const AdminGamesPage = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<Record<string, any>>({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedGame, setSelectedGame] = useState<GameResponse | null>(null);

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

  const handleFilterChange = (name: string, value: any) => {
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
        `${game.favorite_team} vs ${game.underdog_team}`,
    },
    {
      id: "spread",
      label: "Spread",
      align: "center" as const,
      format: (spread: number) =>
        spread > 0 ? `+${spread}` : spread.toString(),
    },
    {
      id: "start_time",
      label: "Start Time",
      format: (dateString: string) => new Date(dateString).toLocaleDateString(),
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
      ? game.favorite_team.toLowerCase().includes(searchTerm.toLowerCase()) ||
        game.underdog_team.toLowerCase().includes(searchTerm.toLowerCase())
      : true;

    const matchesWeek = filters.week
      ? game.week === parseInt(filters.week)
      : true;
    const matchesSeason = filters.season
      ? game.season === parseInt(filters.season)
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
        <AdminActionButtons
          onAdd={() => console.log("Add new game")}
          addLabel="Add Game"
        />
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
        message={`Are you sure you want to delete the game between ${selectedGame?.favorite_team} and ${selectedGame?.underdog_team}? This action cannot be undone.`}
        confirmLabel="Delete"
        severity="error"
      />
    </Box>
  );
};

export default AdminGamesPage;
