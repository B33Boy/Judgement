import { useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { getPlayerName } from "../lib/player";
import { useGame } from "../context/GameContext";

export default function SessionPage() {
  const { sessionId } = useParams();
  const { connect, players, sendMessage } = useGame();
  const playerName = getPlayerName();

  // If no ID or name, go home
  if (!sessionId || !playerName) return <Navigate to="/" />;

  useEffect(() => {
    connect(sessionId, playerName);
  }, [sessionId, playerName, connect]);

  const handleStartGame = () => {
    sendMessage("start_game");
  };

  return (
    <div>
      <h2>Session: {sessionId}</h2>
      <h3>Lobby</h3>

      <div className="player-list">
        <h4>Players ({players.length})</h4>
        <ul>
          {players.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>

      <button
        disabled={players.length < 3 || players.length > 7}
        onClick={handleStartGame}
      >
        Start Game
      </button>
    </div>
  );
}
