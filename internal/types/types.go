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
	// FE -> BE
	MsgStartGame MessageType = "start_game"
	MsgMakeBid   MessageType = "make_bid"
	MsgPlayCard  MessageType = "play_card"

	// BE -> FE
	MsgWelcome       MessageType = "welcome"
	MsgPlayersUpdate MessageType = "players_update"
	MsgGameStarted   MessageType = "game_started"
	MsgGameEnd       MessageType = "game_end"
	MsgPlayerHand    MessageType = "player_hand"
	MsgStateSync     MessageType = "state_sync"
	MsgInvalidAction MessageType = "invalid_action"
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
