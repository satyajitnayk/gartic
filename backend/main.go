package main

import (
	"log"
	"net/http"

	"github.com/satyajitnayk/gartic/config"
	"github.com/satyajitnayk/gartic/handlers"
	"github.com/satyajitnayk/gartic/services"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocketConnections)
	http.HandleFunc("/create-room", handlers.HandleRoomCreation)
	http.HandleFunc("/join", handlers.HandleJoiningRoom)
	http.HandleFunc("/clients", handlers.ListClients)

	go services.HandleMessages()
	go services.CleanEmptyRooms(config.CleanupInterval)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
