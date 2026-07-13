package game

// End-to-end test driving a full 3-player game through the public
// NewGame/HandleGameInput surface, the same way the app package would.

import (
	"context"
	"encoding/json"
	"testing"

	t "github.com/B33Boy/Judgement/internal/types"
)

type fakeSession struct {
	ctx     context.Context
	players map[t.PlayerID]*t.Player
}

func (f *fakeSession) Context() context.Context             { return f.ctx }
func (f *fakeSession) GetPlayers() map[t.PlayerID]*t.Player { return f.players }
func (f *fakeSession) Emit(t.GameOutput)                    {}

func TestFullGamePlaysToCompletion(t2 *testing.T) {
	players := map[t.PlayerID]*t.Player{
		"p1": {ID: "p1", PlayerName: "Alice"},
		"p2": {ID: "p2", PlayerName: "Bob"},
		"p3": {ID: "p3", PlayerName: "Carol"},
	}
	session := &fakeSession{ctx: context.Background(), players: players}

	g := NewGame(session)
	g.Start()

	tricksPlayed := 0
	safety := 0

	for g.sm.state != StateGameOver {
		safety++
		if safety > 5000 {
			t2.Fatalf("game stuck: round=%d state=%s", g.state.Round, g.sm.state)
		}

		switch g.sm.state {
		case StateBid:
			player := g.Players[g.state.TurnPlayer]
			sendBid(t2, g, player.ID, legalBid(g, player))

		case StatePlay:
			player := g.Players[g.state.TurnPlayer]
			if len(g.cardstack) == 0 {
				tricksPlayed++
			}
			sendPlay(t2, g, player.ID, legalCard(t2, g, player))

		default:
			t2.Fatalf("unexpected state %s outside of the client-driven loop", g.sm.state)
		}
	}

	if int(g.state.Round) != len(roundSchedule) {
		t2.Errorf("expected game to finish after %d rounds, got %d", len(roundSchedule), g.state.Round)
	}

	wantTricks := 0
	for _, n := range roundSchedule {
		wantTricks += n
	}
	if tricksPlayed != wantTricks {
		t2.Errorf("expected %d total tricks across the game, got %d", wantTricks, tricksPlayed)
	}

	for id, scores := range g.scores {
		if len(scores) != len(roundSchedule) {
			t2.Errorf("player %s: expected %d score entries, got %d", id, len(roundSchedule), len(scores))
		}
	}

	for id, player := range g.Players {
		if len(player.Cards) != 0 {
			t2.Errorf("player %s finished the game with %d cards still in hand", id, len(player.Cards))
		}
	}
}

func sendBid(t2 *testing.T, g *Game, playerID t.PlayerID, bid Bid) {
	payload, err := json.Marshal(MakeBid{Bid: bid})
	if err != nil {
		t2.Fatalf("marshal bid: %v", err)
	}
	g.HandleGameInput(t.GameInput{
		Player: &t.Player{ID: playerID},
		Env:    t.Envelope{Type: t.MsgMakeBid, Payload: payload},
	})
}

func sendPlay(t2 *testing.T, g *Game, playerID t.PlayerID, card Card) {
	payload, err := json.Marshal(card)
	if err != nil {
		t2.Fatalf("marshal card: %v", err)
	}
	g.HandleGameInput(t.GameInput{
		Player: &t.Player{ID: playerID},
		Env:    t.Envelope{Type: t.MsgPlayCard, Payload: payload},
	})
}

// legalBid picks 0, falling back to 1 if 0 would trip the hook rule (only
// possible for the last bidder in a round) - exactly one bid value is ever
// banned, so 1 is always safe when 0 isn't.
func legalBid(g *Game, player *GamePlayer) Bid {
	if err := g.validateBid(player, 0); err == nil {
		return 0
	}
	return 1
}

// legalCard returns any card from the player's hand that isCardPlayable
// accepts.
func legalCard(t2 *testing.T, g *Game, player *GamePlayer) Card {
	for _, card := range player.Cards {
		if g.isCardPlayable(player, card) {
			return card
		}
	}
	t2.Fatalf("player %s has no legal card among %v", player.ID, player.Cards)
	return Card{}
}
