export interface SessionInfo {
  sessionId: string;
  playerName: string;
}

export interface WSEnvelope {
  type:
    | "welcome"
    | "players_update"
    | "game_started"
    | "error"
    | "start_game"
    | "player_hand"
    | "round_info"
    | "make_bid"
    | "play_card"
    | "state_sync";
  payload?: any;
}

export interface GameState {
  round: number;
  state: "bidding" | "playing" | "resolution" | "gameover";
  turnPlayer: string; // PlayerID
  trumpSuit: string | null;
  table: Record<string, string | undefined>; // PlayerID -> CardID
  bids: Record<string, number>; // PlayerID -> bid
  handsWon: Record<string, number>; // PlayerID -> number hands won this
}

export type PlayerPublic = {
  id: string;
  name: string;
};
export type Players = PlayerPublic[];

export type Scores = Map<string, number[]>;

export type Card = {
  suit: string;
  rank: string;
};

// export type Cards = Card[];

export const suitMap: Record<string, number> = {
  SPADE: 0,
  HEART: 1,
  DIAMOND: 2,
  CLUB: 3,
};

export const rankMap: Record<string, number> = {
  "2": 2,
  "3": 3,
  "4": 4,
  "5": 5,
  "6": 6,
  "7": 7,
  "8": 8,
  "9": 9,
  "10": 10,
  JACK: 11,
  QUEEN: 12,
  KING: 13,
  ACE: 14,
};

export const suitFromValue: Record<string, string> = {
  0: "SPADE",
  1: "HEART",
  2: "DIAMOND",
  3: "CLUB",
};

export const rankFromValue: Record<string, string> = {
  2: "2",
  3: "3",
  4: "4",
  5: "5",
  6: "6",
  7: "7",
  8: "8",
  9: "9",
  10: "10",
  11: "JACK",
  12: "QUEEN",
  13: "KING",
  14: "ACE",
};
