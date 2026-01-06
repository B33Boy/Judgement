package app

import (
	"math/rand"
)

// ======================== Deck ========================
type Deck []string

func newDeck() Deck {
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

func shuffleDeck(cards Deck) {
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

type GamePlayer struct {
	PlayerName string
	Score      int
	Cards      Deck
}

type GameInput struct {
	Player *Player
	Env    Envelope
}

type Game struct {
	emit    func(Envelope)
	Players []*GamePlayer
	Round   int
	// state  GameState
}

func NewGame(session *Session) *Game {
	session.mu.Lock()
	defer session.mu.Unlock()

	deck := newDeck()
	shuffleDeck(deck)

	playerCnt := len(session.Players)
	cards := distributeCards(deck, playerCnt)

	gamePlayers := make([]*GamePlayer, 0, playerCnt)

	i := 0
	for name := range session.Players {
		gamePlayers = append(gamePlayers, &GamePlayer{
			PlayerName: name,
			Score:      0,
			Cards:      cards[i],
		})
		i++
	}

	return &Game{
		emit: func(env Envelope) {
			select {
			case session.outputs <- env:
			case <-session.ctx.Done():
			}
		},
		Players: gamePlayers,
		Round:   0,
	}
}

func (g *Game) Start() {
	// setup game and run state machine
	g.emit(Envelope{Type: MsgGameStarted})

}

func (g *Game) HandleGameInput(input GameInput) {
	// Take input and parse it to an appropriate message
	// perform message action
	// But the message output pipe is contained within the session and the game shouldn't have access to the session

}
