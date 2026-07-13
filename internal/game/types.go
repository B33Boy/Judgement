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

// MinPlayers/MaxPlayers bound how many players a game can seat. MaxPlayers
// is a hard constraint of the round schedule: the largest round deals 7
// cards per player, and 8*7=56 exceeds the 52-card deck.
const (
	MinPlayers = 3
	MaxPlayers = 7
)

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
	TrickContinue   Event = "trick_continue"
	GameDone        Event = "game_done"
	RoundResolved   Event = "round_resolved"
)
