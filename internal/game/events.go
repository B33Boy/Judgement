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

	strHand := getStrHand(player.Cards)

	payload, _ := json.Marshal(struct {
		Cards []string `json:"cards"`
	}{Cards: strHand})

	out := t.GameOutput{
		Players: []t.PlayerID{player.ID},
		Env: t.Envelope{
			Type:    t.MsgPlayerHand,
			Payload: payload,
		},
	}

	g.emit(out)
}

func (g *Game) broadcastGameState() {

	payload, _ := json.Marshal(g.state)

	g.emit(t.GameOutput{
		Players: g.allPlayerIDs(),
		Env: t.Envelope{
			Type:    t.MsgStateSync,
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

func (g *Game) sendInvalidMove(playerID t.PlayerID, message string) {
	payload, _ := json.Marshal(InvalidActionPayload{
		Message: message,
	})

	out := t.GameOutput{
		Players: []t.PlayerID{playerID},
		Env: t.Envelope{
			Type:    t.MsgInvalidAction,
			Payload: payload,
		},
	}

	g.emit(out)
}
