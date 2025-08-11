import React, { useState } from "react";
import axios from "axios";
import { TextField, Button, Alert } from "@mui/material";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleLogin = async () => {
    try {
      await axios.post("/api/login", { email, password });
      window.location.href = "/";
    } catch (error) {
      setError("Invalid email or password");
      console.error("Login failed", error);
    }
  };

  return (
    <div>
      <h1>Login</h1>
      {error && <Alert severity="error">{error}</Alert>}
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
      <Button variant="contained" onClick={handleLogin}>Login</Button>
    </div>
  );
};

export default LoginPage;
