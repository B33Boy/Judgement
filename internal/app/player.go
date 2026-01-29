package app

import (
	"context"

	t "github.com/B33Boy/Judgement/internal/types"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func NewPlayer(playerName string, conn *websocket.Conn) *t.Player {
	ctx, cancel := context.WithCancel(context.Background())

	return &t.Player{
		ID:         t.PlayerID(uuid.NewString()),
		PlayerName: playerName,
		Conn:       conn,
		Send:       make(chan t.Envelope, 100),
		Ctx:        ctx,
		Cancel:     cancel,
	}
}
