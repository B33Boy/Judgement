package app

import (
	"math/rand"
	"sync"
)

type Session struct {
	ID      string             `json:"sessionId"`
	Players map[string]*Player `json:"players"`
	Started bool
	mu      sync.Mutex
}

func NewSession(sessionId string) *Session {
	return &Session{
		ID:      sessionId,
		Players: make(map[string]*Player),
	}
}

func (s *Session) AddPlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Players[player.PlayerName] = player
}

func (s *Session) RemovePlayer(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Players, name)
}

func (s *Session) CopyPlayerList() []*Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	players := make([]*Player, 0, len(s.Players))
	for _, p := range s.Players {
		players = append(players, p)
	}
	return players
}

// func (s *Session) DeletePlayers(playersToDisconnect []*Player) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	for _, p := range playersToDisconnect {
// 		delete(s.Players, p.PlayerName)
// 	}
// }

func (s *Session) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Started = true
}

// ==================================================

type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionStore) GetSession(sessionId string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, exists := s.sessions[sessionId]
	return session, exists
}

func (s *SessionStore) DeleteSession(sessionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionId)
}

func (s *SessionStore) GenerateRandomSession() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	const length = 8
	var id string

	for {
		id = randomString(length)
		if _, exists := s.sessions[id]; !exists {
			break
		}
	}

	session := NewSession(id)
	s.sessions[id] = session
	return session
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
