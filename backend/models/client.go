package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	RoomID string
	Name   string
}

func (c *Client) Send(message interface{}) error {
	return c.Conn.WriteJSON(message)
}
