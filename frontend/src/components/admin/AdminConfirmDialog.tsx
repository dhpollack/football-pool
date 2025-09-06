import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  CircularProgress,
  Alert,
} from "@mui/material";

interface AdminConfirmDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  loading?: boolean;
  error?: string | null;
  confirmLabel?: string;
  cancelLabel?: string;
  severity?: "warning" | "error" | "info";
}

const AdminConfirmDialog = ({
  open,
  onClose,
  onConfirm,
  title,
  message,
  loading = false,
  error = null,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  severity = "warning",
}: AdminConfirmDialogProps) => {
  const getSeverityColor = () => {
    switch (severity) {
      case "error":
        return "error";
      case "info":
        return "info";
      default:
        return "warning";
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <Typography variant="body1">{message}</Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          {cancelLabel}
        </Button>
        <Button
          onClick={onConfirm}
          color={getSeverityColor()}
          variant="contained"
          disabled={loading}
          startIcon={loading && <CircularProgress size={16} />}
        >
          {confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AdminConfirmDialog;
