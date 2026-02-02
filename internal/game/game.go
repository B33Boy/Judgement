package game

// Main game

import (
	"context"
	"log"
	"math/rand"
	"time"

	t "github.com/B33Boy/Judgement/internal/types"
)

// ======================== Game ========================

type GameParams struct {
	maxRounds     Round
	cardsPerRound int
}

type GameState struct {
	Round      Round                `json:"round"`
	State      State                `json:"state"`
	TurnPlayer t.PlayerID           `json:"turnPlayer"`
	TrumpSuit  *Suit                `json:"trumpSuit"`
	Table      map[t.PlayerID]*Card `json:"table"` // Cards currently played
	Bids       map[t.PlayerID]Bid   `json:"bids"`
	HandsWon   map[t.PlayerID]int   `json:"handsWon"`
}

type Game struct {
	// Engine
	ctx    context.Context
	cancel context.CancelFunc
	emit   func(t.GameOutput)
	cycler *PlayerCycler
	sm     *StateMachine

	// Data
	Players   PlayerMap
	params    *GameParams
	state     *GameState
	scores    PlayerScore // historical scores
	cardstack []*Card
}

type SessionView interface {
	Context() context.Context
	GetPlayers() map[t.PlayerID]*t.Player
	Emit(t.GameOutput)
}

func NewGame(session SessionView) *Game {
	players := session.GetPlayers()
	playerCnt := len(players)

	hands := getHands(playerCnt)
	gamePlayers := make(PlayerMap)

	i := 0
	for playerID, player := range players {
		gamePlayers[playerID] = &GamePlayer{
			ID:         playerID,
			PlayerName: player.PlayerName,
			Bid:        nil,
			Cards:      hands[i],
		}
		i++
	}

	ctx, cancel := context.WithCancel(session.Context())

	// Cycler
	cycler := NewPlayerCycler(gamePlayers)

	keys := make([]t.PlayerID, 0, len(gamePlayers))
	for id := range gamePlayers {
		keys = append(keys, id)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	firstPlayerID := keys[rng.Intn(len(keys))]

	err := cycler.StartFrom(firstPlayerID)
	if err != nil {
		log.Println("failed to start cycler:", err)
	}

	// State Machine
	sm := NewStateMachine(StateBid)
	sm.AddTransition(StateBid, BiddingDone, StatePlay)
	sm.AddTransition(StatePlay, PlayingDone, StateResolution)
	sm.AddTransition(StateResolution, PlayingContinue, StateBid)
	sm.AddTransition(StateResolution, PlayingDone, StateGameOver)

	// Params
	params := &GameParams{
		maxRounds:     14,
		cardsPerRound: 7,
	}

	gameState := &GameState{
		Round:      0,
		State:      StateBid,
		TurnPlayer: firstPlayerID,
		TrumpSuit:  nil,
		Table:      make(map[t.PlayerID]*Card),
		Bids:       make(map[t.PlayerID]Bid),
		HandsWon:   make(map[t.PlayerID]int),
	}

	// Scores
	scoreboard := NewScoreboard(playerCnt, gamePlayers, params.maxRounds)

	return &Game{
		ctx:    ctx,
		cancel: cancel,
		emit:   session.Emit,
		cycler: cycler,
		sm:     sm,

		Players:   gamePlayers,
		params:    params,
		state:     gameState,
		scores:    scoreboard,
		cardstack: make([]*Card, 0),
	}
}

func (g *Game) Start() {

	g.sendGameStarted()

	for _, id := range g.allPlayerIDs() {
		g.sendCardsToPlayer(id)
	}

	// Send Round # and Send Turn PlayerId
	g.sendRoundInfo()
}

func (g *Game) HandleGameInput(input t.GameInput) {
	switch g.sm.state {
	case StateBid:
		g.handleBid(input)

	case StatePlay:
		g.handlePlay(input)

	case StateResolution:
		g.handleResolution(input)
	}
}
