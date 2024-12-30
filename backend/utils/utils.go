package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
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

func GetWords(wordCount int) []string {
	if wordCount <= 0 {
		return []string{}
	}
	data, err := os.ReadFile("../data/words.json")
	if err != nil {
		log.Fatalf("Error reading the JSON file: %v", err)
	}

	var words []string
	err = json.Unmarshal(data, &words)
	if err != nil {
		log.Fatalf("Error unmarshalling the JSON data: %v", err)
	}
	var result []string
	var count = 0
	for _, word := range words {
		if count >= wordCount {
			break
		}
		result = append(result, word)
		count++
	}
	return result
}
