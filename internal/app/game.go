package app

import (
	"encoding/json"
	"math/rand"
)

// ======================== Deck ========================
type Deck []string

func NewDeck() Deck {
	cards := make(Deck, 52)
	cardIdx := 0

	for _, suit := range []string{"♠", "♥", "♦", "♣"} {
		for _, rank := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"} {
			cards[cardIdx] = rank + "_" + suit
			cardIdx++
		}
	}
	return cards
}

func ShuffleDeck(cards Deck) {
	for i := range cards {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func distributeCards(deck Deck, playerCnt int) []Deck {
	playerHands := make([]Deck, playerCnt)
	cardsPerPlayer := len(deck) / playerCnt

	for i := range playerCnt {
		start := i * cardsPerPlayer
		end := start + cardsPerPlayer

		// last player gets remaining cards
		if i == playerCnt-1 {
			end = len(deck)
		}

		playerHands[i] = append(Deck(nil), deck[start:end]...)
	}
	return playerHands
}

// ======================== Game ========================
type GameInput struct {
	Player  string
	Type    MessageType
	Payload json.RawMessage
}

type GameEvent struct {
	Type    MessageType
	Payload any
}

type GamePlayer struct {
	PlayerName string
	Score      int
	Cards      Deck
}

type Game struct {
	Players []*GamePlayer
	Turn    int

	inputs chan GameInput
	events chan GameEvent
	// state  GameState
}
