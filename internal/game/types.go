package game

import t "github.com/B33Boy/Judgement/internal/types"

// ================= Game Logic =================
type PlayerMap map[t.PlayerID]*GamePlayer

type GamePlayer struct {
	ID         t.PlayerID
	PlayerName string
	Bid        *Bid
	Cards      Hand
}

// ================= Game Types =================
// Game instance
type Round int
type Score int
type Bid int

// State Machine
type State string
type Event string

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
