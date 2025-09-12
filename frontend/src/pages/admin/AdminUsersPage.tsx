import { useState } from "react";
import { Typography, Box, IconButton, Alert, Chip, TextField, MenuItem } from "@mui/material";
import { Edit, Delete } from "@mui/icons-material";
import { useAdminListUsers } from "../../services/api/user/user";
import AdminDataTable from "../../components/admin/AdminDataTable";
import AdminSearchFilter from "../../components/admin/AdminSearchFilter";
import AdminActionButtons from "../../components/admin/AdminActionButtons";
import AdminConfirmDialog from "../../components/admin/AdminConfirmDialog";
import { UserWithStats } from "../../services/model";

const AdminUsersPage = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState("");
  const [filters, setFilters] = useState<Record<string, any>>({});
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<UserWithStats | null>(null);
  const [editFormData, setEditFormData] = useState({
    name: "",
    email: "",
    role: "",
  });

  const { data, error, isLoading } = useAdminListUsers({
    request: {
      params: {
        search: searchTerm || undefined,
        role: filters.role || undefined,
        page: page + 1,
        limit: rowsPerPage,
      },
    },
  });

  const users = data?.users || [];
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

  const handleDeleteClick = (user: UserWithStats) => {
    setSelectedUser(user);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    // TODO: Implement delete mutation
    setDeleteDialogOpen(false);
    setSelectedUser(null);
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setSelectedUser(null);
  };

  const handleEditClick = (user: UserWithStats) => {
    setSelectedUser(user);
    setEditFormData({
      name: user.name || "",
      email: user.email || "",
      role: user.role || "",
    });
    setEditDialogOpen(true);
  };

  const handleEditCancel = () => {
    setEditDialogOpen(false);
    setSelectedUser(null);
    setEditFormData({
      name: "",
      email: "",
      role: "",
    });
  };

  const handleEditConfirm = () => {
    if (!selectedUser) return;
    
    // TODO: Implement edit mutation
    setEditDialogOpen(false);
    setSelectedUser(null);
    setEditFormData({
      name: "",
      email: "",
      role: "",
    });
  };

  const handleEditFormChange = (field: string, value: string) => {
    setEditFormData(prev => ({ ...prev, [field]: value }));
  };

  const columns = [
    {
      id: "id",
      label: "ID",
      align: "center" as const,
    },
    {
      id: "email",
      label: "Email",
    },
    {
      id: "name",
      label: "Name",
    },
    {
      id: "role",
      label: "Role",
      align: "center" as const,
      format: (user: UserWithStats) => (
        <Chip
          label={user.role}
          color={user.role === "admin" ? "primary" : "default"}
          size="small"
        />
      ),
    },
    {
      id: "pick_count",
      label: "Picks",
      align: "center" as const,
    },
    {
      id: "total_wins",
      label: "Wins",
      align: "center" as const,
    },
    {
      id: "created_at",
      label: "Joined",
      format: (user: UserWithStats) => new Date(user.created_at).toLocaleDateString(),
    },
    {
      id: "actions",
      label: "Actions",
      align: "center" as const,
      format: (user: UserWithStats) => (
        <Box display="flex" gap={1} justifyContent="center">
          <IconButton 
            size="small" 
            color="primary"
            onClick={() => handleEditClick(user)}
          >
            <Edit />
          </IconButton>
          <IconButton
            size="small"
            color="error"
            onClick={() => handleDeleteClick(user)}
          >
            <Delete />
          </IconButton>
        </Box>
      ),
    },
  ];

  const filterFields = [
    {
      name: "role",
      label: "Role",
      type: "select" as const,
      options: [
        { value: "", label: "All Roles" },
        { value: "user", label: "User" },
        { value: "admin", label: "Admin" },
      ],
    },
  ];

  return (
    <Box>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={3}
      >
        <Typography variant="h4">User Management</Typography>
        <AdminActionButtons
          onAdd={() => {}}
          addLabel="Add User"
        />
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Error loading users: {error.message}
        </Alert>
      )}

      <AdminSearchFilter
        searchTerm={searchTerm}
        onSearchChange={handleSearchChange}
        filters={filters}
        onFilterChange={handleFilterChange}
        filterFields={filterFields}
        onClearFilters={handleClearFilters}
        placeholder="Search by email or name..."
      />

      <AdminDataTable
        columns={columns}
        data={users}
        loading={isLoading}
        error={error?.message || null}
        page={page}
        rowsPerPage={rowsPerPage}
        totalCount={totalCount}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        emptyMessage={
          searchTerm || Object.keys(filters).length > 0
            ? "No users match your search criteria"
            : "No users available"
        }
      />

      <AdminConfirmDialog
        open={deleteDialogOpen}
        onClose={handleDeleteCancel}
        onConfirm={handleDeleteConfirm}
        title="Delete User"
        message={`Are you sure you want to delete user ${selectedUser?.email}? This action cannot be undone and will delete all associated picks.`}
        confirmLabel="Delete"
        severity="error"
      />

      {editDialogOpen && (
        <AdminConfirmDialog
          open={editDialogOpen}
          onClose={handleEditCancel}
          onConfirm={handleEditConfirm}
          title="Edit User"
          message={
            <Box component="div">
              <Typography variant="body2" component="div" sx={{ mb: 2 }}>
                Edit user details for {selectedUser?.email}
              </Typography>
              <Box display="flex" flexDirection="column" gap={2} component="div">
                <TextField
                  label="Name"
                  value={editFormData.name}
                  onChange={(e) => handleEditFormChange("name", e.target.value)}
                  fullWidth
                  size="small"
                />
                <TextField
                  label="Email"
                  value={editFormData.email}
                  onChange={(e) => handleEditFormChange("email", e.target.value)}
                  fullWidth
                  size="small"
                />
                <TextField
                  label="Role"
                  value={editFormData.role}
                  onChange={(e) => handleEditFormChange("role", e.target.value)}
                  fullWidth
                  size="small"
                  select
                >
                  <MenuItem value="user">User</MenuItem>
                  <MenuItem value="admin">Admin</MenuItem>
                </TextField>
              </Box>
            </Box>
          }
          confirmLabel="Save"
          severity="info"
        />
      )}
    </Box>
  );
};

export default AdminUsersPage;
