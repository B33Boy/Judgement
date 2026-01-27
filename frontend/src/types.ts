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
    | "player_hand"
    | "round_info"
    | "make_bid";
  payload?: any;
}

export interface PlayerState {
  players: string[];
}

export interface RoundInfo {
  round: number;
  turnPlayer: string;
  state: string;
}

export type Scores = Map<string, number[]>;
