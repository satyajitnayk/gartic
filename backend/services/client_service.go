package services

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/satyajitnayk/gartic/models"
)

var (
	clients      = make(map[*websocket.Conn]*models.Client)
	clientsMutex sync.Mutex
)

// Retrieve all registered clients
func GetClients() []*models.Client {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	result := []*models.Client{}
	for _, client := range clients {
		result = append(result, client)
	}
	return result
}

// Add a client to the global client registry
func AddClient(client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[client.Conn] = client
}

// Remove a client from the global client registry
func RemoveClient(client *models.Client) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, client.Conn)
}
