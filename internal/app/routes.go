package app

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/coder/websocket"
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

		r.Get("/websocket", a.websocketHandler)
	})

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
	json.NewEncoder(w).Encode(session)
}

func (a *App) GetSessionHandler(w http.ResponseWriter, r *http.Request) {

	sessionId := chi.URLParam(r, "sessionId")

	if _, found := a.sessionStore.GetSession(sessionId); !found {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

// func (a *App) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	resp := make(map[string]string)
// 	resp["message"] = "Hello World"

// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}

// 	_, _ = w.Write(jsonResp)
// }
