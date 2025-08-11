import React, { useState } from "react";
import axios from "axios";
import { TextField, Button, Alert } from "@mui/material";
import { Link } from "react-router-dom";

const RegisterPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleRegister = async () => {
    try {
      await axios.post("/api/register", { email, password });
      setSuccess("Registration successful! You can now login.");
      setError("");
    } catch (error) {
      setError("Registration failed. Please try again.");
      setSuccess("");
      console.error("Registration failed", error);
    }
  };

  return (
    <div>
      <h1>Register</h1>
      {error && <Alert severity="error">{error}</Alert>}
      {success && <Alert severity="success">{success}</Alert>}
      <TextField
        label="Email"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        fullWidth
        margin="normal"
      />
      <TextField
        label="Password"
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        fullWidth
        margin="normal"
      />
      <Button variant="contained" onClick={handleRegister}>Register</Button>
      <p>
        Already have an account? <Link to="/login">Login here</Link>
      </p>
    </div>
  );
};

export default RegisterPage;
