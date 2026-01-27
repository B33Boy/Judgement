import { useEffect, useState } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame, type MessageHandler } from "../context/GameContext";
import type { RoundInfo, Scores } from "../types";

export default function GamePage() {
  const { sessionId, playerName } = useParams();
  const {
    isConnected,
    players,
    hand,
    roundInfo,
    scores,
    connect,
    sendMessage,
  } = useGame();

  if (!sessionId || !playerName) {
    return <Navigate to="/" />;
  }

  useEffect(() => {
    if (!isConnected) {
      connect(sessionId, playerName);
    }
  }, [sessionId, playerName, isConnected, connect]);

  const isBiddingTurn =
    roundInfo?.state === "bidding" && roundInfo?.turnPlayer === playerName;

  const isPlaying = roundInfo?.state === "playing";

  return (
    <div className="game-layout">
      <h1 className="title">Judgement</h1>
      <div className="session">
        <SessionBox sessionId={sessionId} playerName={playerName} />
      </div>

      <div className="round">
        <RoundData roundInfo={roundInfo} />
      </div>

      <div className="score">
        <ScoreTable scores={scores} />
      </div>

      <div className="action">
        {isBiddingTurn && <BidBox msgFunction={sendMessage} />}
        {isPlaying && (
          <Table players={players} turnPlayer={roundInfo.turnPlayer} />
        )}
      </div>

      <div className="hand">
        <PlayerHand hand={hand} />
      </div>
    </div>
  );
}

// function PlayerList({ players }: { players: string[] }) {
//   return (
//     <>
//       {players.length === 0 ? (
//         <p>Connecting to game...</p>
//       ) : (
//         <>
//           <h3>Current Players</h3>
//           <ul className="player-list-game">
//             {players.map((p) => (
//               <li key={p}>{p}</li>
//             ))}
//           </ul>
//         </>
//       )}
//     </>
//   );
// }

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

function Table({
  players,
  turnPlayer,
}: {
  players: string[];
  turnPlayer: string;
}) {
  return (
    <div className="table">
      {players.map((p) => (
        <TableEntry key={p} player={p} isCurrentTurn={p === turnPlayer} />
      ))}
    </div>
  );
}

function TableEntry({
  player,
  isCurrentTurn,
}: {
  player: string;
  isCurrentTurn: boolean;
}) {
  return (
    <div
      className={`table-entry-container ${isCurrentTurn ? "current-turn" : ""}`}
    >
      <p className="player-name">{player}</p>
      <div className="table-entry"></div>
    </div>
  );
}

type SessionBoxProps = {
  sessionId: string;
  playerName: string;
};

function SessionBox({ sessionId, playerName }: SessionBoxProps) {
  return (
    <div className="session-box">
      <h3>Session: {sessionId}</h3>
      <h3>{playerName}</h3>
    </div>
  );
}

function ScoreTable({ scores }: { scores: Scores | null }) {
  if (!scores) return;

  const NUM_ROUNDS = 14;
  const entries = Array.from(scores.entries());

  return (
    <table className="score-table">
      <thead>
        <tr>
          <th>Player</th>
          {Array.from({ length: NUM_ROUNDS }, (_, i) => (
            <th key={i}>R{i + 1}</th>
          ))}
        </tr>
      </thead>

      <tbody>
        {entries.map(([player, playerScores]) => (
          <tr key={player}>
            <td>{player}</td>
            {Array.from({ length: NUM_ROUNDS }, (_, i) => (
              <td key={i}>{playerScores?.[i] ?? ""}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
