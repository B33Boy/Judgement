import { useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame } from "../context/GameContext";

export default function GamePage() {
  const { sessionId, playerName } = useParams();
  const { isConnected, players, hand, roundInfo, connect } = useGame();

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
        {<h3>{playerName}'s Hand</h3>}
        <div
          className="playerHand"
          style={{ display: "flex", gap: "15%", flexWrap: "wrap" }}
        >
          {hand.length > 0 ? (
            hand.map((card, index) => (
              <div key={index} className="card-item">
                <Card code={card} />
              </div>
            ))
          ) : (
            <p>Waiting for cards...</p>
          )}
        </div>
        <div className="roundInfo">
          {roundInfo ? (
            <div>
              <p>Round: {roundInfo.round}</p>
              <p>Turn: {roundInfo.turnPlayer}</p>
              <p>State: {roundInfo.state}</p>
            </div>
          ) : (
            <p>Waiting for round info</p>
          )}
        </div>
      </section>
    </div>
  );
}

function Card({ code }: { code: string }) {
  return <img className="card-image" src={`/cards/${code}.svg`} alt={code} />;
}
