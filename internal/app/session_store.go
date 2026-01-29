package app

import (
	"math/rand"
	"sync"
)

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
