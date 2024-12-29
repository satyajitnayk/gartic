package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	Sender  string `json:"sender"`
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

	// If room or name is not provided, close the connection
	if room == "" || name == "" {
		conn.Close()
		return
	}

	client := &Client{Conn: conn, Room: room, Name: name}

	clientsMutex.Lock()
	clients[conn] = client
	clientsMutex.Unlock()

	// Start heartbeat mechanism
	go startHeartBeat(conn)

	welcomeMsg := Message{
		Type:    "welcome",
		Content: fmt.Sprintf("Welcome %s to room: %s", client.Name, client.Room),
		Room:    client.Room,
	}
	conn.WriteJSON(welcomeMsg)

	// Listen for incoming messages
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}
		appendMetaAndBroadcast(conn, msg)
	}
}

func sendMessageToRoom(msg Message) {
	// Send the message to all clients in the same room
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn, client := range clients {
		if client.Room == msg.Room {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				deregisterClient(conn)
			}
		}
	}
}

func handleMessages() {
	for msg := range broadcast {
		sendMessageToRoom(msg)
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

const (
	heartbeatInterval = 10 // Send a ping every 10 seconds
	heartbeatTimeout  = 20 // Disconnect if no pong received within 20 seconds
)

func registerClient(conn *websocket.Conn, r *http.Request) *Client {
	query := r.URL.Query()
	room, name := query.Get("room"), query.Get("name")

	if room == "" || name == "" {
		conn.Close()
		return nil
	}

	client := &Client{Conn: conn, Room: room, Name: name}
	clientsMutex.Lock()
	clients[conn] = client
	clientsMutex.Unlock()

	welcomeMsg := Message{
		Type:    "welcome",
		Content: fmt.Sprintf("Welcome %s to room: %s", client.Name, client.Room),
		Room:    client.Room,
	}
	conn.WriteJSON(welcomeMsg)

	return client
}

func deregisterClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
	conn.Close()
}

func appendMetaAndBroadcast(conn *websocket.Conn, msg Message) {
	// Add metadata (room and username) to the message
	clientsMutex.Lock()
	sender := clients[conn]
	if sender != nil {
		msg.Room = sender.Room
		msg.Sender = sender.Name
	}
	clientsMutex.Unlock()

	// Broadcast the message
	broadcast <- msg
}

func startHeartBeat(conn *websocket.Conn) {
	ticker := time.NewTicker(heartbeatInterval * time.Second)
	defer ticker.Stop()

	// Set read deadline and pong handler
	conn.SetReadDeadline(time.Now().Add(heartbeatTimeout * time.Second))
	conn.SetPongHandler(func(appData string) error {
		fmt.Println("Pong received")
		conn.SetReadDeadline(time.Now().Add(heartbeatTimeout * time.Second)) // Reset deadline
		return nil
	})

	for range ticker.C {
		err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
		if err != nil {
			log.Printf("Heartbeat failed: %v", err)
			deregisterClient(conn)
			break
		}
		fmt.Println("ping sent")
	}
}
