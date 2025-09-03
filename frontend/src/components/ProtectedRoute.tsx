import { ReactNode } from "react";
import { useAuth } from "../hooks/useAuth";
import { CircularProgress, Box, Typography, Container } from "@mui/material";

interface ProtectedRouteProps {
  children: ReactNode;
  requireAdmin?: boolean;
}

const ProtectedRoute = ({ children, requireAdmin = false }: ProtectedRouteProps) => {
  const { isAuthenticated, isAdmin, loading } = useAuth();

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
    return (
      <Container maxWidth="sm" sx={{ mt: 4 }}>
        <Typography variant="h4" color="error" align="center" gutterBottom>
          Authentication Required
        </Typography>
        <Typography variant="body1" align="center">
          Please log in to access this page.
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