import React, { useState, useEffect } from "react";
import axios from "axios";

const PickEntryPage = () => {
  const [games, setGames] = useState([]);
  const [picks, setPicks] = useState({});

  useEffect(() => {
    const fetchGames = async () => {
      const { data } = await axios.get("/api/games?season=2025&week=1");
      setGames(data);
    };
    fetchGames();
  }, []);

  const handlePickChange = (gameId, team) => {
    setPicks({ ...picks, [gameId]: { ...picks[gameId], team } });
  };

  const handleRankChange = (gameId, rank) => {
    setPicks({ ...picks, [gameId]: { ...picks[gameId], rank } });
  };

  const handleQuickPick = () => {
    const newPicks = {};
    const ranks = Array.from({ length: games.length }, (_, i) => i + 1);
    games.forEach((game) => {
      const randomTeam =
        Math.random() > 0.5 ? game.FavoriteTeam : game.UnderdogTeam;
      const randomRankIndex = Math.floor(Math.random() * ranks.length);
      const randomRank = ranks.splice(randomRankIndex, 1)[0];
      newPicks[game.ID] = { team: randomTeam, rank: randomRank };
    });
    setPicks(newPicks);
  };

  const handleSubmit = async () => {
    try {
      const picksToSubmit = Object.keys(picks).map((gameId) => ({
        game_id: parseInt(gameId),
        picked_team: picks[gameId].team,
        rank: parseInt(picks[gameId].rank),
        quick_pick: false, // TODO: determine if quick pick was used
      }));
      await axios.post("/api/picks", picksToSubmit);
    } catch (error) {
      console.error("Failed to submit picks", error);
    }
  };

  return (
    <div>
      <h1>Pick Entry</h1>
      <button onClick={handleQuickPick}>Quick Pick</button>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
      >
        {games.map((game) => (
          <div key={game.ID}>
            <h3>
              {game.FavoriteTeam} vs {game.UnderdogTeam} ({game.Spread})
            </h3>
            <label>
              <input
                type="radio"
                name={`pick-${game.ID}`}
                value={game.FavoriteTeam}
                onChange={() => handlePickChange(game.ID, game.FavoriteTeam)}
                checked={picks[game.ID]?.team === game.FavoriteTeam}
              />
              {game.FavoriteTeam}
            </label>
            <label>
              <input
                type="radio"
                name={`pick-${game.ID}`}
                value={game.UnderdogTeam}
                onChange={() => handlePickChange(game.ID, game.UnderdogTeam)}
                checked={picks[game.ID]?.team === game.UnderdogTeam}
              />
              {game.UnderdogTeam}
            </label>
            <input
              type="number"
              min="1"
              max={games.length}
              placeholder="Rank"
              value={picks[game.ID]?.rank || ""}
              onChange={(e) => handleRankChange(game.ID, e.target.value)}
            />
          </div>
        ))}
        <button type="submit">Submit</button>
      </form>
    </div>
  );
};

export default PickEntryPage;
