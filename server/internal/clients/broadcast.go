package clients

import (
	"context"
	"fmt"
)

func (cls *Clients) Broadcast(ctx context.Context, message string) (bool, error) {
	var err error
	for uuid, client := range cls.list {
		err = client.Conn.Write(ctx, 1, []byte(message))
		if err != nil {
			return false, fmt.Errorf("[error] id=%v error=%w", uuid, err)
		}
	}

	return true, nil
}
