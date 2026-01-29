package game

// Communications

import (
	"encoding/json"

	t "github.com/B33Boy/Judgement/internal/types"
)

func (g *Game) sendGameStarted() {
	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env:     t.Envelope{Type: t.MsgGameStarted},
	})
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

func (g *Game) sendGameFinished() {
	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env: t.Envelope{
			Type: t.MsgGameEnd,
		},
	})
}
