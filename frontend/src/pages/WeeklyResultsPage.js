import React, { useState, useEffect } from "react";
import axios from "axios";

const WeeklyResultsPage = () => {
  const [results, setResults] = useState([]);

  useEffect(() => {
    const fetchResults = async () => {
      const { data } = await axios.get("/api/results/week?season=2025&week=1");
      setResults(data);
    };
    fetchResults();
  }, []);

  return (
    <div>
      <h1>Weekly Results</h1>
      <table>
        <thead>
          <tr>
            <th>Player</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          {results.map((result) => (
            <tr key={result.player_id}>
              <td>{result.player_name}</td>
              <td>{result.score}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default WeeklyResultsPage;
