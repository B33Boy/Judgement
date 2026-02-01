package game

import (
	"math/rand"
)

// ================================== Card ==================================
type Suit int
type Rank int

const (
	Spade Suit = iota
	Heart
	Diamond
	Club
)

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

func (s Suit) String() string {
	return [...]string{"SPADE", "HEART", "DIAMOND", "CLUB"}[s]
}

func (r Rank) String() string {
	return [...]string{
		"", "", "2", "3", "4", "5", "6", "7",
		"8", "9", "10", "JACK", "QUEEN", "KING", "ACE",
	}[r]
}

func (c Card) String() string {
	return c.Suit.String() + "-" + c.Rank.String()
}

func higherRank(a, b Card) bool {
	return a.Rank > b.Rank
}

func sameSuit(a, b Card) bool {
	return a.Suit == b.Suit
}

func (card *Card) greater_than(other *Card) bool {
	return sameSuit(*card, *other) && higherRank(*card, *other)
}

// ================================== Card ==================================

type Deck []Card
type Hand []Card

func getStrHand(hand Hand) []string {
	strCards := make([]string, len(hand))
	for i, card := range hand {
		strCards[i] = card.String()
	}
	return strCards
}

func newDeck() Deck {
	deck := make(Deck, 0, 52)

	for s := Spade; s <= Club; s++ {
		for r := Two; r <= Ace; r++ {
			deck = append(deck, Card{
				Suit: s,
				Rank: r,
			})
		}
	}
	return deck
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

	for i := 0; i < playerCnt; i++ {
		start := i * cardsPerPlayer
		end := start + cardsPerPlayer

		// copy the slice so each hand has its own backing array
		playerHands[i] = append(Hand(nil), deck[start:end]...)
	}
	return playerHands
}

func getHands(playerCount int) []Hand {
	deck := newDeck()
	shuffleDeck(deck)
	return distributeCards(deck, playerCount)
}
