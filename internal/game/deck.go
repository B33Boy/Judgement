package game

import "math/rand"

type Deck []string

func newDeck() Deck {
	cards := make(Deck, 52)
	cardIdx := 0

	for _, suit := range []string{"SPADE", "HEART", "DIAMOND", "CLUB"} {
		for _, rank := range []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING", "ACE"} {
			cards[cardIdx] = suit + "-" + rank
			cardIdx++
		}
	}
	return cards
}

func shuffleDeck(cards Deck) {
	for i := range cards {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func distributeCards(deck Deck, playerCnt int) []Hand {
	playerHands := make([]Hand, playerCnt)
	cardsPerPlayer := 7

	for i := range playerCnt {
		start := i * cardsPerPlayer
		end := start + cardsPerPlayer

		playerHands[i] = Hand(append(Deck(nil), deck[start:end]...))
	}
	return playerHands
}
