package chat

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johndosdos/chat-app/server/internal/database"
)

type Message struct {
	Content []byte
	From    uuid.UUID
}

func DbStoreMessage(ctx context.Context, db database.Queries, recvFromHub chan Message) {
	for {
		select {
		case message := <-recvFromHub:
			_, err := db.CreateMessage(ctx, database.CreateMessageParams{
				UserID:  pgtype.UUID{Bytes: [16]byte(message.From), Valid: true},
				Content: string(message.Content),
			})
			if err != nil {
				log.Printf("[DB error] failed to store message to database: %v", err)
				return
			}
		}
	}
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
			Content: []byte(msg.Content),
		}
	}
}
