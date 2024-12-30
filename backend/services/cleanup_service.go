package services

import (
	"log"
	"time"
)

func CleanEmptyRooms(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for roomID, room := range rooms {
			if room.IsEmpty() {
				log.Printf("Cleaning empty room: %s", roomID)
				delete(rooms, roomID)
			}
		}
	}
}
