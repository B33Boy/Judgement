package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (a *App) wsHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")
	playerName := r.URL.Query().Get("playerName")

	if sessionId == "" || playerName == "" {
		http.Error(w, "missing sessionId or playerName", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*"},
	})
	if err != nil {
		log.Println("websocket accept failed:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	session, exists := a.sessionStore.GetSession(sessionId)
	if !exists {
		log.Println("session does not exist")
		return
	}

	player := NewPlayer(playerName, conn)
	onPlayerJoin(r.Context(), session, player)
	defer onPlayerLeave(r.Context(), session, player)

	for {
		var env Envelope
		if err := wsjson.Read(r.Context(), conn, &env); err != nil {
			log.Println("ws read failed:", err)
			return
		}

		if err := dispatchMessage(r.Context(), session, player, env); err != nil {
			log.Println("dispatch failed:", err)
			return
		}
	}
}

func onPlayerJoin(ctx context.Context, session *Session, player *Player) {
	session.AddPlayer(player)
	broadcastPlayers(ctx, session)
}

func onPlayerLeave(ctx context.Context, session *Session, player *Player) {
	session.RemovePlayer(player.PlayerName)
	broadcastPlayers(ctx, session)
}

func broadcastPlayers(ctx context.Context, session *Session) {
	players := session.CopyPlayerList()

	names := make([]string, 0, len(players))
	for _, p := range players {
		names = append(names, p.PlayerName)
	}

	for _, p := range players {
		if err := send(
			ctx,
			p.Conn,
			MsgPlayersUpdate,
			PlayersUpdatePayload{PlayerNames: names},
		); err != nil {
			p.Conn.CloseNow()
			session.RemovePlayer(p.PlayerName)
		}
	}
}

func broadcastGameStarted(ctx context.Context, session *Session) {
	for _, p := range session.CopyPlayerList() {
		send(ctx, p.Conn, MsgGameStarted, GameStartedPayload{})
	}
}

func dispatchMessage(ctx context.Context, session *Session, player *Player, env Envelope) error {
	switch env.Type {
	case MsgStartGame:
		if success := session.Start(ctx); !success {
			return nil
		}
		broadcastGameStarted(ctx, session)

	case MsgMakeBid, MsgPlayCard:
		session.actions <- Action{Player: player, Message: env}

	default:
		log.Printf("unknown message type: %s", env.Type)
	}

	return nil
}

// ================= Send Abstraction =================

func send(ctx context.Context, conn *websocket.Conn, t MessageType, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Can't marshal the payload!")
		return err
	}
	return wsjson.Write(ctx, conn, Envelope{
		Type:    t,
		Payload: b,
	})
}
