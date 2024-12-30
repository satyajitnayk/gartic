package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/satyajitnayk/gartic/models"
	"github.com/satyajitnayk/gartic/utils"

	"github.com/satyajitnayk/gartic/config"
)

var rooms = make(map[string]*models.Room)

func HandleRoomCreation(w http.ResponseWriter, _ *http.Request) {
	roomID := utils.GenerateRoomID()
	room := models.NewRoom(roomID)

	rooms[roomID] = room

	inviteLink := fmt.Sprintf("%s?room=%s", config.RoomBaseURL, roomID)
	response := map[string]string{"roomID": roomID, "inviteLink": inviteLink}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleJoiningRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	if _, exists := rooms[roomID]; !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Connect to WebSocket endpoint with the room ID"))
}
