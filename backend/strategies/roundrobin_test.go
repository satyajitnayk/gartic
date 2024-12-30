package strategies

import (
	"fmt"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	tests := []struct {
		numParticipants int
		totalTurns      int
		expectedIndexes []int
	}{
		{
			numParticipants: 3,
			totalTurns:      3,
			expectedIndexes: []int{0, 1, 2, 0, 1, 2, 0, 1, 2},
		},
		{
			numParticipants: 2,
			totalTurns:      4,
			expectedIndexes: []int{0, 1, 0, 1, 0, 1, 0, 1},
		},
		{
			numParticipants: 5,
			totalTurns:      1,
			expectedIndexes: []int{0, 1, 2, 3, 4},
		},
		{
			numParticipants: 1,
			totalTurns:      2,
			expectedIndexes: []int{0, 0},
		},
		{
			numParticipants: 4,
			totalTurns:      0,
			expectedIndexes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("RoundRobin_%d_%d", tt.numParticipants, tt.totalTurns), func(t *testing.T) {
			roundRobin := NewRoundRobin(tt.numParticipants, tt.totalTurns)

			var actualIndexes []int
			for i := 0; i < tt.totalTurns*tt.numParticipants; i++ {
				index := roundRobin.GetCurrentIndex()
				if index == -1 {
					break
				}
				actualIndexes = append(actualIndexes, index)
				roundRobin.GetNextIndex()
			}

			if !equal(actualIndexes, tt.expectedIndexes) {
				t.Errorf("expected %v, got %v", tt.expectedIndexes, actualIndexes)
			}
		})
	}

	t.Run("AfterRoundRobinFinished", func(t *testing.T) {
		roundRobin := NewRoundRobin(3, 3)
		turns := []int{roundRobin.GetCurrentIndex()}
		for i := 0; i < 9; i++ {
			turns = append(turns, roundRobin.GetNextIndex())
		}

		// After all turns, GetNextIndex should return -1
		if index := roundRobin.GetNextIndex(); index != -1 {
			t.Errorf("expected -1, got %d after all turns completed", index)
		}
	})
}

// Helper function to compare slices
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
