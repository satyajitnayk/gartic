package models

type Message struct {
	Type    string `json:"type"`    // Type of the message (e.g., "chat", "system", "welcome")
	Content string `json:"content"` // Content of the message
	RoomID  string `json:"room"`    // The room ID associated with the message
	Sender  string `json:"sender"`  // Sender of the message (user's name)
}
