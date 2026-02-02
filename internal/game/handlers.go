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

	g.state.TurnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.sm.Trigger(BiddingDone)
	}

	g.broadcastGameState()
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
	if input.Env.Type != t.MsgPlayCard {
		log.Println("Invalid message type, \"make_bid\" expected")
		return
	}

	// Get card from input
	var playedCard Card
	err := json.Unmarshal(input.Env.Payload, &playedCard)
	if err != nil {
		log.Println("Cannot unmarshal played card")
		return
	}

	log.Printf("%v", playedCard.String())

	// Get player from input and ensure that it is their turn
	curPlayer := g.Players[input.Player.ID]
	if g.verifyPlayerTurn(curPlayer) != nil {
		g.sendInvalidMove(input.Player.ID, "Not your turn")
		return
	}

	// // check if card is playable
	if !g.isCardPlayable(curPlayer, playedCard) {
		g.sendInvalidMove(input.Player.ID, "Card cannot be played")
		return
	}
	// Play card
	// g.playCard(input.Player.ID, input.Card)

	g.state.TurnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.sm.Trigger(PlayingDone)
		g.updateRound()
	}
	// g.broadcastCardPlayed(input.Player.ID, input.Card)
	g.broadcastGameState()
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
	g.state.Round++

	if g.state.Round > g.params.maxRounds {
		// state change to finished game
		g.trigger(PlayingDone)
	}
}

func (g *Game) verifyPlayerTurn(player *GamePlayer) error {
	if player.ID != g.state.TurnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.state.TurnPlayer].PlayerName)
		return errors.New("Incorrect player turn")
	}
	return nil
}

func (g *Game) isCardPlayable(player *GamePlayer, card Card) bool {

	// If there are no cards on the cardstack, any card is playable
	num_cards := len(g.cardstack)
	if num_cards == 0 {
		return true
	}

	curTop := *g.cardstack[num_cards-1]

	// Check if player has other cards in the hand that they can play
	// You can play if same suite or sir, if you don't have either than you can play whatever card
	cards_allowable := make([]Card, 0)
	for _, playerCard := range player.Cards {
		if sameSuit(playerCard, curTop) || playerCard.Suit == *g.state.TrumpSuit {
			// From the allowed cards, if the current is a
			if playerCard == card {
				log.Println("Card is either the same suit or the sir")
				return true
			}
			cards_allowable = append(cards_allowable, playerCard)
		}
	}

	if len(cards_allowable) == 0 {
		log.Println("Card is a fish")
		return true // fish - any card is allowed
	}

	return false
}
