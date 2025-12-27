package app

import (
	"sync"
)

type Session struct {
	ID string `json:"sessionId"`
}

type SessionStore struct {
	sessions map[string]Session
	mu       sync.Mutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionId]
	return session, exists
}

func (s *SessionStore) DeleteSession(sessionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionId)
}
