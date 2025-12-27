import { useNavigate, useParams } from "react-router-dom";
import { useEffect } from "react";

export default function SessionPage() {
  const { sessionId } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    async function validateSession() {
      console.log(`${sessionId}`);
      const res = await fetch(
        `http://localhost:${import.meta.env.VITE_PORT}/api/session/${sessionId}`
      );

      if (!res.ok) {
        alert("Invalid Session Id");
        navigate("/");
      }
    }
    validateSession();
  }, [sessionId, navigate]);

  return <h2>Session ID: {sessionId}</h2>;
}
