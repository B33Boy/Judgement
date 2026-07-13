package game

// Rules (bidding, playing)

import (
	"encoding/json"
	"errors"
	"fmt"
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

	var payload MakeBid
	if err := json.Unmarshal(input.Env.Payload, &payload); err != nil {
		log.Println("Cannot unmarshall MakeBid")
		return
	}

	if err := g.validateBid(curPlayer, payload.Bid); err != nil {
		g.sendInvalidMove(curPlayer.ID, err.Error())
		return
	}

	curPlayer.Bid = &payload.Bid
	g.state.Bids[curPlayer.ID] = payload.Bid

	g.state.TurnPlayer = g.cyclePlayer()

	if g.cycler.CompletedCycle() {
		g.changeState(BiddingDone)
	}

	g.broadcastGameState()
}

// validateBid enforces the bid range and the "hook" rule: the last bidder
// in a round cannot bid a number that makes the total of all bids equal
// the number of cards dealt this round.
func (g *Game) validateBid(curPlayer *GamePlayer, bid Bid) error {
	cardsInHand := len(curPlayer.Cards)

	if bid < 0 || int(bid) > cardsInHand {
		return fmt.Errorf("bid must be between 0 and %d", cardsInHand)
	}

	if g.cycler.WillCompleteNext() {
		sum := Bid(0)
		for _, b := range g.state.Bids {
			sum += b
		}
		if int(sum+bid) == cardsInHand {
			return fmt.Errorf("bid cannot make total bids equal %d", cardsInHand)
		}
	}

	return nil
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
		g.sendInvalidMove(input.Player.ID, "Not your turn")
		log.Printf("HandlePlay: %v", err)
		return
	}

	// For rounds where we start of with no trump suit
	g.handleNoTrumpSuit(playedCard.Suit)

	// check if card is playable
	if !g.isCardPlayable(curPlayer, playedCard) {
		g.sendInvalidMove(input.Player.ID, "Card cannot be played")
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
	g.broadcastGameState()
}

func (g *Game) handleResolution() {
	winnerID := g.determineTrickWinner()
	g.state.HandsWon[winnerID]++
	g.state.TurnPlayer = winnerID

	if !g.roundComplete() {
		// More tricks left this round - lead the next one from the winner.
		// onStateChanged's StatePlay case clears the table/cardstack and
		// restarts the cycler from g.state.TurnPlayer.
		g.changeState(TrickContinue)
		g.broadcastGameState()
		return
	}

	g.finishRound()
}

func (g *Game) verifyPlayerTurn(player *GamePlayer) error {
	if player.ID != g.state.TurnPlayer {
		log.Printf("It is %s's turn!\n", g.Players[g.state.TurnPlayer].PlayerName)
		return errors.New("Incorrect player turn")
	}
	return nil
}

// isCardPlayable enforces follow-suit: a player must play the suit that
// was led if they hold one. Trump has no bearing on legality - it only
// matters for winning the trick - so it's only a legal substitute once the
// player has no card of the led suit at all.
func (g *Game) isCardPlayable(player *GamePlayer, card Card) bool {

	// Leading the trick - any card is playable.
	if len(g.cardstack) == 0 {
		return true
	}
	ledSuit := g.cardstack[0].Suit

	hasLedSuit := false
	for _, playerCard := range player.Cards {
		if playerCard.Suit == ledSuit {
			hasLedSuit = true
			break
		}
	}

	if hasLedSuit {
		return card.Suit == ledSuit
	}

	// No card of the led suit - anything (including trump) is legal.
	return true
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
