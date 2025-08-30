import { useState, useEffect, useId } from "react";
import { TextField, Button, Box, Typography } from "@mui/material";
import { api, ApiError } from "../services/api";

const ProfilePage = () => {
  const nameId = useId();
  const addressId = useId();
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [formErrors, setFormErrors] = useState<{
    name?: string;
  }>({});

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await api.get("/api/users/me");
        setName(data.name);
        setAddress(data.address);
      } catch (error: unknown) {
        if (error instanceof ApiError && error.status === 401) {
          // Authentication error - redirect will be handled by AuthContext
          setError("Please log in to view your profile");
        } else if (error instanceof Error) {
          setError(error.message);
        } else {
          setError("An unknown error occurred");
        }
      }
    };

    fetchProfile();
  }, []);

  const validateForm = () => {
    const errors: {
      name?: string;
    } = {};

    if (!name.trim()) {
      errors.name = "Name is required";
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
      await api.put("/api/users/me", { name, address });
      alert("Profile updated successfully!");
    } catch (error: unknown) {
      if (error instanceof ApiError && error.status === 401) {
        // Authentication error - redirect will be handled by AuthContext
        setError("Please log in to update your profile");
      } else if (error instanceof Error) {
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
      <Typography variant="h4">Profile</Typography>
      {error && <Typography color="error">{error}</Typography>}
      <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id={nameId}
          label="Name"
          name="name"
          autoComplete="name"
          autoFocus
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            if (formErrors.name) {
              setFormErrors({ ...formErrors, name: undefined });
            }
          }}
          error={!!formErrors.name}
          helperText={formErrors.name}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id={addressId}
          label="Address"
          name="address"
          autoComplete="address"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
        />
        <Button
          type="submit"
          fullWidth
          variant="contained"
          sx={{ mt: 3, mb: 2 }}
        >
          Save
        </Button>
      </Box>
    </Box>
  );
};

export default ProfilePage;
