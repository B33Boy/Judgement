package app

import (
	"context"
	"log"
	"sync"

	g "github.com/B33Boy/Judgement/internal/game"
	t "github.com/B33Boy/Judgement/internal/types"
)

// Implement SessionView implicitly
func (s *Session) Context() context.Context {
	return s.ctx
}

func (s *Session) GetPlayers() map[t.PlayerID]*t.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.players
}

func (s *Session) Emit(out t.GameOutput) {
	select {
	case s.outputs <- out:
	case <-s.ctx.Done():
	}
}

type Session struct {
	ID      string `json:"sessionId"`
	players map[t.PlayerID]*t.Player

	inputs  chan t.GameInput
	outputs chan t.GameOutput

	game *g.Game

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.Mutex
}

func NewSession(sessionId string) *Session {

	ctx, cancel := context.WithCancel(context.Background())

	s := &Session{
		ID:      sessionId,
		players: make(map[t.PlayerID]*t.Player),

		inputs:  make(chan t.GameInput, 32),
		outputs: make(chan t.GameOutput, 32),

		game:   nil,
		ctx:    ctx,
		cancel: cancel,
	}

	go s.run()

	return s
}

func (s *Session) AddPlayer(player *t.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prevent duplicate players by kicking out old one first
	if old, ok := s.players[player.ID]; ok {
		old.Cancel()
		close(old.Send)
	}

	s.players[player.ID] = player
}

func (s *Session) RemovePlayer(player *t.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, ok := s.players[player.ID]; ok {
		player.Cancel()    // stop the write loop
		close(player.Send) // close outbound channel
		delete(s.players, player.ID)

		if len(s.players) == 0 {
			s.cancel()
		}
	}
}

func (s *Session) CopyPlayerList() []*t.Player {
	s.mu.Lock()
	defer s.mu.Unlock()

	players := make([]*t.Player, 0, len(s.players))
	for _, p := range s.players {
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

func (s *Session) handleInput(input t.GameInput) {
	switch input.Env.Type {
	case t.MsgStartGame:
		if s.game != nil {
			return // already started
		}
		s.game = g.NewGame(s)
		s.game.Start()
	default:
		if s.game == nil {
			return
		}
		s.game.HandleGameInput(input)
	}
}

func (s *Session) handleOutput(output t.GameOutput) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Route outputs to the specific players given in output
	for _, id := range output.Players {

		// Get Player from id common to Player and GamePlayer
		player := s.players[id]

		select {
		case player.Send <- output.Env:
			// success
		default:
			// slow client, drop or disconnect
			log.Println("Dropping message for slow player with ID:", id)
		}
	}
}
