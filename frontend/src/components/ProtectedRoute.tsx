import type { ReactNode } from "react";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { CircularProgress, Box, Typography, Container } from "@mui/material";

interface ProtectedRouteProps {
  children: ReactNode;
  requireAdmin?: boolean;
}

const ProtectedRoute = ({
  children,
  requireAdmin = false,
}: ProtectedRouteProps) => {
  const { isAuthenticated, isAdmin, loading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      // Redirect to login if not authenticated
      navigate("/login", { replace: true });
    }
  }, [isAuthenticated, loading, navigate]);

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="200px"
      >
        <CircularProgress />
      </Box>
    );
  }

  if (!isAuthenticated) {
    // This will be shown briefly before redirect
    return (
      <Container maxWidth="sm" sx={{ mt: 4 }}>
        <Typography variant="h4" color="error" align="center" gutterBottom>
          Authentication Required
        </Typography>
        <Typography variant="body1" align="center">
          Redirecting to login...
        </Typography>
      </Container>
    );
  }

  if (requireAdmin && !isAdmin) {
    return (
      <Container maxWidth="sm" sx={{ mt: 4 }}>
        <Typography variant="h4" color="error" align="center" gutterBottom>
          Access Denied
        </Typography>
        <Typography variant="body1" align="center">
          You do not have permission to access this area.
        </Typography>
      </Container>
    );
  }

  return <>{children}</>;
};

export default ProtectedRoute;
