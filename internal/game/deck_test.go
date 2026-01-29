package game

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := newDeck()

	// Check total count
	if len(deck) != 52 {
		t.Errorf("Expected 52 cards, got %d", len(deck))
	}

	expectedCard := "SPADE-ACE"
	found := false
	for _, card := range deck {
		if card == expectedCard {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find %s in new deck", expectedCard)
	}
}

func TestShuffleDeck(t *testing.T) {
	deck1 := newDeck()

	// Create a copy to compare against
	deck2 := make(Deck, len(deck1))
	copy(deck2, deck1)

	shuffleDeck(deck1)

	same := true
	for i := range deck1 {
		if deck1[i] != deck2[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("Deck order did not change after shuffle")
	}
}

func TestDistributeCards(t *testing.T) {
	deck := newDeck()
	playerCount := 3
	expectedCardsPerPlayer := 7

	hands := distributeCards(deck, playerCount)

	if len(hands) != playerCount {
		t.Errorf("Expected %d hands, got %d", playerCount, len(hands))
	}

	for i, hand := range hands {
		if len(hand) != expectedCardsPerPlayer {
			t.Errorf("Player %d expected %d cards, got %d", i, expectedCardsPerPlayer, len(hand))
		}
	}
}
