import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

const WS_BASE = `ws://localhost:${import.meta.env.VITE_PORT}`;

interface EnvelopeMessage {
  type: "players_update" | "game_started" | "error";
  payload?: any;
}

interface SessionInfo {
  sessionId: string;
  playerName: string;
}

export default function SessionSocket({ sessionId, playerName }: SessionInfo) {
  const navigate = useNavigate();
  const wsRef = useRef<WebSocket | null>(null);
  const [players, setPlayers] = useState<string[]>([]);

  useEffect(() => {
    const ws = new WebSocket(
      `${WS_BASE}/ws?sessionId=${sessionId}&playerName=${playerName}`
    );

    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data) as EnvelopeMessage;

      if (msg.type === "players_update") {
        setPlayers(msg.payload.players ?? []);
      }

      if (msg.type === "game_started") {
        navigate(`/game/${sessionId}`);
      }

      if (msg.type === "error") {
        alert(msg.payload?.message);
        navigate("/");
      }
    };

    ws.onerror = () => navigate("/");

    wsRef.current = ws;
    return () => ws.close();
  }, [sessionId, playerName, navigate]);

  function startGame() {
    wsRef.current?.send(JSON.stringify({ type: "start_game" }));
  }

  return (
    <>
      <h3>Players</h3>
      <ul>
        {players.map((p) => (
          <li key={p}>{p}</li>
        ))}
      </ul>
      <button disabled={players.length < 2} onClick={startGame}>
        Start Game
      </button>
    </>
  );
}
