import { useEffect, useState } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame, type MessageHandler } from "../context/GameContext";
import type { RoundInfo } from "../types";

export default function GamePage() {
  const { sessionId, playerName } = useParams();
  const { isConnected, players, hand, roundInfo, connect, sendMessage } =
    useGame();

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

      <PlayerList players={players} />
      <PlayerHand hand={hand} />
      <RoundData roundInfo={roundInfo} />

      {roundInfo?.turnPlayer === playerName && (
        <BidBox msgFunction={sendMessage} />
      )}
    </div>
  );
}

function PlayerList({ players }: { players: string[] }) {
  return (
    <>
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
    </>
  );
}

function PlayerHand({ hand }: { hand: string[] }) {
  return (
    <div className="player-hand">
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
  );
}

function Card({ code }: { code: string }) {
  return <img className="card-image" src={`/cards/${code}.svg`} alt={code} />;
}

function RoundData({ roundInfo }: { roundInfo: RoundInfo | null }) {
  return (
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
  );
}

type BidBoxProps = {
  msgFunction: MessageHandler;
};

function BidBox({ msgFunction }: BidBoxProps) {
  const [bidVal, setBidVal] = useState<number | null>(null);

  function handleBid() {
    if (bidVal === null) return;
    msgFunction("make_bid", { bid: bidVal });
  }

  return (
    <div className="bid-box">
      <div className="bid-buttons">
        {Array.from({ length: 8 }, (_, i) => (
          <button
            key={i}
            onClick={() => setBidVal(i)}
            className={bidVal === i ? "selected" : ""}
          >
            {i}
          </button>
        ))}
      </div>
      <button disabled={bidVal === null} onClick={handleBid}>
        Submit Bid
      </button>
    </div>
  );
}
