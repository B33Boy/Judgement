import { useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame } from "../context/GameContext";

export default function GamePage() {
  const { sessionId, playerName } = useParams();
  const { isConnected, players, hand, connect } = useGame();

  if (!sessionId || !playerName) {
    return <Navigate to="/" />;
  }

  useEffect(() => {
    if (!isConnected) {
      connect(sessionId, playerName);
    }
  }, [sessionId, playerName, isConnected, connect]);

  return (
    <div>
      <h1>Judgement</h1>
      <p>Session: {sessionId}</p>

      <section>
        <h3>Players</h3>

        {players.length === 0 ? (
          <p>Connecting to game...</p>
        ) : (
          <>
            <h3>Current Players</h3>
            <ul className="player-list-game">
              {players.map((p) => (
                <li key={p}>{p}</li>
              ))}
            </ul>
          </>
        )}
      </section>

      <hr />
      <section>
        <h3>Your Hand</h3>
        <div style={{ display: "flex", gap: "15%", flexWrap: "wrap" }}>
          {hand.length > 0 ? (
            hand.map((card, index) => (
              <div key={index} className="card-item">
                <span>{card.replace("_", " ")}</span>
              </div>
            ))
          ) : (
            <p>Waiting for cards...</p>
          )}
        </div>
      </section>
    </div>
  );
}
