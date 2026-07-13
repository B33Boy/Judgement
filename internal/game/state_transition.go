package game

// Logic flow based on generic State Machine

import (
	"log"

	t "github.com/B33Boy/Judgement/internal/types"
)

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
		// Clear state from previous trick/round
		g.cardstack = make([]Card, 0)
		g.state.Table = make(map[t.PlayerID]*Card)
		g.cycler.StartFrom(g.state.TurnPlayer)
		log.Println("StatePlay")

	case StateResolution:
		log.Println("StateResolution")
		// Resolution has no client-driven input - run it immediately as
		// part of the transition rather than waiting for a message that
		// will never arrive.
		g.handleResolution()

	case StateGameOver:
		log.Println("StateGameOver")
		g.sendGameFinished()
	}
}
