package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/satyajitnayk/gartic/models"
	"github.com/satyajitnayk/gartic/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade Error: %v", err)
		return
	}
	defer conn.Close()

	roomID := r.URL.Query().Get("room")
	name := r.URL.Query().Get("name")
	if roomID == "" || name == "" {
		conn.WriteJSON(map[string]string{"error": "Room ID and Name are required"})
		return
	}

	client := services.RegisterClient(conn, roomID, name)
	if client == nil {
		conn.WriteJSON(map[string]string{"error": "Unable to join room"})
		return
	}

	go services.StartHeartbeat(client)

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read Error: %v", err)
			services.DeregisterClient(client)
			break
		}
		services.HandleClientMessage(client, msg)
	}
}

func ListClients(w http.ResponseWriter, _ *http.Request) {
	clients := services.GetClients()
	clientDetails := []map[string]interface{}{}

	for _, client := range clients {
		clientDetails = append(clientDetails, map[string]interface{}{
			"roomID":  client.RoomID,
			"name":    client.Name,
			"address": client.Conn.RemoteAddr().String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clientDetails); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
