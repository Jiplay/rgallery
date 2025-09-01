package notify

import (
	"encoding/json"
	"sync"
)

var (
	subscribers = make(map[chan string]struct{})
	mu          sync.Mutex
)

type Notice struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Function to add a subscriber (client waiting for an event)
func AddSubscriber() chan string {

	ch := make(chan string, 1)
	mu.Lock()
	subscribers[ch] = struct{}{}
	mu.Unlock()
	return ch
}

// Function to remove a subscriber
func RemoveSubscriber(ch chan string) {
	mu.Lock()
	_, exists := subscribers[ch]
	if exists {
		delete(subscribers, ch)
		close(ch) // Close only if it exists
	}
	mu.Unlock()
}

// Function to notify all subscribers with a message
func NotifySubscribers(message string, status string) {
	mu.Lock()
	messageJSON, _ := json.Marshal(Notice{Message: message, Status: status})
	for ch := range subscribers {
		select {
		case ch <- string(messageJSON):
		default:
		}
		delete(subscribers, ch)
		close(ch)
	}
	mu.Unlock()
}
