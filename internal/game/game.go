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
	round         Round
	maxRounds     Round
	cardsPerRound int
}

type Game struct {
	ctx        context.Context
	cancel     context.CancelFunc
	emit       func(t.GameOutput)
	Players    PlayerMap
	turnPlayer t.PlayerID
	params     *GameParams
	cycler     *PlayerCycler
	sm         *StateMachine
	scores     PlayerScore
	cardstack  []*Card
	sir        *Suit
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

	// Params
	params := &GameParams{
		round:         0,
		maxRounds:     14,
		cardsPerRound: 7,
	}

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

	// Scores
	scores := NewScoreboard(playerCnt, gamePlayers, params.maxRounds)

	return &Game{
		ctx:        ctx,
		cancel:     cancel,
		emit:       session.Emit,
		Players:    gamePlayers,
		turnPlayer: firstPlayerID,
		params:     params,
		cycler:     cycler,
		sm:         sm,
		scores:     scores,
		cardstack:  make([]*Card, 0),
		sir:        nil,
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
