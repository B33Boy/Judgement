package game

// Logic flow based on generic State Machine

import "log"

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

	case StateGameOver:
		log.Println("StateGameOver")
		g.sendGameFinished()
	}
}
