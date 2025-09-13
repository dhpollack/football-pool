import { Button, Box } from "@mui/material";
import { Add, Edit, Delete, Visibility } from "@mui/icons-material";

interface AdminActionButtonsProps {
  onAdd?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  onView?: () => void;
  addLabel?: string;
  editLabel?: string;
  deleteLabel?: string;
  viewLabel?: string;
  disabled?: {
    add?: boolean;
    edit?: boolean;
    delete?: boolean;
    view?: boolean;
  };
  size?: "small" | "medium" | "large";
}

const AdminActionButtons = ({
  onAdd,
  onEdit,
  onDelete,
  onView,
  addLabel = "Add New",
  editLabel = "Edit",
  deleteLabel = "Delete",
  viewLabel = "View",
  disabled = {},
  size = "medium",
}: AdminActionButtonsProps) => {
  return (
    <Box display="flex" gap={1} flexWrap="wrap">
      {onAdd && (
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={onAdd}
          disabled={disabled.add}
          size={size}
        >
          {addLabel}
        </Button>
      )}
      {onEdit && (
        <Button
          variant="outlined"
          startIcon={<Edit />}
          onClick={onEdit}
          disabled={disabled.edit}
          size={size}
        >
          {editLabel}
        </Button>
      )}
      {onView && (
        <Button
          variant="outlined"
          startIcon={<Visibility />}
          onClick={onView}
          disabled={disabled.view}
          size={size}
        >
          {viewLabel}
        </Button>
      )}
      {onDelete && (
        <Button
          variant="outlined"
          color="error"
          startIcon={<Delete />}
          onClick={onDelete}
          disabled={disabled.delete}
          size={size}
        >
          {deleteLabel}
        </Button>
      )}
    </Box>
  );
};

export default AdminActionButtons;
