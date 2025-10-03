import { useState, useId, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { TextField, Button, Box, Typography } from "@mui/material";
import { useMutation } from "@tanstack/react-query";
import useSignIn from "react-auth-kit/hooks/useSignIn";
import useIsAuthenticated from "react-auth-kit/hooks/useIsAuthenticated";
import { loginUser } from "../services/api/user/user";
import type { LoginRequest } from "../services/model";

const LoginPage = () => {
  const isAuthenticated = useIsAuthenticated();
  const signIn = useSignIn();
  const navigate = useNavigate();
  const emailId = useId();
  const passwordId = useId();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [formErrors, setFormErrors] = useState<{
    email?: string;
    password?: string;
  }>({});

  // Redirect if already authenticated
  useEffect(() => {
    console.log("LoginPage useEffect - isAuthenticated:", isAuthenticated);
    if (isAuthenticated) {
      console.log("User is authenticated, redirecting to home");
      navigate("/", { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const loginMutation = useMutation({
    mutationFn: (loginData: LoginRequest) => loginUser(loginData),
    onMutate: () => {
      setLoading(true);
    },
    onSuccess: (userData) => {
      console.log("Login API response received:", userData);

      // Check if user data has any non-serializable properties
      const userDataForAuth = {
        id: userData.user.id,
        name: userData.user.name,
        email: userData.user.email,
        role: userData.user.role,
        // Omit created_at and updated_at as they might cause issues
      };

      console.log("Prepared user data for auth:", userDataForAuth);

      // Store the authentication state using react-auth-kit
      const signInSuccess = signIn({
        auth: {
          token: userData.token, // Use the actual JWT token from API response
          type: "Bearer",
        },
        userState: userDataForAuth,
      });

      if (signInSuccess) {
        console.log("Sign in successful, waiting for auth state to propagate");
        // Give react-auth-kit time to propagate the authentication state
        // before navigating to avoid race conditions
        setTimeout(() => {
          console.log("Auth state should be propagated, navigating to home");
          navigate("/");
        }, 500); // Increased from 100ms to 500ms
      } else {
        console.error("Sign in failed - react-auth-kit signIn returned false");
        console.error("Token:", userData.token);
        console.error("User data:", userDataForAuth);

        // Try to debug localStorage directly
        const authData = localStorage.getItem("_auth");
        console.error("Current localStorage _auth:", authData);

        // Add a more specific error message for debugging
        setError(
          `Failed to store authentication data. Please check console for details.`,
        );

        // Force an error state that might be visible in tests
        setTimeout(() => {
          // This might help with debugging in test environment
          console.error(
            "AUTH_STORAGE_FAILED: react-auth-kit signIn returned false",
          );
        }, 100);
      }
    },
    onError: (error: Error) => {
      setError(error.message);
    },
    onSettled: () => {
      setLoading(false);
    },
  });

  const validateForm = () => {
    const errors: {
      email?: string;
      password?: string;
    } = {};

    if (!email.trim()) {
      errors.email = "Email is required";
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      errors.email = "Email is invalid";
    }

    if (!password) {
      errors.password = "Password is required";
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);

    if (!validateForm()) {
      return;
    }

    try {
      await loginMutation.mutateAsync({ email, password });
    } catch (error: unknown) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        height: "100vh",
      }}
    >
      <Typography variant="h4">Login</Typography>
      {error && (
        <Typography color="error" data-testid="error-message">
          {error}
        </Typography>
      )}
      <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id={emailId}
          label="Email Address"
          name="email"
          autoComplete="email"
          autoFocus
          value={email}
          onChange={(e) => {
            setEmail(e.target.value);
            if (formErrors.email) {
              setFormErrors({ ...formErrors, email: undefined });
            }
          }}
          error={!!formErrors.email}
          helperText={formErrors.email}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          name="password"
          label="Password"
          type="password"
          id={passwordId}
          autoComplete="current-password"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
            if (formErrors.password) {
              setFormErrors({ ...formErrors, password: undefined });
            }
          }}
          error={!!formErrors.password}
          helperText={formErrors.password}
        />
        <Button
          type="submit"
          fullWidth
          variant="contained"
          sx={{ mt: 3, mb: 2 }}
          disabled={loading}
        >
          {loading ? "Signing In..." : "Sign In"}
        </Button>
      </Box>
    </Box>
  );
};

export default LoginPage;
