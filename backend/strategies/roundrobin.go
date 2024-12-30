package strategies

import (
	"sync"
)

// RoundRobin defines the round-robin algorithm for managing turns.
type RoundRobin struct {
	NumParticipants int // Number of participants (not actual players)
	CurrentIndex    int
	TotalTurns      int
	TurnsTaken      int
	Lock            sync.Mutex
}

func NewRoundRobin(numParticipants, totalTurns int) *RoundRobin {
	return &RoundRobin{
		NumParticipants: numParticipants,
		TotalTurns:      totalTurns,
		CurrentIndex:    0,
		TurnsTaken:      0,
	}
}

func (r *RoundRobin) GetCurrentIndex() int {
	return r.CurrentIndex
}

func (r *RoundRobin) GetNextIndex() int {
	r.Lock.Lock()
	defer r.Lock.Unlock()

	// If we've completed all the turns, return -1
	if r.TurnsTaken >= r.TotalTurns {
		return -1
	}

	// Move to the next participant
	r.CurrentIndex = (r.CurrentIndex + 1) % r.NumParticipants

	// If we've completed one full cycle
	// (all participants have had their turn), increment the TurnsTaken
	if r.CurrentIndex == 0 {
		r.TurnsTaken++
	}
	return r.CurrentIndex
}
