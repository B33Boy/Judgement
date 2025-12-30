import { Navigate, useParams } from "react-router-dom";
import { useSessionValidation } from "../hooks/useSessionValidation";
import { getPlayerName } from "../lib/player";
import SessionSocket from "../components/SessionSocket";

export default function SessionPage() {
  const { sessionId } = useParams();
  const { valid, loading } = useSessionValidation(sessionId);

  if (loading) return <p>Validating session...</p>;
  if (!valid || !sessionId) return <Navigate to="/" />;

  const playerName = getPlayerName();

  return (
    <>
      <h2>Session {sessionId}</h2>
      <h3>Welcome {playerName}</h3>
      <SessionSocket sessionId={sessionId} playerName={playerName} />
    </>
  );
}
