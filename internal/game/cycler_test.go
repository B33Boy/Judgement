package game

import (
	"testing"
)

func TestPlayerCycler(t *testing.T) {
	// Dummy players
	players := PlayerMap{
		"a": &GamePlayer{ID: "a"},
		"b": &GamePlayer{ID: "b"},
		"c": &GamePlayer{ID: "c"},
	}

	cycler := NewPlayerCycler(players)

	// Start from a specific player
	err := cycler.StartFrom("b")
	if err != nil {
		t.Fatalf("StartFrom failed: %v", err)
	}

	// Check that completed cycle is false initially
	if cycler.CompletedCycle() {
		t.Errorf("Expected CompletedCycle false initially")
	}

	// Collect player order
	order := []string{}
	for i := 0; i < len(players); i++ {
		id, err := cycler.Next()
		if err != nil {
			t.Fatalf("Next() failed: %v", err)
		}
		order = append(order, string(id))
	}

	// After one full round, CompletedCycle should be true
	if !cycler.CompletedCycle() {
		t.Errorf("Expected CompletedCycle true after full round")
	}

	// Optional: Check order (it should start from the next player after "b")
	expectedOrder := []string{"c", "a", "b"} // Because Next() increments index first
	for i := range order {
		if order[i] != expectedOrder[i] {
			t.Errorf("Order mismatch at %d: expected %s, got %s", i, expectedOrder[i], order[i])
		}
	}

	// Test StartFrom with invalid player
	err = cycler.StartFrom("x")
	if err == nil {
		t.Errorf("Expected error when starting from non-existent player")
	}
}
