import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Alert,
} from "@mui/material";

const OverallResultsPage = () => {
  const [results, setResults] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const { data } = await axios.get("/api/results/season?season=2025");
        setResults(data);
      } catch (error) {
        setError("Failed to fetch overall results");
        console.error("Failed to fetch overall results", error);
      }
    };
    fetchResults();
  }, []);

  return (
    <div>
      <h1>Overall Results</h1>
      {error && <Alert severity="error">{error}</Alert>}
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Player</TableCell>
              <TableCell>Score</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {results.map((result) => (
              <TableRow key={result.player_id}>
                <TableCell>{result.player_name}</TableCell>
                <TableCell>{result.score}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default OverallResultsPage;
