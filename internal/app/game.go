package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
)

// ======================== Deck ========================
type Deck []string
type Hand []string
type PlayerMap map[PlayerID]*GamePlayer

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
	cardsPerPlayer := 7

	for i := range playerCnt {
		start := i * cardsPerPlayer
		end := start + cardsPerPlayer

		playerHands[i] = Hand(append(Deck(nil), deck[start:end]...))
	}
	return playerHands
}

// ======================== Game Structs ========================

type GamePlayer struct {
	ID         PlayerID
	PlayerName string
	Bid        *bid
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

// ======================== State Machine ========================

type State string
type Event string

type StateMachine struct {
	state       State
	transitions map[State]map[Event]State
}

func NewStateMachine(initial State) *StateMachine {
	return &StateMachine{
		state:       initial,
		transitions: make(map[State]map[Event]State),
	}
}

func (sm *StateMachine) AddTransition(from State, event Event, to State) {
	if sm.transitions[from] == nil {
		sm.transitions[from] = make(map[Event]State)
	}
	sm.transitions[from][event] = to
}

func (sm *StateMachine) Trigger(event Event) error {
	next, ok := sm.transitions[sm.state][event]

	if !ok {
		return fmt.Errorf("invalid transition: %s + %s", sm.state, event)
	}

	sm.state = next
	return nil
}

const (
	// States
	StateBid        State = "bidding"
	StatePlay       State = "playing"
	StateResolution State = "resolution"
	StateGameOver   State = "gameover"

	// Events
	BiddingDone     Event = "bidding_done"
	PlayingContinue Event = "playing_continue"
	PlayingDone     Event = "playing_done"
	RoundResolved   Event = "round_resolved"
)

// ======================== Scoreboard ========================

type round int
type score int
type bid int

type PlayerScore map[PlayerID][]score

// type PlayerBid map[PlayerID][]bid

// ======================== Game ========================

type GameParams struct {
	round         round
	maxRounds     round
	cardsPerRound int
}

type Game struct {
	ctx        context.Context
	cancel     context.CancelFunc
	emit       func(GameOutput)
	Players    PlayerMap
	turnPlayer PlayerID
	params     *GameParams
	cycler     *PlayerCycler
	sm         *StateMachine
	scores     PlayerScore
}

type PlayerCycler struct {
	keys  []PlayerID
	index int
}

func NewPlayerCycler(m PlayerMap) *PlayerCycler {
	keys := make([]PlayerID, 0, len(m))

	for playerID := range m {
		keys = append(keys, playerID)
	}
	return &PlayerCycler{
		keys:  keys,
		index: 0,
	}
}

func (pc *PlayerCycler) Next() (PlayerID, error) {
	if len(pc.keys) == 0 {
		return PlayerID(""), errors.New("0 players to cycle through")
	}

	playerID := pc.keys[pc.index]
	pc.index = (pc.index + 1) % len(pc.keys)

	return playerID, nil
}

func (pc *PlayerCycler) completedCycle() bool {
	if len(pc.keys) == 0 {
		return false
	}
	return pc.index == 0
}

func NewGame(session *Session) *Game {
	session.mu.Lock()
	defer session.mu.Unlock()

	playerCnt := len(session.Players)
	hands := getHands(playerCnt)

	// Game Players
	gamePlayers := make(PlayerMap)

	i := 0
	for playerID, player := range session.Players {
		gamePlayers[playerID] = &GamePlayer{
			ID:         playerID,
			PlayerName: player.PlayerName,
			Bid:        nil,
			Cards:      hands[i],
		}
		i++
	}

	// Params
	params := &GameParams{
		round:         0,
		maxRounds:     14,
		cardsPerRound: 7,
	}

	// Cycler
	cycler := NewPlayerCycler(gamePlayers)
	firstPlayerID, _ := cycler.Next()

	// State Machine
	sm := NewStateMachine(StateBid)

	sm.AddTransition(StateBid, BiddingDone, StatePlay)
	sm.AddTransition(StatePlay, PlayingDone, StateResolution)
	sm.AddTransition(StateResolution, PlayingContinue, StateBid)
	sm.AddTransition(StateResolution, PlayingDone, StateGameOver)

	// Scores
	scores := NewScoreboard(playerCnt, gamePlayers, params.maxRounds)

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
func NewScoreboard(playerCnt int, gamePlayers PlayerMap, maxRounds round) PlayerScore {
	scores := make(PlayerScore, playerCnt)
	for playerId := range gamePlayers {
		scores[playerId] = make([]score, maxRounds)
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
	g.emit(GameOutput{
		Players: g.allPlayerIDs(),
		Env:     Envelope{Type: MsgGameStarted},
	})
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
	payload, _ := json.Marshal(PlayerHandChangePayload{
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

func (g *Game) cyclePlayer() PlayerID {
	startPlayerID, err := g.cycler.Next()
	if err == nil {
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

	g.emit(GameOutput{
		Players: g.allPlayerIDs(),
		Env: Envelope{
			Type:    MsgRoundInfo,
			Payload: payload,
		},
	})
}

func (g *Game) HandleGameInput(input GameInput) {
	switch g.sm.state {
	case StateBid:
		g.handleBid(input)

	case StatePlay:
		g.handlePlay(input)

	case StateResolution:
		g.handleResolution(input)
	}
}

func (g *Game) handleBid(input GameInput) {
	if input.Env.Type != MsgMakeBid {
		log.Println("Invalid message type, \"make_bid\" expected")
		return
	}

	curPlayer := g.Players[input.Player.ID]
	if curPlayer.ID != g.turnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.turnPlayer].PlayerName)
		return
	}

	g.recordBid(curPlayer, input)

	g.turnPlayer = g.cyclePlayer()

	if g.allPlayersBid() {
		g.sm.Trigger(BiddingDone)
		g.updateRound()
	}
	g.sendRoundInfo()
}

func (g *Game) recordBid(curPlayer *GamePlayer, input GameInput) {

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

func (g *Game) allPlayersBid() bool {
	for _, player := range g.Players {
		if player.Bid == nil {
			return false
		}
	}
	return true
}

func (g *Game) handlePlay(input GameInput) {
	// player actions
	// 1) play card 2) reveal sir card
}

func (g *Game) handleResolution(input GameInput) {
	// Update scores and send to frontend
}

func (g *Game) updateRound() {
	// Run this after every player turn, it will only update round when we return back to first player
	completed := g.cycler.completedCycle()
	if !completed {
		return
	}
	g.params.round++

	if g.params.round > g.params.maxRounds {
		// state change to finished game
		g.sendGameFinished()
	}
}

func (g *Game) sendGameFinished() {
	g.emit(GameOutput{
		Players: g.allPlayerIDs(),
		Env: Envelope{
			Type: MsgGameEnd,
		},
	})
}
