package app

import (
	"context"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Player struct {
	ID         PlayerID `json:"id"`
	PlayerName string   `json:"playerName"`
	Conn       *websocket.Conn
	Send       chan Envelope
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewPlayer(playerName string, conn *websocket.Conn) *Player {
	ctx, cancel := context.WithCancel(context.Background())

	return &Player{
		ID:         PlayerID(uuid.NewString()),
		PlayerName: playerName,
		Conn:       conn,
		Send:       make(chan Envelope, 100),
		ctx:        ctx,
		cancel:     cancel,
	}
}
