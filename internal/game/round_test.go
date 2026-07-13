package game

import (
	"testing"

	t "github.com/B33Boy/Judgement/internal/types"
)

func TestBeats(t2 *testing.T) {
	spade := Spade

	cases := []struct {
		name      string
		candidate Card
		current   Card
		ledSuit   Suit
		trump     *Suit
		want      bool
	}{
		{
			name:      "higher card of led suit wins",
			candidate: Card{Suit: Heart, Rank: King},
			current:   Card{Suit: Heart, Rank: Ten},
			ledSuit:   Heart,
			trump:     nil,
			want:      true,
		},
		{
			name:      "lower card of led suit loses",
			candidate: Card{Suit: Heart, Rank: Ten},
			current:   Card{Suit: Heart, Rank: King},
			ledSuit:   Heart,
			trump:     nil,
			want:      false,
		},
		{
			name:      "off-suit discard never wins",
			candidate: Card{Suit: Club, Rank: Ace},
			current:   Card{Suit: Heart, Rank: Two},
			ledSuit:   Heart,
			trump:     nil,
			want:      false,
		},
		{
			name:      "trump beats non-trump led suit",
			candidate: Card{Suit: Spade, Rank: Two},
			current:   Card{Suit: Heart, Rank: Ace},
			ledSuit:   Heart,
			trump:     &spade,
			want:      true,
		},
		{
			name:      "non-trump never beats trump",
			candidate: Card{Suit: Heart, Rank: Ace},
			current:   Card{Suit: Spade, Rank: Two},
			ledSuit:   Heart,
			trump:     &spade,
			want:      false,
		},
		{
			name:      "higher trump beats lower trump",
			candidate: Card{Suit: Spade, Rank: King},
			current:   Card{Suit: Spade, Rank: Two},
			ledSuit:   Heart,
			trump:     &spade,
			want:      true,
		},
	}

	for _, c := range cases {
		t2.Run(c.name, func(t2 *testing.T) {
			got := beats(c.candidate, c.current, c.ledSuit, c.trump)
			if got != c.want {
				t2.Errorf("beats(%v, %v, led=%v, trump=%v) = %v, want %v",
					c.candidate, c.current, c.ledSuit, c.trump, got, c.want)
			}
		})
	}
}

func TestDetermineTrickWinner(t2 *testing.T) {
	trump := Club

	cardA := Card{Suit: Heart, Rank: King} // led suit, high
	cardB := Card{Suit: Heart, Rank: Ace}  // led suit, higher
	cardC := Card{Suit: Club, Rank: Two}   // trump, low rank but still wins
	cardD := Card{Suit: Spade, Rank: King} // off-suit discard

	g := &Game{
		cardstack: []Card{cardA, cardB, cardC, cardD}, // "a" led
		state: &GameState{
			TrumpSuit: &trump,
			Table: map[t.PlayerID]*Card{
				"a": &cardA,
				"b": &cardB,
				"c": &cardC,
				"d": &cardD,
			},
		},
	}

	winner := g.determineTrickWinner()
	if winner != "c" {
		t2.Errorf("expected trump player 'c' to win, got %s", winner)
	}
}

func TestDetermineTrickWinner_NoTrump(t2 *testing.T) {
	cardA := Card{Suit: Heart, Rank: Two}
	cardB := Card{Suit: Heart, Rank: Ace}
	cardC := Card{Suit: Spade, Rank: King} // off-suit, can't win

	g := &Game{
		cardstack: []Card{cardA, cardB, cardC},
		state: &GameState{
			TrumpSuit: nil,
			Table: map[t.PlayerID]*Card{
				"a": &cardA,
				"b": &cardB,
				"c": &cardC,
			},
		},
	}

	winner := g.determineTrickWinner()
	if winner != "b" {
		t2.Errorf("expected 'b' (highest led-suit card) to win, got %s", winner)
	}
}

func TestScoreRound(t2 *testing.T) {
	bidTwo := Bid(2)
	bidOne := Bid(1)

	g := &Game{
		Players: PlayerMap{
			"hit":  {ID: "hit", Bid: &bidTwo},
			"miss": {ID: "miss", Bid: &bidOne},
			"none": {ID: "none", Bid: nil},
		},
		state: &GameState{
			Round: 3,
			HandsWon: map[t.PlayerID]int{
				"hit":  2,
				"miss": 0,
				"none": 1,
			},
		},
		scores: PlayerScore{
			"hit":  make([]Score, 14),
			"miss": make([]Score, 14),
			"none": make([]Score, 14),
		},
	}

	g.scoreRound()

	if got := g.scores["hit"][3]; got != Score(12) {
		t2.Errorf("expected hit player to score 10+2=12, got %d", got)
	}
	if got := g.scores["miss"][3]; got != Score(0) {
		t2.Errorf("expected missed bid to score 0, got %d", got)
	}
	if got := g.scores["none"][3]; got != Score(0) {
		t2.Errorf("expected player with no bid to score 0, got %d", got)
	}
}

func TestRoundComplete(t2 *testing.T) {
	g := &Game{
		Players: PlayerMap{
			"a": {ID: "a", Cards: Hand{}},
			"b": {ID: "b", Cards: Hand{}},
		},
	}
	if !g.roundComplete() {
		t2.Errorf("expected round complete when hands are empty")
	}

	g.Players["a"].Cards = Hand{{Suit: Spade, Rank: Two}}
	if g.roundComplete() {
		t2.Errorf("expected round not complete when a hand still has cards")
	}
}
