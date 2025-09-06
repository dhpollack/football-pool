import {
  Paper,
  TextField,
  Box,
  Button,
  IconButton,
  InputAdornment,
  Grid,
} from "@mui/material";
import { Search, Clear, FilterList } from "@mui/icons-material";

interface FilterField {
  name: string;
  label: string;
  type?: "text" | "number" | "select";
  options?: { value: string; label: string }[];
}

interface AdminSearchFilterProps {
  searchTerm: string;
  onSearchChange: (value: string) => void;
  filters: Record<string, any>;
  onFilterChange: (name: string, value: any) => void;
  filterFields?: FilterField[];
  onClearFilters?: () => void;
  placeholder?: string;
}

const AdminSearchFilter = ({
  searchTerm,
  onSearchChange,
  filters,
  onFilterChange,
  filterFields = [],
  onClearFilters,
  placeholder = "Search...",
}: AdminSearchFilterProps) => {
  const hasActiveFilters =
    Object.values(filters).some(
      (value) => value !== undefined && value !== "" && value !== null,
    ) || searchTerm !== "";

  return (
    <Paper sx={{ p: 2, mb: 3 }}>
      <Grid container spacing={2} alignItems="center">
        <Grid size={{ xs: 12, md: 6 }}>
          <TextField
            fullWidth
            variant="outlined"
            placeholder={placeholder}
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <Search />
                </InputAdornment>
              ),
              endAdornment: searchTerm && (
                <InputAdornment position="end">
                  <IconButton
                    size="small"
                    onClick={() => onSearchChange("")}
                    edge="end"
                  >
                    <Clear />
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
        </Grid>

        {filterFields.length > 0 && (
          <Grid size={{ xs: 12, md: 4 }}>
            <Box display="flex" gap={1} flexWrap="wrap">
              {filterFields.map((field) => (
                <TextField
                  key={field.name}
                  select={field.type === "select"}
                  type={field.type === "number" ? "number" : "text"}
                  label={field.label}
                  value={filters[field.name] || ""}
                  onChange={(e) => onFilterChange(field.name, e.target.value)}
                  size="small"
                  sx={{ minWidth: 120 }}
                  SelectProps={{
                    native: true,
                  }}
                >
                  {field.type === "select" && field.options && (
                    <option value="">All</option>
                  )}
                  {field.type === "select" &&
                    field.options &&
                    field.options.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                </TextField>
              ))}
            </Box>
          </Grid>
        )}

        {(hasActiveFilters || filterFields.length > 0) && (
          <Grid size={{ xs: 12, md: 2 }}>
            <Box display="flex" gap={1} justifyContent={{ md: "flex-end" }}>
              {filterFields.length > 0 && (
                <Button
                  variant="outlined"
                  startIcon={<FilterList />}
                  size="small"
                  disabled={!hasActiveFilters}
                  onClick={onClearFilters}
                >
                  Clear
                </Button>
              )}
            </Box>
          </Grid>
        )}
      </Grid>
    </Paper>
  );
};

export default AdminSearchFilter;
