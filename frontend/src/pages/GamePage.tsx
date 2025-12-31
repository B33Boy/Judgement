import { useParams, Navigate } from "react-router-dom";
import { useSessionValidation } from "../hooks/useSessionValidation";

export default function GamePage() {
  const { sessionId } = useParams();
  const { valid, loading } = useSessionValidation(sessionId);
  if (loading) return <p>Validating session...</p>;
  if (!valid || !sessionId) return <Navigate to="/" />;

  return <h1>Game!</h1>;
}
