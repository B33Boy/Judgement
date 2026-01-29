package types

/*
	Common types
*/

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
)

// ================= App Types =================
type PlayerID string

type Player struct {
	ID         PlayerID `json:"id"`
	PlayerName string   `json:"playerName"`
	Conn       *websocket.Conn
	Send       chan Envelope
	Ctx        context.Context
	Cancel     context.CancelFunc
}

// App Message Types
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

// ================= Transmission Types =================

type Envelope struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type GameInput struct {
	Player *Player
	Env    Envelope
}

type GameOutput struct {
	Players []PlayerID
	Env     Envelope
}
