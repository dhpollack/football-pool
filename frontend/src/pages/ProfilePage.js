import React, { useState, useEffect } from "react";
import axios from "axios";
import { TextField, Button, Alert } from "@mui/material";

const ProfilePage = () => {
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const { data } = await axios.get("/api/users/me");
        setName(data.Name);
        setAddress(data.Address);
      } catch (error) {
        setError("Failed to fetch profile");
        console.error("Failed to fetch profile", error);
      }
    };
    fetchProfile();
  }, []);

  const handleUpdate = async () => {
    try {
      await axios.put("/api/users/me/update", { name, address });
      setSuccess("Profile updated successfully");
    } catch (error) {
      setError("Failed to update profile");
      console.error("Failed to update profile", error);
    }
  };

  return (
    <div>
      <h1>Profile</h1>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      <TextField
        label="Name"
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        fullWidth
        margin="normal"
      />
      <TextField
        label="Address"
        type="text"
        value={address}
        onChange={(e) => setAddress(e.target.value)}
        fullWidth
        margin="normal"
      />
      <Button variant="contained" onClick={handleUpdate}>Update</Button>
    </div>
  );
};

export default ProfilePage;
