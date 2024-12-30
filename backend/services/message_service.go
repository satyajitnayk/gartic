package services

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satyajitnayk/gartic/config"
	"github.com/satyajitnayk/gartic/models"
)

var (
	broadcast  = make(chan models.Message)
	rooms      = make(map[string]*models.Room)
	roomsMutex sync.Mutex
)

// Retrieve the room by ID
func getRoom(roomID string) *models.Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	return rooms[roomID]
}

// Register a client in the specified room
func RegisterClient(conn *websocket.Conn, roomID, name string) *models.Client {
	roomsMutex.Lock()
	room, exists := rooms[roomID]
	roomsMutex.Unlock()

	if !exists {
		return nil
	}

	client := &models.Client{
		Conn:   conn,
		RoomID: roomID,
		Name:   name,
	}

	room.AddClient(client)
	return client
}

// Deregister a client and clean up resources
func DeregisterClient(client *models.Client) {
	room := getRoom(client.RoomID)
	if room != nil {
		room.RemoveClient(client)
	}

	client.Conn.Close()
	log.Printf("Client %s removed from room %s", client.Name, client.RoomID)
}

// Start the heartbeat mechanism for a client
func StartHeartbeat(client *models.Client) {
	ticker := time.NewTicker(config.HeartbeatInterval)
	defer ticker.Stop()

	client.Conn.SetReadDeadline(time.Now().Add(config.HeartbeatTimeout))
	client.Conn.SetPongHandler(func(appData string) error {
		client.Conn.SetReadDeadline(time.Now().Add(config.HeartbeatTimeout))
		return nil
	})

	for range ticker.C {
		if err := client.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(config.HeartbeatInterval)); err != nil {
			log.Printf("Ping failed for client %s: %v", client.Name, err)
			DeregisterClient(client)
			break
		}
	}
	log.Printf("Heartbeat started for client %s", client.Name)
}

func HandleClientMessage(client *models.Client, msg models.Message) {
	msg.RoomID = client.RoomID
	msg.Sender = client.Name
	broadcast <- msg
}

func HandleMessages() {
	for msg := range broadcast {
		room := getRoom(msg.RoomID)
		if room == nil {
			continue
		}

		for client := range room.Clients {
			if err := client.Send(msg); err != nil {
				log.Printf("Broadcast Error: %v", err)
				DeregisterClient(client)
			}
		}
	}
}
