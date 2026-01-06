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

type StartGamePayload struct{}
