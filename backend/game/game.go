package game

import (
	"fmt"
	"log"
	"time"

	"github.com/satyajitnayk/gartic/models"
	"github.com/satyajitnayk/gartic/strategies"
	"github.com/satyajitnayk/gartic/utils"
)

// Game defines the structure of a game.
type Game struct {
	Host          *models.Client
	Players       []*models.Client
	RoundRobin    *strategies.RoundRobin
	Words         []string
	IsStarted     bool
	CurrentTurn   int
	CurrentPlayer *models.Client // Track current player whose turn it is
	TimePerTurn   time.Duration  // Time allowed per player's turn
}

// NewGame initializes a new game with a host and players.
func NewGame(host *models.Client, players []*models.Client, totalTurns int, timePerTurn time.Duration) *Game {
	roundRobin := strategies.NewRoundRobin(len(players), totalTurns)
	wordCounts := len(players) * totalTurns
	words := utils.GetWords(wordCounts) // Ensure utils.GetWords returns the correct number of words
	return &Game{
		Host:        host,
		Players:     players,
		RoundRobin:  roundRobin,
		IsStarted:   false,
		CurrentTurn: 0,
		Words:       words,
		TimePerTurn: timePerTurn, // Store the time per turn here
	}
}

// Start starts the game, setting it to the started state and notifying all players.
func (g *Game) Start() {
	if g.IsStarted {
		log.Println("Game has already started.")
		return
	}
	g.IsStarted = true
	g.notifyPlayers("Game started! Get ready to play!")
}

func (g *Game) StartTurn() {
	// Get the current player's index
	playerIndex := g.RoundRobin.GetCurrentIndex()
	if playerIndex == -1 {
		g.EndGame()
		return
	}

	g.CurrentPlayer = g.Players[playerIndex]
	g.notifyPlayers(fmt.Sprintf("It's %s's turn to guess!", g.CurrentPlayer.Name))

	// Set up a timer for the time per turn
	go g.handleTurnTimeout()
}

func (g *Game) handleTurnTimeout() {
	time.Sleep(g.TimePerTurn)
	// If the game is still running, move to the next player
	if g.IsStarted {
		g.RoundRobin.GetNextIndex()
		g.StartTurn() // Start the next player's turn
	}
}

// notifyPlayers sends a message to all players.
func (g *Game) notifyPlayers(content string) {
	for _, player := range g.Players {
		if player.Conn != nil {
			err := player.Conn.WriteJSON(models.Message{
				Type:    "game_info",
				Content: content,
				RoomID:  player.RoomID,
				Sender:  "system",
			})
			if err != nil {
				log.Printf("Error sending message to player %s: %v", player.Name, err)
			}
		}
	}
}

// NextTurn advances to the next player’s turn in the round-robin cycle.
func (g *Game) NextTurn() {
	g.CurrentTurn++
	if g.CurrentTurn >= len(g.Players)*g.RoundRobin.TotalTurns {
		g.EndGame()
	} else {
		playerIndex := g.RoundRobin.GetCurrentIndex()
		if playerIndex != -1 {
			player := g.Players[playerIndex]
			g.notifyPlayers(fmt.Sprintf("It's %s's turn to guess!", player.Name))
			g.RoundRobin.GetNextIndex()
		}
	}
}

// EndGame ends the game and notifies players.
func (g *Game) EndGame() {
	g.notifyPlayers("Game Over! The game has ended.")

	// Clear the game state or reset values as necessary
	g.IsStarted = false
	g.CurrentTurn = 0
	g.CurrentPlayer = nil
	g.RoundRobin = nil // Reset round-robin so the game can be restarted if necessary
	g.Players = nil    // Clear players (or you can keep them if you want to restart with the same players)
	g.Words = nil
}
