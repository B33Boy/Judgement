export default function GameTable({
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
