import { useState } from "react";
import { Typography, Box, IconButton, Alert } from "@mui/material";
import { Edit, Delete, PlayArrow } from "@mui/icons-material";
import {
  useListWeeks,
  getListWeeksQueryKey,
  useDeleteWeek,
  useActivateWeek,
} from "../../services/api/admin/admin";
import AdminDataTable from "../../components/admin/AdminDataTable";
import AdminSearchFilter from "../../components/admin/AdminSearchFilter";
import AdminActionButtons from "../../components/admin/AdminActionButtons";
import AdminConfirmDialog from "../../components/admin/AdminConfirmDialog";
import WeekForm from "../../components/admin/WeekForm";
import type { WeekResponse } from "../../services/model";
import { useQueryClient } from "@tanstack/react-query";

const AdminWeeksPage = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<
    Record<string, string | number | undefined>
  >({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [weekFormOpen, setWeekFormOpen] = useState(false);
  const [selectedWeek, setSelectedWeek] = useState<WeekResponse | null>(null);

  const queryClient = useQueryClient();
  const { data, error, isLoading } = useListWeeks();
  const deleteWeekMutation = useDeleteWeek();
  const activateWeekMutation = useActivateWeek();

  const weeks = data?.weeks || [];

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

  const handleDeleteClick = (week: WeekResponse) => {
    setSelectedWeek(week);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (selectedWeek) {
      try {
        await deleteWeekMutation.mutateAsync({ id: selectedWeek.id });
        queryClient.invalidateQueries({
          queryKey: getListWeeksQueryKey(),
        });
      } catch (error) {
        console.error("Error deleting week:", error);
      }
    }
    setDeleteDialogOpen(false);
    setSelectedWeek(null);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setSelectedWeek(null);
  };

  const handleActivateWeek = async (week: WeekResponse) => {
    try {
      await activateWeekMutation.mutateAsync({ id: week.id });
      queryClient.invalidateQueries({
        queryKey: getListWeeksQueryKey(),
      });
    } catch (error) {
      console.error("Error activating week:", error);
    }
  };

  const handleAddWeek = () => {
    setSelectedWeek(null);
    setWeekFormOpen(true);
  };

  const handleEditWeek = (week: WeekResponse) => {
    setSelectedWeek(week);
    setWeekFormOpen(true);
  };

  const handleWeekFormClose = () => {
    setWeekFormOpen(false);
    setSelectedWeek(null);
  };

  const handleWeekFormSuccess = () => {
    queryClient.invalidateQueries({
      queryKey: getListWeeksQueryKey(),
    });
    setWeekFormOpen(false);
    setSelectedWeek(null);
  };

  const columns = [
    {
      id: "id",
      label: "ID",
      align: "center" as const,
    },
    {
      id: "week_number",
      label: "Week",
      align: "center" as const,
    },
    {
      id: "season",
      label: "Season",
      align: "center" as const,
    },
    {
      id: "week_start_time",
      label: "Start Time",
      format: (week: WeekResponse) =>
        week.week_start_time
          ? new Date(week.week_start_time).toLocaleDateString()
          : "-",
    },
    {
      id: "week_end_time",
      label: "End Time",
      format: (week: WeekResponse) =>
        week.week_end_time
          ? new Date(week.week_end_time).toLocaleDateString()
          : "-",
    },
    {
      id: "is_active",
      label: "Status",
      align: "center" as const,
      format: (week: WeekResponse) => (
        <Box
          sx={{
            color: week.is_active ? "success.main" : "text.secondary",
            fontWeight: week.is_active ? "bold" : "normal",
          }}
        >
          {week.is_active ? "Active" : "Inactive"}
        </Box>
      ),
    },
    {
      id: "actions",
      label: "Actions",
      align: "center" as const,
      format: (week: WeekResponse) => (
        <Box display="flex" gap={1} justifyContent="center">
          {!week.is_active && (
            <IconButton
              size="small"
              color="success"
              onClick={() => handleActivateWeek(week)}
              disabled={activateWeekMutation.isPending}
            >
              <PlayArrow />
            </IconButton>
          )}
          <IconButton
            size="small"
            color="primary"
            onClick={() => handleEditWeek(week)}
          >
            <Edit />
          </IconButton>
          <IconButton
            size="small"
            color="error"
            onClick={() => handleDeleteClick(week)}
            disabled={deleteWeekMutation.isPending}
          >
            <Delete />
          </IconButton>
        </Box>
      ),
    },
  ];

  const filterFields = [
    {
      name: "week_number",
      label: "Week",
      type: "number" as const,
    },
    {
      name: "season",
      label: "Season",
      type: "number" as const,
    },
    {
      name: "is_active",
      label: "Status",
      type: "select" as const,
      options: [
        { value: "true", label: "Active" },
        { value: "false", label: "Inactive" },
      ],
    },
  ];

  // Apply simple client-side filtering for now
  const filteredWeeks = weeks.filter((week) => {
    const matchesSearch = searchTerm
      ? week.season.toString().includes(searchTerm) ||
        week.week_number.toString().includes(searchTerm)
      : true;

    const matchesWeek = filters.week_number
      ? week.week_number === parseInt(filters.week_number as string, 10)
      : true;
    const matchesSeason = filters.season
      ? week.season === parseInt(filters.season as string, 10)
      : true;
    const matchesStatus =
      filters.is_active !== undefined
        ? week.is_active === (filters.is_active === "true")
        : true;

    return matchesSearch && matchesWeek && matchesSeason && matchesStatus;
  });

  const paginatedWeeks = filteredWeeks.slice(
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
        <Typography variant="h4">Week Management</Typography>
        <AdminActionButtons onAdd={handleAddWeek} addLabel="Add Week" />
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading weeks: {error.message}
        </Alert>
      )}

      <AdminSearchFilter
        searchTerm={searchTerm}
        onSearchChange={handleSearchChange}
        filters={filters}
        onFilterChange={handleFilterChange}
        filterFields={filterFields}
        onClearFilters={handleClearFilters}
        placeholder="Search by week or season..."
      />

      <AdminDataTable
        columns={columns}
        data={paginatedWeeks}
        loading={isLoading}
        error={error?.message || null}
        page={page}
        rowsPerPage={rowsPerPage}
        totalCount={filteredWeeks.length}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        emptyMessage={
          searchTerm || Object.keys(filters).length > 0
            ? "No weeks match your search criteria"
            : "No weeks available"
        }
      />

      <AdminConfirmDialog
        open={deleteDialogOpen}
        onClose={handleDeleteCancel}
        onConfirm={handleDeleteConfirm}
        title="Delete Week"
        message={`Are you sure you want to delete Week ${selectedWeek?.week_number} of Season ${selectedWeek?.season}? This action cannot be undone.`}
        confirmLabel="Delete"
        severity="error"
      />

      <WeekForm
        open={weekFormOpen}
        onClose={handleWeekFormClose}
        onSuccess={handleWeekFormSuccess}
        week={selectedWeek}
      />
    </Box>
  );
};

export default AdminWeeksPage;
