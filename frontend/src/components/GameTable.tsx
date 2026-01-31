import { useDroppable } from "@dnd-kit/core";

function Droppable(props: any) {
  const { isOver, setNodeRef } = useDroppable({ id: props.id });

  const style = {
    color: isOver ? "#fb752da4" : undefined,
    background: isOver ? "#fb752da4" : undefined,
  };

  return (
    <div ref={setNodeRef} style={style}>
      {props.children}
    </div>
  );
}

type GameTableProps = {
  players: string[];
  turnPlayer: string;
};

export default function GameTable({ players, turnPlayer }: GameTableProps) {
  return (
    <div className="table">
      {players.map((p) => (
        <TableEntry key={p} player={p} isCurrentTurn={p === turnPlayer} />
      ))}
    </div>
  );
}

type TableEntryProps = {
  player: string;
  isCurrentTurn: boolean;
};

function TableEntry({ player, isCurrentTurn }: TableEntryProps) {
  const content = (
    <div
      className={`table-entry-container ${isCurrentTurn ? "current-turn" : ""}`}
    >
      <p className="player-name">{player}</p>
      <div className="table-entry"></div>
    </div>
  );

  return isCurrentTurn ? (
    <Droppable id={`table-slot-${player}`}>{content}</Droppable>
  ) : (
    content
  );
}
