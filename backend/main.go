package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Conn *websocket.Conn
	Room string
	Name string
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Room    string `json:"room"`
}

var clientsMutex sync.Mutex
var clients = make(map[*websocket.Conn]*Client)
var broadcast = make(chan Message)

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/clients", listClients)

	go handleMessages()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	query := r.URL.Query()
	room := query.Get("room")
	name := query.Get("name")

	client := &Client{Conn: conn, Room: room, Name: name}
	clientsMutex.Lock()
	clients[conn] = client
	clientsMutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for conn, client := range clients {
			if client.Room == msg.Room {
				err := conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					conn.Close()
					delete(clients, conn)
				}
			}
		}
	}
}

func listClients(w http.ResponseWriter, _ *http.Request) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clientDetails := []map[string]interface{}{}
	for _, client := range clients {
		clientDetails = append(clientDetails, map[string]interface{}{
			"room":    client.Room,
			"name":    client.Name,
			"address": client.Conn.RemoteAddr().String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clientDetails); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
