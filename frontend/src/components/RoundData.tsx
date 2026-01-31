import type { RoundInfo } from "../types";

export default function RoundData({
  roundInfo,
}: {
  roundInfo: RoundInfo | null;
}) {
  return (
    <div className="round-info">
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
