import { useState } from "react";
import {
  Typography,
  Box,
  IconButton,
  Alert,
  Chip,
} from "@mui/material";
import { Delete } from "@mui/icons-material";
import { useAdminListPicks } from "../../services/api/default/default";
import AdminDataTable from "../../components/admin/AdminDataTable";
import AdminSearchFilter from "../../components/admin/AdminSearchFilter";
import AdminConfirmDialog from "../../components/admin/AdminConfirmDialog";
import { PickResponse } from "../../services/model";

const AdminPicksPage = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<Record<string, any>>({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedPick, setSelectedPick] = useState<PickResponse | null>(null);

  const { data, error, isLoading } = useAdminListPicks({
    request: {
      params: {
        user_id: filters.user_id || undefined,
        game_id: filters.game_id || undefined,
        week: filters.week || undefined,
        season: filters.season || undefined,
        page: page + 1,
        limit: rowsPerPage,
      },
    },
  });

  const picks = data?.picks || [];
  const totalCount = data?.pagination?.total || 0;

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

  const handleDeleteClick = (pick: PickResponse) => {
    setSelectedPick(pick);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    // TODO: Implement delete mutation
    console.log("Delete pick:", selectedPick?.id);
    setDeleteDialogOpen(false);
    setSelectedPick(null);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setSelectedPick(null);
  };

  const columns = [
    {
      id: "id",
      label: "ID",
      align: "center" as const,
    },
    {
      id: "user",
      label: "User",
      format: (pick: PickResponse) => pick.user?.email || "Unknown",
    },
    {
      id: "game",
      label: "Game",
      format: (pick: PickResponse) =>
        pick.game ? `${pick.game.favorite_team} vs ${pick.game.underdog_team}` : "Unknown",
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
      id: "created_at",
      label: "Submitted",
      format: (dateString: string) =>
        new Date(dateString).toLocaleDateString(),
    },
    {
      id: "actions",
      label: "Actions",
      align: "center" as const,
      format: (pick: PickResponse) => (
        <Box display="flex" gap={1} justifyContent="center">
          <IconButton
            size="small"
            color="error"
            onClick={() => handleDeleteClick(pick)}
          >
            <Delete />
          </IconButton>
        </Box>
      ),
    },
  ];

  const filterFields = [
    {
      name: "user_id",
      label: "User ID",
      type: "number" as const,
    },
    {
      name: "game_id",
      label: "Game ID",
      type: "number" as const,
    },
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

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Pick Management</Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading picks: {error.message}
        </Alert>
      )}

      <AdminSearchFilter
        searchTerm={searchTerm}
        onSearchChange={handleSearchChange}
        filters={filters}
        onFilterChange={handleFilterChange}
        filterFields={filterFields}
        onClearFilters={handleClearFilters}
        placeholder="Search by user email..."
      />

      <AdminDataTable
        columns={columns}
        data={picks}
        loading={isLoading}
        error={error?.message || null}
        page={page}
        rowsPerPage={rowsPerPage}
        totalCount={totalCount}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        emptyMessage={
          Object.keys(filters).length > 0
            ? "No picks match your filter criteria"
            : "No picks available"
        }
      />

      <AdminConfirmDialog
        open={deleteDialogOpen}
        onClose={handleDeleteCancel}
        onConfirm={handleDeleteConfirm}
        title="Delete Pick"
        message={`Are you sure you want to delete this pick? This action cannot be undone.`}
        confirmLabel="Delete"
        severity="error"
      />
    </Box>
  );
};

export default AdminPicksPage;