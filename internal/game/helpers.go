package game

import (
	"log"

	t "github.com/B33Boy/Judgement/internal/types"
)

func (g *Game) cyclePlayer() t.PlayerID {
	startPlayerID, err := g.cycler.Next()
	if err != nil {
		log.Println("Cannot fetch current player, 0 players in session!")
		g.cancel()
	}
	return startPlayerID
}

func (g *Game) allPlayerIDs() []t.PlayerID {

	all_ids := make([]t.PlayerID, 0, len(g.Players))

	for id := range g.Players {
		all_ids = append(all_ids, id)
	}

	return all_ids
}
