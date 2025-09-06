package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johndosdos/chat-app/server/internal/database"
	"github.com/microcosm-cc/bluemonday"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]uuid.UUID)
	broadcast = make(chan message)
	mu        = sync.Mutex{}

	dbConn    *pgx.Conn
	dbQueries *database.Queries
)

type message struct {
	senderID uuid.UUID
	content  []byte
}

func main() {
	port := ":8080"
	ctx := context.Background()

	// database init start
	var err error
	dbURL := os.Getenv("DB_URL")
	dbConn, err = pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Printf("[Error] cannot connect to postgresql database: %v", err)
		return
	}
	defer dbConn.Close(ctx)

	dbQueries = database.New(dbConn)
	// database init end

	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/ws", handleConnection)

	go handleMessages()

	log.Println("Server starting at port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[error] failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// We need to prevent race conditions when multiple clients are connecting.
	// to our server.
	userID := uuid.New()
	mu.Lock()
	clients[ws] = userID
	mu.Unlock()

	// Send the last 50 messages to the client on new connection.
	msgList, err := dbQueries.ListMessages(context.Background())
	if err != nil {
		log.Printf("[error] failed to load messages from database: %v", err)
		return
	} else {
		for _, msg := range msgList {
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Content))
			if err != nil {
				log.Printf("[Error] failed to load messages to client")
				break
			}
		}
	}

	// Process received data for broadcasting to clients.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			// If the client was disconnected.
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			return
		}

		// Send data to the broadcast channel.
		broadcast <- message{
			senderID: userID,
			content:  p,
		}
	}
}

func handleMessages() {
	p := bluemonday.StrictPolicy()

	for {
		msg := <-broadcast

		// Sanitize the message before broadcasting.
		sanitized := p.SanitizeBytes(msg.content)

		// Create message entry to database.
		_, err := dbQueries.CreateMessage(context.Background(), database.CreateMessageParams{
			UserID: pgtype.UUID{
				Bytes: msg.senderID,
				Valid: true,
			},
			Content: string(sanitized),
		})
		if err != nil {
			log.Printf("[error] failed to create entry to database: %v", err)
			break
		}

		// Acquire a lock to prevent race conditions.
		mu.Lock()
		for conn, clientID := range clients {
			if msg.senderID != clientID {
				err := conn.WriteMessage(websocket.TextMessage, sanitized)
				if err != nil {
					log.Printf("[Error] failed to write data to client")
					delete(clients, conn)
				}
			}
		}
		mu.Unlock()
	}
}
