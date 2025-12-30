package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func (a *App) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", a.HealthHandler)
		r.Post("/session", a.CreateSessionHandler)
		r.Get("/session/{sessionId}", a.GetSessionHandler)
	})
	r.Get("/ws", a.wsHandler)

	return r
}

func (a *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Status": "Connection Healthy",
	})
}

func (a *App) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {

	session := a.sessionStore.GenerateRandomSession()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*session)
}

func (a *App) GetSessionHandler(w http.ResponseWriter, r *http.Request) {

	sessionId := chi.URLParam(r, "sessionId")

	if _, found := a.sessionStore.GetSession(sessionId); !found {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) wsHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")
	playerName := r.URL.Query().Get("playerName")

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*"},
	})
	if err != nil {
		log.Println("websocket accept failed:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	log.Println("websocket connected")

	session, exists := a.sessionStore.GetSession(sessionId)

	if !exists {
		// Throw error
		log.Println("Session doesn't exist")
		return
	}

	session.AddPlayer(NewPlayer(playerName, conn))

	broadcastPlayers(session)

	for {
		var msg Message
		err := wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			session.mu.Lock()
			delete(session.Players, playerName)
			session.mu.Unlock()
			broadcastPlayers(session)
			return
		}

		if msg.Type == "start_game" {
			session.Started = true
			broadcastGameStarted(session)
		}
	}
}

func broadcastPlayers(session *Session) {
	session.mu.Lock()
	defer session.mu.Unlock()

	playerNames := make([]string, 0, len(session.Players))
	for name := range session.Players {
		playerNames = append(playerNames, name)
	}

	msg := Message{
		Type:        "players_update",
		PlayerNames: playerNames,
	}

	for name, player := range session.Players {
		err := wsjson.Write(context.Background(), player.Conn, msg)
		if err != nil {
			// Connection is dead: remove it
			player.Conn.CloseNow()
			delete(session.Players, name)
		}
	}
}

func broadcastGameStarted(session *Session) {
	session.mu.Lock()
	defer session.mu.Unlock()

	msg := Message{
		Type: "game_started",
	}

	for name, player := range session.Players {
		err := wsjson.Write(context.Background(), player.Conn, msg)
		if err != nil {
			player.Conn.CloseNow()
			delete(session.Players, name)
		}
	}
}
