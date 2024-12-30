package utils

import "testing"

func TestGetWords(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name        string
		length      int
		expectedLen int
	}{
		{
			name:        "Get 5 words",
			length:      5,
			expectedLen: 5,
		},
		{
			name:        "Get 3 words",
			length:      3,
			expectedLen: 3,
		},
		{
			name:        "Get 0 word",
			length:      0,
			expectedLen: 0,
		},
		{
			name:        "Get -3 word",
			length:      -3,
			expectedLen: 0,
		},
	}

	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := GetWords(tt.length)
			if len(words) != tt.expectedLen {
				t.Errorf("expected length %d, but got %d", tt.expectedLen, len(words))
			}
		})
	}
}
