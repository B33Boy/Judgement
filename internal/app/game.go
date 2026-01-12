package app

import (
	"context"
	"encoding/json"
	"math/rand"
)

// ======================== Deck ========================
type Deck []string
type Hand []string

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

func distributeCards(deck Deck, playerCnt int) []Hand {
	playerHands := make([]Hand, playerCnt)
	cardsPerPlayer := len(deck) / playerCnt

	for i := range playerCnt {
		start := i * cardsPerPlayer
		end := start + cardsPerPlayer

		// last player gets remaining cards
		if i == playerCnt-1 {
			end = len(deck)
		}

		playerHands[i] = Hand(append(Deck(nil), deck[start:end]...))
	}
	return playerHands
}

// ======================== Game ========================

type GamePlayer struct {
	ID         PlayerID
	PlayerName string
	Score      int
	Cards      Hand
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
	ctx     context.Context
	cancel  context.CancelFunc
	emit    func(GameOutput)
	Players map[PlayerID]*GamePlayer
	Round   int
}

func NewGame(session *Session) *Game {
	session.mu.Lock()
	defer session.mu.Unlock()

	playerCnt := len(session.Players)
	hands := getHands(playerCnt)

	gamePlayers := make(map[PlayerID]*GamePlayer)

	i := 0
	for playerID, player := range session.Players {
		gamePlayers[playerID] = &GamePlayer{
			ID:         playerID,
			PlayerName: player.PlayerName,
			Score:      0,
			Cards:      hands[i],
		}
		i++
	}

	ctx, cancel := context.WithCancel(session.ctx)
	return &Game{
		ctx:    ctx,
		cancel: cancel,
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

func getHands(playerCount int) []Hand {
	deck := newDeck()
	shuffleDeck(deck)
	return distributeCards(deck, playerCount)
}

func (g *Game) Start() {
	g.emit(GameOutput{
		Players: g.allPlayerIDs(),
		Env:     Envelope{Type: MsgGameStarted},
	})

	for _, id := range g.allPlayerIDs() {
		g.sendCardsToPlayer(id)
	}
}

func (g *Game) HandleGameInput(input GameInput) {

}

func (g *Game) allPlayerIDs() []PlayerID {

	all_ids := make([]PlayerID, 0, len(g.Players))

	for id := range g.Players {
		all_ids = append(all_ids, id)
	}

	return all_ids
}

func (g *Game) sendCardsToPlayer(playerID PlayerID) {
	player := g.Players[playerID]
	payload, _ := json.Marshal(PlayerHandChange{
		Cards: player.Cards,
	})

	out := GameOutput{
		Players: []PlayerID{player.ID},
		Env: Envelope{
			Type:    MsgPlayerHand,
			Payload: payload,
		},
	}

	g.emit(out)
}
