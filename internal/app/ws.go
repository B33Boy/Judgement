package app

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Message struct {
	Type        string   `json:"type"`
	PlayerNames []string `json:"players,omitempty"`
}

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
	session.AddPlayer(player)
	broadcastPlayers(session)

	for {
		var msg Message
		if err := wsjson.Read(r.Context(), conn, &msg); err != nil {
			session.RemovePlayer(player.PlayerName)
			broadcastPlayers(session)
			return
		}

		switch msg.Type {
		case "start_game":
			if !session.Started {
				session.Started = true
				broadcastGameStarted(session)
			}
		}
	}
}

func broadcastPlayers(session *Session) {
	players := session.CopyPlayerList()

	names := make([]string, 0, len(players))
	for _, p := range players {
		names = append(names, p.PlayerName)
	}

	msg := Message{
		Type:        "players_update",
		PlayerNames: names,
	}

	for _, p := range players {
		if err := wsjson.Write(context.Background(), p.Conn, msg); err != nil {
			p.Conn.CloseNow()
			session.RemovePlayer(p.PlayerName)
		}
	}
}

func broadcastGameStarted(session *Session) {
	players := session.CopyPlayerList()

	msg := Message{
		Type: "game_started",
	}

	for _, p := range players {
		if err := wsjson.Write(context.Background(), p.Conn, msg); err != nil {
			p.Conn.CloseNow()
			session.RemovePlayer(p.PlayerName)
		}
	}
}
