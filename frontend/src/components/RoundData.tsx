import type { GameState } from "../types";

export default function RoundData({
  gameState,
}: {
  gameState: GameState | null;
}) {
  return (
    <div className="round-info">
      {gameState ? (
        <div>
          <p>Round: {gameState.round}</p>
          <p>Turn: {gameState.turnPlayer}</p>
          <p>State: {gameState.state}</p>
        </div>
      ) : (
        <p>Waiting for round info</p>
      )}
    </div>
  );
}
