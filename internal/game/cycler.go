package game

import (
	"errors"

	t "github.com/B33Boy/Judgement/internal/types"
)

type PlayerCycler struct {
	keys       []t.PlayerID
	index      int
	startIndex int
	started    bool
}

func NewPlayerCycler(m PlayerMap) *PlayerCycler {
	keys := make([]t.PlayerID, 0, len(m))

	for playerID := range m {
		keys = append(keys, playerID)
	}
	return &PlayerCycler{
		keys:       keys,
		index:      0,
		startIndex: 0,
		started:    false,
	}
}

func (pc *PlayerCycler) StartFrom(player t.PlayerID) error {
	for i, id := range pc.keys {
		if id == player {
			pc.index = i
			pc.startIndex = i
			pc.started = false
			return nil
		}
	}
	return errors.New("player not found")
}

func (pc *PlayerCycler) Next() (t.PlayerID, error) {
	if len(pc.keys) == 0 {
		return "", errors.New("0 players to cycle through")
	}

	// 1. Increment the index first to move to the next person
	pc.index = (pc.index + 1) % len(pc.keys)

	// 2. Grab the ID at the new index
	playerID := pc.keys[pc.index]
	pc.started = true

	return playerID, nil
}

func (pc *PlayerCycler) CompletedCycle() bool {
	if len(pc.keys) == 0 || !pc.started {
		return false
	}
	return pc.index == pc.startIndex
}
