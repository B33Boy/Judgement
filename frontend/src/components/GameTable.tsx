import { useDroppable } from "@dnd-kit/core";
import {
  rankFromValue,
  suitFromValue,
  type Card,
  type GameState,
  type PlayerPublic,
  type Players,
} from "../types";
import { CardImg } from "./PlayerHand";

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
  players: Players;
  gameState: GameState;
};

export default function GameTable({ players, gameState }: GameTableProps) {
  return (
    <div className="table">
      {players.map((p: PlayerPublic) => (
        <TableEntry
          key={p.id}
          player={p}
          isCurrentTurn={p.id === gameState.turnPlayer}
          card={gameState.table[p.id]}
        />
      ))}
    </div>
  );
}

type TableEntryProps = {
  player: PlayerPublic;
  isCurrentTurn: boolean;
  card?: Card;
};

function TableEntry({ player, isCurrentTurn, card }: TableEntryProps) {
  const content = (
    <div
      className={`table-entry-container ${isCurrentTurn ? "current-turn" : ""}`}
    >
      <p className="player-name">{player.name}</p>
      <div className="table-entry">
        {card && <CardImg name={cardToImgName(card)} />}
      </div>
    </div>
  );

  return isCurrentTurn ? (
    <Droppable id={`table-slot-${player.id}`}>{content}</Droppable>
  ) : (
    content
  );
}

export function cardToImgName(card: Card): string {
  const suit = suitFromValue[card.suit];
  const rank = rankFromValue[card.rank];

  if (!suit || !rank) {
    throw new Error(`Invalid card: ${JSON.stringify(card)}`);
  }

  return `${suit}-${rank}`;
}
