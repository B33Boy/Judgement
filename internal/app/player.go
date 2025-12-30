package app

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Player struct {
	PlayerName string `json:"playerName"`
	Score      int    `json:"score"`
	Conn       *websocket.Conn
}

func NewPlayer(playerName string, conn *websocket.Conn) *Player {
	return &Player{
		PlayerName: playerName,
		Score:      0,
		Conn:       conn,
	}
}

func (p *Player) Send(msg any) error {
	return wsjson.Write(context.Background(), p.Conn, msg)
}
