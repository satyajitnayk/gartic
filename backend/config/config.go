package config

import "time"

const (
	HeartbeatInterval = 10 * time.Second // Interval to send ping
	HeartbeatTimeout  = 20 * time.Second // Timeout for pong response
	RoomIDLength      = 8
	CleanupInterval   = 10 * time.Second // Empty room cleanup interval
)

var (
	RoomBaseURL = "http://localhost:8080/join"
)
