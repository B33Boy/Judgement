import type { Scores } from "../types";

export default function ScoreTable({ scores }: { scores: Scores | null }) {
  if (!scores) return;

  const NUM_ROUNDS = 14;
  const entries = Array.from(scores.entries());

  return (
    <table className="score-table">
      <thead>
        <tr>
          <th>Player</th>
          {Array.from({ length: NUM_ROUNDS }, (_, i) => (
            <th key={i}>R{i + 1}</th>
          ))}
        </tr>
      </thead>

      <tbody>
        {entries.map(([player, playerScores]) => (
          <tr key={player}>
            <td>{player}</td>
            {Array.from({ length: NUM_ROUNDS }, (_, i) => (
              <td key={i}>{playerScores?.[i] ?? ""}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
