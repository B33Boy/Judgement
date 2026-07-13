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

// roundSchedule defines how many cards are dealt each round: a pyramid
// that shrinks from 7 down to 1 and back up to 7 (14 rounds total).
var roundSchedule = []int{7, 6, 5, 4, 3, 2, 1, 1, 2, 3, 4, 5, 6, 7}

type GameParams struct {
	roundSchedule []int
}

type GameState struct {
	Round      Round                `json:"round"`
	State      State                `json:"state"`
	TurnPlayer t.PlayerID           `json:"turnPlayer"`
	TrumpSuit  *Suit                `json:"trumpSuit"`
	Table      map[t.PlayerID]*Card `json:"table"` // Cards currently played
	Bids       map[t.PlayerID]Bid   `json:"bids"`
	HandsWon   map[t.PlayerID]int   `json:"handsWon"`
	Scores     PlayerScore          `json:"scores"`
}

type Game struct {
	// Engine
	ctx    context.Context
	cancel context.CancelFunc
	emit   func(t.GameOutput)
	cycler *PlayerCycler
	sm     *StateMachine

	// Data
	Players     PlayerMap
	playerOrder []t.PlayerID // fixed seating order, used to rotate the round starter
	dealerIdx   int
	params      *GameParams
	state       *GameState
	scores      PlayerScore // historical scores
	cardstack   []Card
}

type SessionView interface {
	Context() context.Context
	GetPlayers() map[t.PlayerID]*t.Player
	Emit(t.GameOutput)
}

func NewGame(session SessionView) *Game {
	players := session.GetPlayers()
	playerCnt := len(players)

	// Params needs to be created before we use roundSchedule
	params := &GameParams{
		roundSchedule: roundSchedule,
	}

	hands := getHands(playerCnt, params.roundSchedule[0])
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
	firstIdx := rng.Intn(len(keys))
	firstPlayerID := keys[firstIdx]

	err := cycler.StartFrom(firstPlayerID)
	if err != nil {
		log.Println("failed to start cycler:", err)
	}

	// State Machine
	sm := NewStateMachine(StateBid)
	sm.AddTransition(StateBid, BiddingDone, StatePlay)
	sm.AddTransition(StatePlay, PlayingDone, StateResolution)
	sm.AddTransition(StateResolution, TrickContinue, StatePlay)
	sm.AddTransition(StateResolution, PlayingContinue, StateBid)
	sm.AddTransition(StateResolution, GameDone, StateGameOver)

	// Scores
	scoreboard := NewScoreboard(playerCnt, gamePlayers, Round(len(params.roundSchedule)))

	gameState := &GameState{
		Round:      0,
		State:      StateBid,
		TurnPlayer: firstPlayerID,
		TrumpSuit:  nil,
		Table:      make(map[t.PlayerID]*Card),
		Bids:       make(map[t.PlayerID]Bid),
		HandsWon:   make(map[t.PlayerID]int),
		Scores:     scoreboard,
	}

	return &Game{
		ctx:    ctx,
		cancel: cancel,
		emit:   session.Emit,
		cycler: cycler,
		sm:     sm,

		Players:     gamePlayers,
		playerOrder: keys,
		dealerIdx:   firstIdx,
		params:      params,
		state:       gameState,
		scores:      scoreboard,
		cardstack:   make([]Card, 0),
	}
}

func (g *Game) Start() {

	g.sendGameStarted()

	for _, player := range g.Players {
		g.sendCardsToPlayer(player)
	}

	// Send Round # and Send Turn PlayerId
	g.broadcastGameState()
}

func (g *Game) HandleGameInput(input t.GameInput) {
	switch g.sm.state {
	case StateBid:
		g.handleBid(input)

	case StatePlay:
		g.handlePlay(input)

		// StateResolution and StateGameOver are not client-driven - any input
		// received while in those states is simply ignored.
	}
}
