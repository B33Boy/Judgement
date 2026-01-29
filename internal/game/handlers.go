package game

// Rules (bidding, playing)

import (
	"encoding/json"
	"errors"
	"log"

	t "github.com/B33Boy/Judgement/internal/types"
)

func (g *Game) handleBid(input t.GameInput) {
	if input.Env.Type != t.MsgMakeBid {
		log.Println("Invalid message type, \"make_bid\" expected")
		return
	}

	curPlayer := g.Players[input.Player.ID]
	if g.verifyPlayerTurn(curPlayer) != nil {
		return
	}

	g.recordBid(curPlayer, input)

	g.turnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.sm.Trigger(BiddingDone)
	}

	g.sendRoundInfo()
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

func (g *Game) verifyPlayerTurn(player *GamePlayer) error {
	if player.ID != g.turnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.turnPlayer].PlayerName)
		return errors.New("Incorrect player turn")
	}
	return nil
}
