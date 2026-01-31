export default function PlayerHand({ hand }: { hand: string[] }) {
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
