import { Component, type ErrorInfo, type ReactNode } from "react";
import { Button, Box, Typography, Paper, Alert } from "@mui/material";
import axios from "axios";

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onReset?: () => void;
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: ErrorInfo;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Uncaught error:", error, errorInfo);
    this.setState({ errorInfo });

    // Log additional context for React Query errors
    if (error.message.includes("QueryClient")) {
      console.error("React Query error detected:", { error, errorInfo });
    }

    // Log network errors with more context
    if (axios.isAxiosError(error)) {
      console.error("Network error:", {
        status: error.response?.status,
        url: error.config?.url,
        method: error.config?.method,
      });
    }
  }

  private handleReset = () => {
    this.setState({ hasError: false, error: undefined, errorInfo: undefined });
    this.props.onReset?.();
  };

  private getErrorMessage(): string {
    const { error } = this.state;
    if (!error) return "An unknown error occurred";

    if (axios.isAxiosError(error)) {
      const status = error.response?.status;
      switch (status) {
        case 401:
          return "Authentication failed. Please log in again.";
        case 403:
          return "You don't have permission to access this resource.";
        case 404:
          return "The requested resource was not found.";
        case 500:
          return "Server error. Please try again later.";
        default:
          return `Network error: ${error.message}`;
      }
    }

    if (error.message.includes("QueryClient")) {
      return "Data loading error. Please try refreshing the page.";
    }

    return error.message;
  }

  public render() {
    if (this.state.hasError) {
      return (
        this.props.fallback || (
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              minHeight: "50vh",
              padding: 3,
            }}
          >
            <Paper
              elevation={3}
              sx={{
                padding: 4,
                maxWidth: 600,
                width: "100%",
                textAlign: "center",
              }}
            >
              <Alert severity="error" sx={{ mb: 3 }}>
                <Typography variant="h6" gutterBottom>
                  Something went wrong
                </Typography>
                <Typography variant="body1">
                  {this.getErrorMessage()}
                </Typography>
              </Alert>

              <Box sx={{ mb: 2 }}>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={this.handleReset}
                  sx={{ mr: 2 }}
                >
                  Try Again
                </Button>
                <Button
                  variant="outlined"
                  onClick={() => window.location.reload()}
                >
                  Refresh Page
                </Button>
              </Box>

              {import.meta.env.MODE === "development" && this.state.error && (
                <details style={{ marginTop: "16px" }}>
                  <summary style={{ cursor: "pointer", marginBottom: "8px" }}>
                    <Typography variant="body2" component="span">
                      Error details (development only)
                    </Typography>
                  </summary>
                  <Box
                    component="pre"
                    sx={{
                      textAlign: "left",
                      backgroundColor: "grey.100",
                      padding: 2,
                      borderRadius: 1,
                      overflow: "auto",
                      fontSize: "0.75rem",
                      maxHeight: 200,
                    }}
                  >
                    {this.state.error.toString()}
                    {this.state.errorInfo?.componentStack && (
                      <>
                        {"\n\nComponent Stack:"}
                        {this.state.errorInfo.componentStack}
                      </>
                    )}
                  </Box>
                </details>
              )}
            </Paper>
          </Box>
        )
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
