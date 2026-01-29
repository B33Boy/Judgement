package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *App) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	session := a.sessionStore.GenerateRandomSession()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := chi.URLParam(r, "sessionId")

	if _, exists := a.sessionStore.GetSession(sessionId); !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
