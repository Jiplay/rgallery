package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/robbymilo/rgallery/pkg/notify"
)

type Notice = notify.Notice

// ServePoll handles requests for new events.
func ServePoll(w http.ResponseWriter, r *http.Request) {
	// Set a deadline (e.g., 30 seconds from now)
	ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
	defer cancel()

	ch := notify.AddSubscriber()
	defer notify.RemoveSubscriber(ch)

	select {
	case msg := <-ch:
		fmt.Fprintf(w, "%s", msg)
	case <-ctx.Done():
		messageJSON, _ := json.Marshal(Notice{Message: "none", Status: "ok"})
		fmt.Fprint(w, string(messageJSON))
	}
}
