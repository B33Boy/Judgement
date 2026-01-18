import {
  createContext,
  useContext,
  useRef,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";
import type { WSEnvelope, RoundInfo } from "../types";

const WS_BASE = `ws://localhost:${import.meta.env.VITE_PORT}`;

export type MessageHandler = (type: WSEnvelope["type"], payload?: any) => void;

interface GameContextType {
  isConnected: boolean;
  players: string[];
  hand: string[];
  roundInfo: RoundInfo | null;
  connect: (sessionId: string, playerName: string) => void;
  disconnect: () => void;
  sendMessage: MessageHandler;
}

const GameContext = createContext<GameContextType | null>(null);

export function GameProvider({ children }: { children: ReactNode }) {
  // States
  const [isConnected, setIsConnected] = useState(false);
  const [players, setPlayers] = useState<string[]>([]);
  const [hand, setHand] = useState<string[]>([]);
  const [roundInfo, setRoundInfo] = useState<RoundInfo | null>(null);

  // Refs
  const currentPlayerRef = useRef<string>("");
  const wsRef = useRef<WebSocket | null>(null);
  const navigate = useNavigate();

  const connect = useCallback(
    (sessionId: string, playerName: string) => {
      // If we are already connected to this session, don't reconnect
      if (wsRef.current?.readyState === WebSocket.OPEN) return;
      currentPlayerRef.current = playerName;

      // Create new ws conn
      const ws = new WebSocket(
        `${WS_BASE}/ws?sessionId=${sessionId}&playerName=${playerName}`,
      );

      ws.onopen = () => setIsConnected(true);
      ws.onclose = () => {
        setIsConnected(false);
        wsRef.current = null;
      };

      ws.onmessage = (e) => {
        const msg = JSON.parse(e.data) as WSEnvelope;

        switch (msg.type) {
          // Lobby
          case "players_update":
            setPlayers(msg.payload.players ?? []);
            break;

          case "game_started":
            navigate(`/game/${sessionId}/${currentPlayerRef.current}`);
            break;

          // Game Logic
          case "player_hand":
            console.log("Received cards:", msg.payload.cards);
            setHand(msg.payload.cards ?? []);
            break;

          case "round_info":
            console.log("Round Info: ", msg.payload);
            const payload = msg.payload as RoundInfo;
            setRoundInfo(payload);
            break;

          case "error":
            alert(msg.payload.message);
            break;
        }
      };

      wsRef.current = ws;
    },
    [navigate],
  );

  const disconnect = useCallback(() => {
    wsRef.current?.close();
    wsRef.current = null;
  }, []);

  const sendMessage = useCallback((type: WSEnvelope["type"], payload?: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, payload }));
    }
  }, []);

  return (
    <GameContext.Provider
      value={{
        isConnected,
        players,
        hand,
        roundInfo,
        connect,
        disconnect,
        sendMessage,
      }}
    >
      {children}
    </GameContext.Provider>
  );
}

export function useGame() {
  const context = useContext(GameContext);
  if (!context) throw new Error("useGame must be used within a GameProvider");
  return context;
}
