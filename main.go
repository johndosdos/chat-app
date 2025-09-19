package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	components "github.com/johndosdos/chat-app/server/components/chat"
	"github.com/johndosdos/chat-app/server/internal/database"

	"github.com/johndosdos/chat-app/server/internal/handler"
	ws "github.com/johndosdos/chat-app/server/internal/websocket"
)

var (
	dbConn    *pgx.Conn
	dbQueries *database.Queries
)

type message struct {
	UserID  uuid.UUID
	Content string `json:"content"`
}

func main() {
	port := ":8080"
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

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

	// client hub init start
	// hub.Run is our central hub that is always listening for client related
	// events.
	hub := ws.NewHub()
	go hub.Run(ctx, dbQueries)
	// client hub init end

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", templ.Handler(components.Base()))
	http.HandleFunc("/ws", handler.ServeWs(ctx, hub, dbQueries))

	log.Println("Server starting at port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
