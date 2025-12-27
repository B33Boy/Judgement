package app

import (
	"math/rand"
	"sync"
)

type Session struct {
	ID string `json:"sessionId"`
}

type SessionStore struct {
	sessions map[string]Session
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]Session),
	}
}

func (s *SessionStore) CreateSession(sessionId string) Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	newSession := Session{ID: sessionId}
	s.sessions[sessionId] = newSession
	return newSession
}

func (s *SessionStore) GetSession(sessionId string) (Session, bool) {
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

func (s *SessionStore) GenerateRandomSession() Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	const LENGTH = 8
	randId := randomString(LENGTH)

	for _, exists := s.sessions[randId]; exists; {
		randId = randomString(LENGTH)
	}

	newSession := Session{ID: randId}
	s.sessions[randId] = newSession
	return newSession
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
