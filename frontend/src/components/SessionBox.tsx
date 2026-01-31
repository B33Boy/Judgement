type SessionBoxProps = {
  sessionId: string;
  playerName: string;
};

export default function SessionBox({ sessionId, playerName }: SessionBoxProps) {
  return (
    <div className="session-box">
      <h3>Session: {sessionId}</h3>
      <h3>{playerName}</h3>
    </div>
  );
}
