package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	t "github.com/B33Boy/Judgement/internal/types"

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

	session, exists := a.sessionStore.GetSession(sessionId)
	if !exists {
		log.Println("session does not exist")
		return
	}

	player := NewPlayer(playerName, conn)

	defer func() {
		player.Cancel() // stops write loop
		conn.Close(websocket.StatusNormalClosure, "")
		onPlayerLeave(session, player)
	}()

	// ====== Write Loop ======
	go func() {
		for {
			select {
			case <-player.Ctx.Done():
				log.Printf("Player %v exited!", player.PlayerName)
				return
			case <-session.ctx.Done():
				log.Printf("[write loop] Closed session %v", session.ID)
				return
			case env, ok := <-player.Send:
				if !ok {
					return
				}
				if err := wsjson.Write(r.Context(), conn, env); err != nil {
					log.Println("write to websocket failed:", err)
					return
				}
			}
		}
	}()

	onPlayerJoin(session, player)

	// ====== Read Loop ======
	for {
		var env t.Envelope
		if err := wsjson.Read(r.Context(), conn, &env); err != nil {

			status := websocket.CloseStatus(err)

			switch status {
			case websocket.StatusNormalClosure, websocket.StatusGoingAway:
				log.Printf("player %s disconnected", player.PlayerName)

			default:
				log.Printf("ws read error (%s): %v", player.PlayerName, err)
			}
			return
		}

		if err := handleIncomingMessage(session, player, env); err != nil {
			return
		}
	}

}

func onPlayerJoin(session *Session, player *t.Player) {
	session.AddPlayer(player)
	sendWelcome(player, session)
	broadcastPlayersUpdate(session)
	log.Printf("Player (%v) added to session (%v)\n", player.PlayerName, session.ID)
}

func onPlayerLeave(session *Session, player *t.Player) {
	log.Printf("Player (%v) left session (%v)\n", player.PlayerName, session.ID)
	session.RemovePlayer(player)
	broadcastPlayersUpdate(session)
}

func sendWelcome(player *t.Player, session *Session) {
	out := t.GameOutput{
		Players: []t.PlayerID{player.ID},
		Env: t.Envelope{
			Type:    t.MsgWelcome,
			Payload: mustMarshal(player.ID),
		},
	}

	select {
	case session.outputs <- out:
	case <-session.ctx.Done():
		log.Printf("[sendWelcome] Closed session %v", session.ID)
	}
}

func broadcastPlayersUpdate(session *Session) {
	players := session.CopyPlayerList()

	public := make([]PlayerPublic, 0, len(players))
	allIDs := make([]t.PlayerID, 0, len(players))

	for _, p := range players {
		allIDs = append(allIDs, p.ID)
		public = append(public, PlayerPublic{
			ID:   p.ID,
			Name: p.PlayerName,
		})
	}

	out := t.GameOutput{
		Players: allIDs,
		Env: t.Envelope{
			Type:    t.MsgPlayersUpdate,
			Payload: mustMarshal(public),
		},
	}

	select {
	case session.outputs <- out:
	case <-session.ctx.Done():
		log.Printf("[broadcastPlayersUpdate] Closed session %v", session.ID)
	}
}

func handleIncomingMessage(session *Session, player *t.Player, env t.Envelope) error {
	select {
	case session.inputs <- t.GameInput{Player: player, Env: env}:
		return nil

	case <-session.ctx.Done():
		log.Printf("[handleIncomingMessage] Closed session %v", session.ID)
		return errors.New("session closed")
	}
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Panicf("marshal failed: %v", err)
	}
	return b
}
