import React, { useState, useEffect } from "react";
import axios from "axios";

const ResultEntryPage = () => {
  const [games, setGames] = useState([]);
  const [results, setResults] = useState({});

  useEffect(() => {
    const fetchGames = async () => {
      const { data } = await axios.get("/api/games?season=2025&week=1");
      setGames(data);
    };
    fetchGames();
  }, []);

  const handleScoreChange = (gameId, team, score) => {
    setResults({ ...results, [gameId]: { ...results[gameId], [team]: score } });
  };

  const handleSubmit = async (gameId) => {
    try {
      await axios.post("/api/results", { ...results[gameId], game_id: gameId });
    } catch (error) {
      console.error("Failed to submit result", error);
    }
  };

  return (
    <div>
      <h1>Result Entry</h1>
      {games.map((game) => (
        <div key={game.ID}>
          <h3>
            {game.FavoriteTeam} vs {game.UnderdogTeam}
          </h3>
          <input
            type="number"
            placeholder={`${game.FavoriteTeam} Score`}
            onChange={(e) =>
              handleScoreChange(
                game.ID,
                "favorite_score",
                parseInt(e.target.value),
              )
            }
          />
          <input
            type="number"
            placeholder={`${game.UnderdogTeam} Score`}
            onChange={(e) =>
              handleScoreChange(
                game.ID,
                "underdog_score",
                parseInt(e.target.value),
              )
            }
          />
          <button onClick={() => handleSubmit(game.ID)}>Submit</button>
        </div>
      ))}
    </div>
  );
};

export default ResultEntryPage;
