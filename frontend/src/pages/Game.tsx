import { useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame } from "../context/GameContext";

import SessionBox from "../components/SessionBox";
import RoundData from "../components/RoundData";
import ScoreTable from "../components/ScoreTable";
import BidBox from "../components/BidBox";
import GameTable from "../components/GameTable";
import PlayerHand from "../components/PlayerHand";

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
          <GameTable players={players} turnPlayer={roundInfo.turnPlayer} />
        )}
      </div>

      <div className="hand">
        <PlayerHand hand={hand} />
      </div>
    </div>
  );
}
