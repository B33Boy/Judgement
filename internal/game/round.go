package game

// Round and trick lifecycle: determining trick winners, scoring completed
// rounds, and dealing/resetting state for the next round.

import (
	t "github.com/B33Boy/Judgement/internal/types"
)

// determineTrickWinner returns the player who won the current trick.
// The led suit is whatever suit was first played this trick (cardstack is
// cleared at the start of every trick and appended to in play order).
func (g *Game) determineTrickWinner() t.PlayerID {
	ledSuit := g.cardstack[0].Suit
	trump := g.state.TrumpSuit

	var winnerID t.PlayerID
	var winningCard Card
	first := true

	for playerID, cardPtr := range g.state.Table {
		card := *cardPtr
		if first {
			winnerID, winningCard, first = playerID, card, false
			continue
		}
		if beats(card, winningCard, ledSuit, trump) {
			winnerID, winningCard = playerID, card
		}
	}

	return winnerID
}

// beats reports whether candidate wins the trick over the current best
// card, given the suit that was led and the round's trump suit (nil if
// none). Trump beats non-trump; among two trumps or two led-suit cards the
// higher rank wins; a card that is neither trump nor led suit never wins.
func beats(candidate, current Card, ledSuit Suit, trump *Suit) bool {
	candidateIsTrump := trump != nil && candidate.Suit == *trump
	currentIsTrump := trump != nil && current.Suit == *trump

	if candidateIsTrump != currentIsTrump {
		return candidateIsTrump
	}
	if candidateIsTrump && currentIsTrump {
		return candidate.Rank > current.Rank
	}

	// Neither card is trump - only led-suit cards are in contention.
	if candidate.Suit != ledSuit {
		return false
	}
	if current.Suit != ledSuit {
		return true
	}
	return candidate.Rank > current.Rank
}

// roundComplete reports whether every player has played all their cards
// for the round. All hands shrink in lockstep, so checking one is enough.
func (g *Game) roundComplete() bool {
	for _, player := range g.Players {
		return len(player.Cards) == 0
	}
	return true
}

// scoreRound records each player's score for the round currently in
// progress (g.state.Round), before it gets incremented by finishRound.
func (g *Game) scoreRound() {
	for playerID, player := range g.Players {
		var score Score
		if player.Bid != nil && int(*player.Bid) == g.state.HandsWon[playerID] {
			score = Score(10 + int(*player.Bid))
		}
		g.scores[playerID][g.state.Round] = score
	}
}

// startNewRound advances the dealer, deals fresh hands sized by the round
// schedule, and resets all per-round state.
func (g *Game) startNewRound() {
	g.dealerIdx = (g.dealerIdx + 1) % len(g.playerOrder)
	starter := g.playerOrder[g.dealerIdx]

	cardsThisRound := g.params.roundSchedule[g.state.Round]
	hands := getHands(len(g.playerOrder), cardsThisRound)

	for i, playerID := range g.playerOrder {
		player := g.Players[playerID]
		player.Cards = hands[i]
		player.Bid = nil
	}

	g.state.TrumpSuit = nil
	g.state.Bids = make(map[t.PlayerID]Bid)
	g.state.HandsWon = make(map[t.PlayerID]int)
	g.state.Table = make(map[t.PlayerID]*Card)
	g.cardstack = make([]Card, 0)
	g.state.TurnPlayer = starter
	g.cycler.StartFrom(starter)

	for _, player := range g.Players {
		g.sendCardsToPlayer(player)
	}
}

// finishRound scores the round that just ended and either starts the next
// one or ends the game.
func (g *Game) finishRound() {
	g.scoreRound()
	g.state.Round++

	if int(g.state.Round) >= len(g.params.roundSchedule) {
		g.changeState(GameDone)
		g.broadcastGameState()
		return
	}

	g.startNewRound()
	g.changeState(PlayingContinue)
	g.broadcastGameState()
}
