import { useDraggable } from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";

function Draggable(props: any) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: props.id,
  });

  // const style = transform
  //   ? {
  //       transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
  //     }
  //   : undefined;

  const style = {
    transform: CSS.Translate.toString(transform),
  };

  return (
    <button ref={setNodeRef} style={style} {...listeners} {...attributes}>
      {props.children}
    </button>
  );
}

export default function PlayerHand({ hand }: { hand: string[] }) {
  return (
    <div className="player-hand">
      {hand.length > 0 ? (
        hand.map((card) => (
          <Draggable key={card} id={`card-${card}`}>
            <div className="card-item">
              <Card name={card} />
            </div>
          </Draggable>
        ))
      ) : (
        <p>Waiting for cards...</p>
      )}
    </div>
  );
}

function Card({ name }: { name: string }) {
  return <img className="card-image" src={`/cards/${name}.svg`} alt={name} />;
}
