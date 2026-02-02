import { useEffect } from "react";
import { useParams, Navigate } from "react-router-dom";
import { useGame } from "../context/GameContext";

import SessionBox from "../components/SessionBox";
import RoundData from "../components/RoundData";
// import ScoreTable from "../components/ScoreTable";
import BidBox from "../components/BidBox";
import GameTable from "../components/GameTable";
import PlayerHand from "../components/PlayerHand";

import { DndContext, type DragEndEvent } from "@dnd-kit/core";

import { suitMap, rankMap } from "../types.ts";

export default function GamePage() {
  const { sessionId, playerName } = useParams();
  const {
    isConnected,
    playerId,
    players,
    hand,
    gameState,
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
    gameState?.state === "bidding" && gameState?.turnPlayer === playerId;

  const isPlaying = gameState?.state === "playing";

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;

    console.log(`PlayerId: ${playerId}`);
    console.log(`PlayerName: ${playerName}`);

    if (!over || !playerId) return;

    const targetPlayerId = String(over.id).replace("table-slot-", ""); // e.g. "table-slot-Bob"
    console.log(`TargetPlayerID: ${targetPlayerId}`);
    if (targetPlayerId != playerId) return;

    const card = String(active.id).replace("card-", ""); // e.g. "card-CLUB-2"
    const [suit, rank] = card.split("-");

    console.log(`${playerName} played ${card}`);
    sendMessage("play_card", { suit: suitMap[suit], rank: rankMap[rank] });
  }

  return (
    <DndContext onDragEnd={handleDragEnd}>
      <div className="game-layout">
        <h1 className="title">Judgement</h1>
        <div className="session">
          <SessionBox sessionId={sessionId} playerName={playerName} />
        </div>

        <div className="round">
          <RoundData gameState={gameState} />
        </div>

        {/* <div className="score">
          <ScoreTable scores={gameState} />
        </div> */}

        <div className="action">
          {isBiddingTurn && <BidBox msgFunction={sendMessage} />}
          {isPlaying && (
            <GameTable players={players} turnPlayer={gameState.turnPlayer} />
          )}
        </div>

        <div className="hand">
          <PlayerHand hand={hand} />
        </div>
      </div>
    </DndContext>
  );
}
