package game

import "testing"

// Example states and events
const (
	StateA State = "A"
	StateB State = "B"
	StateC State = "C"

	EventX Event = "X"
	EventY Event = "Y"
)

func TestStateMachine_Transitions(t *testing.T) {
	sm := NewStateMachine(StateA)

	// Define transitions
	sm.AddTransition(StateA, EventX, StateB)
	sm.AddTransition(StateB, EventY, StateC)

	// 1. Trigger EventX from StateA → should go to StateB
	next, err := sm.Trigger(EventX)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next != StateB {
		t.Errorf("expected StateB, got %s", next)
	}

	// 2. Trigger EventY from StateB → should go to StateC
	next, err = sm.Trigger(EventY)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next != StateC {
		t.Errorf("expected StateC, got %s", next)
	}

	// 3. Trigger invalid event from StateC → should return error and state stays the same
	next, err = sm.Trigger(EventX)
	if err == nil {
		t.Errorf("expected error for invalid transition, got nil")
	}
	if next != StateC {
		t.Errorf("state should not change on invalid transition, got %s", next)
	}
}
