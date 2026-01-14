package app

import (
	"encoding/json"
)

// ================= MessageTypes =================
type MessageType string

const (
	MsgPlayersUpdate MessageType = "players_update" // BE -> FE
	MsgStartGame     MessageType = "start_game"     // FE -> BE
	MsgGameStarted   MessageType = "game_started"   // BE -> FE

	MsgPlayerHand MessageType = "player_hand" // BE -> FE
	MsgRoundInfo  MessageType = "round_info"  // BE -> FE

	MsgGameEnd  MessageType = "game_end"  // BE -> FE
	MsgMakeBid  MessageType = "make_bid"  // FE -> BE
	MsgPlayCard MessageType = "play_card" // FE -> BE
)

type Envelope struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// ================= Payload Structs =================

type PlayersUpdatePayload struct {
	PlayerNames []string `json:"players"`
}

type PlayerHandChangePayload struct {
	Cards Hand `json:"cards"`
}

type RoundInfoPayload struct {
	Round      round  `json:"round"`
	TurnPlayer string `json:"turnPlayer"`
	State      State  `json:"state"`
}

type MakeBid struct {
	Bid bid `json:"bid"`
}
