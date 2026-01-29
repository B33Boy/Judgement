package game

import (
	"context"
	"encoding/json"
	"errors"
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
	}
}

func getHands(playerCount int) []Hand {
	deck := newDeck()
	shuffleDeck(deck)
	return distributeCards(deck, playerCount)
}
func NewScoreboard(playerCnt int, gamePlayers PlayerMap, maxRounds Round) PlayerScore {
	scores := make(PlayerScore, playerCnt)
	for playerId := range gamePlayers {
		scores[playerId] = make([]Score, maxRounds)
	}
	return scores
}

func (g *Game) Start() {

	g.sendGameStarted()

	for _, id := range g.allPlayerIDs() {
		g.sendCardsToPlayer(id)
	}

	// Send Round # and Send Turn PlayerId
	g.sendRoundInfo()
}

func (g *Game) sendGameStarted() {
	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env:     t.Envelope{Type: t.MsgGameStarted},
	})
}

func (g *Game) allPlayerIDs() []t.PlayerID {

	all_ids := make([]t.PlayerID, 0, len(g.Players))

	for id := range g.Players {
		all_ids = append(all_ids, id)
	}

	return all_ids
}

func (g *Game) sendCardsToPlayer(playerID t.PlayerID) {
	player := g.Players[playerID]
	payload, _ := json.Marshal(PlayerHandChangePayload{
		Cards: player.Cards,
	})

	out := t.GameOutput{
		Players: []t.PlayerID{player.ID},
		Env: t.Envelope{
			Type:    t.MsgPlayerHand,
			Payload: payload,
		},
	}

	g.emit(out)
}

func (g *Game) cyclePlayer() t.PlayerID {
	startPlayerID, err := g.cycler.Next()
	if err != nil {
		log.Println("Cannot fetch current player, 0 players in session!")
		g.cancel()
	}
	return startPlayerID
}

func (g *Game) sendRoundInfo() {
	payload, _ := json.Marshal(RoundInfoPayload{
		Round:      g.params.round,
		TurnPlayer: g.Players[g.turnPlayer].PlayerName,
		State:      g.sm.state,
	})

	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env: t.Envelope{
			Type:    t.MsgRoundInfo,
			Payload: payload,
		},
	})
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

func (g *Game) handleBid(input t.GameInput) {
	if input.Env.Type != t.MsgMakeBid {
		log.Println("Invalid message type, \"make_bid\" expected")
		return
	}

	if input.Player.ID != g.turnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.turnPlayer].PlayerName)
		return
	}

	curPlayer := g.Players[input.Player.ID]
	g.recordBid(curPlayer, input)

	g.turnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.sm.Trigger(BiddingDone)
	}

	g.sendRoundInfo()
}

func (g *Game) verifyPlayerTurn(player *GamePlayer) error {
	if player.ID != g.turnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.turnPlayer].PlayerName)
		return errors.New("Incorrect player turn")
	}
	return nil
}

func (g *Game) recordBid(curPlayer *GamePlayer, input t.GameInput) {

	var payload MakeBid
	err := json.Unmarshal(input.Env.Payload, &payload)

	if err != nil {
		log.Println("Cannot unmarshall MakeBid")
		return
	}

	// Logic to check if bid is possible
	// If not possible send message back

	curPlayer.Bid = &payload.Bid
}

func (g *Game) handlePlay(input t.GameInput) {
	// receive played card, send rejection message if not possible to play
	// get a struct of
	// 1) play card
	if input.Env.Type != t.MsgPlayCard {
		log.Println("Invalid message type, \"make_bid\" expected")
		return
	}

	curPlayer := g.Players[input.Player.ID]
	if g.verifyPlayerTurn(curPlayer) != nil {
		return
	}

	// check if card is valid
	// Play card

	g.turnPlayer = g.cyclePlayer()

	// if g.allPlayersPlayedCard() {
	// 	g.sm.Trigger(PlayingDone)
	// 	g.updateRound()
	// }
	g.sendRoundInfo()
}

func (g *Game) handleResolution(input t.GameInput) {
	// Update scores and send to frontend
	// call StartFrom() to start from the winning player
	// call UpdateRound() here
	g.updateRound()
}

func (g *Game) updateRound() {
	// Run this after every player turn, it will only update round when we return back to first player
	completed := g.cycler.CompletedCycle()
	if !completed {
		log.Println("Player cycle not completed, not updating round yet")
		return
	}
	g.params.round++

	if g.params.round > g.params.maxRounds {
		// state change to finished game
		g.trigger(PlayingDone)
	}
}

func (g *Game) sendGameFinished() {
	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env: t.Envelope{
			Type: t.MsgGameEnd,
		},
	})
}

func (g *Game) trigger(event Event) {
	prev := g.sm.state
	next, err := g.sm.Trigger(event)

	if err != nil {
		log.Printf("trigger state change failed: %v\n", err)
		return
	}

	g.onStateChanged(prev, next)
}

func (g *Game) onStateChanged(from, to State) {
	switch to {

	case StateBid:
		log.Println("StateBid")

	case StatePlay:
		g.cycler.StartFrom(g.turnPlayer)
		log.Println("StatePlay")

	case StateResolution:
		log.Println("StateResolution")

	case StateGameOver:
		log.Println("StateGameOver")
		g.sendGameFinished()
	}
}
