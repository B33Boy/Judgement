package app

import (
	"fmt"
	"testing"

	t "github.com/B33Boy/Judgement/internal/types"
)

func newTestPlayer(id, name string) *t.Player {
	return &t.Player{
		ID:         t.PlayerID(id),
		PlayerName: name,
		Send:       make(chan t.Envelope, 8),
	}
}

func TestAddPlayer_RejectsOnceGameStarted(t2 *testing.T) {
	s := NewSession("test-session")
	defer s.cancel()

	p1 := newTestPlayer("p1", "Alice")
	if err := s.AddPlayer(p1); err != nil {
		t2.Fatalf("expected AddPlayer to succeed before game start, got: %v", err)
	}

	s.gameStarted.Store(true)

	p2 := newTestPlayer("p2", "Bob")
	if err := s.AddPlayer(p2); err == nil {
		t2.Errorf("expected AddPlayer to reject a join once the game has started")
	}

	s.mu.Lock()
	_, added := s.players["p2"]
	s.mu.Unlock()
	if added {
		t2.Errorf("rejected player should not have been added to the session")
	}
}

func TestValidatePlayerCount(t2 *testing.T) {
	s := NewSession("test-session")
	defer s.cancel()

	if err := s.validatePlayerCount(); err == nil {
		t2.Errorf("expected rejection with 0 players")
	}

	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("p%d", i)
		if err := s.AddPlayer(newTestPlayer(id, id)); err != nil {
			t2.Fatalf("AddPlayer failed: %v", err)
		}
	}
	if err := s.validatePlayerCount(); err != nil {
		t2.Errorf("expected 3 players to be valid, got: %v", err)
	}

	for i := 3; i < 8; i++ {
		id := fmt.Sprintf("p%d", i)
		if err := s.AddPlayer(newTestPlayer(id, id)); err != nil {
			t2.Fatalf("AddPlayer failed: %v", err)
		}
	}
	if err := s.validatePlayerCount(); err == nil {
		t2.Errorf("expected rejection with 8 players (max is 7)")
	}
}

func TestHandleInput_StartGame_RejectsBadPlayerCount(t2 *testing.T) {
	s := NewSession("test-session")
	defer s.cancel()

	p1 := newTestPlayer("p1", "Alice")
	if err := s.AddPlayer(p1); err != nil {
		t2.Fatalf("AddPlayer failed: %v", err)
	}

	s.handleInput(t.GameInput{
		Player: p1,
		Env:    t.Envelope{Type: t.MsgStartGame},
	})

	if s.game != nil {
		t2.Errorf("expected game not to start with only 1 player")
	}

	select {
	case env := <-p1.Send:
		if env.Type != t.MsgInvalidAction {
			t2.Errorf("expected invalid_action rejection, got %s", env.Type)
		}
	default:
		t2.Errorf("expected a rejection message to be sent to the requesting player")
	}
}

func TestHandleInput_StartGame_SucceedsWithValidPlayerCount(t2 *testing.T) {
	s := NewSession("test-session")
	defer s.cancel()

	players := []*t.Player{
		newTestPlayer("p1", "Alice"),
		newTestPlayer("p2", "Bob"),
		newTestPlayer("p3", "Carol"),
	}
	for _, p := range players {
		if err := s.AddPlayer(p); err != nil {
			t2.Fatalf("AddPlayer failed: %v", err)
		}
	}

	s.handleInput(t.GameInput{
		Player: players[0],
		Env:    t.Envelope{Type: t.MsgStartGame},
	})

	if s.game == nil {
		t2.Errorf("expected game to start with 3 players")
	}
	if !s.gameStarted.Load() {
		t2.Errorf("expected gameStarted flag to be set once the game begins")
	}

	late := newTestPlayer("late", "Dave")
	if err := s.AddPlayer(late); err == nil {
		t2.Errorf("expected a join after game start to be rejected")
	}
}
