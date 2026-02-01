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
    | "make_bid"
    | "play_card";
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

// export type Card = {
//   suit: string;
//   rank: string;
// };

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
