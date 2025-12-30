import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

const WS_BASE = `ws://localhost:${import.meta.env.VITE_PORT}`;

type PlayersUpdateMsg = {
  type: "players_update";
  players: string[];
};

interface SessionSocketProps {
  sessionId: string;
  playerName: string;
}

export default function SessionSocket({
  sessionId,
  playerName,
}: SessionSocketProps) {
  const navigate = useNavigate();
  const wsRef = useRef<WebSocket | null>(null);
  const [players, setPlayers] = useState<string[]>([]);

  useEffect(() => {
    const ws = new WebSocket(
      `${WS_BASE}/ws?sessionId=${sessionId}&playerName=${playerName}`
    );

    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);

      switch (msg.type) {
        case "players_update":
          setPlayers(msg.players);
          break;
        case "game_started":
          navigate(`/game/${sessionId}`);
          break;
        case "error":
          alert(msg.payload.message);
          ws.close();
          navigate("/");
          break;
      }
    };

    ws.onerror = () => {
      navigate("/");
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [sessionId, playerName, navigate]);

  const startGame = () => {
    wsRef.current?.send(JSON.stringify({ type: "start_game", payload: {} }));
  };

  return (
    <div>
      <h3>Players:</h3>
      <ul>
        {players.map((p) => (
          <li key={p}>{p}</li>
        ))}
      </ul>

      <button onClick={startGame} disabled={players.length < 2}>
        Start Game
      </button>
    </div>
  );
}
