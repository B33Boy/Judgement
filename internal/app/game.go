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
	ID         PlayerID
	PlayerName string
	Score      int
	Cards      Deck
}

type GameInput struct {
	Player *Player
	Env    Envelope
}

type GameOutput struct {
	Players []PlayerID
	Env     Envelope
}

type Game struct {
	emit    func(GameOutput)
	Players map[PlayerID]*GamePlayer
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

	gamePlayers := make(map[PlayerID]*GamePlayer)

	i := 0
	for playerID, player := range session.Players {
		gamePlayers[playerID] = &GamePlayer{
			ID:         playerID,
			PlayerName: player.PlayerName,
			Score:      0,
			Cards:      cards[i],
		}
		i++
	}

	return &Game{
		emit: func(out GameOutput) {
			select {
			case session.outputs <- out:
			case <-session.ctx.Done():
			}
		},
		Players: gamePlayers,
		Round:   0,
	}
}

func (g *Game) Start() {
	out := GameOutput{
		Players: g.allPlayerIDs(),
		Env:     Envelope{Type: MsgGameStarted},
	}
	g.emit(out)
}

func (g *Game) HandleGameInput(input GameInput) {
	// Take input and parse it to an appropriate message
	// perform message action
	// But the message output pipe is contained within the session and the game shouldn't have access to the session

}

func (g *Game) allPlayerIDs() []PlayerID {

	all_ids := make([]PlayerID, 0, len(g.Players))

	for id := range g.Players {
		all_ids = append(all_ids, id)
	}

	return all_ids
}
