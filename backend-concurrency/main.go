package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/instagram/concurrency-service/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow cross-origin requests for web client integration
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := ":8081"
	h := hub.NewHub()
	go h.Run()

	// WebSocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = fmt.Sprintf("guest_%d", time.Now().UnixNano()%100000)
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade websocket: %v", err)
			return
		}

		client := &hub.Client{
			Hub:    h,
			Conn:   conn,
			UserID: userID,
			Send:   make(chan []byte, 256),
		}

		h.Register <- client

		// Spin up reader and writer in independent goroutines
		go client.WritePump()
		go client.ReadPump()
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "UP",
			"service":   "instagram-concurrency-service",
			"timestamp": time.Now().UTC(),
		})
	})

	// Live concurrency metrics
	http.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stats := map[string]interface{}{
			"active_users":       h.ActiveUsersCount(),
			"active_connections": h.ActiveSessionsCount(),
			"goroutines":         runtime.NumGoroutine(),
			"go_version":         runtime.Version(),
			"status":             "healthy",
		}
		json.NewEncoder(w).Encode(stats)
	})

	// Webhook / notification push endpoint (Called by Spring Boot Core or async queues)
	http.HandleFunc("/api/v1/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload hub.MessagePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		payload.Timestamp = time.Now()
		delivered := h.DirectNotify(payload.RecipientID, payload)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"delivered_to_active_ws": delivered,
			"recipient_id":           payload.RecipientID,
		})
	})

	log.Printf("🚀 Instagram Concurrency Service (Go) listening on http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
