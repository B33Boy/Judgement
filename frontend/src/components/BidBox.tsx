import { useState } from "react";
import type { MessageHandler } from "../context/GameContext";

type BidBoxProps = {
  msgFunction: MessageHandler;
};

export default function BidBox({ msgFunction }: BidBoxProps) {
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
