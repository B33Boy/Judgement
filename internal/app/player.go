package app

import (
	"context"

	"github.com/coder/websocket"
)

type Player struct {
	PlayerName string `json:"playerName"`
	Conn       *websocket.Conn
	Send       chan Envelope
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewPlayer(playerName string, conn *websocket.Conn) *Player {
	ctx, cancel := context.WithCancel(context.Background())

	return &Player{
		PlayerName: playerName,
		Conn:       conn,
		Send:       make(chan Envelope, 100),
		ctx:        ctx,
		cancel:     cancel,
	}
}
