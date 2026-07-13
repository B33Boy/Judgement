package game

import (
	"testing"

	t "github.com/B33Boy/Judgement/internal/types"
)

func TestIsCardPlayable_MustFollowLedSuit(t2 *testing.T) {
	trump := Club
	led := Card{Suit: Heart, Rank: Two}

	g := &Game{
		cardstack: []Card{led},
		state:     &GameState{TrumpSuit: &trump},
	}

	player := &GamePlayer{
		ID: "p",
		Cards: Hand{
			{Suit: Heart, Rank: King}, // has led suit
			{Suit: Club, Rank: Ace},   // trump
		},
	}

	// Holding the led suit, a trump card is NOT a legal substitute.
	if g.isCardPlayable(player, Card{Suit: Club, Rank: Ace}) {
		t2.Errorf("expected trump to be illegal when player can follow suit")
	}

	// The led-suit card is legal.
	if !g.isCardPlayable(player, Card{Suit: Heart, Rank: King}) {
		t2.Errorf("expected led-suit card to be legal")
	}
}

func TestIsCardPlayable_LedSuitIsFirstCardNotLastPlayed(t2 *testing.T) {
	// Trick so far: led Hearts, second player had no hearts and discarded
	// a Spade. A third player who *does* have a Heart must still follow
	// the suit that was led (Hearts), not the suit of the last card played.
	g := &Game{
		cardstack: []Card{
			{Suit: Heart, Rank: Two},
			{Suit: Spade, Rank: King},
		},
		state: &GameState{TrumpSuit: nil},
	}

	player := &GamePlayer{
		ID: "p",
		Cards: Hand{
			{Suit: Heart, Rank: Ace},
			{Suit: Club, Rank: Two},
		},
	}

	if !g.isCardPlayable(player, Card{Suit: Heart, Rank: Ace}) {
		t2.Errorf("expected the Heart to be playable (it's the led suit)")
	}
	if g.isCardPlayable(player, Card{Suit: Club, Rank: Two}) {
		t2.Errorf("expected the Club to be illegal - player still holds the led suit (Hearts)")
	}
}

func TestIsCardPlayable_NoLedSuitAnythingGoes(t2 *testing.T) {
	g := &Game{
		cardstack: []Card{{Suit: Heart, Rank: Two}},
		state:     &GameState{TrumpSuit: nil},
	}

	player := &GamePlayer{
		ID:    "p",
		Cards: Hand{{Suit: Club, Rank: Two}, {Suit: Spade, Rank: Ace}},
	}

	if !g.isCardPlayable(player, Card{Suit: Spade, Rank: Ace}) {
		t2.Errorf("expected any card to be playable when the player has none of the led suit")
	}
}

func TestValidateBid_HookRule(t2 *testing.T) {
	players := PlayerMap{
		"a": {ID: "a", Cards: Hand{{}, {}, {}}}, // 3 cards this round
		"b": {ID: "b", Cards: Hand{{}, {}, {}}},
		"c": {ID: "c", Cards: Hand{{}, {}, {}}},
	}
	cycler := NewPlayerCycler(players)
	cycler.StartFrom("a")
	cycler.Next() // -> b
	cycler.Next() // -> c, last bidder

	g := &Game{
		cycler: cycler,
		state: &GameState{
			Bids: map[t.PlayerID]Bid{"a": 1, "b": 1}, // sums to 2 so far
		},
	}

	// 3 cards dealt; a=1, b=1; if c also bids 1, total is exactly 3 -> illegal.
	if err := g.validateBid(players["c"], Bid(1)); err == nil {
		t2.Errorf("expected hook-rule rejection when total bids would equal cards dealt")
	}

	// Any other bid is fine for the last bidder.
	if err := g.validateBid(players["c"], Bid(0)); err != nil {
		t2.Errorf("expected bid of 0 to be valid, got error: %v", err)
	}
	if err := g.validateBid(players["c"], Bid(2)); err != nil {
		t2.Errorf("expected bid of 2 to be valid, got error: %v", err)
	}
}

func TestValidateBid_NotLastBidderNoHookRestriction(t2 *testing.T) {
	players := PlayerMap{
		"a": {ID: "a", Cards: Hand{{}, {}, {}}},
		"b": {ID: "b", Cards: Hand{{}, {}, {}}},
		"c": {ID: "c", Cards: Hand{{}, {}, {}}},
	}
	cycler := NewPlayerCycler(players)
	cycler.StartFrom("a") // "a" is not the last bidder

	g := &Game{
		cycler: cycler,
		state:  &GameState{Bids: map[t.PlayerID]Bid{}},
	}

	// Even a bid that would (hypothetically) sum to the hand size is fine
	// for a non-last bidder.
	if err := g.validateBid(players["a"], Bid(3)); err != nil {
		t2.Errorf("expected non-last bidder to bid freely, got error: %v", err)
	}
}

func TestValidateBid_OutOfRange(t2 *testing.T) {
	players := PlayerMap{
		"a": {ID: "a", Cards: Hand{{}, {}}},
	}
	cycler := NewPlayerCycler(players)
	cycler.StartFrom("a")

	g := &Game{
		cycler: cycler,
		state:  &GameState{Bids: map[t.PlayerID]Bid{}},
	}

	if err := g.validateBid(players["a"], Bid(-1)); err == nil {
		t2.Errorf("expected negative bid to be rejected")
	}
	if err := g.validateBid(players["a"], Bid(3)); err == nil {
		t2.Errorf("expected bid greater than hand size to be rejected")
	}
}
