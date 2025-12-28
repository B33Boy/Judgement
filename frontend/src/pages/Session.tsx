import { Navigate, useNavigate, useParams } from "react-router-dom";
// import { useEffect, useState, type SetStateAction } from "react";
import { useSessionValidation } from "../hooks/useSessionValidation";

export default function SessionPage() {
  const { sessionId } = useParams();
  const navigate = useNavigate();

  // const [currentPlayers, setCurrentPlayers] = useState();
  const { valid, loading } = useSessionValidation(sessionId, navigate);

  if (loading) return <p> Validating Session</p>;

  if (!valid) {
    return <Navigate to="/" />;
  }

  return (
    <>
      <h2>Session ID: {sessionId}</h2>
    </>
  );
}
