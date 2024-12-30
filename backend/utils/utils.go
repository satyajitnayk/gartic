package utils

import (
	"math/rand"
	"time"

	"github.com/satyajitnayk/gartic/config"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRoomID() string {
	rand.Seed(time.Now().UnixNano())
	id := make([]rune, config.RoomIDLength)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}
