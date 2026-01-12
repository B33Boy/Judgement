export interface SessionInfo {
  sessionId: string;
  playerName: string;
}

export interface WSEnvelope {
  type:
    | "players_update"
    | "game_started"
    | "error"
    | "start_game"
    | "player_hand";
  payload?: any;
}

export interface PlayerState {
  players: string[];
}
