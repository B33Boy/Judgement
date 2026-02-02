import {
  createContext,
  useContext,
  useRef,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";
import type { WSEnvelope, GameState, Players } from "../types";

const WS_BASE = `ws://localhost:${import.meta.env.VITE_PORT}`;

export type MessageHandler = (type: WSEnvelope["type"], payload?: any) => void;

interface GameContextType {
  isConnected: boolean;
  playerId: string | null;
  players: Players;
  hand: string[];
  gameState: GameState | null;
  connect: (sessionId: string, playerName: string) => void;
  disconnect: () => void;
  sendMessage: MessageHandler;
}

const GameContext = createContext<GameContextType | null>(null);

export function GameProvider({ children }: { children: ReactNode }) {
  // States
  const [isConnected, setIsConnected] = useState(false);
  const [playerId, setPlayerId] = useState<string | null>(null);
  const [players, setPlayers] = useState<Players>([]);
  const [hand, setHand] = useState<string[]>([]);
  // const [roundInfo, setRoundInfo] = useState<RoundInfo | null>(null);
  // const [scores, setScores] = useState<Scores | null>(null);

  // Inside your Provider:
  const [gameState, setGameState] = useState<GameState | null>(null);

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
          case "welcome":
            setPlayerId(msg.payload);
            break;

          case "players_update":
            setPlayers(msg.payload ?? []);
            break;

          case "game_started":
            navigate(`/game/${sessionId}/${currentPlayerRef.current}`);
            break;

          // Game Logic
          case "player_hand":
            console.log("Received cards:", msg.payload.cards);
            setHand(msg.payload.cards ?? []);
            break;

          case "state_sync":
            console.log("State Sync: ", msg.payload);
            setGameState(msg.payload);
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
        playerId,
        players,
        hand,
        gameState,
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
