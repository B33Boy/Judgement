package app

import (
	"context"
	"log"
	"math/rand"
	"sync"
)

type Session struct {
	ID      string             `json:"sessionId"`
	Players map[string]*Player `json:"players"`

	inputs  chan GameInput
	outputs chan Envelope

	game *Game

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.Mutex
}

func NewSession(sessionId string) *Session {

	ctx, cancel := context.WithCancel(context.Background())

	s := &Session{
		ID:      sessionId,
		Players: make(map[string]*Player),

		inputs:  make(chan GameInput, 32),
		outputs: make(chan Envelope, 32),

		game:   nil,
		ctx:    ctx,
		cancel: cancel,
	}

	go s.run()

	return s
}

func (s *Session) AddPlayer(player *Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prevent duplicate players by kicking out old one first
	if old, ok := s.Players[player.PlayerName]; ok {
		old.cancel()
		close(old.Send)
	}

	s.Players[player.PlayerName] = player
}

func (s *Session) RemovePlayer(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, ok := s.Players[name]; ok {
		player.cancel()    // stop the write loop
		close(player.Send) // close outbound channel
		delete(s.Players, name)

		if len(s.Players) == 0 {
			s.cancel()
		}
	}
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

func (s *Session) run() {
	for {
		select {
		case <-s.ctx.Done():
			return

		case input := <-s.inputs:
			s.handleInput(input)

		case output := <-s.outputs:
			s.handleOutput(output)
		}
	}
}

func (s *Session) handleInput(input GameInput) {
	switch input.Env.Type {
	case MsgStartGame:
		if s.game != nil {
			return // already started
		}
		s.game = NewGame(s)
		s.game.Start()
	default:
		if s.game == nil {
			return
		}
		s.game.HandleGameInput(input)
	}
}

func (s *Session) handleOutput(env Envelope) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, player := range s.Players {
		select {
		case player.Send <- env:
			// success
		default:
			// slow client, drop or disconnect
			log.Println("Dropping message for slow player:", name)
		}
	}
}

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
