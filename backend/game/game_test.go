package game

import (
	"testing"
	"time"

	"github.com/satyajitnayk/gartic/models"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	// Prepare sample data
	host := &models.Client{Name: "Host"}
	players := []*models.Client{
		{Name: "Player1"},
		{Name: "Player2"},
	}
	totalTurns := 3
	timePerTurn := 30 * time.Second

	// Create the game
	game := NewGame(host, players, totalTurns, timePerTurn)

	// Test if game is initialized correctly
	assert.NotNil(t, game)
	assert.Equal(t, len(players), len(game.Players))
	assert.Equal(t, totalTurns, game.RoundRobin.TotalTurns)
	assert.Equal(t, 0, game.CurrentTurn)
	assert.Equal(t, 0, game.RoundRobin.TurnsTaken)
	assert.NotNil(t, game.RoundRobin)
}

func TestStartGame(t *testing.T) {
	host := &models.Client{Name: "Host"}
	players := []*models.Client{
		{Name: "Player1"},
		{Name: "Player2"},
	}
	totalTurns := 3
	timePerTurn := 30 * time.Second

	game := NewGame(host, players, totalTurns, timePerTurn)

	// Ensure the game is not started initially
	assert.False(t, game.IsStarted)

	// Start the game
	game.Start()

	// Check if the game has started
	assert.True(t, game.IsStarted)
	assert.NotNil(t, game.RoundRobin)
}

func TestNextTurn(t *testing.T) {
	host := &models.Client{Name: "Host"}
	players := []*models.Client{
		{Name: "Player1"},
		{Name: "Player2"},
	}
	totalTurns := 2
	timePerTurn := 30 * time.Second

	game := NewGame(host, players, totalTurns, timePerTurn)
	game.Start() // Start the game

	// Ensure the first turn is for the host (assuming random start)
	assert.Equal(t, game.RoundRobin.CurrentIndex, 0)

	// Move to next turn
	game.NextTurn()

	// Check if the current player is the next player in the round-robin
	assert.Equal(t, game.RoundRobin.CurrentIndex, 1)

	// Move again
	game.NextTurn()
	assert.Equal(t, game.RoundRobin.CurrentIndex, 0) // Loop back to player 1
}

func TestEndGame(t *testing.T) {
	host := &models.Client{Name: "Host"}
	players := []*models.Client{
		{Name: "Player1"},
		{Name: "Player2"},
	}
	totalTurns := 1
	timePerTurn := 30 * time.Second

	game := NewGame(host, players, totalTurns, timePerTurn)
	game.Start()

	// End the game after 1 turn
	game.EndGame()

	// Ensure game state is reset
	assert.False(t, game.IsStarted)
	assert.Equal(t, 0, game.CurrentTurn)
	assert.Nil(t, game.Players)
	assert.Nil(t, game.Words)
}
