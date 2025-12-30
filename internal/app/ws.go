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
			session.mu.Lock()
			delete(session.Players, playerName)
			session.mu.Unlock()
			broadcastPlayers(session)
			return
		}

		if msg.Type == "start_game" {
			session.mu.Lock()
			session.Started = true
			session.mu.Unlock()
			broadcastGameStarted(session)
		}
	}
}

func broadcastPlayers(session *Session) {
	session.mu.Lock()
	defer session.mu.Unlock()

	names := make([]string, 0, len(session.Players))
	for name := range session.Players {
		names = append(names, name)
	}

	msg := Message{
		Type:        "players_update",
		PlayerNames: names,
	}

	for name, player := range session.Players {
		if err := wsjson.Write(context.Background(), player.Conn, msg); err != nil {
			player.Conn.CloseNow()
			delete(session.Players, name)
		}
	}
}

func broadcastGameStarted(session *Session) {
	session.mu.Lock()
	defer session.mu.Unlock()

	msg := Message{Type: "game_started"}

	for name, player := range session.Players {
		if err := wsjson.Write(context.Background(), player.Conn, msg); err != nil {
			player.Conn.CloseNow()
			delete(session.Players, name)
		}
	}
}
