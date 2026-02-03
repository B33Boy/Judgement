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
	if err := g.verifyPlayerTurn(curPlayer); err != nil {
		log.Printf("HandleBid: %v", err)
		return
	}

	g.recordBid(curPlayer, input)

	g.state.TurnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.changeState(BiddingDone)
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
	if err := g.verifyPlayerTurn(curPlayer); err != nil {
		// g.sendInvalidMove(input.Player.ID, "Not your turn")
		log.Printf("HandlePlay: %v", err)
		return
	}

	// For rounds where we start of with no trump suit
	g.handleNoTrumpSuit(playedCard.Suit)

	// // check if card is playable
	if !g.isCardPlayable(curPlayer, playedCard) {
		// g.sendInvalidMove(input.Player.ID, "Card cannot be played")
		log.Printf("Card not playable: %v", playedCard)
		return
	}

	// Play card
	g.playCard(curPlayer, playedCard)
	g.sendCardsToPlayer(curPlayer)

	g.state.TurnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.changeState(PlayingDone)
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
		g.changeState(GameDone)
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
	if len(g.cardstack) == 0 {
		return true
	}
	curTop := g.cardstack[len(g.cardstack)-1]

	trump := g.state.TrumpSuit
	hasTrump := trump != nil

	hasLegalAlternative := false

	for _, playerCard := range player.Cards {
		if sameSuit(playerCard, curTop) || (hasTrump && playerCard.Suit == *trump) {

			hasLegalAlternative = true

			// If the played card is one of them, allow
			if playerCard.Equals(card) {
				log.Println("Card follows suit or trump")
				return true
			}
		}
	}

	// If no legal alternatives exist, player can play anything
	if !hasLegalAlternative {
		log.Println("No matching suit or trump, any card allowed")
		return true
	}

	return false // had legal option but player didn't use it
}

func (g *Game) playCard(player *GamePlayer, card Card) {
	g.removeCardFromPlayer(player, card)
	g.addCardToTable(player, card)
}

func (g *Game) removeCardFromPlayer(player *GamePlayer, playedCard Card) {
	for i, handCard := range player.Cards {
		if handCard.Equals(playedCard) {
			player.Cards = append(player.Cards[:i], player.Cards[i+1:]...)
			return
		}
	}
}

func (g *Game) addCardToTable(player *GamePlayer, card Card) {
	g.cardstack = append(g.cardstack, card)
	g.state.Table[player.ID] = &card
}

func (g *Game) handleNoTrumpSuit(suit Suit) {
	trump := g.state.TrumpSuit
	if trump == nil {
		// Make current card (initial) the trump suit
		g.state.TrumpSuit = &suit
	}
}
