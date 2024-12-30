package models

import "sync"

type Room struct {
	ID      string
	Clients map[*Client]struct{}
	Mutex   sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[*Client]struct{}),
	}
}

func (r *Room) AddClient(client *Client) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Clients[client] = struct{}{}
}

func (r *Room) RemoveClient(client *Client) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Clients, client)
}

func (r *Room) IsEmpty() bool {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	return len(r.Clients) == 0
}
