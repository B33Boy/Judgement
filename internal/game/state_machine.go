package game

import (
	"fmt"
)

type StateMachine struct {
	state       State
	transitions map[State]map[Event]State
}

func NewStateMachine(initial State) *StateMachine {
	return &StateMachine{
		state:       initial,
		transitions: make(map[State]map[Event]State),
	}
}

func (sm *StateMachine) AddTransition(from State, event Event, to State) {
	if sm.transitions[from] == nil {
		sm.transitions[from] = make(map[Event]State)
	}
	sm.transitions[from][event] = to
}

func (sm *StateMachine) Trigger(event Event) (State, error) {
	next, ok := sm.transitions[sm.state][event]

	if !ok {
		return sm.state, fmt.Errorf("invalid transition: %s + %s", sm.state, event)
	}

	sm.state = next
	return next, nil
}
