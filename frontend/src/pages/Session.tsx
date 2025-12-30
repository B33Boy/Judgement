import { Navigate, useNavigate, useParams } from "react-router-dom";
import { useSessionValidation } from "../hooks/useSessionValidation";
import SessionSocket from "../components/SessionSocket";

export default function SessionPage() {
  const navigate = useNavigate();

  const { sessionId } = useParams();
  const { valid, loading } = useSessionValidation(sessionId, navigate);

  const playerName = localStorage.getItem("playerName") || "";

  if (loading) return <p>Validating session...</p>;
  if (!valid) return <Navigate to="/" />;
  if (!sessionId) return <Navigate to="/" />;

  return (
    <>
      <h2>Session ID: {sessionId}</h2>
      <h3>Welcome player: {playerName}</h3>
      <SessionSocket sessionId={sessionId} playerName={playerName} />
    </>
  );
}
