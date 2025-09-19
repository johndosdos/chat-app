package chat

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/johndosdos/chat-app/internal/database"
)

type Message struct {
	Content string `json:"content"`
	From    uuid.UUID
}

func DbLoadChatHistory(ctx context.Context, recv chan Message, db *database.Queries) {
	// Send the last 50 messages to the client on new connection.
	dbMessageList, err := db.ListMessages(ctx)
	if err != nil {
		log.Printf("[error] failed to load messages from database: %v", err)
		return
	}

	// Use the hub's accept channel.
	for _, msg := range dbMessageList {
		recv <- Message{
			From:    msg.UserID.Bytes,
			Content: msg.Content,
		}
	}
}
