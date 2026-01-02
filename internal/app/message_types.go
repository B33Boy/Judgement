package app

import (
	"encoding/json"
)

// ================= MessageTypes =================
type MessageType string

const (
	// Session-level
	MsgPlayersUpdate MessageType = "players_update"
	MsgStartGame     MessageType = "start_game"

	// Game-level
	MsgGameStarted MessageType = "game_started"
	MsgGameEnd     MessageType = "game_end"
	MsgMakeBid     MessageType = "make_bid"
	MsgPlayCard    MessageType = "play_card"
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
